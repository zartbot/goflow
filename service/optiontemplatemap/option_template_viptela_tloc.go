package optiontemplatemap

import (
	"errors"
	"fmt"
	"net"
	"sync"
)

/*
  _____________________________________________________________________________
  |                 Field                   |    ID | Ent.ID | Offset |  Size |
  -----------------------------------------------------------------------------
  | TLOC TABLE OVERLAY SESSION ID           | 12435 |      9 |      0 |     4 |
  | tloc local color                        | 12437 |      9 |      4 |    16 |
  | tloc remote color                       | 12439 |      9 |     20 |    16 |
  | tloc tunnel protocol                    | 12440 |      9 |     36 |     8 |
  | tloc local system ip address            | 12436 |      9 |     44 |     4 |
  | tloc remote system ip address           | 12438 |      9 |     48 |     4 |
  | sdwan bfd-loss                          | 12512 |      9 |     52 |     4 |
  | sdwan bfd-jitter                        | 12513 |      9 |     56 |     4 |
  | sdwan bfd-latency                       | 12514 |      9 |     60 |     4 |
  -----------------------------------------------------------------------------


12432,9,overlay_session_id_input,unsigned32
12433,9,overlay_session_id_output,unsigned32
12434,9,routing_vrf_service,unsigned32
12435,9,tloc_table_overlay_session_id,unsigned32
12436,9,tloc_local_system_ip_address,ipv4Address
12437,9,tloc_local_color,string
12438,9,tloc_remote_system_ip_address,ipv4Address
12439,9,tloc_remote_color,string
12440,9,tloc_tunnel_protocol,string

12512,9,tloc_bfd_loss,unsigned32
12513,9,tloc_bfd_jitter,unsigned32
12514,9,tloc_bfd_latency,unsigned32
*/

type TLOCData struct {
	Protocol    string
	LocalIP     net.IP
	LocalColor  string
	RemoteIP    net.IP
	RemoteColor string
	Loss        uint32
	Jitter      uint32
	Latency     uint32
}

func (t *TLOCData) String() string {
	return fmt.Sprintf("[%s]-Local:%s:%s--Remote:%s:%s", t.Protocol, t.LocalIP, t.LocalColor, t.RemoteIP, t.RemoteColor)
}

type TLOCMapKey struct {
	AgentID   string
	SessionID uint32
}

var ViptelaTLocTable sync.Map

func UpdateTLOCDatabase(AgentID string, sessionID uint32, tloc *TLOCData) error {
	key := TLOCMapKey{
		AgentID:   AgentID,
		SessionID: sessionID,
	}
	data := TLOCData{
		Protocol:    tloc.Protocol,
		LocalIP:     tloc.LocalIP,
		LocalColor:  tloc.LocalColor,
		RemoteIP:    tloc.RemoteIP,
		RemoteColor: tloc.RemoteColor,
		Loss:        tloc.Loss,
		Jitter:      tloc.Jitter,
		Latency:     tloc.Latency,
	}
	ViptelaTLocTable.Store(key, data)
	return nil
}

func FetchTLOCDatabaseByIndex(AgentID string, sessionID uint32) (TLOCData, error) {
	var err error
	var i TLOCData
	key := TLOCMapKey{
		AgentID:   AgentID,
		SessionID: sessionID,
	}

	value, ok := ViptelaTLocTable.Load(key)
	if ok {
		i, valid := value.(TLOCData)
		if !valid {
			err = errors.New("invalid type assertion")
			return i, err
		} else {
			return i, nil
		}
	} else {
		err = errors.New("TLOC-Session-ID does not found in database")
		return i, err
	}
}

func ShowTLOCDatabase() {
	ViptelaTLocTable.Range(func(k, v interface{}) bool {
		value := v.(TLOCData)
		key := k.(TLOCMapKey)
		fmt.Printf("TLOC Session Map [%10s] | %d  %s \n", key.AgentID, key.SessionID, value.String())
		return true
	})
}

func UpdateViptelaTLOCMap(d map[string]interface{}, AgentID string) error {
	var err error
	sessionid, has_sessionid := d["tloc_table_overlay_session_id"]
	local_ip, has_local_ip := d["tloc_local_system_ip_address"]
	local_color, has_local_color := d["tloc_local_color"]

	if has_sessionid && has_local_ip && has_local_color {
		//is Option Template, update db
		remote_ip, _ := d["tloc_remote_system_ip_address"]
		remote_color, _ := d["tloc_remote_color"]
		protocol_id, _ := d["tloc_tunnel_protocol"]

		tlocdata := &TLOCData{
			RemoteIP:    remote_ip.(net.IP),
			RemoteColor: remote_color.(string),
			LocalIP:     local_ip.(net.IP),
			LocalColor:  local_color.(string),
			Protocol:    protocol_id.(string),
			Latency:     65535,
			Jitter:      65535,
			Loss:        100,
		}

		latencyT, exist := d["tloc_bfd_latency"]
		if exist {
			tlocdata.Latency = latencyT.(uint32)
		}
		jitterT, exist := d["tloc_bfd_jitter"]
		if exist {
			tlocdata.Jitter = jitterT.(uint32)
		}
		lossT, exist := d["tloc_bfd_loss"]
		if exist {
			tlocdata.Loss = lossT.(uint32)
		}

		err = UpdateTLOCDatabase(AgentID, sessionid.(uint32), tlocdata)
	}

	inputID, hasInputID := d["overlay_session_id_input"]
	if hasInputID {
		r, err := FetchTLOCDatabaseByIndex(AgentID, inputID.(uint32)+10000)
		if err == nil {
			d["input_TLOC_Session"] = r
		}
	}

	outputID, hasOutputID := d["overlay_session_id_output"]
	if hasOutputID {
		r, err := FetchTLOCDatabaseByIndex(AgentID, outputID.(uint32)+10000)
		if err == nil {
			d["output_TLOC_Session"] = r
		}
	}
	return err
}
