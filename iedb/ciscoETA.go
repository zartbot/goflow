package iedb

import (
	"encoding/binary"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// ParseCiscoETAIDP : Parse Cisco ETA IDP packet
func ParseCiscoETA_IDP(b []byte) *gopacket.Packet {
	ethP := gopacket.NewPacket(b, layers.LayerTypeIPv4, gopacket.Default)
	//logrus.Warn("ETA-IDP", ethP)
	return &ethP
}

type CiscoETASPLT struct {
	Length   [10]uint16
	Interval [10]uint16
}

//ParseCiscoETA_SPLT : Parse Cisco ETA SPLT packet
func ParseCiscoETA_SPLT(b []byte) *CiscoETASPLT {
	r := new(CiscoETASPLT)
	for i := 0; i < 10; i++ {
		start := 2 * i
		end := 2*i + 2
		r.Length[i] = uint16(binary.BigEndian.Uint16(b[start:end]))
	}
	for i := 0; i < 10; i++ {
		start := 2*i + 20
		end := 2*i + 22
		r.Interval[i] = uint16(binary.BigEndian.Uint16(b[start:end]))
	}
	return r
}
