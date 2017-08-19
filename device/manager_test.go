package device

import (
	"github.com/hemtjanst/hemtjanst/messaging"
	"io/ioutil"
	"log"
	"testing"
)

func init() {
	log.SetFlags(0)
	log.SetOutput(ioutil.Discard)
}

func TestNewManager(t *testing.T) {
	c := &messaging.TestingMQTTClient{}
	m := messaging.NewTestingMessenger(c)
	mn := NewManager(m, nil)

	if mn.client != m {
		t.Error("Manager is missing a PublishSubsciber")
	}
	if mn.devices == nil {
		t.Error("Expected devices to be initialised to a map")
	}
	if len(mn.handlers) > 0 {
		t.Error("Did not expect any device handlers")
	}
}

func TestManagerAdd(t *testing.T) {
	c := &messaging.TestingMQTTClient{}
	m := messaging.NewTestingMessenger(c)
	mn := NewManager(m, nil)

	mn.Add("lightbulb/kitchen", []byte(`{}`))
	if len(mn.devices) != 1 {
		t.Error("Expected only 1 device, have ", len(mn.devices))
	}

	// Add the same device again, nothing should happen
	mn.Add("lightbulb/kitchen", []byte(`{}`))
	if len(mn.devices) != 1 {
		t.Error("Expected only 1 device, have ", len(mn.devices))
	}
}

func TestManagerGet(t *testing.T) {
	c := &messaging.TestingMQTTClient{}
	m := messaging.NewTestingMessenger(c)
	mn := NewManager(m, nil)

	mn.Add("lightbulb/kitchen", []byte(`{}`))
	_, err := mn.Get("lightbulb/kitchen")
	if err != nil {
		t.Error("Expected to find device")
	}
	_, err = mn.Get("contactSensor/kitchen")
	if err == nil {
		t.Error("Expected to not find device")
	}
}

func TestManagerGetAll(t *testing.T) {
	c := &messaging.TestingMQTTClient{}
	m := messaging.NewTestingMessenger(c)
	mn := NewManager(m, nil)

	devs := mn.GetAll()
	if len(devs) > 0 {
		t.Error("Expected 0 devices, got ", len(devs))
	}

	mn.Add("lightbulb/kitchen", []byte(`{}`))
	mn.Add("contactSensor/kitchen", []byte(`{}`))

	devs = mn.GetAll()
	if len(devs) != 2 {
		t.Error("Expected 2 devices, got ", len(devs))
	}
}

func TestManagerRemove(t *testing.T) {
	c := &messaging.TestingMQTTClient{}
	m := messaging.NewTestingMessenger(c)
	mn := NewManager(m, nil)

	mn.Add("lightbulb/kitchen", []byte(`{}`))
	mn.Add("contactSensor/kitchen", []byte(`{}`))

	mn.Add("lightbulb/kitchen", []byte(""))
	if len(mn.devices) != 1 {
		t.Error("Expected 1 device, got ", len(mn.devices))
	}

	mn.Add("contactSensor/bathroom", []byte(`{}`))
	mn.devices["contactSensor/bathroom"].LastWillID = "ted"
	mn.Leave("ted")

	if mn.devices["contactSensor/bathroom"].Reachable {
		t.Error("Expected contactSensor/bathroom to be unreachable")
	}
}

func TestManagerAddHandler(t *testing.T) {
	c := &messaging.TestingMQTTClient{}
	m := messaging.NewTestingMessenger(c)
	mn := NewManager(m, nil)

	mn.Add("lightbulb/kitchen", []byte(`{}`))
	h := &TestingDeviceHandler{}
	mn.AddHandler(h)

	if len(mn.handlers) != 1 {
		t.Error("Expected 1 device handler, got ", len(mn.handlers))
	}
}
