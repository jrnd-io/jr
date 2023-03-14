package jr

import (
	"net"
	"time"
)

func unixTimeStamp(days int) int64 {
	unixEpoch := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	now := time.Now()
	first := now.AddDate(0, 0, -days).Sub(unixEpoch).Seconds()
	last := now.Sub(unixEpoch).Seconds()
	return Random.Int63n(int64(last-first)) + int64(first)
}

func ipKnownPorts() string {
	ports := []string{"80", "81", "443", "22", "631"}
	return ports[Random.Intn(len(ports))]
}

func ipKnownProtocols() string {
	protocols := []string{"TCP", "UDP", "ICMP", "FTP", "HTTP", "SFTP"}
	return protocols[Random.Intn(len(protocols))]
}

func ip(cidr string) string {

GENERATE:

	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return "0.0.0.0"
	}

	ones, _ := ipnet.Mask.Size()
	quotient := ones / 8
	remainder := ones % 8

	r := make([]byte, 4)
	Random.Read(r)

	for i := 0; i <= quotient; i++ {
		if i == quotient {
			shifted := byte(r[i]) >> remainder
			r[i] = ^ipnet.IP[i] & shifted
		} else {
			r[i] = ipnet.IP[i]
		}
	}
	ip = net.IPv4(r[0], r[1], r[2], r[3])

	if ip.Equal(ipnet.IP) /*|| ip.Equal(broadcast) */ {
		goto GENERATE
	}
	return ip.String()
}
