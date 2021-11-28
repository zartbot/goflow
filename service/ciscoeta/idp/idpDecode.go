package idp

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func isIPv4(d map[string]interface{}) bool {
	if _, valid := d["sourceIPv4Address"]; valid {
		return true
	} else {
		return false
	}
}

func DecodeIDPField(d map[string]interface{}) error {
	var packet gopacket.Packet
	if idpRaw, valid := d["ETA_IDP"]; valid && idpRaw != nil {
		//v4v6
		if isIPv4(d) {
			packet = gopacket.NewPacket(idpRaw.([]byte), layers.LayerTypeIPv4, gopacket.Default)
		} else {
			packet = gopacket.NewPacket(idpRaw.([]byte), layers.LayerTypeIPv6, gopacket.Default)
		}

		ParseDNSPacket(d, packet)

	}
	return nil
}
