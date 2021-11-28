package optiontemplatemap

import (
	"errors"
	"fmt"
	"sync"
)

type FWEventMapKey struct {
	AgentID string
	EventID uint32
}

var FWEvent sync.Map

func UpdateFWEvent(AgentID string, eventid uint32, Description string) error {
	key := FWEventMapKey{
		AgentID: AgentID,
		EventID: eventid,
	}
	FWEvent.Store(key, Description)

	return nil
}

func FetchFWEventByIndex(AgentID string, eventid uint32) (string, error) {
	var err error
	var i string
	key := FWEventMapKey{
		AgentID: AgentID,
		EventID: eventid,
	}

	value, ok := FWEvent.Load(key)
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

func ShowFWEvent() {
	FWEvent.Range(func(k, v interface{}) bool {
		value := v.(string)
		key := k.(FWEventMapKey)
		fmt.Printf("FW Event Map[%10s:%-6d] | Description: %-50s \n", key.AgentID, key.EventID, value)
		return true
	})
}

func UpdateFWEventMap(d map[string]interface{}, AgentID string) error {
	var err error
	eventID, hasEventID := d["policy_firewall_event_extended"]
	if hasEventID {
		//isTemplate ?
		eventDescription, hasEventDescription := d["policy_firewall_event_extended_description"]

		if hasEventID && hasEventDescription {
			err = UpdateFWEvent(AgentID, eventID.(uint32), eventDescription.(string))
		} else {
			//data record contains interface field
			r, err := FetchFWEventByIndex(AgentID, eventID.(uint32))
			if err == nil {
				//var ir fieldvalue.FieldValue
				//ir.Name = "ingressInterfaceName"
				//ir.Value = r.Description
				//ir.Type = "string"
				d["firewall_event"] = r
			}
		}
	}
	return err
}
