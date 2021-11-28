package ciscoavc

import (
	"github.com/zartbot/go_utils/map2value"
)

/*
Field Definition Link:
https://www.cisco.com/c/en/us/td/docs/ios/solutions_docs/avc/guide/avc-user-guide/avc_config.html
https://www.cisco.com/c/en/us/td/docs/routers/access/ISRG2/AVC/api/guide/AVC_Metric_Definition_Guide/5_AVC_Metric_Def.html
https://www.cisco.com/c/en/us/td/docs/ios/solutions_docs/avc/ios_xe3_9/avc_soln_guide_iosxe3_9/avc_app_exported_fields.html

conn_delay_app_sum [milliseconds]
	CLI: collect connection delay app sum

conn_delay_response_to_server_sum [milliseconds]
	* CLI: collect connection delay response to-server sum
	* Response time is the time between the client request and the corresponding first response packet
	* from the server, as observed at the observation point. The value of this information element is
	* the sum of all response times observed for the responses of this flow. For the average, this field
	* must be divided by numRespsCountDelta (42060).

conn_delay_response_client_to_server_sum [milliseconds]
	* CLI: collect connection delay response client-to-server sum
	* Total delay is the time between the client request and the first response packet from the server,
	* as seen by the client. This is the sum of all total delays observed for the responses of this flow.
	*  For the average, this field must be divided by numRespsCountDelta (42060)

*/
type ResponseDelay struct {
	App           float32
	Client2Server float32
	ToServer      float32
}

func meanAppDelay(d map[string]interface{}) *ResponseDelay {
	r := &ResponseDelay{
		App:           0,
		Client2Server: 0,
		ToServer:      0,
	}
	SrvRspCnt, ErrSrvRspCnt := map2value.MapToUInt32(d, "conn_server_response_cnt")
	DelayAppSum, ErrDelayAppSum := map2value.MapToUInt32(d, "conn_delay_app_sum")
	DelayRspC2SSum, ErrDelayRspC2SSum := map2value.MapToUInt32(d, "conn_delay_response_client_to_server_sum")
	DelayRsp2SSum, ErrDelayRsp2SSum := map2value.MapToUInt32(d, "conn_delay_response_to_server_sum")
	if ErrSrvRspCnt == nil && SrvRspCnt > 0 {
		var SrvRspCntFloat32 = float32(SrvRspCnt)
		if ErrDelayAppSum == nil {
			r.App = float32(DelayAppSum) / SrvRspCntFloat32
		}

		if ErrDelayRspC2SSum == nil {
			r.Client2Server = float32(DelayRspC2SSum) / SrvRspCntFloat32
		}

		if ErrDelayRsp2SSum == nil {
			r.ToServer = float32(DelayRsp2SSum) / SrvRspCntFloat32
		}
	}
	return r
}

type NetWorkDelay struct {
	Client float32
	Server float32
	Total  float32 //client2server
}

type Duration struct {
	ConnTotal  float32
	ConnMean   float32
	TransTotal float32
	TransMean  float32
}

func meanNetDelay(d map[string]interface{}) (*NetWorkDelay, *Duration) {

	network := &NetWorkDelay{
		Client: 0,
		Server: 0,
		Total:  0,
	}
	duration := &Duration{
		ConnTotal:  0,
		ConnMean:   0,
		TransTotal: 0,
		TransMean:  0,
	}

	//Trans duration mean
	NumTrans, ErrNumTrans := map2value.MapToUInt32(d, "conn_transaction_cnt_complete")
	TransDuraSum, ErrTransDuraSum := map2value.MapToUInt32(d, "conn_transaction_duration_sum")
	if (ErrNumTrans == nil) && (ErrTransDuraSum == nil) {
		if NumTrans > 0 {
			duration.TransMean = float32(TransDuraSum) / float32(NumTrans)
			duration.TransTotal = float32(TransDuraSum)
		}
	}

	NumConn, ErrNumConn := map2value.MapToUInt32(d, "newConnectionDeltaCount")
	DuraSum, ErrDuraSum := map2value.MapToUInt64(d, "connectionSumDurationSeconds")
	CNetDelaySum, ErrCNetDelaySum := map2value.MapToUInt32(d, "conn_delay_network_to_client_sum")
	SNetDelaySum, ErrSNetDelaySum := map2value.MapToUInt32(d, "conn_delay_network_to_server_sum")
	C2SNetDelaySum, ErrC2SNetDelaySum := map2value.MapToUInt32(d, "conn_delay_network_client_to_server_sum")

	if ErrNumConn == nil && NumConn > 0 {
		var NumConnFloat32 = float32(NumConn)
		if ErrDuraSum == nil {
			duration.ConnMean = float32(float64(DuraSum) / float64(NumConn))
			duration.ConnTotal = float32(DuraSum)
		}

		if ErrCNetDelaySum == nil {
			network.Client = float32(float32(CNetDelaySum) / NumConnFloat32)

		}
		if ErrSNetDelaySum == nil {
			network.Server = float32(float32(SNetDelaySum) / NumConnFloat32)
		}

		if ErrC2SNetDelaySum == nil {
			network.Total = float32(float32(C2SNetDelaySum) / NumConnFloat32)
		}
	}

	tn := network.Client + network.Server
	if tn > network.Total {
		network.Total = tn
	}
	return network, duration
}

//Long Live Record Mode
func meanLLNetDelay(d map[string]interface{}) (*NetWorkDelay, *Duration) {
	network := new(NetWorkDelay)
	duration := new(Duration)

	//Conn duration mean
	NumConn, ErrNumConn := map2value.MapToUInt32(d, "newConnectionDeltaCount")
	DuraSum, ErrDuraSum := map2value.MapToUInt64(d, "connectionSumDurationSeconds")
	if ErrNumConn == nil && NumConn > 0 {
		if ErrDuraSum == nil {
			duration.ConnMean = float32(float64(DuraSum) / float64(NumConn))
			duration.ConnTotal = float32(DuraSum)
		}
	}

	//Trans duration mean
	NumTrans, ErrNumTrans := map2value.MapToUInt32(d, "conn_transaction_cnt_complete")
	TransDuraSum, ErrTransDuraSum := map2value.MapToUInt32(d, "conn_transaction_duration_sum")
	if (ErrNumTrans == nil) && (ErrTransDuraSum == nil) {
		if NumTrans > 0 {
			duration.TransMean = float32(TransDuraSum) / float32(NumTrans)
			duration.TransTotal = float32(TransDuraSum)
		}
	}

	//client
	NumC, ErrNumC := map2value.MapToUInt32(d, "conn_delay_network_to_client_num_samples")
	CNetDelaySum, ErrCNetDelaySum := map2value.MapToUInt32(d, "conn_delay_network_long_lived_to_client")
	if (ErrNumC == nil) && (ErrCNetDelaySum == nil) {
		if NumC > 0 {
			network.Client = float32(CNetDelaySum) / float32(NumC)
		} else {
			network.Client = 0
		}
	}

	//server
	NumS, ErrNumS := map2value.MapToUInt32(d, "conn_delay_network_to_server_num_samples")
	SNetDelaySum, ErrSNetDelaySum := map2value.MapToUInt32(d, "conn_delay_network_long_lived_to_server")
	if (ErrNumS == nil) && (ErrSNetDelaySum == nil) {
		if NumS > 0 {
			network.Server = float32(SNetDelaySum) / float32(NumS)
		} else {
			network.Server = 0
		}
	}

	//client-to-server
	NumC2S, ErrNumC2S := map2value.MapToUInt32(d, "conn_delay_network_client_to_server_num_samples")
	C2SNetDelaySum, ErrC2SNetDelaySum := map2value.MapToUInt32(d, "conn_delay_network_long_lived_client_to_server")
	if (ErrNumC2S == nil) && (ErrC2SNetDelaySum == nil) {
		if NumC2S > 0 {
			network.Total = float32(C2SNetDelaySum) / float32(NumC2S)
		}
	}

	//sometime to_client sample counter is zero, need to fix from C2SNetDelay and SNetDelay
	if network.Client == 0 && network.Server != 0 {
		delta := network.Total - network.Server
		if delta > 0 {
			network.Client = delta
		}
	}

	//sometime to_server sample counter is zero, need to fix from C2SNetDelay and CNetDelay
	if network.Client != 0 && network.Server == 0 {
		delta := network.Total - network.Client
		if delta > 0 {
			network.Server = delta
		}
	}
	return network, duration
}
