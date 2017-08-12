package flagmqtt

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	mq "github.com/eclipse/paho.mqtt.golang"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type ClientConfig struct {
	WillTopic            string
	WillPayload          string
	WillQoS              int
	WillRetain           bool
	ClientID             string
	ConnectTimeout       time.Duration
	KeepAlive            time.Duration
	MaxReconnectInterval time.Duration
	PingTimeout          time.Duration
	WriteTimeout         time.Duration
}

func NewPersistantMqtt(config ClientConfig) (mqttClient mq.Client, err error) {
	useTls := false

	if val, ok := os.LookupEnv("MQTT_TLS"); ok {
		useTls = val != "0" && val != "" && strings.ToLower(val) != "false"
	}
	if *mqttTLSFlag {
		useTls = true
	}

	caPath := envOrFlagStr(*mqttCAFlag, "MQTT_CA_PATH", "")
	certPath := envOrFlagStr(*mqttCertFlag, "MQTT_CERT_PATH", "")
	keyPath := envOrFlagStr(*mqttKeyFlag, "MQTT_KEY_PATH", "")
	address := envOrFlagStr(*mqttAddressFlag, "MQTT_ADDRESS", "localhost:1883")
	username := envOrFlagStr(*mqttUsernameFlag, "MQTT_USERNAME", "")
	password := envOrFlagStr(*mqttPasswordFlag, "MQTT_PASSWORD", "")

	var tlsCfg *tls.Config
	if useTls {
		tlsCfg, err = setupTls(caPath, certPath, keyPath)
		if err != nil {
			return
		}
	}

	clientId := config.ClientID

	if clientId == "" {
		var rndUuid uuid.UUID
		rndUuid = uuid.NewV4()
		clientId = rndUuid.String()
	}

	if config.ConnectTimeout <= 0 {
		config.ConnectTimeout = 30 * time.Second
	}
	if config.MaxReconnectInterval <= 0 {
		config.MaxReconnectInterval = 1 * time.Minute
	}
	if config.KeepAlive <= 0 {
		config.KeepAlive = 1 * time.Minute
	}
	if config.PingTimeout <= 0 {
		config.PingTimeout = 30 * time.Second
	}
	if config.WriteTimeout <= 0 {
		config.WriteTimeout = 30 * time.Second
	}

	opts := mq.NewClientOptions().
		AddBroker(fmt.Sprintf("tcp://%s", address)).
		SetClientID(clientId).
		SetConnectTimeout(config.ConnectTimeout).
		SetKeepAlive(config.KeepAlive).
		SetMaxReconnectInterval(config.MaxReconnectInterval).
		SetMessageChannelDepth(100).
		SetPingTimeout(config.PingTimeout).
		SetProtocolVersion(4).
		SetWriteTimeout(config.WriteTimeout)

	if config.WillTopic != "" {
		opts.SetWill(
			config.WillTopic,
			config.WillPayload,
			byte(config.WillQoS&0xff),
			config.WillRetain,
		)
	}

	if username != "" {
		opts.SetUsername(username)
	}
	if password != "" {
		opts.SetPassword(password)
	}
	if useTls {
		opts.SetTLSConfig(tlsCfg)
	}

	return mq.NewClient(opts), nil
}

func setupTls(caPath, certPath, keyPath string) (*tls.Config, error) {
	tlsCfg := &tls.Config{}
	if caPath != "" {
		caPem, err := ioutil.ReadFile(caPath)
		if err != nil {
			return nil, err
		}
		tlsCfg.RootCAs = x509.NewCertPool()
		tlsCfg.RootCAs.AppendCertsFromPEM(caPem)
	}

	if certPath != "" && keyPath != "" {
		if keyPath == "" {
			return nil, errors.New("Certificate path specified, but key path missing")
		}
		cert, err := tls.LoadX509KeyPair(certPath, keyPath)
		if err != nil {
			return nil, err
		}
		tlsCfg.Certificates = []tls.Certificate{cert}
		tlsCfg.BuildNameToCertificate()
	}

	return tlsCfg, nil
}
