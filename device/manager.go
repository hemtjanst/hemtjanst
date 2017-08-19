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
	sync.RWMutex
}

func NewManager(c messaging.PublishSubscriber) *Manager {
	return &Manager{
		client:   c,
		devices:  make(map[string]*Device, 10),
		handlers: []Handler{},
	}
}

func (m *Manager) Add(topic string, meta []byte) {
	var dev *Device
	var existing bool
	// If meta is empty, device should be deleted permanently
	del := len(meta) == 0
	if dev, existing = m.devices[topic]; !existing {
		log.Print("Got announce for new device ", topic)
		dev = &Device{Topic: topic, transport: m.client}
	}
	if del {
		if existing {
			log.Print("Got remove for device ", topic)
			// Got empty payload, remove device
			m.Lock()
			defer m.Unlock()
			m.forHandler(func(handler Handler) {
				handler.Removed(dev)
			})

			// Loop through all topics and add to slice first
			// instead of calling unsubscribe() on every feature
			topics := make([]string, len(dev.Features))
			for _, ft := range dev.Features {
				if ft.GetTopic != "" {
					log.Print("Unsubscribing from ", topic)
					topics = append(topics, ft.GetTopic)
				}
			}
			if len(topics) > 0 {
				m.client.Unsubscribe(topics...)
			}

			// Forget device
			delete(m.devices, topic)
			return
		}
		log.Print("Got remove for (non-existing) device ", topic)
		return
	}
	log.Print("Processing meta for device ", topic)

	err := json.Unmarshal(meta, dev)
	dev.Reachable = true
	if err != nil {
		log.Print(err)
		return
	}
	for name, ft := range dev.Features {
		if ft.SetTopic == "" {
			ft.SetTopic = fmt.Sprintf("%s/%s/set", topic, name)
		}
		if ft.GetTopic == "" {
			ft.GetTopic = fmt.Sprintf("%s/%s/get", topic, name)
		}
		ft.devRef = dev
	}

	go m.forHandler(func(handler Handler) {
		handler.Updated(dev)
	})

	if !existing {
		m.Lock()
		defer m.Unlock()
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
