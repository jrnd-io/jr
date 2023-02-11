package functions

import (
	"math/rand"
	"net"
	"time"
)

func unixTimeStamp(days int) int64 {
	unixEpoch := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	now := time.Now()
	first := now.AddDate(0, 0, -days).Sub(unixEpoch).Seconds()
	last := now.Sub(unixEpoch).Seconds()
	return rand.Int63n(int64(last-first)) + int64(first)
}

func ipKnownPorts() string {
	ports := []string{"80", "81", "443", "22", "631"}
	return ports[rand.Intn(len(ports))]
}

func ipKnownProtocols() string {
	protocols := []string{"TCP", "UDP", "ICMP", "FTP", "HTTP", "SFTP"}
	return protocols[rand.Intn(len(protocols))]
}

func ip(cidr string) string {

GENERATE:

	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return "0.0.0.0"
	}

	// The number of leading 1s in the mask
	ones, _ := ipnet.Mask.Size()
	quotient := ones / 8
	remainder := ones % 8

	// create random 4-byte byte slice
	r := make([]byte, 4)
	rand.Read(r)

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
		// we got unlucky. The host portion of our ipv4 address was
		// either all 0s (the network address) or all 1s (the broadcast address)
		goto GENERATE
	}
	return ip.String()
}
