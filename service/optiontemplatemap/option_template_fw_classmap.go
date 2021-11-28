package optiontemplatemap

import (
	"errors"
	"fmt"
	"sync"
)

type FWClassMapMapKey struct {
	AgentID string
	classid uint8
}

var FWClassMap sync.Map

func UpdateFWClassMap(AgentID string, classid uint8, Description string) error {
	key := FWClassMapMapKey{
		AgentID: AgentID,
		classid: classid,
	}
	FWClassMap.Store(key, string(Description))

	return nil
}

func FetchFWClassMapByIndex(AgentID string, classid uint8) (string, error) {
	var err error
	var i string
	key := FWClassMapMapKey{
		AgentID: AgentID,
		classid: classid,
	}

	value, ok := FWClassMap.Load(key)
	if ok {
		i, valid := value.(string)
		if !valid {
			err = errors.New("invalid type assertion")
			return i, err
		} else {
			return i, nil
		}
	} else {
		err = errors.New("Event does not found in database")
		return i, err
	}
}

func ShowFWClassMap() {
	FWClassMap.Range(func(k, v interface{}) bool {
		value := v.(string)
		key := k.(FWClassMapMapKey)
		fmt.Printf("FW Class Map[%10s:%-6d] | Description: %-50s \n", key.AgentID, key.classid, value)
		return true
	})
}

func UpdateFWClassMapMap(d map[string]interface{}, AgentID string) error {
	var err error
	classid, hasclassid := d["classId"]
	if hasclassid {
		//isTemplate ?
		className, hasclassName := d["className"]

		if hasclassid && hasclassName {
			err = UpdateFWClassMap(AgentID, classid.(uint8), className.(string))
		} else {
			//data record contains interface field
			r, err := FetchFWClassMapByIndex(AgentID, classid.(uint8))
			if err == nil {
				//var ir fieldvalue.FieldValue
				//ir.Name = "ingressInterfaceName"
				//ir.Value = r.Description
				//ir.Type = "string"
				d["fw_class_map"] = r
			}
		}
	}
	return err
}
