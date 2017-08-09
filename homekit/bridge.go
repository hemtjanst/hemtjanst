package homekit

import (
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/service"
)

type Bridge struct {
	*accessory.Accessory
	BridgingState *service.BridgingState
}

func NewBridge() *Bridge {
	acc := Bridge{}
	info := accessory.Info{
		Name:         "hemtjanst",
		SerialNumber: "000-000",
		Manufacturer: "hemtjanst",
		Model:        "v0",
	}
	acc.Accessory = accessory.New(info, accessory.TypeBridge)
	acc.BridgingState = service.NewBridgingState()
	acc.BridgingState.Category.SetValue(1)
	acc.BridgingState.Reachable.SetValue(true)
	acc.BridgingState.LinkQuality.SetValue(100)

	acc.AddService(acc.BridgingState.Service)

	return &acc
}
