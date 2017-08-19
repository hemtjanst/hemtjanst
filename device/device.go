package device

import (
	"encoding/json"
	"github.com/hemtjanst/hemtjanst/messaging"
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

func (d *Device) HasFeature(feature string) bool {
	if _, ok := d.Features[feature]; ok {
		return true
	}
	return false
}

func (d *Device) PublishMeta(prefix string) error {
	js, err := json.Marshal(d)
	if err != nil {
		return err
	}
	d.transport.Publish(prefix + d.Topic, js, 1, true)
	return nil
}

func (f *Feature) Set(value string) {
	f.devRef.transport.Publish(f.SetTopic, []byte(value), 1, false)
}

func (f *Feature) OnSet(callback func(msg messaging.Message)) {
	f.devRef.transport.Subscribe(f.SetTopic, 1, callback)
}

func (f *Feature) Update(value string) {
	f.devRef.transport.Publish(f.GetTopic, []byte(value), 1, true)
}

func (f *Feature) OnUpdate(callback func(msg messaging.Message)) {
	f.devRef.transport.Subscribe(f.GetTopic, 1, callback)
}
