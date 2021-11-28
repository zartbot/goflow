package ciscoavc

import (
	"github.com/zartbot/go_utils/map2value"
)

/*
MOS	Quality	Impairment
5	Excellent	Imperceptible
4	Good	Perceptible but not annoying
3	Fair	Slightly annoying
2	Poor	Annoying
1	Bad    	>Very annoying

MOS Calculation

  effective_latency = latency +2*jitter + 10.0

For effective_latency < 160.0 ms:
  R = 93.2 - (effective_latency)/40.0

For effective_latency >= 160.0 ms:

  R = 93.2 - (effective_latency - 120.0)/10.0

If the effective latency is less than 160.0 ms the overall impact to the voice quality is moderate. For larger values, the voice quality drops more significantly, which is why R is penalized more.
The packet loss (in percentage points) is taken into consideration as follows:

  R = R - 2.5 * packet_loss

For R < 0:
  MOS = 1.0

For 0 < R < 100.0:
  MOS = 1 + 0.035*R + 0.000007*R*(R-60)*(100-R)

For R >= 100.0:
  MOS = 4.5
*/

type RTPStats struct {
	AvgJitter float64
	Loss      float64
	MOS       float64
	Score     float64
}

func MOS(delay float64, jitter float64, loss float64) (float64, float64) {
	var R, Score, MOS float64
	effectiveLatency := delay + 2*jitter + 10
	if effectiveLatency < 160 {
		R = 93.2 - effectiveLatency/40.0
	} else {
		R = 93.2 - (effectiveLatency-120.0)/10
	}
	R = R - 2.5*loss
	if R < 0 {
		MOS = 1.0
		Score = 0
	} else if R < 100 {
		MOS = 1 + 0.035*R + 0.000007*R*(R-60)*(100-R)
		Score = R
	} else {
		MOS = 4.5
		Score = 100
	}
	return Score * Score * Score / 10000, MOS
}

func RTPRecordExtend(d map[string]interface{}) *RTPStats {
	r := new(RTPStats)
	_, ErrRTPPayloadType := map2value.MapToUInt8(d, "trans_rtp_payload_type")
	if ErrRTPPayloadType == nil {
		AvgRTPJitter, ErrAvgRTPJitter := map2value.MapToUInt32(d, "trans_rtp_jitter_mean")
		if ErrAvgRTPJitter == nil {
			if AvgRTPJitter > 4294967290 {
				AvgRTPJitter = 0
				d["trans_rtp_jitter_mean"] = 0
			}
		}
		MaxRTPJitter, ErrMaxRTPJitter := map2value.MapToUInt32(d, "trans_rtp_jitter_maximum")
		if ErrMaxRTPJitter == nil {
			if MaxRTPJitter > 4294967290 {
				MaxRTPJitter = 0
				d["trans_rtp_jitter_maximum"] = 0
			}
		}

		MinRTPJitter, ErrMinRTPJitter := map2value.MapToUInt32(d, "trans_rtp_jitter_minimum")
		if ErrMinRTPJitter == nil {
			if MinRTPJitter > 4294967290 {
				MinRTPJitter = 0
				d["trans_rtp_jitter_minimum"] = 0
			}
		}

		SumRTPJitter, ErrSumRTPJitter := map2value.MapToUInt64(d, "trans_rtp_jitter_mean_sum")
		if ErrSumRTPJitter == nil {
			if SumRTPJitter > 4294967290 {
				SumRTPJitter = 0
				d["trans_rtp_jitter_mean_sum"] = 0
			}
		}
		var pktLostCnt uint64
		if _v, ok := d["trans_pkts_lost_cnt_long"]; ok {
			if v, valid := _v.(uint32); valid {
				if v > 4294967290 {
					pktLostCnt = uint64(0)
				} else {
					pktLostCnt = uint64(v)
				}
			} else if v, valid := _v.(uint64); valid {
				if v > 4294967290 {
					pktLostCnt = 0
				}
			}
		}
		d["trans_pkts_lost_cnt_long"] = pktLostCnt
		d["total_pkts_loss_sum"] = pktLostCnt
		Totalpkts, ErrTotalPkts := map2value.MapToUInt64(d, "total_pkts")
		if ErrTotalPkts == nil {
			r.Loss = float64(pktLostCnt) / float64(Totalpkts)
			d["rtp_loss_rate"] = float32(r.Loss)
			r.AvgJitter = float64(SumRTPJitter) / float64(Totalpkts)
			d["rtp_jitter_mean"] = r.AvgJitter
			//hardcode for network latency to 20ms
			//since currently ART record does not contain delay
			r.Score, r.MOS = MOS(20, r.AvgJitter, r.Loss)
			d["performance_score"] = r.Score
			d["rtp_mos"] = r.MOS
		}
		FixRTPRecordDirection(d)
	}
	return r
}
