package flagmqtt

import "flag"

var (
	MqttAddressFlag          = flag.String("mqtt.address", "localhost:1883", "Address to MQTT endpoint")
	MqttUsernameFlag         = flag.String("mqtt.username", "", "MQTT Username")
	MqttPasswordFlag         = flag.String("mqtt.password", "", "MQTT Password")
	MqttTLSFlag              = flag.Bool("mqtt.tls", false, "Enable TLS")
	MqttCAFlag               = flag.String("mqtt.ca", "", "Path to CA certificate")
	MqttCertFlag             = flag.String("mqtt.cert", "", "Path to Client certificate")
	MqttKeyFlag              = flag.String("mqtt.key", "", "Path to Client certificate key")
	MqttConnectionTimeout    = flag.Int("mqtt.connection-timeout", 10, "Connection timeout in seconds")
	MqttKeepAlive            = flag.Int("mqtt.keepalive", 5, "Time in seconds between each PING packet")
	MqttMaxReconnectInterval = flag.Int("mqtt.max-reconnect-interval", 2, "Maximum time in minutes to wait between reconnect attemps")
	MqttPingTimeout          = flag.Int("mqtt.ping-timeout", 10, "Time in seconds after which a ping times out")
	MqttWriteTimeout         = flag.Int("mqtt.write-timeout", 5, "Time in seconds after which a write will time out")
)
