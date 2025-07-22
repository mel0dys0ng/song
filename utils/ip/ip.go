package ip

import (
	"errors"
	"net"
)

var (
	NoValidIPv4AddressFound = errors.New("no valid IPv4 address found")
)

/*
GetLocalIP 获取本机非环回地址的 IPv4 地址
返回值:

	ip - 找到的第一个有效 IPv4 地址
	err - 错误信息（未找到有效IP或接口错误）
*/
func GetLocalIP() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		// 跳过环回接口和非活动接口
		if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}

			ip := ipNet.IP.To4()
			if ip != nil && !ip.IsLoopback() {
				return ip.String(), nil
			}
		}
	}

	return "", NoValidIPv4AddressFound
}

/*
GetOutboundIP 获取本机出口 IPv4 地址
返回值:

	ip - 出口 IP 地址
	err - 网络连接错误
*/
func GetOutboundIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80") // 使用 Google DNS
	if err != nil {
		return "", err
	}

	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}
