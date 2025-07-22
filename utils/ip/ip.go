package ip

import (
	"net"

	"github.com/mel0dys0ng/song/utils/strjngs"
)

// GetLocalHost 获取本地对外IP地址
func GetLocalHost() (ip string) {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		return
	}

	localAddr, ok := conn.LocalAddr().(*net.UDPAddr)
	if ok {
		ip = strjngs.IndexOfSplit(localAddr.String(), ":", 0)
	}

	return
}
