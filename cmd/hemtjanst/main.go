package main

import (
	"flag"
	"fmt"
	"github.com/brutella/hc/accessory"
	"github.com/hemtjanst/hemtjanst/device"
	"github.com/hemtjanst/hemtjanst/homekit"
	"github.com/hemtjanst/hemtjanst/homekit/bridge"
	"github.com/hemtjanst/hemtjanst/messaging"
	"github.com/hemtjanst/hemtjanst/messaging/flagmqtt"
	"github.com/hemtjanst/hemtjanst/web"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var (
	name     = flag.String("name", "hemtjanst", "Name of bridge instance")
	addr     = flag.String("address", "127.0.0.1", "IP or hostname for Hemtjänst to bind on")
	port     = flag.String("port", "12345", "Port for Hemtjänst to bind on")
	pin      = flag.String("pin", "01020304", "Pairing pin for the HomeKit bridge")
	startWeb = flag.Bool("web.ui", false, "Start the built-in web UI")
	wAddr    = flag.String("web.addr", ":8080", "IP/host:port to bind the webinterface to")
	dbPath   = flag.String("db.path", "./db", "Path to store the database with HomeKit key pairs etc.")
)

const (
	announceTopicPrefix = "announce/"
	leaveTopic          = "leave"
	discoverTopic       = "discover"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Parameters:\n\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
	}
	flag.Parse()

	log.Print("Initialing Hemtjänst")
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	announce := make(chan messaging.Message)
	leave := make(chan messaging.Message)
	handlerInit := make(chan bool)
	managerInit := make(chan bool)

	bridgeConfig := bridge.Config{
		Pin:         *pin,
		Port:        *port,
		IP:          *addr,
		StoragePath: *dbPath,
	}
	bridgeInfo := accessory.Info{
		Name:         *name,
		SerialNumber: "12345",
		Manufacturer: "BEDS Inc.",
		Model:        "v0.1",
	}

	log.Print("Attempting to connect to MQTT broker")
	handler := &messaging.Handler{
		Ann:           announce,
		Leave:         leave,
		AnnounceTopic: announceTopicPrefix + "#",
		LeaveTopic:    leaveTopic,
		DiscoverTopic: discoverTopic,
		// Wait 5 seconds before sending discover to allow
		// persistent announcements to be fetched
		DiscoverDelay: 5 * time.Second,
		DiscoverStart: handlerInit,
	}
	cID := flagmqtt.NewUniqueIdentifier()
	conf := flagmqtt.ClientConfig{
		ClientID:                fmt.Sprintf("hemtjanst-%s", cID),
		OnConnectHandler:        handler.OnConnect,
		OnConnectionLostHandler: handler.OnConnectionLost,
		WillTopic:               "leave",
		WillPayload:             fmt.Sprintf("hemtjanst-%s", cID),
		WillRetain:              false,
		WillQoS:                 0,
	}

	c, err := flagmqtt.NewPersistentMqtt(conf)
	if err != nil {
		log.Fatal("Could not configure the MQTT client: ", err)
	}

	go func() {
		if token := c.Connect(); token.Wait() && token.Error() != nil {
			log.Fatal("Failed to establish connection with broker: ", token.Error())
		}
	}()

	hkBridge, err := bridge.NewBridge(bridgeConfig, bridgeInfo)
	if err != nil {
		log.Fatal("Could not start HomeKit bridge: ", err)
	}

	manager := device.NewManager(messaging.NewMQTTMessenger(c), managerInit)
	log.Print("Started device manager")

	hk := homekit.NewHomekit(hkBridge, manager)

	go func() {
		// Wait for handler to have sent its discover
		<-handlerInit

		// Tell manager that we're initialised
		managerInit <- true

		// Wait a few more seconds before starting bridge
		<-time.After(5 * time.Second)
		log.Print("Starting HomeKit bridge")
		manager.AddHandler(hk)
		hkBridge.Start()
	}()

	if *startWeb {
		go func() {
			<-time.After(5 * time.Second)
			web.Serve(manager, *wAddr)
		}()
		log.Print("Started web interface")
	}

loop:
	for {
		select {
		case sig := <-quit:
			log.Printf("Received signal: %s, proceeding to shutdown", sig)
			break loop
		case msg := <-announce:
			newReg := msg.Topic()
			log.Print("New announcement: ", newReg)
			if len(newReg) <= len(announceTopicPrefix) || newReg[0:len(announceTopicPrefix)] != announceTopicPrefix {
				// Announcement doesn't start with prefix
				log.Printf("Ignoring: Announcement doesn't start with " + announceTopicPrefix)
				continue
			}
			newReg = newReg[len(announceTopicPrefix):]
			if !strings.Contains(newReg, "/") {
				// We expect topics we care about to contain at least 1 /
				log.Printf("Ignoring: Malformed topic %s", newReg)
				continue
			}
			if len(msg.Payload()) == 0 {
				// Empty payload, remove device
				go manager.Remove(newReg)
				continue
			}

			go manager.Add(newReg, msg.Payload())
		case msg := <-leave:
			go manager.Leave(string(msg.Payload()))
		}
	}

	// When the MQTT lib is connecting but hasn't establish a conneciton yet
	// the IsConnected() method returns true. As a consequence, b/c it believes
	// it is connected the call to .Disconnect() will panic if we terminate
	// before we've managed to establish a connection to the broker, as it
	// tries to close one of its own channels that are currently still nil.
	//
	// To avoid getting an ugly panic printed for what is arguably a bug in the
	// library defer a recover that does nothing and exit normally.
	defer func() {
		recover()
	}()

	c.Disconnect(250)
	hkBridge.Stop()
	log.Print("Disconnected from broker. Bye!")
	os.Exit(0)
}
