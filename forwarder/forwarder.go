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
