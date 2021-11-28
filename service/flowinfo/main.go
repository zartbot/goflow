package flowinfo

import (
	"hash/fnv"
)

func UpdateFlowInfo(d map[string]interface{}, EnableNAT bool) {
	var isFlowRecord, direction bool

	NormalizeForwardingStats(d)
	//flowinfo lookup
	flowrecord, isFlowRecord := SrcDstLookUp(d, EnableNAT)
	if isFlowRecord == false {
		flowrecord, _ = BidiLookUp(d, EnableNAT)
	}

	//flow direction
	bidirectionT, bDirectionValid := d["biflowDirection"]
	flowdirectionT, fDirectionValid := d["flowDirection"]
	if fDirectionValid && flowdirectionT != nil {
		if flowdirectionT.(uint8) == 0 {
			//FlowDirection is INPUT
			/*-----OUT-----^observationPoint^|==Device==|---IN------*/
			/*-----OUT-----(inputInterface)>>>>>>>>>>>>>>---IN------*/
			direction = true
			flowrecord.IPAddress_A, flowrecord.IPAddress_B = flowrecord.IPAddress_B, flowrecord.IPAddress_A
			flowrecord.Port_A, flowrecord.Port_B = flowrecord.Port_B, flowrecord.Port_A
			flowrecord.Direction = !flowrecord.Direction
			if isFlowRecord {
				d["out2in_bytes"] = d["total_bytes"]
				d["out2in_pkts"] = d["total_pkts"]
			}
		} else {
			//FlowDirection is OUTPUT
			/*-----OUT-----^observationPoint^|==Device==|---IN------*/
			/*-----OUT-----(egressIntf)<<<<<<<<<<<<<<<<<<---IN------*/
			if isFlowRecord {
				d["in2out_bytes"] = d["total_bytes"]
				d["in2out_pkts"] = d["total_pkts"]
			}
		}
		if _, ValidObvid := d["observationPointId"]; ValidObvid == false {
			//add Observable Point Name
			if direction {
				//FlowDirection is INPUT
				d["observationInterface"] = d["ingressInterfaceName"]
				d["observationPointId"] = d["ingressInterface"]
				d["observationFlowDirection"] = "input"
			} else {
				d["observationInterface"] = d["egressInterfaceName"]
				d["observationPointId"] = d["egressInterface"]
				d["observationFlowDirection"] = "output"
			}
		}
	}

	if bDirectionValid && bidirectionT != nil {
		// this is a very tricky detection since Router currently only report 0,1,2
		// all biflowDirection == 0 flow has src=0.0.0.0 , dst=0.0.0.0 ...
		if bidirectionT.(uint8) == 2 {
			flowrecord.IPAddress_A, flowrecord.IPAddress_B = flowrecord.IPAddress_B, flowrecord.IPAddress_A
			flowrecord.Port_A, flowrecord.Port_B = flowrecord.Port_B, flowrecord.Port_A
			flowrecord.Direction = !flowrecord.Direction
		}
	}

	if (flowrecord.IPAddress_A != nil) && (flowrecord.IPAddress_B != nil) {
		h32a := fnv.New32a()
		h32b := fnv.New32a()
		h32a.Write(flowrecord.IPAddress_A)
		h32b.Write(flowrecord.IPAddress_B)
		flowrecord.Hash_A = h32a.Sum32()
		flowrecord.Hash_B = h32b.Sum32()
		flowrecord.Hash_A_Mod = h32a.Sum32() % 10000
		flowrecord.Hash_B_Mod = h32b.Sum32() % 10000

		d["flowinfo"] = flowrecord

	}
}
