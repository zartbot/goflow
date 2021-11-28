package identity

import (
	"fmt"
	"time"

	"github.com/zartbot/go_utils/map2value"

	"github.com/zartbot/goflow/lib/iputil"
)

//Type: 64 UserInfo
//Type: 8  ServerInfo
//Type :1  DNS Record

func LookupAddr(ipaddr string, ptrLookup bool) (map[string]interface{}, int) {
	if ua, utype := LookupUserAddr(ipaddr); utype == 64 {
		//fmt.Println(ua, ustr, utype)
		return ua, utype
	} else if sa, stype := LookupServerAddr(ipaddr); stype == 8 {
		return sa, stype
	} else {
		da, dtype := LookupDNSAddr(ipaddr, ptrLookup)
		return da, dtype
	}
}

func SecurityLookUp(d map[string]interface{}) (int, string, string, int, int) {
	var flow interface{}
	var hasflow bool
	var nameA string
	var nameB string
	var riskscore, riskscoreA, riskscoreB int
	flow, hasflow = d["flowinfo"]
	if hasflow {
		f := flow.(iputil.FlowInfo)
		repA, repCodeA := LookupIPReputation(f.IPAddress_A.String())
		repB, repCodeB := LookupIPReputation(f.IPAddress_B.String())

		var botA, botB string
		var botCodeA, botCodeB int

		/*
			botA, botCodeA := LookupBotNetAddr(f.IPAddress_A.String())
			botB, botCodeB := LookupBotNetAddr(f.IPAddress_B.String())
		*/
		riskscoreA = repCodeA + botCodeA
		if riskscoreA > 0 {
			nameA = fmt.Sprintf("%s:%s", botA, repA)
		}
		riskscoreB = repCodeB + botCodeB
		if riskscoreB > 0 {
			nameB = fmt.Sprintf("%s:%s", botB, repB)
		}
		riskscore = riskscoreA + riskscoreB
	}
	return riskscore, nameA, nameB, riskscoreA, riskscoreB
}

func AddressLookUP(d map[string]interface{}, ptrLookup bool) error {
	var err error
	var flow interface{}
	var hasflow bool

	if d["Type"] == "nat" {
		flowrisk, A, B, _, _ := SecurityLookUp(d)
		d["flowrisk"] = flowrisk
		if flowrisk > 0 {
			d["flowname"] = fmt.Sprintf("%s<-->%s", A, B)
		}
		return err
	}

	//update record for ssl-common-name
	if sslcn, ok := d["ssl_common-name"]; ok && sslcn != nil {
		if serverIP, serverValid := iputil.FetchIPAddress(d, "server"); serverValid {
			Service["DNS"].Store(serverIP.String(), sslcn.(string), time.Now())
		}
	}

	flow, hasflow = d["flowinfo"]
	if hasflow {
		f := flow.(iputil.FlowInfo)
		nameA, typeA := LookupAddr(f.IPAddress_A.String(), ptrLookup)
		nameB, typeB := LookupAddr(f.IPAddress_B.String(), ptrLookup)
		d["hosta"] = nameA
		d["hostb"] = nameB

		flowrisk, secA, secB, riska, riskb := SecurityLookUp(d)
		d["flowrisk"] = flowrisk
		nameA["risk"] = riska
		nameB["risk"] = riskb

		HostNameA, _ := map2value.MapToString(nameA, "HostName")
		HostNameB, _ := map2value.MapToString(nameB, "HostName")
		if riska > 0 {
			HostNameA = fmt.Sprintf("[%s]%s", secA, HostNameA)
		}
		if riskb > 0 {
			HostNameB = fmt.Sprintf("[%s]%s", secB, HostNameB)
		}

		d["flowname"] = fmt.Sprintf("%s<-->%s", HostNameA, HostNameB)
		d["flowtype"] = typeA + typeB

	}
	return err
}
