package optiontemplatemap

import (
	"encoding/binary"
	"errors"
	"fmt"
	"sync"
)

// InterfaceData : Interface Template
// _____________________________________________________________________________
// |                 Field                   |    ID | Ent.ID | Offset |  Size |
// -----------------------------------------------------------------------------
// | INTERFACE INPUT SNMP                    |    10 |        |      0 |     4 |
// | interface name short                    |    82 |        |      4 |    32 |
// | interface name long                     |    83 |        |     36 |    64 |
// -----------------------------------------------------------------------------

type InterfaceData struct {
	Name        string
	Description string
}

type InterfaceMapKey struct {
	AgentID string
	Ifindex uint32
}

var InterfaceTable sync.Map

func UpdateInterfaceDatabase(AgentID string, ifindex uint32, name string, Description string) error {
	key := InterfaceMapKey{
		AgentID: AgentID,
		Ifindex: ifindex,
	}
	data := InterfaceData{
		Name:        name,
		Description: Description,
	}
	/*
		//store will cause map lock, so try to load and compare before store
		olddata, _ := FetchInterfaceDatabaseByIndex(AgentID, ifindex)
		if olddata != data {
			InterfaceTable.Store(key, data)
		}
	*/
	InterfaceTable.Store(key, data)

	return nil
}

func FetchInterfaceDatabaseByIndex(AgentID string, ifindex uint32) (InterfaceData, error) {
	var err error
	var i InterfaceData
	key := InterfaceMapKey{
		AgentID: AgentID,
		Ifindex: ifindex,
	}

	value, ok := InterfaceTable.Load(key)
	if ok {
		i, valid := value.(InterfaceData)
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

func ShowInterfaceDatabase() {
	InterfaceTable.Range(func(k, v interface{}) bool {
		value := v.(InterfaceData)
		key := k.(InterfaceMapKey)
		fmt.Printf("Interface Map[%10s:%-6d] | Name: %-20s | Description: %-50s \n", key.AgentID, key.Ifindex, value.Name, value.Description)
		return true
	})
}

func UpdateInterfaceMap(d map[string]interface{}, AgentID string) error {
	var err error
	ingressIntf, hasInputIntf := d["ingressInterface"]
	if hasInputIntf {
		//isTemplate ?
		intfName, hasNameShort := d["interfaceName"]
		intfDescription, hasNameLong := d["interfaceDescription"]
		if hasNameShort && hasNameLong {
			err = UpdateInterfaceDatabase(AgentID, ingressIntf.(uint32), intfName.(string), intfDescription.(string))
		} else {
			//data record contains interface field
			r, err := FetchInterfaceDatabaseByIndex(AgentID, ingressIntf.(uint32))
			if err == nil {
				//var ir fieldvalue.FieldValue
				//ir.Name = "ingressInterfaceName"
				//ir.Value = r.Description
				//ir.Type = "string"
				d["ingressInterfaceName"] = r.Description
			}
		}
	}
	egressIntf, egressOk := d["egressInterface"]
	if egressOk {
		r, err := FetchInterfaceDatabaseByIndex(AgentID, egressIntf.(uint32))
		if err == nil {
			//var er fieldvalue.FieldValue
			//er.Name = "egressInterfaceName"
			//er.Type = "string"
			//er.Value = r.Description
			d["egressInterfaceName"] = r.Description
		}
	}
	if _obvID, ok := d["observationPointId"]; ok && _obvID != nil {
		if obvID, valid := _obvID.([]byte); valid {
			value := binary.BigEndian.Uint32(obvID[4:8])
			r, rerr := FetchInterfaceDatabaseByIndex(AgentID, value)
			if rerr == nil {
				d["observationInterface"] = r.Description
				d["observationPointId"] = value
			} else {
				d["observationInterface"] = "Unknown"
				d["observationPointId"] = value
			}
		}
	}
	return err
}
