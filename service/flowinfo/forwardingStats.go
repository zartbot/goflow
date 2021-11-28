package flowinfo

import (
	"github.com/zartbot/go_utils/map2value"
)

func NormalizeForwardingStats(d map[string]interface{}) {
	//rename pkt/bytes [long] to total_pkt/bytes
	if pkts, ok := d["packetDeltaCount"]; ok {
		var total_pkts uint64
		if v, valid := pkts.(uint64); valid {
			total_pkts = v
		} else if v, valid := pkts.(uint32); valid {
			total_pkts = uint64(v)
		}
		d["total_pkts"] = total_pkts
		//update taildrop

		if total_pkts != 0 {
			winSize, errWin := map2value.MapToUInt64(d, "trans_tcp_window_size_sum")
			if errWin == nil {
				d["trans_tcp_window_size_mean"] = float64(winSize) / float64(total_pkts)
			}

			queue_depth_taildrop_counter, errtail := map2value.MapToUInt64(d, "queue_depth_taildrop_counter")
			if errtail == nil {
				d["taildrop_rate_mean"] = float64(queue_depth_taildrop_counter) / float64(total_pkts)
			}

			queue_depth_sum, errQdpSum := map2value.MapToUInt64(d, "queue_depth_sum")
			if errQdpSum == nil {
				d["queue_depth_mean"] = float64(queue_depth_sum) / float64(total_pkts)
			}
		}
	}
	if bs, ok := d["octetDeltaCount"]; ok {
		if v, valid := bs.(uint64); valid {
			d["total_bytes"] = v
		} else if v, valid := bs.(uint32); valid {
			d["total_bytes"] = uint64(v)
		}
	}
}
