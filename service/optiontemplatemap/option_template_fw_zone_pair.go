package optiontemplatemap

import (
	"errors"
	"fmt"
	"sync"
)

type FWZonePairMapKey struct {
	AgentID    string
	zonepairid uint32
}

var FWZonePair sync.Map

func UpdateFWZonePair(AgentID string, zonepairid uint32, Description []byte) error {
	key := FWZonePairMapKey{
		AgentID:    AgentID,
		zonepairid: zonepairid,
	}
	FWZonePair.Store(key, string(Description))

	return nil
}

func FetchFWZonePairByIndex(AgentID string, zonepairid uint32) (string, error) {
	var err error
	var i string
	key := FWZonePairMapKey{
		AgentID:    AgentID,
		zonepairid: zonepairid,
	}

	value, ok := FWZonePair.Load(key)
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

func ShowFWZonePair() {
	FWZonePair.Range(func(k, v interface{}) bool {
		value := string(v.([]byte))
		key := k.(FWZonePairMapKey)
		fmt.Printf("FW Zonepair Map[%10s:%-6d] | Description: %-50s \n", key.AgentID, key.zonepairid, value)
		return true
	})
}

func UpdateFWZonePairMap(d map[string]interface{}, AgentID string) error {
	var err error
	zonepairid, haszonepairid := d["policy_firewall_zone_pair_id"]
	if haszonepairid {
		//isTemplate ?
		zonePairName, haszonePairName := d["policy_firewall_zone_pair_name"]

		if haszonepairid && haszonePairName {
			err = UpdateFWZonePair(AgentID, zonepairid.(uint32), zonePairName.([]byte))
		} else {
			//data record contains interface field
			r, err := FetchFWZonePairByIndex(AgentID, zonepairid.(uint32))
			if err == nil {
				//var ir fieldvalue.FieldValue
				//ir.Name = "ingressInterfaceName"
				//ir.Value = r.Description
				//ir.Type = "string"
				d["fw_zone_pair"] = r
			}
		}
	}
	return err
}
