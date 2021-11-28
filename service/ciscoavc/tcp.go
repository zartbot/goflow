package ciscoavc

import (
	"math"
	"strings"

	"github.com/zartbot/go_utils/map2value"
)

func TCPWindowSize(d map[string]interface{}) float64 {
	var avgWinSize float64
	winSizeSum, ErrWinSize := map2value.MapToUInt64(d, "trans_tcp_window_size_sum")
	pkts, ErrPkts := map2value.MapToUInt64(d, "total_pkts")
	if (ErrWinSize == nil) && (ErrPkts == nil) {
		avgWinSize = float64(winSizeSum) / float64(pkts)
	}
	return avgWinSize
}

//TCPThrouhgput is a function to calculate theoretical network limit
//based on the following formula
// if has pkt loss  rate <= (MSS/RTT)*(C/sqrt(Loss))
// if does not have pkt loss rate <= TCPBufferSize /RTT
//MSS: Bytes uint64
//RTT: ms float64
//Loss: % float64
//WindowSize:  Bytes float64
func TCPTheoreticalThrouhgput(rtt float64, loss float64, winsize float64, mss uint64) float64 {
	if loss == 0 {
		if winsize == 0 {
			return math.Min(16000/rtt, float64(mss)/(rtt*math.Sqrt(0.0001))) * 8000
		}
		return winsize / rtt * 8000
	}
	return float64(mss) / (rtt * math.Sqrt(loss)) * 8000
}

type TCPBenchmark struct {
	AvgWindowSize  float64
	TheoreticalBW  float64
	LatencyScore   float64
	BandwidthScore float64
	RefBW          float64
}

func InterfaceBW(d map[string]interface{}) float64 {
	intf, ErrIntf := map2value.MapToString(d, "observationInterface")
	//default bw is 1G
	var intfbw float64 = 1000000000
	if ErrIntf == nil {
		if strings.Contains(intf, "Tunnel") {
			intfbw = 1000000000
		} else if strings.Contains(intf, "TenGigabitEthernet") {
			intfbw = 10000000000
		} else if strings.Contains(intf, "FastEthernet") {
			intfbw = 100000000
		}
	}
	return intfbw
}

func TCPPerformance(d map[string]interface{}, rtt float64, loss float64, brtt float64, bloss float64) *TCPBenchmark {
	const MSS uint64 = 1460
	const DefaultWindowSize float64 = 64000
	t := new(TCPBenchmark)
	t.AvgWindowSize = TCPWindowSize(d)
	if rtt < 1 {
		rtt = 1
	}
	if loss > 1 {
		loss = 1
	}
	//intfbw := InterfaceBW(d)
	t.TheoreticalBW = TCPTheoreticalThrouhgput(rtt, loss, t.AvgWindowSize, MSS)
	t.RefBW = TCPTheoreticalThrouhgput(brtt, bloss, DefaultWindowSize, MSS)
	t.BandwidthScore = math.Min(20*math.Log(t.TheoreticalBW/t.RefBW)+100, 100)
	if t.BandwidthScore < 0 {
		t.BandwidthScore = 0
	}
	t.LatencyScore = math.Min(100*math.Log(2*brtt/rtt)+32, 100)
	if t.LatencyScore < 0 {
		t.LatencyScore = 0
	}
	//t.LatencyScore = 1 / (1 + math.Exp(-rtt/brtt))
	return t
}
