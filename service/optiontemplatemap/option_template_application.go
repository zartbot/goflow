package optiontemplatemap

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/zartbot/goflow/iedb"
)

type ApplicationMapKey struct {
	AgentID string
	AppID   uint32
}

// ApplicationName :
// _____________________________________________________________________________
// |                 Field                   |    ID | Ent.ID | Offset |  Size |
// -----------------------------------------------------------------------------
// | APPLICATION ID                          |    95 |        |      0 |     4 |
// | application name                        |    96 |        |      4 |    24 |
// | application description                 |    94 |        |     28 |    55 |
// -----------------------------------------------------------------------------
type ApplicationNameDataType struct {
	Name        string
	Description string
}

var ApplicationNameTable sync.Map

func UpdateAppNameDatabase(AgentID string, appid uint32, appName ApplicationNameDataType) error {

	key := ApplicationMapKey{
		AgentID: AgentID,
		AppID:   appid,
	}
	/*
		//store will cause map lock, so try to load and compare before store
		olddata, _ := FetchAppNameDatabaseByID(AgentID, appid)
		if olddata != appName {
			ApplicationNameTable.Store(key, appName)
		}
	*/

	ApplicationNameTable.Store(key, appName)
	return nil
}

func FetchAppNameDatabaseByID(AgentID string, appid uint32) (ApplicationNameDataType, error) {
	var err error
	var i ApplicationNameDataType
	key := ApplicationMapKey{
		AgentID: AgentID,
		AppID:   appid,
	}

	value, ok := ApplicationNameTable.Load(key)
	if ok {
		i, valid := value.(ApplicationNameDataType)
		if !valid {
			err = errors.New("invalid type assertion")
			return i, err
		} else {
			return i, nil
		}
	} else {
		err = errors.New("Application does not found in database")
		return i, err
	}
}

func ShowApplicationNameDatabase() {
	ApplicationNameTable.Range(func(k, v interface{}) bool {
		value := v.(ApplicationNameDataType)
		key := k.(ApplicationMapKey)
		fmt.Printf("Application Map[%10s:%-6d] | %v \n", key.AgentID, key.AppID, value)
		return true
	})
}

// ApplicationAttribute : Application Template
// _____________________________________________________________________________
// |                 Field                   |    ID | Ent.ID | Offset |  Size |
// -----------------------------------------------------------------------------
// | APPLICATION ID                          |    95 |        |      0 |     4 |
// | application category name               | 12232 |      9 |      4 |    32 |
// | application sub category name           | 12233 |      9 |     36 |    32 |
// | application group name                  | 12234 |      9 |     68 |    32 |
// | application traffic-class               | 12243 |      9 |    100 |    32 |
// | application business-relevance          | 12244 |      9 |    132 |    32 |
// | p2p technology                          |   288 |        |    164 |    10 |
// | tunnel technology                       |   289 |        |    174 |    10 |
// | encrypted technology                    |   290 |        |    184 |    10 |
// | application set name                    | 12231 |      9 |    194 |    32 |
// | application family name                 | 12230 |      9 |    226 |    32 |
// -----------------------------------------------------------------------------

type ApplicationAttributeDataType struct {
	Category          string
	SubCategory       string
	GroupName         string
	TrafficClass      string
	BusinessRelevance string
	p2pTech           string
	TunnelTech        string
	EncryptTech       string
	SetName           string
	FamilyName        string
}

var ApplicationAttributeTable sync.Map

func UpdateAppAttributeDatabase(AgentID string, appid uint32, appAttr ApplicationAttributeDataType) error {

	key := ApplicationMapKey{
		AgentID: AgentID,
		AppID:   appid,
	}
	/*
		store will cause map lock, so try to load and compare before store
		olddata, _ := FetchAppAttributeDatabaseByID(AgentID, appid)
		if olddata != appAttr {
			ApplicationAttributeTable.Store(key, appAttr)
		}
	*/
	ApplicationAttributeTable.Store(key, appAttr)

	return nil

}

func FetchAppAttributeDatabaseByID(AgentID string, appid uint32) (ApplicationAttributeDataType, error) {
	var err error
	var i ApplicationAttributeDataType
	key := ApplicationMapKey{
		AgentID: AgentID,
		AppID:   appid,
	}

	value, ok := ApplicationAttributeTable.Load(key)
	if ok {
		i, valid := value.(ApplicationAttributeDataType)
		if !valid {
			err = errors.New("invalid type assertion")
			return i, err
		} else {
			return i, nil
		}
	} else {
		err = errors.New("AppAttribute does not found in database")
		return i, err
	}
}

func ShowApplicationAttributeDatabase() {
	ApplicationAttributeTable.Range(func(k, v interface{}) bool {
		value := v.(ApplicationAttributeDataType)
		key := k.(ApplicationMapKey)
		fmt.Printf("App Attribute Map[%10s:%-6d] | %v \n", key.AgentID, key.AppID, value)
		return true
	})
}

//
// Client: Option sub-application-table
// Template layout
// _____________________________________________________________________________
// |                 Field                   |    ID | Ent.ID | Offset |  Size |
// -----------------------------------------------------------------------------
// | APPLICATION ID                          |    95 |        |      0 |     4 |
// | SUB APPLICATION TAG                     |    97 |        |      4 |     4 |
// | sub application name                    |   109 |        |      8 |    80 |
// | sub application description             |   110 |        |     88 |    80 |
// -----------------------------------------------------------------------------
type SubApplicationMapKey struct {
	AgentID   string
	AppID     uint32
	SubAppTag uint32
}

var SubApplicationNameTable sync.Map

func UpdateSubAppNameDatabase(AgentID string, appid uint32, subapptag uint32, subAppName string) error {

	key := SubApplicationMapKey{
		AgentID:   AgentID,
		AppID:     appid,
		SubAppTag: subapptag,
	}

	SubApplicationNameTable.Store(key, subAppName)
	return nil
}

func FetchSubAppNameDatabaseByID(AgentID string, appid uint32, subapptag uint32) (string, error) {
	var err error
	var i string
	key := SubApplicationMapKey{
		AgentID:   AgentID,
		AppID:     appid,
		SubAppTag: subapptag,
	}

	value, ok := SubApplicationNameTable.Load(key)
	if ok {
		i, valid := value.(string)
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

func ShowSubApplicationNameDatabase() {
	SubApplicationNameTable.Range(func(k, v interface{}) bool {
		value := v.(string)
		key := k.(SubApplicationMapKey)
		fmt.Printf("Interface Map[%10s:%-6d:%-6d] | %v \n", key.AgentID, key.AppID, key.SubAppTag, value)
		return true
	})
}

func UpdateAppMap(d map[string]interface{}, AgentID string) error {
	var err error
	appid, hasappid := d["applicationId"]
	if hasappid {
		// applicationAttributeTable
		Category, has1 := d["applicationCategory"]
		SubCategory, has2 := d["applicationSubCategory"]

		// applicationNameTable
		appName, has11 := d["applicationName"]
		appDescription, has12 := d["applicationDescription"]
		// subApplicationTable
		subAppTag, has21 := d["subApplicationTag"]
		subAppName, has22 := d["subApplicationName"]

		if has1 && has2 {
			GroupName, has3 := d["applicationGroup"]
			TrafficClass, has4 := d["applicationTrafficClass"]
			BusinessRelevance, has5 := d["applicationBusinessRelevance"]
			p2pTech, has6 := d["p2pTechnology"]
			TunnelTech, has7 := d["tunnelTechnology"]
			EncryptTech, has8 := d["encryptedTechnology"]
			SetName, has9 := d["applicationSet"]
			FamilyName, has10 := d["applicationFamilyName"]
			if has3 && has4 && has5 && has6 && has7 && has8 && has9 && has10 {
				var a ApplicationAttributeDataType
				a.Category = Category.(string)
				a.SubCategory = SubCategory.(string)
				a.GroupName = GroupName.(string)
				a.TrafficClass = TrafficClass.(string)
				a.BusinessRelevance = BusinessRelevance.(string)
				a.p2pTech = p2pTech.(string)
				a.TunnelTech = TunnelTech.(string)
				a.EncryptTech = EncryptTech.(string)
				a.SetName = SetName.(string)
				a.FamilyName = FamilyName.(string)
				err = UpdateAppAttributeDatabase(AgentID, appid.(uint32), a)
			}
		} else if has11 && has12 {
			var a ApplicationNameDataType
			a.Name = appName.(string)
			a.Description = appDescription.(string)
			err = UpdateAppNameDatabase(AgentID, appid.(uint32), a)

		} else if has21 && has22 {
			err = UpdateSubAppNameDatabase(AgentID, appid.(uint32), subAppTag.(uint32), subAppName.(string))
		} else {
			//data record contains interface field
			r, err := FetchAppAttributeDatabaseByID(AgentID, appid.(uint32))
			if err == nil {
				d["applicationAttribute"] = r
			}
			r1, err1 := FetchAppNameDatabaseByID(AgentID, appid.(uint32))
			if err1 == nil {
				d["applicationName"] = r1
			}

		}
	}
	return err
}

func UpdateCiscoVarString(d map[string]interface{}, AgentID string) error {
	var err error
	ciscovarstr, valid := d["applicationVarString"]
	if valid {
		valueList := ciscovarstr.(map[iedb.CiscoAppVarStringKey]string)
		for key, value := range valueList {
			if len(value) > 0 {
				fieldname, err := FetchSubAppNameDatabaseByID(AgentID, key.AppID, key.SubAppTag)
				if err == nil {
					fieldname = strings.Replace(fieldname, " ", "_", -1)
					d[fieldname] = value
				}

			}

		}
		//delete varstring from map
		delete(d, "applicationVarString")

	}
	return err
}
