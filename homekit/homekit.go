package homekit

import (
	"git.neotor.se/daenney/hemtjanst/device"
	"git.neotor.se/daenney/hemtjanst/homekit/bridge"
	"git.neotor.se/daenney/hemtjanst/homekit/util"
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

func (h *Homekit) DeviceUpdated(d *device.Device) {
	h.lock.Lock()
	defer h.lock.Unlock()
	if val, ok := h.devices[d.Topic]; ok {
		val.device = d

		oldAcc := val.accessory
		newAcc := createAccessoryFromDevice(d)
		util.SetReachability(newAcc, util.GetReachability(oldAcc))

		if !newAcc.Equal(oldAcc) {
			val.accessory = newAcc
			h.bridge.ReplaceAccessory(oldAcc, newAcc)
		}

	} else {
		newDev, err := newDeviceHolder(d)
		if err != nil {
			return
		}
		if newDev.accessory != nil {
			util.SetReachability(newDev.accessory, true)
			h.bridge.AddAccessory(newDev.accessory)
		}
		h.devices[d.Topic] = newDev
	}
}

func (h *Homekit) DeviceLeave(d *device.Device) {
	if val, ok := h.devices[d.Topic]; ok {
		h.lock.Lock()
		defer h.lock.Unlock()
		if val.accessory != nil {
			util.SetReachability(val.accessory, false)
		}
	}
}

func (h *Homekit) DeviceRemoved(d *device.Device) {
	if val, ok := h.devices[d.Topic]; ok {
		h.lock.Lock()
		defer h.lock.Unlock()
		if val.accessory != nil {
			h.bridge.RemoveAccessory(val.accessory)
		}
		delete(h.devices, d.Topic)
	}
}
