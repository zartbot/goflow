package ciscoavc

import "github.com/zartbot/go_utils/map2value"

func Swap(d map[string]interface{}, a string, b string) {
	d[a], d[b] = d[b], d[a]
}

func FixConnRecordDirection(d map[string]interface{}) error {
	//Flow direction detect based on Observation Point view
	//
	//     Ref: RFC:5103 Bidirectional Flow Export Using IP Flow Information Export (IPFIX)
	//     +-------+------------------+----------------------------------------+
	//     | Value | Name             | Description                            |
	//     +-------+------------------+----------------------------------------+
	//     | 0x00  | arbitrary        | Direction was assigned arbitrarily.    |
	//     | 0x01  | initiator        | The Biflow Source is the flow          |
	//     |       |                  | initiator, as determined by the        |
	//     |       |                  | Metering Process' best effort to       |
	//     |       |                  | detect the initiator.                  |
	//     | 0x02  | reverseInitiator | The Biflow Destination is the flow     |
	//     |       |                  | initiator, as determined by the        |
	//     |       |                  | Metering Process' best effort to       |
	//     |       |                  | detect the initiator.  This value is   |
	//     |       |                  | provided for the convenience of        |
	//     |       |                  | Exporting Processes to revise an       |
	//     |       |                  | initiator estimate without re-encoding |
	//     |       |                  | the Biflow Record.                     |
	//     | 0x03  | perimeter        | The Biflow Source is the endpoint      |
	//     |       |                  | outside of a defined perimeter.  The   |
	//     |       |                  | perimeter's definition is implicit in  |
	//     |       |                  | the set of Biflow Source and Biflow    |
	//     |       |                  | Destination addresses exported in the  |
	//     |       |                  | Biflow Records.                        |
	//     +-------+------------------+----------------------------------------+
	bidirectionT, bDirectionValid := d["biflowDirection"]
	flowdirectionT, fDirectionValid := d["flowDirection"]
	if bDirectionValid && fDirectionValid && bidirectionT != nil && flowdirectionT != nil {
		//only process tcp packet with latency field
		ProtocolID, _ := map2value.MapToUInt8(d, "protocolIdentifier")
		//_, NetworkDelayExist := d["conn_delay_network_to_server_mean"]
		dataType, _ := map2value.MapToString(d, "Type")

		PktsRetransExist := false
		BytesRetransExist := false
		server_retrans_pkt_cnt, err_S_R_P_CNT := map2value.MapToUInt32(d, "conn_server_cnt_pkts_retransmit")
		client_retrans_pkt_cnt, err_C_R_P_CNT := map2value.MapToUInt32(d, "conn_client_cnt_pkts_retransmit")
		if (err_S_R_P_CNT == nil) && (err_C_R_P_CNT == nil) {
			PktsRetransExist = true
		}

		server_retrans_bytes_cnt, err_S_R_B_CNT := map2value.MapToUInt32(d, "conn_server_cnt_bytes_retransmit")
		client_retrans_bytes_cnt, err_C_R_B_CNT := map2value.MapToUInt32(d, "conn_client_cnt_bytes_retransmit")
		if (err_S_R_B_CNT == nil) && (err_C_R_B_CNT == nil) {
			BytesRetransExist = true
		}

		if flowdirectionT.(uint8) == 0 {
			//FlowDirection is INPUT

			/*Client(Initiator)-----OUT-----^observationPoint^|==Device==|---IN------Server(Responder) */
			if (ProtocolID == 6) && dataType == "art" {
				d["in_network_delay_mean"] = d["conn_delay_network_to_server_mean"]
				d["out_network_delay_mean"] = d["conn_delay_network_to_client_mean"]
			}

			if PktsRetransExist {
				d["in2out_pkts_loss_sum"] = uint64(server_retrans_pkt_cnt)
				d["out2in_pkts_loss_sum"] = uint64(client_retrans_pkt_cnt)
				d["in2out_pkts_loss_rate"] = d["conn_server_pkts_retransmit_rate"]
				d["out2in_pkts_loss_rate"] = d["conn_client_pkts_retransmit_rate"]
			}

			if BytesRetransExist {
				d["in2out_bytes_loss_sum"] = uint64(server_retrans_bytes_cnt)
				d["out2in_bytes_loss_sum"] = uint64(client_retrans_bytes_cnt)
				d["in2out_bytes_loss_rate"] = d["conn_server_bytes_retransmit_rate"]
				d["out2in_bytes_loss_rate"] = d["conn_client_bytes_retransmit_rate"]
			}

			d["in2out_bytes"] = d["conn_server_to_client_bytes"]
			d["out2in_bytes"] = d["conn_client_to_server_bytes"]
			d["in2out_pkts"] = d["responderPackets"]
			d["out2in_pkts"] = d["initiatorPackets"]
		} else {
			//FlowDirection is OUTPUT
			/*Server(Responder)-----OUT-----^observationPoint^|==Device==|---IN------Client(Initiator) */
			if (ProtocolID == 6) && dataType == "art" {
				d["out_network_delay_mean"] = d["conn_delay_network_to_server_mean"]
				d["in_network_delay_mean"] = d["conn_delay_network_to_client_mean"]
			}
			if PktsRetransExist {
				d["out2in_pkts_loss_sum"] = uint64(server_retrans_pkt_cnt)
				d["in2out_pkts_loss_sum"] = uint64(client_retrans_pkt_cnt)
				d["out2in_pkts_loss_rate"] = d["conn_server_pkts_retransmit_rate"]
				d["in2out_pkts_loss_rate"] = d["conn_client_pkts_retransmit_rate"]
			}
			if BytesRetransExist {
				d["out2in_bytes_loss_sum"] = uint64(server_retrans_bytes_cnt)
				d["in2out_bytes_loss_sum"] = uint64(client_retrans_bytes_cnt)
				d["out2in_bytes_loss_rate"] = d["conn_server_bytes_retransmit_rate"]
				d["in2out_bytes_loss_rate"] = d["conn_client_bytes_retransmit_rate"]
			}
			d["out2in_bytes"] = d["conn_server_to_client_bytes"]
			d["in2out_bytes"] = d["conn_client_to_server_bytes"]
			d["out2in_pkts"] = d["responderPackets"]
			d["in2out_pkts"] = d["initiatorPackets"]
		}
		// this is a very tricky detection since Router currently only report 0,1,2
		// all biflowDirection == 0 flow has src=0.0.0.0 , dst=0.0.0.0 ...
		if bidirectionT.(uint8) == 2 {
			if (ProtocolID == 6) && dataType == "art" {
				Swap(d, "in_network_delay_mean", "out_network_delay_mean")
			}
			if PktsRetransExist {
				Swap(d, "out2in_pkts_loss_sum", "in2out_pkts_loss_sum")
				Swap(d, "out2in_pkts_loss_rate", "in2out_pkts_loss_rate")
			}
			if BytesRetransExist {
				Swap(d, "out2in_bytes_loss_sum", "in2out_bytes_loss_sum")
				Swap(d, "out2in_bytes_loss_rate", "in2out_bytes_loss_rate")
			}
			Swap(d, "out2in_bytes", "in2out_bytes")
			Swap(d, "out2in_pkts", "in2out_pkts")
		}
	}
	return nil
}

func FixRTPRecordDirection(d map[string]interface{}) error {
	bidirectionT, bDirectionValid := d["biflowDirection"]
	flowdirectionT, fDirectionValid := d["flowDirection"]
	if bDirectionValid && fDirectionValid && bidirectionT != nil && flowdirectionT != nil {
		if flowdirectionT.(uint8) == 0 {
			//FlowDirection is INPUT
			d["out2in_pkts_loss_rate"] = d["rtp_loss_rate"]
			d["out2in_rtp_jitter_mean"] = d["rtp_jitter_mean"]
			d["out2in_rtp_mos"] = d["rtp_mos"]
			d["out2in_performance_score"] = d["performance_score"]
		} else {
			//FLowDirection is OUTPUT
			d["in2out_pkts_loss_rate"] = d["rtp_loss_rate"]
			d["in2out_rtp_jitter_mean"] = d["rtp_jitter_mean"]
			d["in2out_rtp_mos"] = d["rtp_mos"]
			d["in2out_performance_score"] = d["performance_score"]
		}
	}
	return nil
}
