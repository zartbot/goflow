package iedb

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type CiscoFARecord struct {
	PacketNum     uint8
	FlowDirection [16]uint8
	TCPFlag       [16]byte
	Interval      [16]uint16
}

func (fa *CiscoFARecord) String() string {

	buf := fmt.Sprintf("Packet Number:%d\n", fa.PacketNum)
	for i := uint8(0); i < fa.PacketNum; i++ {

		if fa.FlowDirection[i] == 0 {
			buf += fmt.Sprintf(" -> Flags:%20s\tLatency:%d\n", PrintTCPFlags(fa.TCPFlag[i]), fa.Interval[i])
		} else {
			buf += fmt.Sprintf(" <- Flags:%20s\tLatency:%d\n", PrintTCPFlags(fa.TCPFlag[i]), fa.Interval[i])
		}

	}
	return buf
}

//ParseCiscoFA : Parse Cisco FA packet
func ParseCiscoFA(b []byte) *CiscoFARecord {
	r := new(CiscoFARecord)

	r.PacketNum = uint8(b[0])

	if b[1]&1 == 1 {
		r.FlowDirection[8] = 1
	}

	if b[1]&2 == 2 {
		r.FlowDirection[9] = 1
	}

	if b[1]&4 == 4 {
		r.FlowDirection[10] = 1
	}

	if b[1]&8 == 8 {
		r.FlowDirection[11] = 1
	}

	if b[1]&16 == 16 {
		r.FlowDirection[12] = 1
	}

	if b[1]&32 == 32 {
		r.FlowDirection[13] = 1
	}
	if b[1]&64 == 64 {
		r.FlowDirection[14] = 1
	}
	if b[1]&128 == 128 {
		r.FlowDirection[15] = 1
	}

	if b[2]&1 == 1 {
		r.FlowDirection[0] = 1
	}

	if b[2]&2 == 2 {
		r.FlowDirection[1] = 1
	}

	if b[2]&4 == 4 {
		r.FlowDirection[2] = 1
	}

	if b[2]&8 == 8 {
		r.FlowDirection[3] = 1
	}

	if b[2]&16 == 16 {
		r.FlowDirection[4] = 1
	}

	if b[2]&32 == 32 {
		r.FlowDirection[5] = 1
	}
	if b[2]&64 == 64 {
		r.FlowDirection[6] = 1
	}
	if b[2]&128 == 128 {
		r.FlowDirection[7] = 1
	}

	for i := 0; i < 16; i++ {
		r.TCPFlag[i] = b[3+i]

	}

	for i := 0; i < 16; i++ {
		r.Interval[i] = uint16(binary.BigEndian.Uint16(b[19+2*i : 21+2*i]))
	}

	//fmt.Println(hex.Dump(b))
	//fmt.Println(r)
	/*
		for i := 0; i < 10; i++ {
			start := 2 * i
			end := 2*i + 2
			r.Length[i] = uint16(binary.BigEndian.Uint16(b[start:end]))
		}
		for i := 0; i < 10; i++ {
			start := 2*i + 20
			end := 2*i + 22
			r.Interval[i] = uint16(binary.BigEndian.Uint16(b[start:end]))
		}*/

	return r
}

func PrintTCPFlags(b byte) string {

	var buffer bytes.Buffer

	if b&1 == 1 {
		buffer.WriteString("FIN ")
	}
	if b&2 == 2 {
		buffer.WriteString("SYN ")
	}
	if b&4 == 4 {
		buffer.WriteString("RST ")
	}
	if b&8 == 8 {
		buffer.WriteString("PSH ")
	}
	if b&16 == 16 {
		buffer.WriteString("ACK ")
	}
	if b&32 == 32 {
		buffer.WriteString("URG ")
	}
	return buffer.String()
}
