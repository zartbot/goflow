package flowinfo

import (
	"net"

	"github.com/zartbot/goflow/lib/iputil"
	"github.com/zartbot/goflow/service/identity"
)

type NATinfo struct {
	Direction  bool
	IPAddressA net.IP
	PortA      uint16
	IPAddressB net.IP
	PortB      uint16
}

/*
natEvent :1 Create  :2 Delete
[sourceIPv4Address][sourceTransportPort] <->[destinationIPv4Address][destinationTransportPort]
[postNATSourceIPv4Address][postNAPTSourceTransportPort] <->[postNATDestinationIPv4Address][postNAPTDestinationTransportPort]
*/

func FetchNATAddress(f iputil.FlowInfo) (iputil.FlowInfo, bool) {
	var v iputil.FlowInfo
	natlog, ok := identity.Service["NAT"].Load(f.FlowIDWithPort)
	if ok {
		v := natlog.(iputil.FlowInfo)
		if f.Direction != v.Direction {
			//preNAT src/dst also need to swap
			tempIP := v.IPAddress_A
			tempPort := v.Port_A
			v.IPAddress_A = v.IPAddress_B
			v.Port_A = v.Port_B
			v.IPAddress_B = tempIP
			v.Port_B = tempPort
		}
		return v, true
	} else {
		return v, false
	}
}
