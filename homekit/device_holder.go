package homekit

import (
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/characteristic"
	"github.com/brutella/hc/service"
	"github.com/gosexy/to"
	"github.com/hemtjanst/hemtjanst/device"
	"github.com/hemtjanst/hemtjanst/homekit/util"
	"github.com/hemtjanst/hemtjanst/messaging"
	"log"
	"net"
	"reflect"
)

type deviceHolder struct {
	device    *device.Device
	accessory *accessory.Accessory
}

func createAccessoryFromDevice(d *device.Device) *accessory.Accessory {

	info := accessory.Info{
		Name:         d.Name,
		Manufacturer: d.Manufacturer,
		Model:        d.Model,
		SerialNumber: d.SerialNumber,
	}

	dType := util.AccessoryType(d.Type)
	sType := util.ServiceType(d.Type)

	a := accessory.New(info, dType)
	a.ID = util.TopicToInt64(d.Topic)

	svc := service.New(sType)
	chCount := 0

	for name, feature := range d.Features {
		ch := util.CharacteristicType(name)
		if ch == nil {
			log.Printf("Ignoring unknown characteristic '%s' (from %s)", name, d.Topic)
			continue
		}

		switch ch.Format {
		case characteristic.FormatBool:
			break
		case characteristic.FormatData:
			break
		case characteristic.FormatFloat:
			if feature.Max > 0 {
				ch.MaxValue = float64(feature.Max)
			}
			if feature.Min > 0 {
				ch.MinValue = float64(feature.Min)
			}
			if feature.Step > 0 {
				ch.StepValue = float64(feature.Step)
			}
			break
		case characteristic.FormatInt32, characteristic.FormatUInt8, characteristic.FormatUInt16, characteristic.FormatUInt32, characteristic.FormatUInt64:
			if feature.Max > 0 {
				ch.MaxValue = feature.Max
			}
			if feature.Min > 0 {
				ch.MinValue = feature.Min
			}
			if feature.Step > 0 {
				ch.StepValue = feature.Step
			}
			break

		case characteristic.FormatString:
			break
		case characteristic.FormatTLV8:
			break
		}
		svc.AddCharacteristic(ch)
		chCount++

		ch.OnValueUpdateFromConn(func(conn net.Conn, c *characteristic.Characteristic, newValue, oldValue interface{}) {
			var out string

			if converted, err := to.Convert(newValue, reflect.TypeOf(out).Kind()); err == nil {
				out = converted.(string)
			}
			if b, ok := newValue.(bool); ok {
				if b {
					out = "1"
				} else {
					out = "0"
				}
			}

			if out != "" {
				feature.Set(out)
			}
		})
		feature.OnUpdate(func(msg messaging.Message) {
			ch.UpdateValue(string(msg.Payload()))
		})

	}

	if chCount > 0 {
		a.AddService(svc)
		for _, s := range a.GetServices() {
			// There should never be multiple instances with the same type added to a device
			// so it should be safe to set ID of service/characteristics to its type
			s.ID = util.HexToInt64(s.Type, s.ID)
			for _, c := range s.GetCharacteristics() {
				c.ID = util.HexToInt64(c.Type, c.ID)
			}
		}
	}
	return a
}

func newDeviceHolder(d *device.Device) (*deviceHolder, error) {

	newDev := &deviceHolder{
		device:    d,
		accessory: createAccessoryFromDevice(d),
	}
	return newDev, nil
}
