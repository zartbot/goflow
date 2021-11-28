package iedb

import "bytes"

type TcpCtrlBits struct {
	//FlagMap    map[string]int
	FlagName    string
	RawCtrlBits byte
}

func ParseTcpControlBits(b []byte) TcpCtrlBits {
	var result TcpCtrlBits

	var buffer bytes.Buffer
	//result.FlagMap = make(map[string]int)
	if len(b) == 1 {
		result.RawCtrlBits = b[0]
	} else if len(b) == 2 {
		result.RawCtrlBits = b[1]
	} else {
		return result
	}
	if result.RawCtrlBits&1 == 1 {
		//result.FlagMap["FIN"] = 1
		buffer.WriteString("FIN ")
	}
	if result.RawCtrlBits&2 == 2 {
		//result.FlagMap["SYN"] = 1
		buffer.WriteString("SYN ")
	}
	if result.RawCtrlBits&4 == 4 {
		//result.FlagMap["RST"] = 1
		buffer.WriteString("RST ")
	}
	if result.RawCtrlBits&8 == 8 {
		//result.FlagMap["PSH"] = 1
		buffer.WriteString("PSH ")
	}
	if result.RawCtrlBits&16 == 16 {
		//result.FlagMap["ACK"] = 1
		buffer.WriteString("ACK ")
	}
	if result.RawCtrlBits&32 == 32 {
		//result.FlagMap["URG"] = 1
		buffer.WriteString("URG ")
	}
	if result.RawCtrlBits&64 == 64 {
		//result.FlagMap["ECE"] = 1
		buffer.WriteString("ECE ")
	}
	if result.RawCtrlBits&128 == 128 {
		//result.FlagMap["CWR"] = 1
		buffer.WriteString("CWR ")
	}
	result.FlagName = buffer.String()
	return result
}
