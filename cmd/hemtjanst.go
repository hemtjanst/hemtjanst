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
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var (
	addr = flag.String("address", "127.0.0.1", "IP or hostname for Hemtjänst to bind on")
	port = flag.String("port", "12345", "Port for Hemtjänst to bind on")
	pin  = flag.String("pin", "01020304", "Pairing pin for the HomeKit bridge")
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
	announce := make(chan []byte)
	leave := make(chan []byte)

	bridgeConfig := bridge.Config{
		Pin:         *pin,
		Port:        *port,
		StoragePath: "./db",
	}
	bridgeInfo := accessory.Info{
		Name:         "Hemtjänst",
		SerialNumber: "12345",
		Manufacturer: "BEDS Inc.",
		Model:        "v0.1",
	}

	log.Print("Attempting to connect to MQTT broker")
	handler := &messaging.Handler{
		Ann: announce,
		Leave: leave,
	}
	conf := flagmqtt.ClientConfig{
		ClientID:                "hemtjanst",
		OnConnectHandler:        handler.OnConnect,
		OnConnectionLostHandler: handler.OnConnectionLost,
		WillTopic:               "leave",
		WillPayload:             "hemtjanst",
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

	manager := device.NewManager(messaging.NewMQTTMessenger(c))
	log.Print("Started device manager")

	hk := homekit.NewHomekit(hkBridge, manager)
	manager.AddHandler(hk)

	go func() {
		<-time.After(2 * time.Second)
		hkBridge.Start()
	}()
	log.Print("Started HomeKit bridge")

loop:
	for {
		select {
		case sig := <-quit:
			log.Printf("Received signal: %s, proceeding to shutdown", sig)
			break loop
		case msg := <-announce:
			newReg := string(msg)
			log.Print("New announcement: ", newReg)
			if !strings.Contains(newReg, "/") {
				// We expect topics we care about to contain at least 1 /
				break
			}
			manager.Add(newReg)
		case msg := <-leave:
			manager.Remove(string(msg))
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
