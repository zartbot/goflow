package flowinfo

import (
	"time"

	"github.com/zartbot/goflow/lib/iputil"
	"github.com/zartbot/goflow/service/identity"
)

func SrcDstLookUp(d map[string]interface{}, EnableNAT bool) (iputil.FlowInfo, bool) {
	var flowrecord iputil.FlowInfo
	//flow mode
	srcIP, srcValid := iputil.FetchIPAddress(d, "src")
	dstIP, dstValid := iputil.FetchIPAddress(d, "dst")
	if srcValid && dstValid {
		srcport := iputil.FetchPort(d, "src")
		dstport := iputil.FetchPort(d, "dst")
		flowrecord = NewFlowInfo(srcIP, srcport, dstIP, dstport)

		//processNAT HSL
		if EnableNAT {
			//validate NAT HSL
			natsrcIP, natSrcValid := iputil.FetchIPAddress(d, "postnatsrc")
			natdstIP, natDstValid := iputil.FetchIPAddress(d, "postnatdst")
			if natSrcValid && natDstValid {
				natsrcport := iputil.FetchPort(d, "postnatsrc")
				natdstport := iputil.FetchPort(d, "postnatdst")
				v := NewFlowInfo(natsrcIP, natsrcport, natdstIP, natdstport)
				d["postNATflowinfo"] = v
				if natEvent, ok := d["natEvent"]; ok {
					//fmt.Println(natEvent)
					currentTime := time.Now()
					if natEvent.(uint8) == 1 {
						//create session in NAT table
						identity.Service["NAT"].Store(v.FlowIDWithPort, flowrecord, currentTime)
					} else if natEvent.(uint8) == 2 {
						//delete session in NAT table , consider delete race condition
						//we just set the expiretime is Time.now - map.Timeout + 180s
						//to make sure the session clean by GC
						identity.Service["NAT"].Store(v.FlowIDWithPort, flowrecord, currentTime.Add(time.Duration(180-identity.Service["NAT"].Timeout)*time.Second))
					}
				}
			} else {
				if origflowinfo, ok := FetchNATAddress(flowrecord); ok {
					flowrecord = origflowinfo
				}
			}
		}
		return flowrecord, true
	}
	return flowrecord, false
}

func BidiLookUp(d map[string]interface{}, EnableNAT bool) (iputil.FlowInfo, bool) {
	var flowrecord iputil.FlowInfo
	//connection mode
	clientIP, clientValid := iputil.FetchIPAddress(d, "client")
	serverIP, serverValid := iputil.FetchIPAddress(d, "server")
	if clientValid && serverValid {
		clientport := iputil.FetchPort(d, "client")
		serverport := iputil.FetchPort(d, "server")
		flowrecord = NewFlowInfo(clientIP, clientport, serverIP, serverport)
		if EnableNAT {
			if origflowinfo, ok := FetchNATAddress(flowrecord); ok {
				flowrecord = origflowinfo
			}
		}
		return flowrecord, true
	}
	return flowrecord, false
}
