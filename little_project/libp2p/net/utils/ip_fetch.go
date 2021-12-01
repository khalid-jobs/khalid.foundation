package main

import (
	"fmt"
	"github.com/libp2p/go-netroute"
	"net"
)

func main()  {
	rsl := make([]net.IP, 0)

	// try to use the default ipv4/6 address
	r, err := netroute.New()
	if err != nil {
		panic(err)
	}
	_, _, localIPv4, err := r.Route(net.IPv4zero)
	if err != nil {
		panic(err)
	}

	if localIPv4.IsGlobalUnicast() {
		rsl = append(rsl, localIPv4)
	}
	fmt.Println(rsl)

	if _, _, localIPv6, err := r.Route(net.IPv6unspecified); err != nil {
		fmt.Println("---\n", err)
	} else if localIPv6.IsGlobalUnicast() {
		rsl = append(rsl, localIPv6)
	}

	// resolve the interface addresses
	if addrs, err := net.InterfaceAddrs(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(addrs)
		for _, addr := range addrs {
			ip := net.ParseIP(addr.String())
			if !isIp6LinkLocal(ip) && ip.IsLoopback() {
				rsl = append(rsl, ip)
			}
		}
	}
   fmt.Println("result")

	for _, ip := range rsl {
		fmt.Println(ip.String())
	}

}

func isIp6LinkLocal(ip net.IP) bool {
	return ip.IsLinkLocalMulticast() || ip.IsLinkLocalUnicast()
}

