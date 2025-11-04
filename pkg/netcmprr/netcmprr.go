package netcmprr

import (
	"fmt"
	"net"
)

type Host string

func IsTrustedSubnet(trustedSubnet, requestIP string) (bool, error) {
	_, ipNet, err := net.ParseCIDR(trustedSubnet)
	if err != nil {
		return false, fmt.Errorf("error parsing trusted subnet %s: %s", trustedSubnet, err)
	}
	reqIP := net.ParseIP(requestIP)
	return ipNet.Contains(reqIP), nil
}
