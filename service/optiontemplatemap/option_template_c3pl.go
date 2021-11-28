package optiontemplatemap

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"sync"
)

/*
Client: Option classmap option table
_____________________________________________________________________________
|                 Field                   |    ID | Ent.ID | Offset |  Size |
-----------------------------------------------------------------------------
| C3PL CLASS CCE-ID                       |  8233 |      9 |      0 |     4 |
| c3pl class name                         |  8234 |      9 |      4 |   512 |
| c3pl class type                         |  8235 |      9 |    516 |   256 |
-----------------------------------------------------------------------------

*/

type C3plClassData struct {
	Name string
	Type string
}

type C3plClassKey struct {
	AgentID string
	CCEID   uint32
}

var C3plClassTable sync.Map

func UpdateC3plClassDatabase(AgentID string, cceid uint32, name string, t string) error {
	key := C3plClassKey{
		AgentID: AgentID,
		CCEID:   cceid,
	}
	data := C3plClassData{
		Name: name,
		Type: t,
	}
	C3plClassTable.Store(key, data)
	return nil
}

func FetchC3plClassDatabaseByIndex(AgentID string, cceid uint32) (C3plClassData, error) {
	var err error
	var i C3plClassData
	key := C3plClassKey{
		AgentID: AgentID,
		CCEID:   cceid,
	}

	value, ok := C3plClassTable.Load(key)
	if ok {
		i, valid := value.(C3plClassData)
		if !valid {
			err = errors.New("invalid type assertion")
			return i, err
		} else {
			return i, nil
		}
	} else {
		err = errors.New("Interface does not found in database")
		return i, err
	}
}

func ShowClassMapDatabase() {
	C3plClassTable.Range(func(k, v interface{}) bool {
		value := v.(C3plClassData)
		key := k.(C3plClassKey)
		fmt.Printf("C3PL Class Map[%10s:%-6d] | Name: %-20s | Type: %-50s \n", key.AgentID, key.CCEID, value.Name, value.Type)
		return true
	})
}

/*
Client: Option policymap option table
_____________________________________________________________________________
|                 Field                   |    ID | Ent.ID | Offset |  Size |
-----------------------------------------------------------------------------
| C3PL POLICY CCE-ID                      |  8236 |      9 |      0 |     4 |
| c3pl policy name                        |  8237 |      9 |      4 |   512 |
| c3pl policy type                        |  8238 |      9 |    516 |   256 |
-----------------------------------------------------------------------------
*/
type C3plPolicyData struct {
	Name string
	Type string
}

type C3plPolicyKey struct {
	AgentID string
	CCEID   uint32
}

var C3plPolicyTable sync.Map

func UpdateC3plPolicyDatabase(AgentID string, cceid uint32, name string, t string) error {
	key := C3plPolicyKey{
		AgentID: AgentID,
		CCEID:   cceid,
	}
	data := C3plPolicyData{
		Name: name,
		Type: t,
	}
	C3plPolicyTable.Store(key, data)
	return nil
}

func FetchC3plPolicyDatabaseByIndex(AgentID string, cceid uint32) (C3plPolicyData, error) {
	var err error
	var i C3plPolicyData
	key := C3plPolicyKey{
		AgentID: AgentID,
		CCEID:   cceid,
	}

	value, ok := C3plPolicyTable.Load(key)
	if ok {
		i, valid := value.(C3plPolicyData)
		if !valid {
			err = errors.New("invalid type assertion")
			return i, err
		} else {
			return i, nil
		}
	} else {
		err = errors.New("Interface does not found in database")
		return i, err
	}
}

func ShowPolicyMapDatabase() {
	C3plPolicyTable.Range(func(k, v interface{}) bool {
		value := v.(C3plPolicyData)
		key := k.(C3plPolicyKey)
		fmt.Printf("C3PL Policy Map[%10s:%-6d] | Name: %-20s | Type: %-50s \n", key.AgentID, key.CCEID, value.Name, value.Type)
		return true
	})
}

func ParseClassificationHierarchy(b []byte, AgentID string) string {
	var result string
	for idx := 0; idx < len(b); idx = idx + 4 {
		value := binary.BigEndian.Uint32(b[idx : idx+4])
		if value != 0 {
			if idx == 0 {
				r, err := FetchC3plPolicyDatabaseByIndex(AgentID, value)
				if err == nil {
					result = r.Name
				} else {
					return ""
				}
			} else {
				r, err := FetchC3plClassDatabaseByIndex(AgentID, value)
				if err == nil {
					result = result + "|->" + r.Name
				} else {
					return ""
				}
			}
		} else {
			break
		}
	}
	return result
}

func UpdateC3PLMap(d map[string]interface{}, AgentID string) error {
	var err error

	/*
		v, errv := map2value.MapToUInt32(d, "policy_qos_queue_index")
		cd, errcd := map2value.MapToUInt64(d, "policy_qos_queue_drops")

		if errv == nil && errcd == nil {
			logrus.Warn(v, cd)
		}
	*/
	classID, hasClassID := d["c3pl_class_cce_id"]
	if hasClassID {
		//isTemplate ?
		className, hasName := d["c3pl_class_name"]
		classType, hasType := d["c3pl_class_type"]
		if hasName && hasType {
			n := bytes.IndexByte(className.([]byte), 0)
			namestr := string(className.([]byte)[:n])

			tn := bytes.IndexByte(classType.([]byte), 0)
			typestr := string(classType.([]byte)[:tn])

			err = UpdateC3plClassDatabase(AgentID, classID.(uint32), namestr, typestr)
			//fmt.Println("DEBUG----->:", AgentID, classID.(uint32), className.(string), classType.(string))
		} else {
			//data record contains class field
			r, err := FetchC3plClassDatabaseByIndex(AgentID, classID.(uint32))
			if err == nil {
				d["class_map_name"] = r.Name
				d["class_map_type"] = r.Type
			}
		}
	}
	policyID, hasPolicyID := d["c3pl_policy_cce_id"]
	if hasPolicyID {
		//isTemplate ?
		policyName, hasName := d["c3pl_policy_name"]
		policyType, hasType := d["c3pl_policy_type"]
		if hasName && hasType {
			n := bytes.IndexByte(policyName.([]byte), 0)
			namestr := string(policyName.([]byte)[:n])

			tn := bytes.IndexByte(policyType.([]byte), 0)
			typestr := string(policyType.([]byte)[:tn])

			err = UpdateC3plPolicyDatabase(AgentID, policyID.(uint32), namestr, typestr)
			//fmt.Println("DEBUG--POLICY-MAP--->:", AgentID, policyID.(uint32), policyName.(string), policyType.(string))
		} else {
			//data record contains policy map field
			r, err := FetchC3plPolicyDatabaseByIndex(AgentID, policyID.(uint32))
			if err == nil {
				d["policy_map_name"] = r.Name
				d["policy_map_type"] = r.Type
			}
		}
	}

	if policyarray, ok := d["policy_qos_classification_hierarchy"]; ok {
		d["policy_hierarchy_array"] = ParseClassificationHierarchy(policyarray.([]byte), AgentID)
		delete(d, "policy_qos_classification_hierarchy")
		//fmt.Printf("DEBUG:::------->:::::::POLICYARRAY:%+v\t %+v\n", policyarray.([]byte), r)
	}

	/*
		if a, ok := d["policy_qos_queue_drops"]; ok {
			fmt.Println(a)
		}
	*/

	return err
}
