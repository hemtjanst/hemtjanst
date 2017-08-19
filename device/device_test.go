package device

import (
	"bytes"
	"github.com/hemtjanst/hemtjanst/messaging"
	"reflect"
	"testing"
)

func TestNewDevice(t *testing.T) {
	d := NewDevice("test", &messaging.TestingMessenger{})
	if d.Topic != "test" {
		t.Errorf("Expected topic of %s, got %s", "test", d.Topic)
	}

	if d.HasFeature("") {
		t.Error("Expected false, got ", d.HasFeature(""))
	}

	f := &Feature{GetTopic: "lightbulb/on/get", SetTopic: "lightbulb/on/set", devRef: d}
	d.Features = map[string]*Feature{}
	d.Features["on"] = f

	if !d.HasFeature("on") {
		t.Error("Expected true, got ", d.HasFeature("on"))
	}
}

func TestPublishMeta(t *testing.T) {
	m := &messaging.TestingMessenger{}
	d := NewDevice("lightbulb/kitchen", m)
	err := d.PublishMeta()
	if err != nil {
		t.Error("Expected to successfully publish meta, got ", err)
	}

	if m.Action != "publish" {
		t.Error("Expected to publish, but tried to ", m.Action)
	}
	if !reflect.DeepEqual(m.Topic, []string{"announce/lightbulb/kitchen"}) {
		t.Error("Expected topic to be announce/lightbulb/kitchen, got ", m.Topic)
	}
	if m.Qos != 1 {
		t.Error("Expected QoS of 1, got ", m.Qos)
	}
	if !m.Persist {
		t.Error("Expected persist, got ", m.Persist)
	}
	msg := `{"topic":"lightbulb/kitchen","name":"","manufacturer":"","model":"","serialNumber":"","type":"","feature":null}`
	if !bytes.Equal(m.Message, []byte(msg)) {
		t.Errorf("Expected %s, got %s", msg, string(m.Message))
	}
}

func TestFeatureSet(t *testing.T) {
	m := &messaging.TestingMessenger{}
	d := NewDevice("lightbulb", m)
	f := &Feature{GetTopic: "lightbulb/on/get", SetTopic: "lightbulb/on/set", devRef: d}
	d.Features = map[string]*Feature{}
	d.Features["on"] = f
	f.Set("1")

	if m.Action != "publish" {
		t.Error("Expected to publish, but instead tried to ", m.Action)
	}
	if !reflect.DeepEqual(m.Topic, []string{"lightbulb/on/set"}) {
		t.Error("Expected topic to be lightbulb/on/set, got ", m.Topic)
	}
	if m.Qos != 1 {
		t.Error("Expected QoS of 1, got ", m.Qos)
	}
	if m.Persist {
		t.Error("Expected message without persist, got ", m.Persist)
	}
	if !bytes.Equal(m.Message, []byte("1")) {
		t.Error("Expected message of 1, got ", string(m.Message))
	}
}

func TestFeatureOnSet(t *testing.T) {
	m := &messaging.TestingMessenger{}
	d := NewDevice("lightbulb", m)
	f := &Feature{GetTopic: "lightbulb/on/get", SetTopic: "lightbulb/on/set", devRef: d}
	d.Features = map[string]*Feature{}
	d.Features["on"] = f

	f.OnSet(func(messaging.Message) {
		return
	})
	if m.Action != "subscribe" {
		t.Error("Expected to subscribe, but instead tried to ", m.Action)
	}
	if !reflect.DeepEqual(m.Topic, []string{"lightbulb/on/set"}) {
		t.Error("Expected topic to be lightbulb/on/set, got ", m.Topic)
	}
	if m.Qos != 1 {
		t.Error("Expected QoS of 1, got ", m.Qos)
	}
	if m.Callback == nil {
		t.Error("Expected a callback, got nil")
	}
}

func TestFeatureUpdate(t *testing.T) {
	m := &messaging.TestingMessenger{}
	d := NewDevice("lightbulb", m)
	f := &Feature{GetTopic: "lightbulb/on/get", SetTopic: "lightbulb/on/set", devRef: d}
	d.Features = map[string]*Feature{}
	d.Features["on"] = f
	f.Update("1")

	if m.Action != "publish" {
		t.Error("Expected to publish, but instead tried to ", m.Action)
	}
	if !reflect.DeepEqual(m.Topic, []string{"lightbulb/on/get"}) {
		t.Error("Expected topic to be lightbulb/on/get, got ", m.Topic)
	}
	if m.Qos != 1 {
		t.Error("Expected QoS of 1, got ", m.Qos)
	}
	if !m.Persist {
		t.Error("Expected message to persist, got ", m.Persist)
	}
	if !bytes.Equal(m.Message, []byte("1")) {
		t.Error("Expected message of 1, got ", string(m.Message))
	}
}

func TestFeatureOnUpdate(t *testing.T) {
	m := &messaging.TestingMessenger{}
	d := NewDevice("lightbulb", m)
	f := &Feature{GetTopic: "lightbulb/on/get", SetTopic: "lightbulb/on/set", devRef: d}
	d.Features = map[string]*Feature{}
	d.Features["on"] = f

	f.OnUpdate(func(messaging.Message) {
		return
	})
	if m.Action != "subscribe" {
		t.Error("Expected to subscribe, but instead tried to ", m.Action)
	}
	if !reflect.DeepEqual(m.Topic, []string{"lightbulb/on/get"}) {
		t.Error("Expected topic to be lightbulb/on/get, got ", m.Topic)
	}
	if m.Qos != 1 {
		t.Error("Expected QoS of 1, got ", m.Qos)
	}
	if m.Callback == nil {
		t.Error("Expected a callback, got nil")
	}
}
