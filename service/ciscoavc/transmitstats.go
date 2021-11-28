package ciscoavc

import (
	"github.com/zartbot/go_utils/map2value"
)

type TransmitStats struct {
	Bytes struct {
		ClientToServer uint64
		ServerToClient uint64
		Total          uint64
	}
	Packets struct {
		ClientToServer uint64
		ServerToClient uint64
		Total          uint64
	}
}
type RetransmitRate struct {
	Bytes struct {
		IsValid        bool
		TotalCount     uint64
		ClientToServer float32
		ServerToClient float32
		Total          float32
	}
	Packets struct {
		IsValid        bool
		TotalCount     uint64
		ClientToServer float32
		ServerToClient float32
		Total          float32
	}
}

func transmitStats(d map[string]interface{}) (*TransmitStats, *RetransmitRate) {
	ts := new(TransmitStats)
	rs := new(RetransmitRate)

	ClientPkts, ErrClientPkts := map2value.MapToUInt64(d, "initiatorPackets")
	ServerPkts, ErrServerPkts := map2value.MapToUInt64(d, "responderPackets")
	ClientBytes, ErrClientBytes := map2value.MapToUInt64(d, "initiatorOctets")
	ServerBytes, ErrServerBytes := map2value.MapToUInt64(d, "responderOctets")
	if ErrClientPkts == nil {
		ts.Packets.ClientToServer = ClientPkts
	}
	if ErrServerPkts == nil {
		ts.Packets.ServerToClient = ServerPkts
	}
	if ErrClientBytes == nil {
		ts.Bytes.ClientToServer = ClientBytes
	}
	if ErrServerBytes == nil {
		ts.Bytes.ServerToClient = ServerBytes
	}

	//if longlive mode enabled, it will has new field
	LLClientBytes, ErrLLClientBytes := map2value.MapToUInt64(d, "conn_client_cnt_bytes_network_long")
	LLServerBytes, ErrLLServerBytes := map2value.MapToUInt64(d, "conn_server_cnt_bytes_network_long")
	if ErrLLClientBytes == nil {
		ts.Bytes.ClientToServer = LLClientBytes
	}
	if ErrLLServerBytes == nil {
		ts.Bytes.ServerToClient = LLServerBytes
	}
	//fix layer2 bytes
	ts.Bytes.ClientToServer = ts.Bytes.ClientToServer + ts.Packets.ClientToServer*14
	ts.Bytes.ServerToClient = ts.Bytes.ServerToClient + ts.Packets.ServerToClient*14

	ts.Bytes.Total = ts.Bytes.ClientToServer + ts.Bytes.ServerToClient
	ts.Packets.Total = ts.Packets.ClientToServer + ts.Packets.ServerToClient

	CRetransCnt, ErrCRetransCnt := map2value.MapToUInt32(d, "conn_client_cnt_pkts_retransmit")
	SRetransCnt, ErrSRetransCnt := map2value.MapToUInt32(d, "conn_server_cnt_pkts_retransmit")
	if (ErrCRetransCnt == nil) && (ts.Packets.ClientToServer > 0) {
		rs.Packets.ClientToServer = float32(CRetransCnt) / float32(ts.Packets.ClientToServer)
	}
	if (ErrSRetransCnt == nil) && (ts.Packets.ServerToClient > 0) {
		rs.Packets.ServerToClient = float32(SRetransCnt) / float32(ts.Packets.ServerToClient)
	}
	if (ErrCRetransCnt == nil) && (ErrSRetransCnt == nil) && (ts.Packets.Total > 0) {
		rs.Packets.Total = float32(CRetransCnt+SRetransCnt) / float32(ts.Packets.Total)
		rs.Packets.TotalCount = uint64(CRetransCnt) + uint64(SRetransCnt)
		rs.Packets.IsValid = true
	}

	CRetransBytes, ErrCRetransBytes := map2value.MapToUInt32(d, "conn_client_cnt_bytes_retransmit")
	SRetransBytes, ErrSRetransBytes := map2value.MapToUInt32(d, "conn_server_cnt_bytes_retransmit")

	if (ErrCRetransBytes == nil) && (ts.Bytes.ClientToServer > 0) {
		rs.Bytes.ClientToServer = float32(CRetransBytes) / float32(ts.Bytes.ClientToServer)
	}
	if (ErrSRetransBytes == nil) && (ts.Bytes.ServerToClient > 0) {
		rs.Bytes.ServerToClient = float32(SRetransBytes) / float32(ts.Bytes.ServerToClient)
	}
	if (ErrCRetransCnt == nil) && (ErrSRetransCnt == nil) && (ts.Bytes.Total > 0) {
		rs.Bytes.Total = float32(CRetransBytes+SRetransBytes) / float32(ts.Bytes.Total)
		rs.Bytes.TotalCount = uint64(CRetransBytes + SRetransBytes)
		rs.Bytes.IsValid = true
	}

	return ts, rs
}
