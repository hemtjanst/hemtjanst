package device

import (
	"encoding/json"
	"fmt"
	"github.com/hemtjanst/hemtjanst/messaging"
	"github.com/satori/go.uuid"
)

type Device struct {
	Topic        string
	Name         string              `json:"name"`
	Manufacturer string              `json:"manufacturer"`
	Model        string              `json:"model"`
	SerialNumber string              `json:"serialNumber"`
	Type         string              `json:"device"`
	LastWillID   uuid.UUID           `json:"lastWillID,omitempty"`
	Features     map[string]*Feature `json:"feature"`
	transport    messaging.PublishSubscriber
}

type Feature struct {
	Min      int    `json:"min,omitempty"`
	Max      int    `json:"max,omitempty"`
	Step     int    `json:"step,omitempty"`
	GetTopic string `json:"getTopic,omitempty"`
	SetTopic string `json:"setTopic,omitempty"`
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

func (d *Device) Set(feature, value string) error {
	if !d.HasFeature(feature) {
		return fmt.Errorf("Feature %s not found on device %s", feature, d.Topic)
	}
	ft := d.Features[feature]
	d.transport.Publish(ft.SetTopic,
		[]byte(value), 1, false)
	return nil
}

func (d *Device) Watch(feature string, callback func(msg messaging.Message)) error {
	if !d.HasFeature(feature) {
		return fmt.Errorf("Feature %s not found on device %s", feature, d.Topic)
	}
	ft := d.Features[feature]
	d.transport.Subscribe(ft.GetTopic,
		1, callback)
	return nil
}

func (d *Device) Update(feature, value string) error {
	if !d.HasFeature(feature) {
		return fmt.Errorf("Feature %s not found on device %s", feature, d.Topic)
	}
	ft := d.Features[feature]
	d.transport.Publish(ft.GetTopic,
		[]byte(value), 1, true)
	return nil
}

func (d *Device) OnSet(feature string, callback func(msg messaging.Message)) error {
	if !d.HasFeature(feature) {
		return fmt.Errorf("Feature %s not found on device %s", feature, d.Topic)
	}
	ft := d.Features[feature]
	d.transport.Subscribe(ft.SetTopic,
		1, callback)
	return nil
}

func (d *Device) PublishMeta() error {
	js, err := json.Marshal(d)
	if err != nil {
		return err
	}
	d.transport.Publish(d.Topic+"/meta", js, 1, true)
	return nil
}
