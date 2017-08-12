package flagmqtt

import (
	"os"
	"strconv"
)

func envOrFlagStr(flagVal, envName, defVal string) (ret string) {
	ret = defVal

	if val, ok := os.LookupEnv(envName); ok {
		ret = val
	}
	if flagVal != defVal {
		ret = flagVal
	}
	return
}

func envOrFlagInt(flagVal int, envName string, defVal int) (ret int) {
	ret = defVal

	if val, ok := os.LookupEnv(envName); ok {
		ret, _ = strconv.Atoi(val)
	}
	if flagVal != defVal {
		ret = flagVal
	}
	return
}
