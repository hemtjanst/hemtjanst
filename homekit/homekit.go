package homekit

import (
	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	lg "github.com/brutella/hc/log"
	"log"
)

func New() {
	debug := lg.Debug
	debug.Enable()
	bridge := NewBridge()
	config := hc.Config{Pin: "01020304", Port: "12345", StoragePath: "./db"}

	t, err := hc.NewIPTransport(config, bridge.Accessory)
	if err != nil {
		log.Fatal("Could not start HomeKit bridge: ", err)
	}

	hc.OnTermination(func() {
		t.Stop()
	})

	t.Start()
}
