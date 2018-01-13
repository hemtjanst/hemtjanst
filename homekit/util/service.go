package util

import (
	"github.com/brutella/hc/service"
	"strings"
)

func ServiceType(t string) string {
	switch strings.ToLower(t) {
	case "accessoryinformation":
		return service.TypeAccessoryInformation
	case "airpurifier":
		return service.TypeAirPurifier
	case "airqualitysensor":
		return service.TypeAirQualitySensor
	case "batteryservice":
		return service.TypeBatteryService
	case "bridgeconfiguration":
		return service.TypeBridgeConfiguration
	case "bridgingstate":
		return service.TypeBridgingState
	case "cameracontrol":
		return service.TypeCameraControl
	case "camerartpstreammanagement":
		return service.TypeCameraRTPStreamManagement
	case "carbondioxidesensor":
		return service.TypeCarbonDioxideSensor
	case "carbonmonoxidesensor":
		return service.TypeCarbonMonoxideSensor
	case "contactsensor":
		return service.TypeContactSensor
	case "door":
		return service.TypeDoor
	case "doorbell":
		return service.TypeDoorbell
	case "fan":
		return service.TypeFan
	case "fanv2":
		return service.TypeFanV2
	case "filtermaintenance":
		return service.TypeFilterMaintenance
	case "garagedooropener":
		return service.TypeGarageDoorOpener
	case "heatercooler":
		return service.TypeHeaterCooler
	case "humidifierdehumidifier":
		return service.TypeHumidifierDehumidifier
	case "humiditysensor":
		return service.TypeHumiditySensor
	case "leaksensor":
		return service.TypeLeakSensor
	case "lightsensor":
		return service.TypeLightSensor
	case "lightbulb":
		return service.TypeLightbulb
	case "lockmanagement":
		return service.TypeLockManagement
	case "lockmechanism":
		return service.TypeLockMechanism
	case "microphone":
		return service.TypeMicrophone
	case "motionsensor":
		return service.TypeMotionSensor
	case "occupancysensor":
		return service.TypeOccupancySensor
	case "outlet":
		return service.TypeOutlet
	case "securitysystem":
		return service.TypeSecuritySystem
	case "slat":
		return service.TypeSlat
	case "smokesensor":
		return service.TypeSmokeSensor
	case "speaker":
		return service.TypeSpeaker
	case "statefulprogrammableswitch":
		return service.TypeStatefulProgrammableSwitch
	case "statelessprogrammableswitch":
		return service.TypeStatelessProgrammableSwitch
	case "switch":
		return service.TypeSwitch
	case "temperaturesensor":
		return service.TypeTemperatureSensor
	case "thermostat":
		return service.TypeThermostat
	case "timeinformation":
		return service.TypeTimeInformation
	case "tunneledbtleaccessoryservice":
		return service.TypeTunneledBTLEAccessoryService
	case "window":
		return service.TypeWindow
	case "windowcovering":
		return service.TypeWindowCovering
	default:
		return ""
	}
}
