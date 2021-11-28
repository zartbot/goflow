package ipfix

import (
	"fmt"
	"net"
	"time"

	"github.com/zartbot/goflow/datarecord"
	"github.com/zartbot/goflow/reader"
)

//Packet : This is a struct for IPFIX Decoder
type Packet struct {
	RemoteAddr *net.UDPAddr
	LocalPort  int
	Data       *reader.Reader
}

// NewPacket construction func
func NewPacket(remoteAddr *net.UDPAddr, localport int, b []byte, n int) *Packet {
	return &Packet{remoteAddr, localport, reader.NewReader(b, n)}
}

//MessageHeader :  RFC 7011 - Message Header Format
//
// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |       Version Number          |            Length             |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                           Export Time                         |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                       Sequence Number                         |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                    Observation Domain ID                      |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
type MessageHeader struct {
	Version    uint16
	Length     uint16
	ExportTime uint32
	SeqNo      uint32
	DomainID   uint32
}

// Decode : A method to parse MessageHeader
func (m *MessageHeader) Decode(r *reader.Reader) error {
	var err error
	if m.Version, err = r.ReadUint16(); err != nil {
		return err
	}
	if m.Version != 10 {
		err = fmt.Errorf("Invalid ipfix version (%d),expect version=10", m.Version)
		return err
	}
	if m.Length, err = r.ReadUint16(); err != nil {
		return err
	}
	if m.ExportTime, err = r.ReadUint32(); err != nil {
		return err
	}
	if m.SeqNo, err = r.ReadUint32(); err != nil {
		return err
	}
	if m.DomainID, err = r.ReadUint32(); err != nil {
		return err
	}
	return nil
}

// FlowSetHeader : RFC 7011 - 3.3.2. Set Header Format
// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |          Set ID               |          Length               |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
type FlowSetHeader struct {
	SetID      uint16
	Length     uint16
	IsTemplate bool //add istemplate for template cache handling...
}

// Decode : A method to parse FlowSetHeader
func (s *FlowSetHeader) Decode(r *reader.Reader) error {
	var err error

	if s.SetID, err = r.ReadUint16(); err != nil {
		return err
	}

	if s.Length, err = r.ReadUint16(); err != nil {
		return err
	}

	if s.SetID == 2 || s.SetID == 3 {
		s.IsTemplate = true
	} else {
		s.IsTemplate = false
	}

	return nil
}

func (m *MessageHeader) String() string {
	var result string
	result = fmt.Sprintf("Message Header: Version:%4d | Length: %6d | ExportTime: %20s | Seq: %10d | DomainID: %-4d\n", m.Version, m.Length, time.Unix(int64(m.ExportTime), 0), m.SeqNo, m.DomainID)
	return result
}

func (s *FlowSetHeader) String() string {
	var result string
	result = fmt.Sprintf("SetID:%6d | Length: %6d | IsTemplate: %t", s.SetID, s.Length, s.IsTemplate)
	return result
}

func PacketParser(r *RawMessageUDP, outputChan chan *datarecord.DataFrame) error {
	var err error
	n := len(r.body)
	packet := NewPacket(r.remoteAddr, r.localport, r.body, n)
	var msgHeader MessageHeader
	err = msgHeader.Decode(packet.Data)
	if err != nil {
		return err
	}
	flowset := 0
	for {
		var flowsetHeader FlowSetHeader
		err = flowsetHeader.Decode(packet.Data)
		if err != nil {
			return err
		}
		if flowsetHeader.IsTemplate {
			//logrus.Warn("AgentID:",packet.RemoteAddr,"|Domain-ID:",msgHeader.DomainID)
			err = HandleTemplateRecord(packet, &msgHeader, &flowsetHeader)
		} else {
			//Read RecordData to buffer
			d, err := packet.Data.ReadN(flowsetHeader.Length - 4)
			if err != nil {
				return err
			}
			//decode DataFrame
			var df datarecord.DataFrame
			rr := reader.NewReader(d, len(d))

			df, err = DataRecordDecode(rr, packet.RemoteAddr.IP.String(), packet.LocalPort, &msgHeader, &flowsetHeader)
			if err != nil {
				return err
			}
			//send to output channel
			outputChan <- &df

		}
		if packet.Data.Pos >= uint16(n) {
			break
		} else {
			flowset++
		}
	}
	return nil
}
