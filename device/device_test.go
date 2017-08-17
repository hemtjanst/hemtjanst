package device

import (
	"github.com/hemtjanst/hemtjanst/messaging"
	"testing"
)

func TestNewDevice(t *testing.T) {
	d := NewDevice("test", &messaging.TestingMessenger{})
	if d.Topic != "test" {
		t.Errorf("Expected topic of %s, got %s", "test", d.Topic)
	}

	if d.HasFeature("") {
		t.Error("Expected false, got ", d.HasFeature(""))
	}

}
