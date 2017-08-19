package homekit

import (
	"github.com/hemtjanst/hemtjanst/device"
	"github.com/hemtjanst/hemtjanst/homekit/bridge"
	"github.com/hemtjanst/hemtjanst/homekit/util"
	"sync"
)

type Homekit struct {
	lock    sync.RWMutex
	bridge  bridge.Bridge
	manager *device.Manager
	devices map[string]*deviceHolder
}

func NewHomekit(bridge bridge.Bridge, manager *device.Manager) *Homekit {
	return &Homekit{
		bridge:  bridge,
		manager: manager,
		devices: map[string]*deviceHolder{},
	}
}

func (h *Homekit) Updated(d *device.Device) {
	h.lock.Lock()
	defer h.lock.Unlock()
	if val, ok := h.devices[d.Topic]; ok {
		val.deviceUpdate(d)
	} else {
		newDev, err := newDeviceHolder(d)
		if err != nil {
			return
		}
		if newDev.accessory != nil {
			util.SetReachability(newDev.accessory, d.Reachable)
			h.bridge.AddAccessory(newDev.accessory)
		}
		h.devices[d.Topic] = newDev
	}
}

func (h *Homekit) Removed(d *device.Device) {
	if val, ok := h.devices[d.Topic]; ok {
		h.lock.Lock()
		defer h.lock.Unlock()
		if val.accessory != nil {
			h.bridge.RemoveAccessory(val.accessory)
		}
		delete(h.devices, d.Topic)
	}
}
