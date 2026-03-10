package ip

import (
	"errors"
	"net"
	"sync"
)

var (
	NoValidIPv4AddressFound = errors.New("no valid IPv4 address found")
)

// 全局变量用于缓存IP地址和相关锁
var (
	cachedOutboundIp string
	ipMutex          sync.RWMutex
)

/*
GetLocalIp 获取本机非环回地址的 IPv4 地址
返回值:

	ip - 找到的第一个有效 IPv4 地址
	err - 错误信息（未找到有效IP或接口错误）
*/
func GetLocalIp() (string, error) {
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

// GetOutboundIp 获取本机对外网络连接使用的IP地址
// 通过建立一个UDP连接到远程地址(8.8.8.8:80)，然后获取本地绑定的IP地址
// 这个方法可以获取到系统用于出站连接的实际IP地址
// 首次调用时会执行网络操作获取IP，后续调用将直接从内存缓存中读取
// 当网络连接失败时，最多会重试3次

// 返回值:
//   - string: 本机对外IP地址
//   - error: 错误信息，如果获取失败则返回错误
func GetOutboundIp() (string, error) {
	// 首先尝试使用读锁获取已缓存的IP
	ipMutex.RLock()

	if cachedOutboundIp != "" {
		defer ipMutex.RUnlock()
		return cachedOutboundIp, nil
	}

	ipMutex.RUnlock()

	// 缓存中没有IP，需要获取，这时使用写锁
	ipMutex.Lock()
	defer ipMutex.Unlock()

	// 双重检查，防止在等待写锁期间有其他goroutine已经设置了缓存
	if cachedOutboundIp != "" {
		return cachedOutboundIp, nil
	}

	// 最多重试3次
	var lastErr error
	for range 3 {
		// 建立一个UDP连接到Google DNS服务器，目的是获取本地绑定的IP地址
		conn, err := net.Dial("udp", "8.8.8.8:80") // 使用 Google DNS
		if err != nil {
			lastErr = err
			continue // 继续下一次尝试
		}

		defer conn.Close()

		// 从连接中提取本地地址并转换为字符串格式返回
		localAddr := conn.LocalAddr().(*net.UDPAddr)
		ip := localAddr.IP.String()

		// 将获取到的IP保存到缓存中
		cachedOutboundIp = ip
		return ip, nil
	}

	// 如果所有尝试都失败了，返回最后一次的错误
	return "", lastErr
}
