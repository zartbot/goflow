package ciscoavc

import (
	"math"

	"github.com/zartbot/go_utils/map2value"
)

func ExtendField(d map[string]interface{}) error {
	dataType, _ := map2value.MapToString(d, "Type")

	if dataType == "art" {
		ts, rs := transmitStats(d)
		d["conn_client_to_server_pkts"] = ts.Packets.ClientToServer
		d["conn_server_to_client_pkts"] = ts.Packets.ServerToClient
		d["total_pkts"] = ts.Packets.Total

		d["conn_client_to_server_bytes"] = ts.Bytes.ClientToServer
		d["conn_server_to_client_bytes"] = ts.Bytes.ServerToClient
		d["total_bytes"] = ts.Bytes.Total

		if rs.Bytes.IsValid {
			d["conn_server_bytes_retransmit_rate"] = rs.Bytes.ServerToClient
			d["conn_client_bytes_retransmit_rate"] = rs.Bytes.ClientToServer
			d["total_bytes_retransmit_rate"] = rs.Bytes.Total
			d["total_bytes_loss_sum"] = rs.Bytes.TotalCount
		}
		if rs.Packets.IsValid {
			d["conn_server_pkts_retransmit_rate"] = rs.Packets.ServerToClient
			d["conn_client_pkts_retransmit_rate"] = rs.Packets.ClientToServer
			d["total_pkts_loss_rate"] = rs.Packets.Total
			d["total_pkts_loss_sum"] = rs.Packets.TotalCount
		}

		ProtocolID, _ := map2value.MapToUInt8(d, "protocolIdentifier")

		//only process tcp packet
		if ProtocolID == 6 {
			_, AppDExist := d["conn_delay_app_sum"]
			_, ArtExist := d["conn_delay_network_to_client_sum"]
			_, LLArtExist := d["conn_delay_network_long_lived_to_client"]

			if AppDExist && (ArtExist || LLArtExist) {
				network := new(NetWorkDelay)
				duration := new(Duration)

				r := meanAppDelay(d)
				d["conn_delay_app_mean"] = r.App
				d["conn_delay_client_to_server_mean"] = r.Client2Server
				d["conn_delay_to_server_mean"] = r.ToServer
				if ArtExist {
					network, duration = meanNetDelay(d)
					d["conn_duration_mean"] = duration.ConnMean
					d["transaction_duration_mean"] = duration.TransMean
					d["conn_delay_network_to_client_mean"] = network.Client
					d["conn_delay_network_to_server_mean"] = network.Server
					d["conn_delay_network_mean"] = network.Total
				} else if LLArtExist {
					network, duration = meanLLNetDelay(d)
					d["conn_duration_mean"] = duration.ConnMean
					d["transaction_duration_mean"] = duration.TransMean
					d["conn_delay_network_to_client_mean"] = network.Client
					d["conn_delay_network_to_server_mean"] = network.Server
					d["conn_delay_network_mean"] = network.Total
				}

				rtt := math.Max(float64(network.Total+r.App), float64(r.Client2Server)) * 2
				d["round_trip_time"] = float32(rtt)
				d["tcp_windowsize_mean"] = float32(TCPWindowSize(d))

			}
		}
		FixConnRecordDirection(d)
	}
	RTPRecordExtend(d)
	return nil
}

/*
old performance calc

				//TODO: Get Reference loss/RTT by Application or Application Category
				//Here is a hardcode :
				refRTT := float64(300)
				refLoss := float64(0.005)
				tcpbm := TCPPerformance(d, rtt, float64(rs.Packets.Total), refRTT, refLoss)

				d["score_latency"] = tcpbm.LatencyScore
				d["score_bw"] = tcpbm.BandwidthScore
				d["tcp_theoretical_bw"] = tcpbm.TheoreticalBW
				d["tcp_reference_bw"] = tcpbm.RefBW
				d["performance_score2"] = math.Min(tcpbm.LatencyScore, tcpbm.BandwidthScore)
				//logrus.Warn("Theoretical_BW", tcpbm.TheoreticalBW, "REF:", tcpbm.RefBW, "SCORE:", tcpbm.BandwidthScore)
*/
