package utils

import "../libzt"

const NETWORK_ID = "8056c2e21c000001"

func InitZT() {
	libzt.SimpleStart("./zt", NETWORK_ID)

	log.Infof("ipv4 = %s \n", libzt.GetIpv4Address(NETWORK_ID))
	log.Infof("ipv6 = %s \n", libzt.GetIpv6Address(NETWORK_ID))
}
