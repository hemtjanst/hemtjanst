package util

import (
	"strings"

	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/characteristic"
	"github.com/brutella/hc/service"
	ht_char "github.com/hemtjanst/hemtjanst/homekit/characteristic"
)

func FindCharacteristic(a *accessory.Accessory, cType string) *characteristic.Characteristic {
	for _, s := range a.GetServices() {
		for _, c := range s.GetCharacteristics() {
			if c.Type == cType {
				return c
			}
		}
	}
	return nil
}

func SetReachability(a *accessory.Accessory, value bool) {
	c := FindCharacteristic(a, characteristic.TypeReachable)
	if c == nil {
		c = characteristic.NewReachable().Characteristic
		a.AddService(&service.Service{
			Type:            service.TypeBridgingState,
			Characteristics: []*characteristic.Characteristic{c},
		})
	}
	c.UpdateValue(value)
}

func GetReachability(a *accessory.Accessory) bool {
	c := FindCharacteristic(a, characteristic.TypeReachable)
	if c != nil {
		if v, ok := c.Value.(bool); ok {
			return v
		}
	}
	SetReachability(a, false)
	return false
}

func CharacteristicType(t string) *characteristic.Characteristic {
	switch strings.ToLower(t) {
	// Allow both contactSensorState and state
	case "accessoryflags":
		return characteristic.NewAccessoryFlags().Characteristic
	case "accessoryidentifier":
		return characteristic.NewAccessoryIdentifier().Characteristic
	case "active":
		return characteristic.NewActive().Characteristic
	case "administratoronlyaccess":
		return characteristic.NewAdministratorOnlyAccess().Characteristic
	case "airparticulatedensity":
		return characteristic.NewAirParticulateDensity().Characteristic
	case "airparticulatesize":
		return characteristic.NewAirParticulateSize().Characteristic
	case "airquality":
		return characteristic.NewAirQuality().Characteristic
	case "appmatchingidentifier":
		return characteristic.NewAppMatchingIdentifier().Characteristic
	case "audiofeedback":
		return characteristic.NewAudioFeedback().Characteristic
	case "batterylevel":
		return characteristic.NewBatteryLevel().Characteristic
	case "brightness":
		return characteristic.NewBrightness().Characteristic
	case "carbondioxidedetected":
		return characteristic.NewCarbonDioxideDetected().Characteristic
	case "carbondioxidelevel":
		return characteristic.NewCarbonDioxideLevel().Characteristic
	case "carbondioxidepeaklevel":
		return characteristic.NewCarbonDioxidePeakLevel().Characteristic
	case "carbonmonoxidedetected":
		return characteristic.NewCarbonMonoxideDetected().Characteristic
	case "carbonmonoxidelevel":
		return characteristic.NewCarbonMonoxideLevel().Characteristic
	case "carbonmonoxidepeaklevel":
		return characteristic.NewCarbonMonoxidePeakLevel().Characteristic
	case "category":
		return characteristic.NewCategory().Characteristic
	case "chargingstate":
		return characteristic.NewChargingState().Characteristic
	case "colortemperature":
		return ht_char.NewColorTemperature().Characteristic
	case "configurebridgedaccessory":
		return characteristic.NewConfigureBridgedAccessory().Characteristic
	case "configurebridgedaccessorystatus":
		return characteristic.NewConfigureBridgedAccessoryStatus().Characteristic
	case "contactsensorstate", "state":
		return characteristic.NewContactSensorState().Characteristic
	case "coolingthresholdtemperature":
		return characteristic.NewCoolingThresholdTemperature().Characteristic
	case "currentairpurifierstate":
		return characteristic.NewCurrentAirPurifierState().Characteristic
	case "currentambientlightlevel":
		return characteristic.NewCurrentAmbientLightLevel().Characteristic
	case "currentdoorstate":
		return characteristic.NewCurrentDoorState().Characteristic
	case "currentfanstate":
		return characteristic.NewCurrentFanState().Characteristic
	case "currentheatercoolerstate":
		return characteristic.NewCurrentHeaterCoolerState().Characteristic
	case "currentheatingcoolingstate":
		return characteristic.NewCurrentHeatingCoolingState().Characteristic
	case "currenthorizontaltiltangle":
		return characteristic.NewCurrentHorizontalTiltAngle().Characteristic
	case "currenthumidifierdehumidifierstate":
		return characteristic.NewCurrentHumidifierDehumidifierState().Characteristic
	case "currentposition":
		return characteristic.NewCurrentPosition().Characteristic
	case "currentrelativehumidity":
		return characteristic.NewCurrentRelativeHumidity().Characteristic
	case "currentslatstate":
		return characteristic.NewCurrentSlatState().Characteristic
	case "currenttemperature":
		return characteristic.NewCurrentTemperature().Characteristic
	case "currenttiltangle":
		return characteristic.NewCurrentTiltAngle().Characteristic
	case "currenttime":
		return characteristic.NewCurrentTime().Characteristic
	case "currentverticaltiltangle":
		return characteristic.NewCurrentVerticalTiltAngle().Characteristic
	case "dayoftheweek":
		return characteristic.NewDayOfTheWeek().Characteristic
	case "digitalzoom":
		return characteristic.NewDigitalZoom().Characteristic
	case "discoverbridgedaccessories":
		return characteristic.NewDiscoverBridgedAccessories().Characteristic
	case "discoveredbridgedaccessories":
		return characteristic.NewDiscoveredBridgedAccessories().Characteristic
	case "filterchangeindication":
		return characteristic.NewFilterChangeIndication().Characteristic
	case "filterlifelevel":
		return characteristic.NewFilterLifeLevel().Characteristic
	case "firmwarerevision":
		return characteristic.NewFirmwareRevision().Characteristic
	case "hardwarerevision":
		return characteristic.NewHardwareRevision().Characteristic
	case "heatingthresholdtemperature":
		return characteristic.NewHeatingThresholdTemperature().Characteristic
	case "holdposition":
		return characteristic.NewHoldPosition().Characteristic
	case "hue":
		return characteristic.NewHue().Characteristic
	case "identify":
		return characteristic.NewIdentify().Characteristic
	case "imagemirroring":
		return characteristic.NewImageMirroring().Characteristic
	case "imagerotation":
		return characteristic.NewImageRotation().Characteristic
	case "inuse":
		return characteristic.NewInUse().Characteristic
	case "isconfigured":
		return characteristic.NewIsConfigured().Characteristic
	case "leakdetected":
		return characteristic.NewLeakDetected().Characteristic
	case "linkquality":
		return characteristic.NewLinkQuality().Characteristic
	case "lockcontrolpoint":
		return characteristic.NewLockControlPoint().Characteristic
	case "lockcurrentstate":
		return characteristic.NewLockCurrentState().Characteristic
	case "locklastknownaction":
		return characteristic.NewLockLastKnownAction().Characteristic
	case "lockmanagementautosecuritytimeout":
		return characteristic.NewLockManagementAutoSecurityTimeout().Characteristic
	case "lockphysicalcontrols":
		return characteristic.NewLockPhysicalControls().Characteristic
	case "locktargetstate":
		return characteristic.NewLockTargetState().Characteristic
	case "logs":
		return characteristic.NewLogs().Characteristic
	case "manufacturer":
		return characteristic.NewManufacturer().Characteristic
	case "model":
		return characteristic.NewModel().Characteristic
	case "motiondetected":
		return characteristic.NewMotionDetected().Characteristic
	case "mute":
		return characteristic.NewMute().Characteristic
	case "name":
		return characteristic.NewName().Characteristic
	case "nightvision":
		return characteristic.NewNightVision().Characteristic
	case "nitrogendioxidedensity":
		return characteristic.NewNitrogenDioxideDensity().Characteristic
	case "obstructiondetected":
		return characteristic.NewObstructionDetected().Characteristic
	case "occupancydetected":
		return characteristic.NewOccupancyDetected().Characteristic
	case "on":
		return characteristic.NewOn().Characteristic
	case "opticalzoom":
		return characteristic.NewOpticalZoom().Characteristic
	case "outletinuse":
		return characteristic.NewOutletInUse().Characteristic
	case "ozonedensity":
		return characteristic.NewOzoneDensity().Characteristic
	case "pairsetup":
		return characteristic.NewPairSetup().Characteristic
	case "pairverify":
		return characteristic.NewPairVerify().Characteristic
	case "pairingfeatures":
		return characteristic.NewPairingFeatures().Characteristic
	case "pairingpairings":
		return characteristic.NewPairingPairings().Characteristic
	case "pm10density":
		return characteristic.NewPM10Density().Characteristic
	case "pm2_5density":
		return characteristic.NewPM2_5Density().Characteristic
	case "positionstate":
		return characteristic.NewPositionState().Characteristic
	case "programmableswitchevent":
		return characteristic.NewProgrammableSwitchEvent().Characteristic
	case "programmableswitchoutputstate":
		return characteristic.NewProgrammableSwitchOutputState().Characteristic
	case "programmode":
		return characteristic.NewProgramMode().Characteristic
	case "reachable":
		return characteristic.NewReachable().Characteristic
	case "relativehumiditydehumidifierthreshold":
		return characteristic.NewRelativeHumidityDehumidifierThreshold().Characteristic
	case "relativehumidityhumidifierthreshold":
		return characteristic.NewRelativeHumidityHumidifierThreshold().Characteristic
	case "remainingduration":
		return characteristic.NewRemainingDuration().Characteristic
	case "resetfilterindication":
		return characteristic.NewResetFilterIndication().Characteristic
	case "rotationdirection":
		return characteristic.NewRotationDirection().Characteristic
	case "rotationspeed":
		return characteristic.NewRotationSpeed().Characteristic
	case "saturation":
		return characteristic.NewSaturation().Characteristic
	case "securitysystemalarmtype":
		return characteristic.NewSecuritySystemAlarmType().Characteristic
	case "securitysystemcurrentstate":
		return characteristic.NewSecuritySystemCurrentState().Characteristic
	case "securitysystemtargetstate":
		return characteristic.NewSecuritySystemTargetState().Characteristic
	case "selectedrtpstreamconfiguration":
		return characteristic.NewSelectedRTPStreamConfiguration().Characteristic
	case "selectedstreamconfiguration":
		return characteristic.NewSelectedStreamConfiguration().Characteristic
	case "serialnumber":
		return characteristic.NewSerialNumber().Characteristic
	case "servicelabelindex":
		return characteristic.NewServiceLabelIndex().Characteristic
	case "servicelabelnamespace":
		return characteristic.NewServiceLabelNamespace().Characteristic
	case "setduration":
		return characteristic.NewSetDuration().Characteristic
	case "setupendpoints":
		return characteristic.NewSetupEndpoints().Characteristic
	case "slattype":
		return characteristic.NewSlatType().Characteristic
	case "smokedetected":
		return characteristic.NewSmokeDetected().Characteristic
	case "softwarerevision":
		return characteristic.NewSoftwareRevision().Characteristic
	case "statusactive":
		return characteristic.NewStatusActive().Characteristic
	case "statusfault":
		return characteristic.NewStatusFault().Characteristic
	case "statusjammed":
		return characteristic.NewStatusJammed().Characteristic
	case "statuslowbattery":
		return characteristic.NewStatusLowBattery().Characteristic
	case "statustampered":
		return characteristic.NewStatusTampered().Characteristic
	case "streamingstatus":
		return characteristic.NewStreamingStatus().Characteristic
	case "sulphurdioxidedensity":
		return characteristic.NewSulphurDioxideDensity().Characteristic
	case "supportedaudiostreamconfiguration":
		return characteristic.NewSupportedAudioStreamConfiguration().Characteristic
	case "supportedrtpconfiguration":
		return characteristic.NewSupportedRTPConfiguration().Characteristic
	case "supportedvideostreamconfiguration":
		return characteristic.NewSupportedVideoStreamConfiguration().Characteristic
	case "swingmode":
		return characteristic.NewSwingMode().Characteristic
	case "targetairpurifierstate":
		return characteristic.NewTargetAirPurifierState().Characteristic
	case "targetairquality":
		return characteristic.NewTargetAirQuality().Characteristic
	case "targetdoorstate":
		return characteristic.NewTargetDoorState().Characteristic
	case "targetfanstate":
		return characteristic.NewTargetFanState().Characteristic
	case "targetheatercoolerstate":
		return characteristic.NewTargetHeaterCoolerState().Characteristic
	case "targetheatingcoolingstate":
		return characteristic.NewTargetHeatingCoolingState().Characteristic
	case "targethorizontaltiltangle":
		return characteristic.NewTargetHorizontalTiltAngle().Characteristic
	case "targethumidifierdehumidifierstate":
		return characteristic.NewTargetHumidifierDehumidifierState().Characteristic
	case "targetposition":
		return characteristic.NewTargetPosition().Characteristic
	case "targetrelativehumidity":
		return characteristic.NewTargetRelativeHumidity().Characteristic
	case "targetslatstate":
		return characteristic.NewTargetSlatState().Characteristic
	case "targettemperature":
		return characteristic.NewTargetTemperature().Characteristic
	case "targettiltangle":
		return characteristic.NewTargetTiltAngle().Characteristic
	case "targetverticaltiltangle":
		return characteristic.NewTargetVerticalTiltAngle().Characteristic
	case "temperaturedisplayunits":
		return characteristic.NewTemperatureDisplayUnits().Characteristic
	case "timeupdate":
		return characteristic.NewTimeUpdate().Characteristic
	case "tunnelconnectiontimeout":
		return characteristic.NewTunnelConnectionTimeout().Characteristic
	case "tunneledaccessoryadvertising":
		return characteristic.NewTunneledAccessoryAdvertising().Characteristic
	case "tunneledaccessoryconnected":
		return characteristic.NewTunneledAccessoryConnected().Characteristic
	case "tunneledaccessorystatenumber":
		return characteristic.NewTunneledAccessoryStateNumber().Characteristic
	case "valvetype":
		return characteristic.NewValveType().Characteristic
	case "version":
		return characteristic.NewVersion().Characteristic
	case "vocdensity":
		return characteristic.NewVOCDensity().Characteristic
	case "volume":
		return characteristic.NewVolume().Characteristic
	case "waterlevel":
		return characteristic.NewWaterLevel().Characteristic
	default:
		return nil
	}
}
