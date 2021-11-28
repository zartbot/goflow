package ipfix

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/zartbot/goflow/iedb"
	"github.com/zartbot/goflow/reader"
)

var TEMPLATE_DEBUG bool = false

// TemplateHeader : RFC 7011 - 3.4.1 Template Record Format
// [Set ID == 2]
// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |      Template ID (> 255)      |         Field Count           |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//
// OptionTemplateHeader : RFC 7011 - 3.4.2.2 Option Template Record Format
// [Set ID == 3]
// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |         Template ID (> 255)   |         Field Count           |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |      Scope Field Count        |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
type TemplateHeader struct {
	TemplateID      uint16
	FieldCount      uint16
	ScopeFieldCount uint16
}

// Decode : A method to parse TemplateHeader
func (t *TemplateHeader) Decode(r *reader.Reader, SetID uint16) error {
	var err error
	if t.TemplateID, err = r.ReadUint16(); err != nil {
		return err
	}

	if t.FieldCount, err = r.ReadUint16(); err != nil {
		return err
	}

	if SetID == 3 {
		if t.ScopeFieldCount, err = r.ReadUint16(); err != nil {
			return err
		}
	}
	return nil
}

// TemplateFieldSpecifier : RFC 7011 - 3.2. Field Specifier Format
//
// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |0|  Information Element ident. |        Field Length           |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//
// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |1|  Information Element ident. |        Field Length           |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                      Enterprise Number                        |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
type TemplateFieldSpecifier struct {
	Key    iedb.ElementKey
	Length uint16
	Name   string
	Type   string
	Dtype  iedb.DataType
}

// Decode : TemplateFiledSpecifier
func (t *TemplateFieldSpecifier) Decode(r *reader.Reader) error {
	var err error

	if t.Key.ElementID, err = r.ReadUint16(); err != nil {
		return err
	}

	if t.Length, err = r.ReadUint16(); err != nil {
		return err
	}

	if t.Key.ElementID > 0x8000 {
		t.Key.ElementID = t.Key.ElementID & 0x7fff
		if t.Key.EnterpriseNo, err = r.ReadUint32(); err != nil {
			return err
		}
	}
	value, ok := iedb.IEDatabase.Load(t.Key)
	if ok {
		ie, valid := value.(iedb.InformationElement)
		if !valid {
			err = errors.New("invalid type assertion")
			return err
		} else {
			t.Name = ie.Name
			t.Type = ie.Type
			t.Dtype = ie.Dtype
		}
	}

	return nil
}

// TemplateRecord : This is an Exmaple
// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |          Set ID = 2           |          Length               |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |      Template ID = 256        |         Field Count = N       |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |1| Information Element id. 1.1 |        Field Length 1.1       |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                    Enterprise Number  1.1                     |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |0| Information Element id. 1.2 |        Field Length 1.2       |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |             ...               |              ...              |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |1| Information Element id. 1.N |        Field Length 1.N       |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                    Enterprise Number  1.N                     |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |      Template ID = 257        |         Field Count = M       |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |0| Information Element id. 2.1 |        Field Length 2.1       |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |1| Information Element id. 2.2 |        Field Length 2.2       |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                    Enterprise Number  2.2                     |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |             ...               |              ...              |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |1| Information Element id. 2.M |        Field Length 2.M       |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                    Enterprise Number  2.M                     |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                          Padding (opt)                        |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// OptionTemplate Example : RFC 7011 - 3.4.2.2 Option Template Record Format
//
// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |          Set ID = 3           |          Length               |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |         Template ID = 258     |         Field Count = N + M   |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |     Scope Field Count = N     |0|  Scope 1 Infor. Element id. |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |     Scope 1 Field Length      |0|  Scope 2 Infor. Element id. |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |     Scope 2 Field Length      |             ...               |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |            ...                |1|  Scope N Infor. Element id. |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |     Scope N Field Length      |   Scope N Enterprise Number  ...
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//...  Scope N Enterprise Number   |1| Option 1 Infor. Element id. |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |    Option 1 Field Length      |  Option 1 Enterprise Number  ...
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//... Option 1 Enterprise Number   |              ...              |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |             ...               |0| Option M Infor. Element id. |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |     Option M Field Length     |      Padding (optional)       |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//
//   The example shows an Options Template Set with mixed IANA-assigned
//   and enterprise-specific Information Elements.  It consists of a
//   Set Header, an Options Template Header, and several Field Specifiers.
type TemplateRecord struct {
	DomainID         uint32
	SetID            uint16
	Header           TemplateHeader
	FieldRecord      []TemplateFieldSpecifier
	ScopeFieldRecord []TemplateFieldSpecifier
}

// Decode : This is a method to decode template and option template set
func (t *TemplateRecord) Decode(r *reader.Reader, m *MessageHeader, s *FlowSetHeader) error {
	var err error
	initpos := r.Pos //used to calculate length to validate Padding Field.

	t.DomainID = m.DomainID
	//t.ExportTime = m.ExportTime
	t.SetID = s.SetID

	err = t.Header.Decode(r, t.SetID)
	if err != nil {
		return err
	}

	if t.SetID == 3 {
		for fieldid := uint16(0); fieldid < t.Header.ScopeFieldCount; fieldid++ {
			var fi TemplateFieldSpecifier
			fi.Decode(r)
			t.ScopeFieldRecord = append(t.ScopeFieldRecord, fi)
		}
	}
	for fieldid := uint16(0); fieldid < t.Header.FieldCount-t.Header.ScopeFieldCount; fieldid++ {
		var fi TemplateFieldSpecifier
		fi.Decode(r)
		t.FieldRecord = append(t.FieldRecord, fi)
	}

	// calculate padding size by length
	// SetLength - 4(SetHeader) - ReadBytes(r.Pos-initpos)
	paddingSize := s.Length - 4 - (r.Pos - initpos)
	if paddingSize > 0 {
		_, err = r.ReadN(paddingSize)
		if err != nil {
			return err
		}
	}
	return nil
}

type TemplateMapKey struct {
	AgentID   string
	LocalPort int
	DomainID  uint32
	SetID     uint16
}

func NewTemplateMapKey(AgentID string, LocalPort int, DomainID uint32, SetID uint16) TemplateMapKey {
	return TemplateMapKey{
		AgentID:   AgentID,
		LocalPort: LocalPort,
		DomainID:  DomainID,
		SetID:     SetID,
	}
}

func FetchTemplateMapByKey(key TemplateMapKey) (TemplateRecord, error) {
	var err error
	var t TemplateRecord
	value, ok := TemplateMap.Load(key)
	if ok {
		t, valid := value.(TemplateRecord)
		if !valid {
			err = errors.New("invalid type assertion")
			return t, err
		} else {
			return t, nil
		}
	} else {
		err = errors.New("value not exist in map")
		return t, err
	}
}

func FetchTemplateMap(AgentID string, LocalPort int, DomainID uint32, SetID uint16) (TemplateRecord, error) {
	key := NewTemplateMapKey(AgentID, LocalPort, DomainID, SetID)
	v, err := FetchTemplateMapByKey(key)
	return v, err
}

func isEqualFS(a []TemplateFieldSpecifier, b []TemplateFieldSpecifier) bool {
	if (a == nil) != (b == nil) {
		return false
	}
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func HandleTemplateRecord(packet *Packet, msghdr *MessageHeader, flowsethdr *FlowSetHeader) error {
	var err error

	initpos := packet.Data.Pos
	var byteremain uint16

	for {
		//var oldts TemplateRecord
		var ts TemplateRecord
		err = ts.Decode(packet.Data, msghdr, flowsethdr)
		if err != nil {
			logrus.Warn("Decode Template Error:", err)
			return err
		}

		if TEMPLATE_DEBUG {
			fmt.Println("DEBUG-TEMPLATE---------------------------")
			fmt.Println(ts.String())
		}

		key := NewTemplateMapKey(packet.RemoteAddr.IP.String(), packet.LocalPort, ts.DomainID, ts.Header.TemplateID)
		TemplateMap.Store(key, ts)
		//fmt.Println(time.Now(), " Template Recieved-->", key)
		/*
			key := NewTemplateMapKey(packet.RemoteAddr.IP.String(), ts.DomainID, ts.Header.TemplateID)
			oldts, err = FetchTemplateMapByKey(key)
			if (err != nil) || (ts.DomainID != oldts.DomainID) || (ts.SetID != oldts.SetID) || !isEqualFS(ts.FieldRecord, oldts.FieldRecord) || !isEqualFS(ts.ScopeFieldRecord, oldts.ScopeFieldRecord) {
				TemplateMap.Store(key, ts)
			}
		*/

		byteremain = flowsethdr.Length - 4 - (packet.Data.Pos - initpos)
		//fmt.Printf("DEBUG-HANDLETEMPLATE:::SetLen:%d\tRPOS:%d\tINITPOS:%d\n", flowsethdr.Length, packet.Data.Pos, initpos)
		if byteremain < 4 {
			//read padding...
			if byteremain > 0 {
				_, err = packet.Data.ReadN(byteremain)
				if err != nil {
					return err
				}
			}
			break
		}
	}
	return nil
}

func (t *TemplateRecord) String() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("[IPFIX] | Template: ScopeFieldCount:%4d | FieldCount:%4d\n", t.Header.ScopeFieldCount, t.Header.FieldCount))
	for idx := uint16(0); idx < t.Header.ScopeFieldCount; idx++ {
		r := t.ScopeFieldRecord[idx]
		buf.WriteString(fmt.Sprintf("[IPFIX] | Template: Scope FieldRecord[%3d]: | TemplateID:DomainID: [%4d:%-4d] | ElementID:[%2d:%-6d] | Length:%4d | Name: %s\n", idx, t.Header.TemplateID, t.DomainID, r.Key.EnterpriseNo, r.Key.ElementID, r.Length, r.Name))
	}
	for idx := uint16(0); idx < t.Header.FieldCount-t.Header.ScopeFieldCount; idx++ {
		r := t.FieldRecord[idx]
		buf.WriteString(fmt.Sprintf("[IPFIX] | Template:       FieldRecord[%3d]: | TemplateID:DomainID: [%4d:%-4d] | ElementID:[%2d:%-6d] | Length:%4d | Name: %s\n", idx, t.Header.TemplateID, t.DomainID, r.Key.EnterpriseNo, r.Key.ElementID, r.Length, r.Name))
	}
	return buf.String()

}
