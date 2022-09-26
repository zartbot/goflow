package datarecord

var RecordTypeList = []string{"nat", "security", "queuedrop", "nimbleflow", "urlstats", "art", "conn", "rtp", "flow", "wireless", "viptela", "viptelabfd", "zbfw"}

func (d *DataFrame) TypeAssertion() {
	if len(d.Record) < 1 {
		d.Type = "NULL"
	} else {
		r := d.Record[0]
		if IsOptionTemplate(r) {
			d.Type = "OptionTemplate"
			return
		} else if IsFW(r) {
			d.Type = "zbfw"
			return
		} else if IsNATRecord(r) {
			d.Type = "nat"
			return
		} else if IsETARecord(r) {
			d.Type = "security"
			return
		} else if isNimbleFlowRecord(r) {
			d.Type = "nimbleflow"
			return
		} else if isUrlStatsRecord(r) {
			d.Type = "urlstats"
			return
		} else if isArtRecord(r) {
			d.Type = "art"
			return
		} else if isConnRecord(r) {
			d.Type = "conn"
		} else if isRTPRecord(r) {
			d.Type = "rtp"
		} else if isViptelaRecord(r) {
			d.Type = "viptela"
			return
		} else if isViptelaBFDRecord(r) {
			d.Type = "viptelabfd"
			return
		} else if isQueueDropRecord(r) {
			d.Type = "queuedrop"
			return
		} else {
			d.Type = "flow"
			return
		}
	}

}

func IsOptionTemplate(d map[string]interface{}) bool {
	// applicationAttributeTable
	_, has1 := d["applicationCategory"]
	_, has2 := d["applicationSubCategory"]

	if has1 && has2 {
		return true
	}

	// applicationNameTable
	_, has11 := d["applicationName"]
	_, has12 := d["applicationDescription"]
	if has11 && has12 {
		return true
	}

	// subApplicationTable
	_, has21 := d["subApplicationTag"]
	_, has22 := d["subApplicationName"]
	if has21 && has22 {
		return true
	}

	// interfaceTemplate
	_, has31 := d["interfaceName"]
	_, has32 := d["interfaceDescription"]
	if has31 && has32 {
		return true
	}

	// VRF template
	_, has41 := d["ingressVRFID"]
	_, has42 := d["VRFname"]
	if has41 && has42 {
		return true
	}

	// ClassMap template
	_, has51 := d["c3pl_class_name"]
	_, has52 := d["c3pl_class_type"]
	if has51 && has52 {
		return true
	}

	// PolicyMap template
	_, has61 := d["c3pl_policy_name"]
	_, has62 := d["c3pl_policy_type"]
	if has61 && has62 {
		return true
	}
	// sample template
	_, has71 := d["samplerMode"]
	_, has72 := d["samplerName"]
	if has71 && has72 {
		return true
	}

	// drop cause template

	_, has81 := d["drop_cause_name"]
	if has81 {
		return true
	}

	// viptela tloc
	_, has91 := d["tloc_table_overlay_session_id"]
	_, has92 := d["tloc_local_system_ip_address"]
	_, has93 := d["tloc_remote_system_ip_address"]
	if has91 && has92 && has93 {
		return true
	}
	return false
}

func IsSampleRecord(d map[string]interface{}) bool {
	_, has71 := d["samplerMode"]
	_, has72 := d["samplerName"]
	if has71 && has72 {
		return true
	}
	return false
}

func IsFW(d map[string]interface{}) bool {
	_, has1 := d["policy_firewall_event_extended"]
	_, has2 := d["policy_firewall_zone_pair_id"]
	if has1 && has2 {
		return true
	}
	return false
}

func IsNATRecord(d map[string]interface{}) bool {
	_, has1 := d["sourceTransportPort"]
	_, has2 := d["postNAPTSourceTransportPort"]
	if has1 && has2 {
		return true
	}
	/*
		_, has11 := d["destinationTransportPort"]
		_, has12 := d["postNAPTDestinationTransportPort"]

		if has11 && has12 {
			return true
		}
	*/
	return false
}

func IsETARecord(d map[string]interface{}) bool {
	_, has1 := d["sourceTransportPort"]
	_, has2 := d["destinationTransportPort"]

	if has1 && has2 {
		_, has11 := d["ETA_IDP"]
		if has11 {
			return true
		}

		_, has12 := d["ETA_SPLT"]
		if has12 {
			return true
		}

		_, has13 := d["ETA_SALT"]
		if has13 {
			return true
		}

		_, has14 := d["ETA_BD"]
		if has14 {
			return true
		}

		_, has15 := d["TLS_Record"]
		_, has16 := d["TLS_Cipher"]
		if has15 && has16 {
			return true
		}
	}
	return false
}

func isNimbleFlowRecord(d map[string]interface{}) bool {
	_, a := d["queue_depth_sum"]
	_, b := d["queue_depth_taildrop_counter"]
	if a && b {
		return true
	}
	return false
}

func isUrlStatsRecord(d map[string]interface{}) bool {
	_, uris := d["app_http_uri_statistics"]
	if uris {
		return true
	}
	return false
}

func isArtRecord(d map[string]interface{}) bool {
	_, cnd := d["conn_delay_network_to_client_sum"]
	_, snd := d["conn_delay_network_to_server_sum"]
	if cnd && snd {
		return true
	}
	_, llcnd := d["conn_delay_network_long_lived_to_client"]
	_, llsnd := d["conn_delay_network_long_lived_to_server"]
	if llcnd && llsnd {
		return true
	}
	return false
}

func isConnRecord(d map[string]interface{}) bool {
	_, c4valid := d["conn_client_ipv4_address"]
	_, s4valid := d["conn_server_ipv4_address"]
	if c4valid && s4valid {
		return true
	}
	_, c6valid := d["conn_client_ipv6_address"]
	_, s6valid := d["conn_server_ipv6_address"]
	if c6valid && s6valid {
		return true
	}
	return false
}

func isRTPRecord(d map[string]interface{}) bool {
	_, rtppt := d["trans_rtp_payload_type"]
	_, biflow := d["biflowDirection"]
	if rtppt && biflow {
		return true
	}

	return false
}

func isViptelaRecord(d map[string]interface{}) bool {
	_, sessionIDinput := d["overlay_session_id_input"]
	_, sessionIDoutput := d["overlay_session_id_output"]
	if sessionIDinput || sessionIDoutput {
		return true
	}
	return false
}

func isViptelaBFDRecord(d map[string]interface{}) bool {
	_, a := d["overlay_session_id"]
	_, b := d["cisco_sdwan_bfd_update_ts"]
	if a && b {
		return true
	}
	return false
}

func isQueueDropRecord(d map[string]interface{}) bool {
	_, hasidx := d["policy_qos_queue_index"]
	_, hascount := d["policy_qos_queue_drops"]
	if hasidx && hascount {
		return true
	}
	return false
}
