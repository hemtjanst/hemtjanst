package device

import (
	"encoding/json"
	"fmt"
	"github.com/hemtjanst/hemtjanst/messaging"
	"log"
	"sync"
)

type DeviceHandler interface {
	DeviceUpdated(*Device)
	DeviceLeave(*Device)
	DeviceRemoved(*Device)
}

type Manager struct {
	devices  map[string]*Device
	handlers []DeviceHandler
	client   messaging.PublishSubscriber
	sync.RWMutex
}

func NewManager(c messaging.PublishSubscriber) *Manager {
	return &Manager{
		client:   c,
		devices:  make(map[string]*Device, 10),
		handlers: []DeviceHandler{},
	}
}

func (m *Manager) Add(topic string) {
	if _, ok := m.devices[topic]; ok {
		log.Print("Got announce for existing device ", topic)
		return
	}
	log.Print("Going to add device ", topic)
	dev := &Device{Topic: topic, transport: m.client}
	m.client.Subscribe(fmt.Sprintf("%s/meta", topic), 1, func(msg messaging.Message) {
		err := json.Unmarshal(msg.Payload(), dev)
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
		}

		go m.forHandler(func(handler DeviceHandler) {
			handler.DeviceUpdated(dev)
		})
	})

	m.Lock()
	defer m.Unlock()
	m.devices[topic] = dev
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

func (m *Manager) Remove(msg string) {
	log.Print("Attempting to remove device ", msg)
	m.Lock()
	defer m.Unlock()
	if val, ok := m.devices[msg]; ok {
		log.Print("Found device, unsubscribing and removing")
		m.client.Unsubscribe(fmt.Sprintf("%s/meta", msg))
		for _, ft := range val.Features {
			m.client.Unsubscribe(ft.GetTopic)
		}

		go m.forHandler(func(handler DeviceHandler) {
			handler.DeviceLeave(val)
		})

		delete(m.devices, msg)
		return
	}
	for _, d := range m.devices {
		if d.LastWillID == msg {
			log.Print("Found device match for LastWillUID, calling Remove")
			m.Remove(d.Topic)
		}
	}
}

func (m *Manager) forHandler(f func(handler DeviceHandler)) {
	for _, h := range m.handlers {
		f(h)
	}
}

func (m *Manager) AddHandler(handler DeviceHandler) {
	m.Lock()
	defer m.Unlock()
	m.handlers = append(m.handlers, handler)
	go func() {
		for _, device := range m.devices {
			handler.DeviceUpdated(device)
		}
	}()
}
