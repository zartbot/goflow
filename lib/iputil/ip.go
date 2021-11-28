package iputil

import (
	"net"
)

var privateIPBlocks []*net.IPNet

func init() {
	for _, cidr := range []string{
		"127.0.0.0/8",    // IPv4 loopback
		"10.0.0.0/8",     // RFC1918
		"172.16.0.0/12",  // RFC1918
		"192.168.0.0/16", // RFC1918
		"169.254.0.0/16", // RFC5735
		"::1/128",        // IPv6 loopback
		"fe80::/10",      // IPv6 link-local
	} {
		_, block, _ := net.ParseCIDR(cidr)
		privateIPBlocks = append(privateIPBlocks, block)
	}
}

func IsPrivateIP(ip net.IP) bool {
	for _, block := range privateIPBlocks {
		if block.Contains(ip) {
			return true
		}
	}
	return false
}

type FlowInfo struct {
	IPAddress_A    net.IP
	Hash_A         uint32
	Hash_B         uint32
	Hash_A_Mod     uint32
	Hash_B_Mod     uint32
	Port_A         uint16
	IPAddress_B    net.IP
	Port_B         uint16
	FlowID         string
	FlowIDWithPort string
	Direction      bool
}

func FetchIPAddress(d map[string]interface{}, mode string) (net.IP, bool) {
	var v4fieldname, v6fieldname string
	switch mode {
	case "src":
		v4fieldname = "sourceIPv4Address"
		v6fieldname = "sourceIPv6Address"
	case "dst":
		v4fieldname = "destinationIPv4Address"
		v6fieldname = "destinationIPv6Address"
	case "client":
		v4fieldname = "conn_client_ipv4_address"
		v6fieldname = "conn_client_ipv6_address"
	case "server":
		v4fieldname = "conn_server_ipv4_address"
		v6fieldname = "conn_server_ipv6_address"
	case "postnatsrc":
		v4fieldname = "postNATSourceIPv4Address"
		v6fieldname = "postNATSourceIPv6Address"
	case "postnatdst":
		v4fieldname = "postNATDestinationIPv4Address"
		v6fieldname = "postNATDestinationIPv6Address"
	}
	if v4, valid4 := d[v4fieldname]; valid4 {
		return v4.(net.IP), true
	} else if v6, valid6 := d[v6fieldname]; valid6 {
		return v6.(net.IP), true
	} else {
		return nil, false
	}
}

func FetchPort(d map[string]interface{}, mode string) uint16 {
	var fieldname string
	switch mode {
	case "src":
		fieldname = "sourceTransportPort"
	case "dst":
		fieldname = "destinationTransportPort"
	case "client":
		fieldname = "conn_client_trans_port"
	case "server":
		fieldname = "conn_server_trans_port"
	case "postnatsrc":
		fieldname = "postNAPTSourceTransportPort"
	case "postnatdst":
		fieldname = "postNAPTDestinationTransportPort"
	}
	v, valid := d[fieldname]
	if valid {
		return v.(uint16)
	}
	return 0
}
