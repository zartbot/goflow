package iedb

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

var IEDatabase sync.Map

/* CSV Format
r[0]//ElementID
r[1]//EnterpriseNo
r[2]//FieldType(name)
r[3]//DataType
*/

func ReadCiscoIE(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Println("[IEDB Parser] Error:", err)
		return
	}
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	for idx := 1; idx < len(records); idx++ {
		r := records[idx]
		var key ElementKey
		var nfv9Key ElementKey
		var ie InformationElement
		ieid, _ := strconv.ParseUint(r[0], 0, 16)
		enno, _ := strconv.ParseUint(r[1], 0, 32)
		key.ElementID = uint16(ieid)
		key.EnterpriseNo = uint32(enno)
		//support NFv9
		nfv9Key.ElementID = uint16(ieid) + 32768
		nfv9Key.EnterpriseNo = 0

		ie.Name = strings.TrimSpace(r[2])
		ie.Type = strings.TrimSpace(r[3])
		ie.Dtype = DataTypeMap[strings.TrimSpace(r[3])]
		IEDatabase.Store(key, ie)
		IEDatabase.Store(nfv9Key, ie)

		/* Temp fix ISR1100 Enterprise Number endian issue
		if key.EnterpriseNo == 9 {
			var EndianIssue ElementKey
			EndianIssue.EnterpriseNo = 150994944
			EndianIssue.ElementID = key.ElementID
			IEDatabase.Store(EndianIssue, ie)
		}*/

	}
}

func ReadIANA(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Println("[IEDB Parser] Error:", err)
		return
	}
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	for idx := 1; idx < len(records); idx++ {
		r := records[idx]
		var key ElementKey
		var ie InformationElement
		ieid, _ := strconv.ParseUint(r[0], 0, 16)
		enno, _ := strconv.ParseUint(r[1], 0, 32)
		key.ElementID = uint16(ieid)
		key.EnterpriseNo = uint32(enno)
		ie.Name = strings.TrimSpace(r[2])
		ie.Type = strings.TrimSpace(r[3])

		ie.Dtype = DataTypeMap[strings.TrimSpace(r[3])]
		IEDatabase.Store(key, ie)

	}
}

func ShowIEDB() {
	log.SetPrefix("[IE Database]:")
	IEDatabase.Range(func(k, v interface{}) bool {
		value := v.(InformationElement)
		key := k.(ElementKey)
		log.Printf("MapKey[%2d:%-6d] | Name: %-50s | Type: %-20s \n", key.EnterpriseNo, key.ElementID, value.Name, value.Type)
		return true
	})
}
