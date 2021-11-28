package ipfix

import (
	"github.com/zartbot/goflow/datarecord"
	"github.com/zartbot/goflow/iedb"
	"github.com/zartbot/goflow/reader"
)

// FieldValueDecode : a fv decode function
func FieldValueDecode(r *reader.Reader, tr *TemplateFieldSpecifier) (interface{}, error) {
	var f interface{}
	var err error
	//f.Key = tr.Key
	//f.Name = tr.Name
	//f.Type = tr.Type
	//  Variable Length Field
	//	0                   1                   2                   3
	//	0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	//  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//  | Length (< 255)|          Information Element                  |
	//  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//  |                      ... continuing as needed                 |
	//  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//		  Figure R: Variable-Length Information Element (IE)
	//						 (Length < 255 Octets)
	//  0                   1                   2                   3
	//  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	//  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//  |      255      |      Length (0 to 65535)      |       IE      |
	//  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//  |                      ... continuing as needed                 |
	//  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	//          Figure S: Variable-Length Information Element (IE)
	//                      (Length 0 to 65535 Octets)
	var tempValue []byte
	if tr.Length == 65535 {
		vlen, err := r.ReadUint8()
		if err != nil {
			return f, err
		}
		if vlen == 255 {
			//Figure.S Length more than 255
			var vlen2 uint16
			vlen2, err = r.ReadUint16()
			if err != nil {
				return f, err
			}
			tempValue, err = r.ReadN(vlen2)
		} else {
			//Figure.R Length less than 255
			if vlen > 0 {
				tempValue, err = r.ReadN(uint16(vlen))
			} else {
				tempValue = []byte{}
			}
		}
		if err != nil {
			return f, err
		}
	} else {
		//Normal fixed Length Field Reader
		tempValue, err = r.ReadN(tr.Length)
		if err != nil {
			return f, err
		}
	}

	f = iedb.ConvertDataType(&tempValue, tr.Dtype)
	return f, nil
}

// DataRecordDecode : decode function
func DataRecordDecode(r *reader.Reader, AgentID string, localport int, m *MessageHeader, s *FlowSetHeader) (datarecord.DataFrame, error) {
	//d := datarecord.NewDataFrame()

	var d datarecord.DataFrame
	var err error

	var t TemplateRecord

	d.AgentID = AgentID
	d.DomainID = m.DomainID
	d.ExportTime = m.ExportTime
	d.SetID = s.SetID

	t, err = FetchTemplateMap(d.AgentID, localport, d.DomainID, d.SetID)
	if err != nil {
		//Template does not found, drop record
		return d, nil
	}

	//debug processing time
	//now := time.Now().Format("2006-01-02 15:04:05")

	initpos := r.Pos
	var byteremain uint16
	for {
		data := make(map[string]interface{})
		//debug processing time
		//data["GoFlow Processing Time"] = now
		for idx := uint16(0); idx < t.Header.ScopeFieldCount; idx++ {
			tr := t.ScopeFieldRecord[idx]
			var fv interface{}
			fv, err = FieldValueDecode(r, &tr)
			if err != nil {
				return d, err
			}

			if tr.Type == "ciscoappvarstring" {
				//multiple same type element should be merge together
				existElement, ok := data[tr.Name]
				if ok {
					e := existElement.(map[iedb.CiscoAppVarStringKey]string)
					result := fv.(map[iedb.CiscoAppVarStringKey]string)
					for k, v := range e {
						result[k] = v
					}
					fv = result
				}
			}
			data[tr.Name] = fv
		}
		for idx := uint16(0); idx < t.Header.FieldCount-t.Header.ScopeFieldCount; idx++ {
			tr := t.FieldRecord[idx]
			var fv interface{}
			fv, err = FieldValueDecode(r, &tr)
			if err != nil {
				return d, err
			}

			if tr.Type == "ciscoappvarstring" {
				//multiple same type element should be merge together
				existElement, ok := data[tr.Name]
				if ok {
					e := existElement.(map[iedb.CiscoAppVarStringKey]string)
					result := fv.(map[iedb.CiscoAppVarStringKey]string)
					for k, v := range e {
						result[k] = v
					}
					fv = result
				}
			}
			data[tr.Name] = fv
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
					return d, err
				}
			}
			break
		}
	}
	return d, nil
}
