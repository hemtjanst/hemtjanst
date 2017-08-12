package util

import (
	"github.com/brutella/hc/accessory"
	"strings"
	"crypto/sha256"
	"encoding/binary"
)

func TopicToInt64(topic string) int64 {
	sum := sha256.Sum256([]byte(topic))
	return int64(binary.BigEndian.Uint64(append([]byte{0x0F}, sum[:7]...)))
}

func AccessoryType(t string) accessory.AccessoryType {
	switch strings.ToLower(t) {
	case "other":
		return accessory.TypeOther
	case "bridge":
		return accessory.TypeBridge
	case "fan":
		return accessory.TypeFan
	case "garagedooropener":
		return accessory.TypeGarageDoorOpener
	case "lightbulb":
		return accessory.TypeLightbulb
	case "doorlock":
		return accessory.TypeDoorLock
	case "outlet":
		return accessory.TypeOutlet
	case "switch":
		return accessory.TypeSwitch
	case "thermostat":
		return accessory.TypeThermostat
	case "sensor":
		return accessory.TypeSensor
	case "securitysystem":
		return accessory.TypeSecuritySystem
	case "door":
		return accessory.TypeDoor
	case "window":
		return accessory.TypeWindow
	case "windowcovering":
		return accessory.TypeWindowCovering
	case "programmableswitch":
		return accessory.TypeProgrammableSwitch
	case "ipcamera":
		return accessory.TypeIPCamera
	case "videodoorbell":
		return accessory.TypeVideoDoorbell
	case "airpurifier":
		return accessory.TypeAirPurifier
	case "heater":
		return accessory.TypeHeater
	case "airconditioner":
		return accessory.TypeAirConditioner
	case "humidifer":
		return accessory.TypeHumidifer
	case "dehumidifier":
		return accessory.TypeDehumidifier
	default:
		return accessory.TypeUnknown
	}
}
