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
	WillTopic               string
	WillPayload             string
	WillQoS                 int
	WillRetain              bool
	ClientID                string
	OnConnectHandler        func(mq.Client)
	OnConnectionLostHandler func(mq.Client, error)
}

func NewPersistantMqtt(config ClientConfig) (mqttClient mq.Client, err error) {
	useTls := false

	if val, ok := os.LookupEnv("MQTT_TLS"); ok {
		useTls = val != "0" && val != "" && strings.ToLower(val) != "false"
	}
	if *MqttTLSFlag {
		useTls = true
	}

	caPath := envOrFlagStr(*MqttCAFlag, "MQTT_CA_PATH", "")
	certPath := envOrFlagStr(*MqttCertFlag, "MQTT_CERT_PATH", "")
	keyPath := envOrFlagStr(*MqttKeyFlag, "MQTT_KEY_PATH", "")
	address := envOrFlagStr(*MqttAddressFlag, "MQTT_ADDRESS", "localhost:1883")
	username := envOrFlagStr(*MqttUsernameFlag, "MQTT_USERNAME", "")
	password := envOrFlagStr(*MqttPasswordFlag, "MQTT_PASSWORD", "")
	connectionTimeout := envOrFlagInt(*MqttConnectionTimeout, "MQTT_CONNECTION_TIMEOUT", 10)
	keepAlive := envOrFlagInt(*MqttKeepAlive, "MQTT_KEEPALIVE", 5)
	maxReconnectInterval := envOrFlagInt(*MqttMaxReconnectInterval, "MQTT_MAX_RECONNECT_INTERVAL", 2)
	pingTimeout := envOrFlagInt(*MqttPingTimeout, "MQTT_PING_TIMEOUT", 10)
	writeTimeout := envOrFlagInt(*MqttWriteTimeout, "MQTT_WRITE_TIMEOUT", 5)

	var tlsCfg *tls.Config
	if useTls {
		tlsCfg, err = setupTLS(caPath, certPath, keyPath)
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

	opts := mq.NewClientOptions().
		AddBroker(fmt.Sprintf("tcp://%s", address)).
		SetClientID(clientId).
		SetConnectTimeout(time.Duration(connectionTimeout) * time.Second).
		SetKeepAlive(time.Duration(keepAlive) * time.Second).
		SetMaxReconnectInterval(time.Duration(maxReconnectInterval) * time.Minute).
		SetMessageChannelDepth(100).
		SetPingTimeout(time.Duration(pingTimeout) * time.Second).
		SetProtocolVersion(4).
		SetWriteTimeout(time.Duration(writeTimeout) * time.Second)

	if config.OnConnectHandler != nil {
		opts.SetOnConnectHandler(config.OnConnectHandler)
	}

	if config.OnConnectionLostHandler != nil {
		opts.SetConnectionLostHandler(config.OnConnectionLostHandler)
	}

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

func setupTLS(caPath, certPath, keyPath string) (*tls.Config, error) {
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
