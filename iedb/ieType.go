package iedb

import (
	"bytes"
	"encoding/binary"
	"math"
	"net"
	"strings"
)

type ElementKey struct {
	EnterpriseNo uint32
	ElementID    uint16
}

type InformationElement struct {
	Name  string
	Type  string
	Dtype DataType
}

type DataType int

const (
	//https://www.iana.org/assignments/ipfix/ipfix-information-element-data-types.csv
	OctetArray DataType = iota
	Unsigned8
	Unsigned16
	Unsigned32
	Unsigned64
	Signed8
	Signed16
	Signed32
	Signed64
	Float32
	Float64
	Boolean
	MacAddress
	String
	DateTimeSeconds
	DateTimeMilliseconds
	DateTimeMicroseconds
	DateTimeNanoseconds
	Ipv4Address
	Ipv6Address
	BasicList
	SubTemplateList
	SubTemplateMultiList
	CiscoAppVarString
	CiscoURLHits
	TcpFlag
	CiscoETA_SPLT
	//CiscoETA_IDP
	CiscoFA
	dyanmic2B4B
)

var DataTypeMap = map[string]DataType{
	"octetArray":           OctetArray,
	"unsigned8":            Unsigned8,
	"unsigned16":           Unsigned16,
	"unsigned32":           Unsigned32,
	"unsigned64":           Unsigned64,
	"signed8":              Signed8,
	"signed16":             Signed16,
	"signed32":             Signed32,
	"signed64":             Signed64,
	"float32":              Float32,
	"float64":              Float64,
	"boolean":              Boolean,
	"macAddress":           MacAddress,
	"string":               String,
	"dateTimeSeconds":      DateTimeSeconds,
	"dateTimeMilliseconds": DateTimeMilliseconds,
	"dateTimeMicroseconds": DateTimeMicroseconds,
	"dateTimeNanoseconds":  DateTimeNanoseconds,
	"ipv4Address":          Ipv4Address,
	"ipv6Address":          Ipv6Address,
	"basicList":            BasicList,
	"subTemplateList":      SubTemplateList,
	"subTemplateMultiList": SubTemplateMultiList,
	"ciscoappvarstring":    CiscoAppVarString,
	"ciscourlhits":         CiscoURLHits,
	"tcpflag":              TcpFlag,
	"CiscoETA_SPLT":        CiscoETA_SPLT,
	"ciscofa":              CiscoFA,
	"dyanmic2B4B":          dyanmic2B4B,
	//"ciscoETA_IDP":          CiscoETA_IDP,
}

func ConvertDataType(b *[]byte, t DataType) interface{} {
	length := len(*b)
	if length < t.minLen() {
		return *b
	}

	switch t {
	case Boolean:
		return (*b)[0] == 1
	case Unsigned8:
		return (*b)[0]
	case Unsigned16:
		return binary.BigEndian.Uint16(*b)
	case Unsigned32:
		return binary.BigEndian.Uint32(*b)
	case Unsigned64:
		if length == 4 {
			return binary.BigEndian.Uint32(*b)
		} else if length == 8 {
			return binary.BigEndian.Uint64(*b)
		} else {
			return *b
		}
	case Signed8:
		return int8((*b)[0])
	case Signed16:
		return int16(binary.BigEndian.Uint16(*b))
	case Signed32:
		return int32(binary.BigEndian.Uint32(*b))
	case Signed64:
		return int64(binary.BigEndian.Uint64(*b))
	case Float32:
		return math.Float32frombits(binary.BigEndian.Uint32(*b))
	case Float64:
		return math.Float64frombits(binary.BigEndian.Uint64(*b))
	case MacAddress:
		return net.HardwareAddr(*b)
	case String:
		return strings.TrimSpace(string(bytes.Trim(*b, "\x00")))
	case Ipv4Address, Ipv6Address:
		return net.IP(*b)
	case DateTimeSeconds:
		return binary.BigEndian.Uint32(*b)
	case DateTimeMilliseconds, DateTimeMicroseconds, DateTimeNanoseconds:
		if length == 4 {
			return binary.BigEndian.Uint32(*b)
		} else if length == 8 {
			return binary.BigEndian.Uint64(*b)
		} else {
			return *b
		}
	case dyanmic2B4B:
		if length == 2 {
			return uint32(binary.BigEndian.Uint16(*b))
		} else if length == 4 {
			return uint32(binary.BigEndian.Uint32(*b))
		} else {
			return *b
		}
	case OctetArray:
		return *b
	case CiscoAppVarString:
		return ParseCiscoAppVarString(*b)
	case CiscoURLHits:
		return ParseCiscoURLHitString(*b)
	case TcpFlag:
		return ParseTcpControlBits(*b)
	case CiscoETA_SPLT:
		return ParseCiscoETA_SPLT(*b)
	//case CiscoETA_IDP:
	//	return ParseCiscoETA_IDP(*b)
	case CiscoFA:
		return ParseCiscoFA(*b)
	}

	return *b
}

func (t DataType) minLen() int {
	switch t {
	case Boolean:
		return 1
	case Unsigned8, Signed8:
		return 1
	case Signed16, Unsigned16:
		return 2
	case Unsigned32, Signed32, Float32:
		return 4
	case Unsigned64: //bugfix for cisco mixed collect counter bytes/packet [long]
		return 4
	case Signed64, Float64:
		return 8
	case DateTimeSeconds, DateTimeMilliseconds, DateTimeMicroseconds, DateTimeNanoseconds:
		return 4
	case MacAddress:
		return 6
	case Ipv4Address:
		return 4
	case Ipv6Address:
		return 16
	case CiscoAppVarString:
		return 6
	case TcpFlag:
		return 1
	case CiscoETA_SPLT:
		return 40
	//case CiscoETA_IDP:
	//	return 1
	default:
		return 0
	}
}
