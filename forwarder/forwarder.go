package forwarder

import (
	"os"
)

func GetOtherIP() string {
	if len(os.Args) >= 2 {
		return os.Args[1]
	} else {
		return ""
	}
}

func LocalPortToForward() string {
	if len(os.Args) >= 3 {
		return os.Args[2]
	} else {
		return "2222"
	}
}
