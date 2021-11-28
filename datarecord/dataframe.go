package datarecord

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// FieldValue : define a common struct
// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |   Set ID = Template ID        |          Length               |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |   Record 1 - Field Value 1    |   Record 1 - Field Value 2    |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
type FieldValue struct {
	//Key   iedb.ElementKey
	//Name  string
	Type  string
	Value interface{}
}

// DataFrame :  Containing Data Records
// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |   Set ID = Template ID        |          Length               |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |   Record 1 - Field Value 1    |   Record 1 - Field Value 2    |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |   Record 1 - Field Value 3    |             ...               |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |   Record 2 - Field Value 1    |   Record 2 - Field Value 2    |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |   Record 2 - Field Value 3    |             ...               |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |   Record 3 - Field Value 1    |   Record 3 - Field Value 2    |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |   Record 3 - Field Value 3    |             ...               |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |              ...              |      Padding (optional)       |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

type DataFrame struct {
	AgentID    string
	DomainID   uint32
	ExportTime uint32
	SetID      uint16
	Type       string
	Record     []map[string]interface{} //use hashmap could be easily add new record.
}

func (d *DataFrame) Print(prefix string) {
	fmt.Printf("[%s] | ==========================================================DataFrame Start========================================================\n", prefix)

	fmt.Printf("[%s] | AgentIP: %-20s | ExportTime: %s | DomainID/SetID: [%6d:%-6d] \n", prefix, d.AgentID, time.Unix(int64(d.ExportTime), 0), d.DomainID, d.SetID)
	for idx := 0; idx < len(d.Record); idx++ {
		fvr := d.Record[idx]
		fmt.Printf("[%s] | Record:[%4d]--------------------------------------------------------------------------------------------------------------------\n", prefix, idx)
		for idy := range fvr {
			f := fvr[idy]
			//fmt.Printf("[%s] | Name: %-35s | Type: %-20s | Value: %+v\n", prefix, idy, f.Type, f.Value)
			fmt.Printf("[%s] | Name: %-35s | Value: %+v\n", prefix, idy, f)
		}

	}
	fmt.Printf("[%s] | ==========================================================DataFrame  End ========================================================\n\n\n", prefix)
}

func (d *DataFrame) ToJSON() (string, error) {
	r, err := json.Marshal(d)
	if err != nil {
		log.Printf("[JSON-Marshal]Error:%s", err)
		return "", err
	} else {
		return string(r), err
	}
}

//RecordList : generate record list for bulk output
func (d *DataFrame) RecordList() []map[string]interface{} {
	var recordList []map[string]interface{}
	for _, r := range d.Record {

		recordList = append(recordList, r)
	}
	return recordList
}
