package utils

type IPProto string

const UDP IPProto = "udp"
const TCP IPProto = "tcp"

func (c IPProto) GetName() string {
	return string(c)
}

func GetIPProto(isUDP bool) IPProto {
	if isUDP {
		return UDP
	} else {
		return TCP
	}
}
