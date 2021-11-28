package packet

// IPv4HDR  : The struct of IPv4 Packet Header
// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |Version|  IHL  |Type of Service|          Total Length         |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |         Identification        |Flags|      Fragment Offset    |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |  Time to Live |    Protocol   |         Header Checksum       |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                       Source Address                          |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                    Destination Address                        |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                    Options                    |    Padding    |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

type IPv4HDR struct {
	Version        uint8
	IHL            uint8
	TOS            uint8
	Length         uint16
	ID             uint16
	Flags          uint8
	FragmentOffset uint16
	TTL            uint8
	Protocol       uint8
	Checksum       uint16
	SrcIP          []byte
	DstIP          []byte
}

/*

const (
	IP_ICMP = 1
	IP_INIP = 4
	IP_TCP  = 6
	IP_UDP  = 17
)

type CiscoIDPElement struct {
	Version  uint8
	Length   uint16
	Protocol uint8
	SrcIP    []byte
	DstIP    []byte
}

func ParseCiscoIDP(pkt []byte) CiscoIDPElement {
	var r CiscoIDPElement

	//pos := 0
	r.Version = uint8(pkt[0]) >> 4
	if r.Version == 4 {
		r.Length = binary.BigEndian.Uint16(pkt[2:4])
		r.Protocol = pkt[9]
		r.SrcIP = pkt[12:16]
		r.DstIP = pkt[16:20]
		pos = 20
	}
	if r.Version == 6 {
		r.Length = binary.BigEndian.Uint16(pkt[4:6])
		r.Protocol = pkt[6] //NextHeader
		r.SrcIP = pkt[8:24]
		r.DstIP = pkt[24:40]
		//pos = 40
	}
	switch r.Protocol {
	case IP_TCP:

	}

	return r

}
*/
