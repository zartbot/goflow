package output

import (
	"time"

	"github.com/zartbot/goflow/datarecord"
	"github.com/zartbot/goflow/service/ciscoavc"
	"github.com/zartbot/goflow/service/flowinfo"
	"github.com/zartbot/goflow/service/geoipmap"
	"github.com/zartbot/goflow/service/optiontemplatemap"
)

func RecordMap(d *datarecord.DataFrame) {
	d.TypeAssertion()

	createAt := time.Unix(int64(d.ExportTime), 0)
	for _, value := range d.Record {
		value["CreateAt"] = createAt
		value["AgentID"] = d.AgentID
		value["Type"] = d.Type
		optiontemplatemap.UpdateInterfaceMap(value, d.AgentID)
		optiontemplatemap.UpdateAppMap(value, d.AgentID)
		optiontemplatemap.UpdateCiscoVarString(value, d.AgentID)
		flowinfo.UpdateFlowInfo(value, true)
		ciscoavc.ExtendField(value)

	}
}

//RecordStreamProccesor : stream process worker....
func RecordStreamProccesor(dfchan chan *datarecord.DataFrame, fanoutList []chan *datarecord.DataFrame) {
	for {
		d := <-dfchan
		RecordMap(d)
		if (d.Type != "NULL") && (d.Type != "OptionTemplate") {
			for _, item := range fanoutList {
				item <- d
			}
		}
	}
}

func RecordMapWithPrediction(d *datarecord.DataFrame, g *geoipmap.GeoIPCollector) {
	d.TypeAssertion()

	createAt := time.Unix(int64(d.ExportTime), 0)
	for _, value := range d.Record {
		value["CreateAt"] = createAt
		value["AgentID"] = d.AgentID
		value["Type"] = d.Type
		optiontemplatemap.UpdateInterfaceMap(value, d.AgentID)
		optiontemplatemap.UpdateAppMap(value, d.AgentID)
		optiontemplatemap.UpdateCiscoVarString(value, d.AgentID)
		optiontemplatemap.UpdateC3PLMap(value, d.AgentID)
		flowinfo.UpdateFlowInfo(value, true)
		geoipmap.UpdateGeoLocationInfo(value, g)
		ciscoavc.ExtendField(value)
	}
}
