package optiontemplatemap

import (
	"errors"
	"fmt"
	"sync"
)

/*
_____________________________________________________________________________
|                 Field                   |    ID | Ent.ID | Offset |  Size |
-----------------------------------------------------------------------------
| DROP CAUSE ID                           | 12442 |      9 |      0 |     2 |
| drop cause name                         | 12447 |      9 |      2 |    40 |
| drop cause desc                         | 12448 |      9 |     42 |    40 |
-----------------------------------------------------------------------------
*/

type DropCauseDataType struct {
	Name string
	//	Description string
}

type DropCauseMapKey struct {
	AgentID     string
	DropCauseID uint16
}

var DropCauseTable sync.Map

func UpdateDropCauseDatabase(AgentID string, dropcauseid uint16, name string) error {
	key := DropCauseMapKey{
		AgentID:     AgentID,
		DropCauseID: dropcauseid,
	}
	data := DropCauseDataType{
		Name: name,
		//	Description: Description,
	}
	/*
		//store will cause map lock, so try to load and compare before store
		olddata, _ := FetchDropCauseDataTypebaseByIndex(AgentID, dropcauseid)
		if olddata != data {
			DropCauseTable.Store(key, data)
		}
	*/
	DropCauseTable.Store(key, data)

	return nil
}

func FetchDropCauseDataTypebaseByIndex(AgentID string, dropcauseid uint16) (DropCauseDataType, error) {
	var err error
	var i DropCauseDataType
	key := DropCauseMapKey{
		AgentID:     AgentID,
		DropCauseID: dropcauseid,
	}

	value, ok := DropCauseTable.Load(key)
	if ok {
		i, valid := value.(DropCauseDataType)
		if !valid {
			err = errors.New("invalid type assertion")
			return i, err
		} else {
			return i, nil
		}
	} else {
		err = errors.New("DropCause does not found in database")
		return i, err
	}
}

func ShowDropCauseDataTypebase() {
	DropCauseTable.Range(func(k, v interface{}) bool {
		value := v.(DropCauseDataType)
		key := k.(DropCauseMapKey)
		fmt.Printf("DropCause Map[%10s:%-6d] | Name: %-20s | Description: %-50s \n", key.AgentID, key.DropCauseID, value.Name, value.Name)
		return true
	})
}

func UpdateDropCauseMap(d map[string]interface{}, AgentID string) error {
	var err error
	dropID, valid := d["drop_cause_id"]
	if valid {
		//isTemplate ?
		dropName, hasNameShort := d["drop_cause_name"]
		//dropDescription, hasNameLong := d["drop_cause_desc"]
		if hasNameShort {
			err = UpdateDropCauseDatabase(AgentID, dropID.(uint16), dropName.(string))
		} else {
			//data record contains DropCause field
			r, err := FetchDropCauseDataTypebaseByIndex(AgentID, dropID.(uint16))
			if err == nil {
				d["drop_cause"] = r.Name
			}
		}
	}
	return err
}
