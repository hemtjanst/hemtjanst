package device

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hemtjanst/hemtjanst/messaging"
	"sync"
)

var (
	devRefError = errors.New(`Feature is missing reference to device. Use .AddFeature()
to add a feature to a device`)
)

type Device struct {
	Topic        string              `json:"topic"`
	Name         string              `json:"name"`
	Manufacturer string              `json:"manufacturer"`
	Model        string              `json:"model"`
	SerialNumber string              `json:"serialNumber"`
	Type         string              `json:"type"`
	LastWillID   string              `json:"lastWillID,omitempty"`
	Features     map[string]*Feature `json:"feature"`
	Reachable    bool                `json:"-"`
	transport    messaging.PublishSubscriber
	sync.RWMutex
}

type Feature struct {
	Min      int    `json:"min,omitempty"`
	Max      int    `json:"max,omitempty"`
	Step     int    `json:"step,omitempty"`
	GetTopic string `json:"getTopic,omitempty"`
	SetTopic string `json:"setTopic,omitempty"`
	devRef   *Device
}

func NewDevice(topic string, client messaging.PublishSubscriber) *Device {
	return &Device{Topic: topic, transport: client}
}

// UnmarshalJSON unmarshals the JSON representation of a device. This is a
// custom implementation so that we can correctly call .AddFeature() for
// every feature. This guarantees that the devRef is properly set and that
// the GetTopic and SetTopic are always set.
func (d *Device) UnmarshalJSON(b []byte) error {
	var objmap map[string]*json.RawMessage
	err := json.Unmarshal(b, &objmap)
	if err != nil {
		return errors.New("Failed to decode device object")
	}
	if val, ok := objmap["topic"]; ok {
		json.Unmarshal(*val, &d.Topic)
	}
	if val, ok := objmap["name"]; ok {
		json.Unmarshal(*val, &d.Name)
	}
	if val, ok := objmap["type"]; ok {
		json.Unmarshal(*val, &d.Type)
	}
	if val, ok := objmap["manufacturer"]; ok {
		json.Unmarshal(*val, &d.Manufacturer)
	}
	if val, ok := objmap["model"]; ok {
		json.Unmarshal(*val, &d.Model)
	}
	if val, ok := objmap["serialNumber"]; ok {
		json.Unmarshal(*val, &d.SerialNumber)
	}
	if val, ok := objmap["lastWillID"]; ok {
		json.Unmarshal(*val, &d.LastWillID)
	}
	if val, ok := objmap["feature"]; ok {
		// We have features, lets add them
		var ftmap map[string]*json.RawMessage
		err = json.Unmarshal(*val, &ftmap)
		if err != nil {
			return errors.New("Failed to decode feature map")
		}
		for name, settings := range ftmap {
			ftr := &Feature{}
			err = json.Unmarshal(*settings, ftr)
			if err != nil {
				return errors.New("Failed to decode settings map")
			}
			d.AddFeature(name, ftr)
		}
	}
	return nil
}

func (d *Device) HasFeature(feature string) bool {
	d.RLock()
	defer d.RUnlock()
	if _, ok := d.Features[feature]; ok {
		return true
	}
	return false
}

// AddFeature adds a feature to the device. It ensures the devRef,
// GetTopic and SetTopic are correctly populated
func (d *Device) AddFeature(name string, ft *Feature) {
	d.Lock()
	defer d.Unlock()
	if d.Features == nil {
		d.Features = map[string]*Feature{}
	}
	d.Features[name] = ft
	d.Features[name].devRef = d
	if ft.GetTopic == "" {
		ft.GetTopic = fmt.Sprintf("%s/%s/%s", d.Topic, name, "get")
	}
	if ft.SetTopic == "" {
		ft.SetTopic = fmt.Sprintf("%s/%s/%s", d.Topic, name, "set")
	}
}

// GetFeature returns a *Feature if a feature by that name is found
// on the device.
func (d *Device) GetFeature(feature string) (*Feature, error) {
	if !d.HasFeature(feature) {
		return nil, fmt.Errorf("Device has no feature: %s", feature)
	}
	d.RLock()
	defer d.RUnlock()
	return d.Features[feature], nil
}

// RemoveFeature removes a feature by that name from the device.
// It returns an error if you try to remove a feature that does
// not exist.
func (d *Device) RemoveFeature(name string) error {
	if !d.HasFeature(name) {
		return fmt.Errorf("No feature found on device named %s", name)
	}
	d.Lock()
	defer d.Unlock()
	delete(d.Features, name)
	return nil
}

func (d *Device) PublishMeta() error {
	d.RLock()
	defer d.RUnlock()
	js, err := json.Marshal(d)
	if err != nil {
		return err
	}
	d.transport.Publish("announce/"+d.Topic, js, 1, true)
	return nil
}

func (f *Feature) Set(value string) error {
	if f.devRef == nil {
		return devRefError
	}
	f.devRef.transport.Publish(f.SetTopic, []byte(value), 1, false)
	return nil
}

func (f *Feature) OnSet(callback func(msg messaging.Message)) error {
	if f.devRef == nil {
		return devRefError
	}
	f.devRef.transport.Subscribe(f.SetTopic, 1, callback)
	return nil
}

func (f *Feature) Update(value string) error {
	if f.devRef == nil {
		return devRefError
	}
	f.devRef.transport.Publish(f.GetTopic, []byte(value), 1, true)
	return nil
}

func (f *Feature) OnUpdate(callback func(msg messaging.Message)) error {
	if f.devRef == nil {
		return devRefError
	}
	f.devRef.transport.Subscribe(f.GetTopic, 1, callback)
	return nil
}
