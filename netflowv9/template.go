package netflowv9

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/zartbot/goflow/iedb"
	"github.com/zartbot/goflow/reader"
)

// TemplateHeader : RFC 3954 - 5.2 Template FlowSet Format
// [Set ID == 0]
// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |         Template ID           |         Field Count           |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//
// OptionTemplateHeader : RFC 7011 - 3.4.2.2 Option Template Record Format
// [Set ID == 1]
// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |         Template ID           |      Option Scope Length      |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |        Option Length          |
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

	if SetID == 1 {
		t.ScopeFieldCount = t.FieldCount / 4
		if t.FieldCount, err = r.ReadUint16(); err != nil {
			return err
		}
		t.FieldCount = t.FieldCount / 4
	}

	//logrus.Warn(t.TemplateID, "--------------", t.FieldCount, "SCOPE:", t.ScopeFieldCount)

	return nil
}

// TemplateFieldSpecifier : RFC 7011 - 3.2. Field Specifier Format
//
// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |   Information Element ident.  |        Field Length           |
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

	t.Key.EnterpriseNo = 0
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

// TemplateRecord : 5.2 Template FlowSet Format
//
// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |       FlowSet ID = 0          |          Length               |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |      Template ID 256          |         Field Count           |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |        Field Type 1           |         Field Length 1        |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |        Field Type 2           |         Field Length 2        |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |             ...               |              ...              |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |        Field Type N           |         Field Length N        |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |      Template ID 257          |         Field Count           |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |        Field Type 1           |         Field Length 1        |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |        Field Type 2           |         Field Length 2        |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |             ...               |              ...              |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |        Field Type M           |         Field Length M        |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |             ...               |              ...              |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |        Template ID K          |         Field Count           |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |             ...               |              ...              |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//
//
// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |       FlowSet ID = 1          |          Length               |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |         Template ID           |      Option Scope Length      |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |        Option Length          |       Scope 1 Field Type      |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |     Scope 1 Field Length      |               ...             |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |     Scope N Field Length      |      Option 1 Field Type      |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |     Option 1 Field Length     |             ...               |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |     Option M Field Length     |           Padding             |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

type TemplateRecord struct {
	DomainID         uint32
	ExportTime       uint32
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
	t.ExportTime = m.ExportTime
	t.SetID = s.SetID

	err = t.Header.Decode(r, t.SetID)
	if err != nil {
		return err
	}

	if t.SetID == 1 {
		for fieldid := uint16(0); fieldid < t.Header.ScopeFieldCount; fieldid++ {
			var fi TemplateFieldSpecifier
			fi.Decode(r)
			t.ScopeFieldRecord = append(t.ScopeFieldRecord, fi)
		}
	}
	for fieldid := uint16(0); fieldid < t.Header.FieldCount; fieldid++ {
		var fi TemplateFieldSpecifier
		fi.Decode(r)
		t.FieldRecord = append(t.FieldRecord, fi)
	}

	// calculate padding size by length
	// SetLength - 4(SetHeader) - ReadBytes(r.Pos-initpos)
	//fmt.Printf("DEBUG:::SetLen:%d\tRPOS:%d\tINITPOS:%d\n", s.Length, r.Pos, initpos)

	paddingSize := s.Length - 4 - (r.Pos - initpos)
	/* Debug Padding
	logrus.Warn("PaddingSize  s.Length - 4 - (r.Pos - initpos)")
	fmt.Printf("DEBUG:::Padding Size:%d\t\tSetLen:%d\tRPOS:%d\tINITPOS:%d\n", paddingSize, s.Length, r.Pos, initpos)
	*/
	if (paddingSize < 4) && (paddingSize > 0) {
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
			logrus.Warn("[Template Parse Error]:=", err)
			logrus.Warn("[Template Parse Error]:=Template Record", ts)
			return err
		}

		key := NewTemplateMapKey(packet.RemoteAddr.IP.String(), packet.LocalPort, ts.DomainID, ts.Header.TemplateID)
		TemplateMap.Store(key, ts)
		//fmt.Println("DEBUG-NFv9-HandleTemplate-->", key)
		/*
			seems compare record may take more CPU resource than force update...
			oldts, err = FetchTemplateMapByKey(key)
			if (err != nil) || (ts.DomainID != oldts.DomainID) || (ts.SetID != oldts.SetID) || !isEqualFS(ts.FieldRecord, oldts.FieldRecord) || !isEqualFS(ts.ScopeFieldRecord, oldts.ScopeFieldRecord) {
				TemplateMap.Store(key, ts)
			}

				fmt.Printf("DEBUG OLD----------------------%+v\n", key)
				fmt.Println(oldts.String())
				fmt.Println("DEBUG New----------------------")
				fmt.Println(ts.String())
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
	buf.WriteString(fmt.Sprintf("[NFv9 ] | Template: ScopeFieldCount:%4d | FieldCount:%4d\n", t.Header.ScopeFieldCount, t.Header.FieldCount))
	for idx := uint16(0); idx < t.Header.ScopeFieldCount; idx++ {
		r := t.ScopeFieldRecord[idx]
		buf.WriteString(fmt.Sprintf("[NFv9 ] | Template: Scope FieldRecord[%3d]: | TemplateID:DomainID: [%4d:%-4d] | ElementID:[%2d:%-6d] | Length:%4d | Name: %s\n", idx, t.Header.TemplateID, t.DomainID, r.Key.EnterpriseNo, r.Key.ElementID, r.Length, r.Name))
	}
	for idx := uint16(0); idx < t.Header.FieldCount; idx++ {
		r := t.FieldRecord[idx]
		buf.WriteString(fmt.Sprintf("[NFv9 ] | Template:       FieldRecord[%3d]: | TemplateID:DomainID: [%4d:%-4d] | ElementID:[%2d:%-6d] | Length:%4d | Name: %s\n", idx, t.Header.TemplateID, t.DomainID, r.Key.EnterpriseNo, r.Key.ElementID, r.Length, r.Name))
	}
	return buf.String()

}
