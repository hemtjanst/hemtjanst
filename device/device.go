package device

import (
	"encoding/json"
	"fmt"
	"github.com/hemtjanst/hemtjanst/messaging"
)

type Device struct {
	Topic        string
	Name         string              `json:"name"`
	Manufacturer string              `json:"manufacturer"`
	Model        string              `json:"model"`
	SerialNumber string              `json:"serialNumber"`
	Type         string              `json:"device"`
	LastWillID   string              `json:"lastWillID,omitempty"`
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

func (d *Device) GetTopic(feature string) string {
	if ft, ok := d.Features[feature]; ok && ft.GetTopic != "" {
		return ft.GetTopic
	}
	return d.Topic + "/" + feature + "/get"
}
func (d *Device) SetTopic(feature string) string {
	if ft, ok := d.Features[feature]; ok && ft.SetTopic != "" {
		return ft.SetTopic
	}
	return d.Topic + "/" + feature + "/set"
}

func (d *Device) Set(feature, value string) error {
	if !d.HasFeature(feature) {
		return fmt.Errorf("Feature %s not found on device %s", feature, d.Topic)
	}
	d.transport.Publish(d.SetTopic(feature),
		[]byte(value), 1, false)
	return nil
}

func (d *Device) Watch(feature string, callback func(msg messaging.Message)) error {
	if !d.HasFeature(feature) {
		return fmt.Errorf("Feature %s not found on device %s", feature, d.Topic)
	}
	d.transport.Subscribe(d.GetTopic(feature),
		1, callback)
	return nil
}

func (d *Device) Update(feature, value string) error {
	if !d.HasFeature(feature) {
		return fmt.Errorf("Feature %s not found on device %s", feature, d.Topic)
	}
	d.transport.Publish(d.GetTopic(feature),
		[]byte(value), 1, true)
	return nil
}

func (d *Device) OnSet(feature string, callback func(msg messaging.Message)) error {
	if !d.HasFeature(feature) {
		return fmt.Errorf("Feature %s not found on device %s", feature, d.Topic)
	}
	d.transport.Subscribe(d.SetTopic(feature),
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
