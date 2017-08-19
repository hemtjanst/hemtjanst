package device

import (
	"encoding/json"
	"fmt"
	"github.com/hemtjanst/hemtjanst/messaging"
	"log"
	"sync"
)

type Handler interface {
	Updated(*Device)
	Removed(*Device)
}

type Manager struct {
	devices  map[string]*Device
	handlers []Handler
	client   messaging.PublishSubscriber
	init     bool
	sync.RWMutex
}

func NewManager(c messaging.PublishSubscriber, initChan chan bool) *Manager {
	m := &Manager{
		client:   c,
		devices:  make(map[string]*Device, 10),
		handlers: []Handler{},
		// Set init to true if initChan is missing, otherwise wait for a init signal
		init: initChan == nil,
	}
	if initChan != nil {
		go func() {
			<-initChan
			m.init = true
		}()
	}

	return m
}

func (m *Manager) Add(topic string, meta []byte) {
	var dev *Device
	var existing bool
	if dev, existing = m.devices[topic]; !existing {
		log.Print("Got announce for new device ", topic)
		dev = &Device{Topic: topic, transport: m.client}
	}
	log.Print("Processing meta for device ", topic)

	m.Lock()
	defer m.Unlock()
	err := json.Unmarshal(meta, dev)
	dev.Reachable = m.init
	if err != nil {
		log.Print(err)
		return
	}

	go m.forHandler(func(handler Handler) {
		handler.Updated(dev)
	})

	if !existing {
		m.devices[topic] = dev
	}
}

func (m *Manager) Get(id string) (*Device, error) {
	log.Print("Looking for device ", id)
	m.RLock()
	defer m.RUnlock()
	if val, ok := m.devices[id]; ok {
		return val, nil
	}
	return nil, fmt.Errorf("Unknown device %s", id)
}

func (m *Manager) GetAll() map[string]*Device {
	m.RLock()
	defer m.RUnlock()
	return m.devices
}

func (m *Manager) Remove(msg string) {
	m.Lock()
	defer m.Unlock()
	var dev *Device
	var existing bool
	if dev, existing = m.devices[msg]; !existing {
		log.Print("Got remove for (non-existing) device ", msg)
		return
	}
	log.Print("Got remove for device ", msg)
	// Got empty payload, remove device
	go m.forHandler(func(handler Handler) {
		handler.Removed(dev)
	})

	// Loop through all topics and add to slice first
	// instead of calling unsubscribe() on every feature
	topics := make([]string, len(dev.Features))
	for _, ft := range dev.Features {
		if ft.GetTopic != "" {
			log.Print("Unsubscribing from ", ft.GetTopic)
			topics = append(topics, ft.GetTopic)
		}
	}
	if len(topics) > 0 {
		m.client.Unsubscribe(topics...)
	}

	// Forget device
	delete(m.devices, msg)
	return
}

func (m *Manager) Leave(msg string) {
	log.Print("Attempting to remove device ", msg)
	m.Lock()
	defer m.Unlock()
	for _, d := range m.devices {
		if d.LastWillID == msg || d.Topic == msg {
			log.Printf("Found: %s, setting unreachable", d.Topic)
			d.Reachable = false
			dev := d
			go m.forHandler(func(handler Handler) {
				handler.Updated(dev)
			})
		}
	}
}

func (m *Manager) forHandler(f func(handler Handler)) {
	m.RLock()
	defer m.RUnlock()
	for _, h := range m.handlers {
		f(h)
	}
}

func (m *Manager) AddHandler(handler Handler) {
	m.Lock()
	defer m.Unlock()
	m.handlers = append(m.handlers, handler)
	go func() {
		for _, device := range m.devices {
			handler.Updated(device)
		}
	}()
}

// TestingDeviceHandler is a noop device handler. It is meant to
// be used in tests.
type TestingDeviceHandler struct{}

func (t *TestingDeviceHandler) Updated(*Device) {}
func (t *TestingDeviceHandler) Removed(*Device) {}
