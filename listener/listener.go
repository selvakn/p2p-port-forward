package listener

import "os"

func LocalPortToListen() string {
	if len(os.Args) >= 2 {
		return os.Args[1]
	} else {
		return "22"
	}
}
