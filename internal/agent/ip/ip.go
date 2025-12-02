package ip

import (
	"fmt"
	"net"
)

func GetLocalIP() (string, error) {
	// Получаем все сетевые интерфейсы
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		// Пропускаем неактивные интерфейсы
		if iface.Flags&net.FlagUp == 0 {
			continue
		}
		// Пропускаем loopback интерфейсы
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		// Получаем адреса интерфейса
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			// Пропускаем IPv6 адреса (по желанию)
			if ip == nil || ip.IsLoopback() {
				continue
			}

			// Преобразуем в IPv4
			ip = ip.To4()
			if ip == nil {
				continue // не IPv4
			}

			return ip.String(), nil
		}
	}

	return "", fmt.Errorf("could not determine local IP")
}
