package flagmqtt

import "os"

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
