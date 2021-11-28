package flowinfo

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"

	"github.com/zartbot/goflow/lib/iputil"
)

/*
var FlowTable ksyncmap.Map

func init() {
	FlowTable.CheckFreq = 10
	FlowTable.Timeout = 20
	go FlowTable.Run()
}
*/

/* FlowInfo used mark flow
   in many Netflow/IPFIX implementation may use multiple field
   mark the SRC/DST IPAddr or Connection Initial/Response Addr
   it's very hard for future flow correlation.
   This package is trying to compare address pair ,and use the
   inner address in IPAddress_A ,outer in IPAddressB */

func ip2int(ip net.IP) uint32 {
	if len(ip) == 16 {
		return binary.BigEndian.Uint32(ip[12:16])
	}
	return binary.BigEndian.Uint32(ip)
}

func Less(IPA net.IP, PortA uint16, IPB net.IP, PortB uint16) bool {
	a := ip2int(IPA)
	b := ip2int(IPB)

	if a < b {
		return true
	} else if a > b {
		return false
	} else {
		// should not be IPA == IPB case in NF record, but still
		// write code here to handle this case.
		if PortA < PortB {
			return true
		} else {
			return false
		}
	}
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func NewFlowInfo(IPA net.IP, PortA uint16, IPB net.IP, PortB uint16) iputil.FlowInfo {
	var f iputil.FlowInfo
	f.Direction = Less(IPA, PortA, IPB, PortB)
	if f.Direction {
		f.FlowID = fmt.Sprintf("%s-%s", IPA, IPB)
		f.FlowIDWithPort = fmt.Sprintf("%s:%d-%s:%d", IPA, PortA, IPB, PortB)

	} else {
		f.FlowID = fmt.Sprintf("%s-%s", IPB, IPA)
		f.FlowIDWithPort = fmt.Sprintf("%s:%d-%s:%d", IPB, PortB, IPA, PortA)

	}
	//this field used to find user/server identity(Geo/DNS/uid...),especially to handle NAT case
	//src in flowbased, client in conn based..
	f.IPAddress_A = IPA
	f.Port_A = PortA
	//dst in flowbased, server in conn based..
	f.IPAddress_B = IPB
	f.Port_B = PortB
	return f
}
