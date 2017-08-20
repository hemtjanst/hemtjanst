package homekit

import (
	"fmt"
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
	device          *device.Device
	accessory       *accessory.Accessory
	mainService     *service.Service
	characteristics map[string]*characteristic.Characteristic
}

func newDeviceHolder(d *device.Device) (*deviceHolder, error) {
	newDev := &deviceHolder{
		device:          d,
		accessory:       nil,
		mainService:     nil,
		characteristics: map[string]*characteristic.Characteristic{},
	}
	err := newDev.createAccessory()
	if err != nil {
		return nil, err
	}
	return newDev, nil
}

func (h *deviceHolder) onHomekitUpdate(c string, value interface{}) {
	log.Printf("onHomeKitUpdate(%s, %v) on device %s\n", c, value, h.device.Topic)
	log.Print(h.device)
	if feature, ok := h.device.Features[c]; ok {
		log.Print(feature)
		var out string

		if converted, err := to.Convert(value, reflect.TypeOf(out).Kind()); err == nil {
			out = converted.(string)
		}
		if b, ok := value.(bool); ok {
			if b {
				out = "1"
			} else {
				out = "0"
			}
		}

		if out != "" {
			feature.Set(out)
		}
		return
	}
}

func (h *deviceHolder) onUpdate(c, value string) {
	log.Printf("onUpdate(%s, %s) on device %s\n", c, value, h.device.Topic)
	if ch, ok := h.characteristics[c]; ok {
		log.Print("Found characteristic: ", c)
		ch.UpdateValue(value)
	}
}

func (h *deviceHolder) deviceUpdate(d *device.Device) {
	h.device = d
	util.SetReachability(h.accessory, d.Reachable)

	// TODO
	// h.updateAccessory()
}

func (h *deviceHolder) createAccessory() (err error) {
	if h.accessory != nil {
		return fmt.Errorf("Accessory already created for device %s", h.device.Topic)
	}

	info := accessory.Info{
		Name:         h.device.Name,
		Manufacturer: h.device.Manufacturer,
		Model:        h.device.Model,
		SerialNumber: h.device.SerialNumber,
	}

	dType := util.AccessoryType(h.device.Type)
	a := accessory.New(info, dType)
	h.accessory = a
	a.ID = util.TopicToInt64(h.device.Topic)

	return h.updateAccessory()
}
func (h *deviceHolder) updateAccessory() (err error) {
	sType := util.ServiceType(h.device.Type)

	// TODO: Compare with current service/characteristics if any are set
	//       instead of creating new ones.
	svc := service.New(sType)
	h.mainService = svc
	chCount := 0

	for name, feature := range h.device.Features {
		chName := name
		ch := util.CharacteristicType(name)

		if ch == nil {
			log.Printf("Ignoring unknown characteristic '%s' (from %s)", name, h.device.Topic)
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
		h.characteristics[chName] = ch
		svc.AddCharacteristic(ch)
		chCount++

		ch.OnValueUpdateFromConn(func(conn net.Conn, c *characteristic.Characteristic, newValue, oldValue interface{}) {
			h.onHomekitUpdate(chName, newValue)
		})
		feature.OnUpdate(func(msg messaging.Message) {
			h.onUpdate(chName, string(msg.Payload()))
		})

	}

	if chCount > 0 {
		h.accessory.AddService(svc)
		for _, s := range h.accessory.GetServices() {
			// There should never be multiple instances with the same type added to a device
			// so it should be safe to set ID of service/characteristics to its type
			s.ID = util.HexToInt64(s.Type, s.ID)
			for _, c := range s.GetCharacteristics() {
				c.ID = util.HexToInt64(c.Type, c.ID)
			}
		}
	}

	return
}
