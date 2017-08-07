package device

import (
	"encoding/json"
	"fmt"
	"git.neotor.se/daenney/hemtjanst/messaging"
	"log"
	"sync"
)

type Manager struct {
	devices map[string]*Device
	client  messaging.PublishSubscriber
	sync.RWMutex
}

func NewManager(c messaging.PublishSubscriber) *Manager {
	return &Manager{
		client:  c,
		devices: make(map[string]*Device, 10),
	}
}

func (m *Manager) Add(topic string) {
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
		for f, _ := range val.Features {
			m.client.Unsubscribe(fmt.Sprintf("%s/%s/get", msg, f))
		}
		delete(m.devices, msg)
		return
	}
	for _, d := range m.devices {
		if d.LastWillID.String() == msg {
			log.Print("Found device match for LastWillUID, calling Remove")
			m.Remove(d.Topic)
		}
	}
}
