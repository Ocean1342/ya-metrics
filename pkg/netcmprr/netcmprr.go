package netcmprr

import (
	"fmt"
	"net"
)

func IsTrustedSubnet(trustedSubnet, requestIP string) (bool, error) {
	_, ipNet, err := net.ParseCIDR(trustedSubnet)
	if err != nil {
		return false, fmt.Errorf("error parsing trusted subnet %s: %s", trustedSubnet, err)
	}
	reqIP := net.ParseIP(requestIP)

	return ipNet.Contains(reqIP), nil
	// убедиться, что айпи находится в этой же подсети
}
