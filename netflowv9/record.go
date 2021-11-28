package netflowv9

import (
	"github.com/zartbot/goflow/datarecord"
	"github.com/zartbot/goflow/iedb"
	"github.com/zartbot/goflow/reader"
)

// FieldValueDecode :  FV decode function
func FieldValueDecode(r *reader.Reader, tr *TemplateFieldSpecifier) (interface{}, error) {
	var f interface{}
	var err error
	//f.Key = tr.Key
	//f.Name = tr.Name
	//f.Type = tr.Type

	var tempValue []byte
	tempValue, err = r.ReadN(tr.Length)
	if err != nil {
		return f, err
	}

	f = iedb.ConvertDataType(&tempValue, tr.Dtype)
	return f, nil
}

/*
func (f *FieldValue) String() string {
	return fmt.Sprintf("ElementKey:+%v\t\tName:\t%s\t          \tType:\t%s\t           \t%+v\n", f.Key, f.Name, f.Type, f.Value)
}
*/

// DataRecordDecode : A method to parse DataRecord
func DataRecordDecode(r *reader.Reader, AgentID string, LocalPort int, m *MessageHeader, s *FlowSetHeader) (datarecord.DataFrame, error) {
	//d := datarecord.NewDataFrame()
	d := &datarecord.DataFrame{}
	var err error
	var t TemplateRecord

	d.AgentID = AgentID
	d.DomainID = m.DomainID
	d.ExportTime = m.ExportTime
	d.SetID = s.SetID

	t, err = FetchTemplateMap(d.AgentID, LocalPort, d.DomainID, d.SetID)
	if err != nil {
		//Template does not found, drop record
		//logrus.Warn("Template does not found")
		return *d, nil
	}

	initpos := r.Pos
	var byteremain uint16
	for {
		data := make(map[string]interface{})

		for idx := uint16(0); idx < t.Header.ScopeFieldCount; idx++ {
			tr := t.ScopeFieldRecord[idx]
			var fv interface{}
			fv, err := FieldValueDecode(r, &tr)
			if err != nil {
				return *d, err
			}
			data[tr.Name] = fv
		}
		for idx := uint16(0); idx < t.Header.FieldCount; idx++ {
			tr := t.FieldRecord[idx]
			var fv interface{}
			fv, err := FieldValueDecode(r, &tr)
			if err != nil {
				return *d, err
			}
			data[tr.Name] = fv
		}
		//handle NFv9 SysUpTime First/Last
		FirstTime, ok1 := data["flowStartSysUpTime"]
		LastTime, ok2 := data["flowEndSysUpTime"]

		if ok1 && ok2 && (LastTime != nil) && (FirstTime != nil) {
			first := FirstTime.(uint32)
			last := LastTime.(uint32)
			data["flowEndMilliseconds"] = uint64(last) + m.BootTime
			if first > last {
				data["flowStartMilliseconds"] = m.BootTime - uint64(100000000) + uint64(first)
			} else {
				data["flowStartMilliseconds"] = m.BootTime + uint64(first)
			}
		}

		d.Record = append(d.Record, data)

		//fmt.Printf("Recode Finished......\tFlowSetLength:%d\tRead_Pos:%d\tInitPos:%d\tAlreadyRead:%d\t\n", s.Length, r.Pos, initpos, r.Pos-initpos)

		byteremain = s.Length - 4 - (r.Pos - initpos)
		// RFC 7011 Section 3.3.1. Set Format
		// Padding
		//
		// The Exporting Process MAY insert some padding octets, so that the
		// subsequent Set starts at an aligned boundary.  For security
		// reasons, the padding octet(s) MUST be composed of octets with
		// value zero (0).  The padding length MUST be shorter than any
		// allowable record in this Set.  If padding of the IPFIX Message is
		// desired in combination with very short records, then the padding
		// Information Element 'paddingOctets' can be used for padding
		// records such that their length is increased to a multiple of 4 or
		// 8 octets.  Because Template Sets are always 4-octet aligned by
		// definition, padding is only needed in the case of other
		// alignments, e.g., on 8-octet boundaries.

		if byteremain < 4 {
			//read padding...
			if byteremain > 0 {
				_, err = r.ReadN(byteremain)
				if err != nil {
					return *d, err
				}
			}
			break
		}
	}

	return *d, nil
}
