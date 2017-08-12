package flagmqtt

import "flag"

var (
	mqttAddressFlag  = flag.String("mqtt.address", "localhost:1883", "Address to MQTT endpoint")
	mqttUsernameFlag = flag.String("mqtt.username", "", "MQTT Username")
	mqttPasswordFlag = flag.String("mqtt.password", "", "MQTT Password")
	mqttTLSFlag      = flag.Bool("mqtt.tls", false, "Enable TLS")
	mqttCAFlag       = flag.String("mqtt.ca", "", "Path to CA certificate")
	mqttCertFlag     = flag.String("mqtt.cert", "", "Path to Client certificate")
	mqttKeyFlag      = flag.String("mqtt.key", "", "Path to Client certificate key")
)
