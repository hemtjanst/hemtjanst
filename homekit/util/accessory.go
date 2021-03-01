package util

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"strings"

	"github.com/brutella/hc/accessory"
)

func TopicToUint64(topic string) uint64 {
	sum := sha256.Sum256([]byte(topic))
	return binary.BigEndian.Uint64(append([]byte{0x0F}, sum[:7]...))
}

func HexToUint64(hexStr string, def uint64) uint64 {
	if b, err := hex.DecodeString(hexStr); err == nil {
		var ret uint64 = 0
		for i := len(b) - 1; i >= 0; i-- {
			ret = ret<<8 + uint64(b[i])
		}
		// Static addition to avoid collision with auto-assigned ID:s
		return ret + 0x100000000
	}
	return def
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
	case "humidifier":
		return accessory.TypeHumidifier
	case "dehumidifier":
		return accessory.TypeDehumidifier
	case "sprinklers":
		return accessory.TypeSprinklers
	case "faucets":
		return accessory.TypeFaucets
	case "showersystems":
		return accessory.TypeShowerSystems
	default:
		return accessory.TypeUnknown
	}
}
