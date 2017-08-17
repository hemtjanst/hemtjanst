package characteristic

import "github.com/brutella/hc/characteristic"

const TypeColorTemperature = "CE"

type ColorTemperature struct {
	*characteristic.Int
}

func NewColorTemperature() *ColorTemperature {
	char := characteristic.NewInt(TypeColorTemperature)
	char.Format = characteristic.FormatUInt32
	char.Perms = []string{characteristic.PermRead, characteristic.PermWrite, characteristic.PermEvents}
	char.SetMinValue(50)
	char.SetMaxValue(400)
	char.SetStepValue(1)
	char.SetValue(0)

	return &ColorTemperature{char}
}
