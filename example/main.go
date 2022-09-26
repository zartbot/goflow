package main

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/zartbot/goflow/service/ciscoavc"
	"github.com/zartbot/goflow/service/ciscoeta"
	"github.com/zartbot/goflow/service/flowinfo"
	"github.com/zartbot/goflow/service/geoipmap"
	"github.com/zartbot/goflow/service/identity"
	"github.com/zartbot/goflow/service/optiontemplatemap"

	"github.com/zartbot/goflow/datarecord"
	"github.com/zartbot/goflow/iedb"
	"github.com/zartbot/goflow/ipfix"
	"github.com/zartbot/goflow/netflowv9"
)

func main() {

	ipaddr := "0.0.0.0"

	iedb.ReadIANA("./config/database/informationElement/iana_ie.csv")
	iedb.ReadCiscoIE("./config/database/informationElement/cisco_ie.csv")
	iedb.ReadCiscoIE("./config/database/informationElement/key_ie.csv")
	iedb.ShowIEDB()

	logrus.Info("Starting Identification service....")
	identity.LoadServerInfo("./config/database/server/server.csv")
	identity.LoadReputationInfo("./config/database/security/reputation.csv")
	identity.LoadBotnetInfo("./config/database/security/spam_bot_net.csv")
	go identity.Run()

	//g, geoServiceErr := geoipmap.NewGeoIPCollector("./config/database/geolocation/city.mmdb", "./config/database/geolocation/asn.mmdb", 31.123, 111.11)

	logrus.Info("Starting Netflow Collectors.....")

	//IPFIXDataFrameChan is dataframe output channel for future streaming process.
	var IPFIXDataFrameChan = make(chan *datarecord.DataFrame, 10000)
	var NFv9DataFrameChan = make(chan *datarecord.DataFrame, 10000)

	//input Collector
	var ipfixCollector []*ipfix.IPFIXCollector
	for _, port := range []int{4739, 4740} {
		c := ipfix.NewIPFIXCollector(ipaddr, port, IPFIXDataFrameChan, 1)
		ipfixCollector = append(ipfixCollector, c)
	}

	var nfv9Collector []*netflowv9.NetflowV9Collector
	for _, port := range []int{2055, 2056} {
		c := netflowv9.NewNFV9Collector(ipaddr, port, NFv9DataFrameChan, 1)
		nfv9Collector = append(nfv9Collector, c)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	for _, c := range ipfixCollector {
		go c.Run()
	}

	for _, c := range nfv9Collector {
		go c.Run()
	}

	for {
		select {
		case e1 := <-IPFIXDataFrameChan:
			RecordMap(e1)
			e1.Print("IPFIX")

		case e2 := <-NFv9DataFrameChan:
			RecordMap(e2)
			e2.Print("NFv9")
		}
	}

	wg.Wait()

}

func RecordMap(d *datarecord.DataFrame) {
	d.TypeAssertion()
	createAt := time.Unix(int64(d.ExportTime), 0)
	for _, value := range d.Record {
		value["CreateAt"] = createAt
		value["AgentID"] = d.AgentID
		value["Type"] = d.Type
		//Optional Template Correlation
		optiontemplatemap.UpdateInterfaceMap(value, d.AgentID)
		optiontemplatemap.UpdateAppMap(value, d.AgentID)
		optiontemplatemap.UpdateCiscoVarString(value, d.AgentID)
		optiontemplatemap.UpdateC3PLMap(value, d.AgentID)
		optiontemplatemap.UpdateDropCauseMap(value, d.AgentID)
		optiontemplatemap.UpdateViptelaTLOCMap(value, d.AgentID)
		optiontemplatemap.UpdateFWEventMap(value, d.AgentID)
		optiontemplatemap.UpdateFWZonePairMap(value, d.AgentID)
		optiontemplatemap.UpdateFWClassMapMap(value, d.AgentID)

		//Additional Service
		flowinfo.UpdateFlowInfo(value, true)
		identity.AddressLookUP(value, true)
		ciscoavc.ExtendField(value)
		ciscoeta.ExtendField(value)
	}
}

func GeoRecordMap(d *datarecord.DataFrame, g *geoipmap.GeoIPCollector) {
	d.TypeAssertion()
	createAt := time.Unix(int64(d.ExportTime), 0)
	for _, value := range d.Record {
		value["CreateAt"] = createAt
		value["AgentID"] = d.AgentID
		value["Type"] = d.Type
		//Optional Template Correlation
		optiontemplatemap.UpdateInterfaceMap(value, d.AgentID)
		optiontemplatemap.UpdateAppMap(value, d.AgentID)
		optiontemplatemap.UpdateCiscoVarString(value, d.AgentID)
		optiontemplatemap.UpdateC3PLMap(value, d.AgentID)
		optiontemplatemap.UpdateDropCauseMap(value, d.AgentID)
		optiontemplatemap.UpdateViptelaTLOCMap(value, d.AgentID)
		optiontemplatemap.UpdateFWEventMap(value, d.AgentID)
		optiontemplatemap.UpdateFWZonePairMap(value, d.AgentID)
		optiontemplatemap.UpdateFWClassMapMap(value, d.AgentID)

		//Additional Service
		flowinfo.UpdateFlowInfo(value, true)
		identity.AddressLookUP(value, true)
		geoipmap.UpdateGeoLocationInfo(value, g)
		ciscoavc.ExtendField(value)
		ciscoeta.ExtendField(value)
	}
}
