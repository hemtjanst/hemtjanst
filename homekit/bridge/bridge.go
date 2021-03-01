package bridge

import (
	"github.com/brutella/hc/accessory"
	"github.com/hemtjanst/hemtjanst/homekit/util"
)

// Bridge
type Bridge interface {
	AddAccessory(a *accessory.Accessory)
	RemoveAccessory(a *accessory.Accessory)
	ReplaceAccessory(old, new *accessory.Accessory)
	Start()
	Stop()
}

type bridge struct {
	*accessory.Accessory
	transport *ipTransport
}

// NewBridge creates a new bridge
func NewBridge(config Config, info accessory.Info) (_ Bridge, err error) {
	acc := &bridge{}
	acc.Accessory = accessory.New(info, accessory.TypeBridge)
	acc.transport, err = NewIPTransport(config, acc.Accessory)
	if err != nil {
		return
	}

	return acc, nil
}

func (b *bridge) AddAccessory(a *accessory.Accessory) {
	// Make sure accessory has reachable characteristic
	util.GetReachability(a)

	if b.transport != nil {
		id := a.ID
		b.transport.addAccessory(a)
		if id > 0 {
			a.ID = id
		}
		b.transport.updateConfig()
	}
}

func (b *bridge) RemoveAccessory(a *accessory.Accessory) {
	if b.transport.container != nil {
		b.transport.container.RemoveAccessory(a)
	}
	b.transport.updateConfig()
}

func (b *bridge) ReplaceAccessory(old, new *accessory.Accessory) {
	if b.transport.container == nil {
		return
	}
	var id uint64
	if old != nil {
		id = old.ID
		b.transport.container.RemoveAccessory(old)
	}
	if new.ID > 0 {
		id = new.ID
	}
	b.transport.addAccessory(new)
	if id > 0 {
		new.ID = id
	}
	b.transport.updateConfig()
}

func (b *bridge) Start() {
	b.transport.Start()
}

func (b *bridge) Stop() {
	b.transport.Stop()
}
