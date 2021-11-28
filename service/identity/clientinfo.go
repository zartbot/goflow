package identity

import (
	"net"
)

type WirelessClientInfo struct {
	MAC      string
	UserName string
	HostType string
	Location string `json:"location"`
	Speed    int    `json:"speed"`
	Device   struct {
		HostName   string
		DeviceType string `json:"type"`
		OS         string
		Version    string
	} `json:"deviceinfo"`
	Network struct {
		Ipv4    net.IP `json:"ipv4"`
		Ipv6    net.IP `json:"ipv6"`
		Gateway net.IP `json:"gateway"`
		//DNS     []net.IP `json:"dns-servers"`
	} `json:"network"`
	Wireless struct {
		APMAC         string  `json:"ap-mac"`
		RSSI          int     `json:"rssi"`
		SNR           int     `json:"snr"`
		SpatialStream int     `json:"spatial-stream"`
		PowerSaveMode int     `json:"power-save"`
		MaxSpeed      float64 `json:"max-data-rate"`
		RadioSlot     int     `json:"radio-slot"`
		WlanID        int     `json:"wlan-id"`
		SSID          string  `json:"ssid"`
		State         int     `json:"state"`
		AssocState    int     `json:"assoc-state"`
		ConnectedTime uint32  `json:"connected-time"`

		APNAME            string `json:"ap-name"`
		Channel           int    `json:"channel"`
		Interference      int    `json:"interference"`
		Utilization       int    `json:"utilization"`
		APTxPower         int    `json:"ap-tx-power"`
		PowerChangeReason int    `json:"pwr-change-reason"`
		APTrafficLoad     int    `json:"ap-traffic-load"`
		APDataRate        int    `json:"ap-traffic-data-rate"`
		MinAirQuality     int    `json:"min-air-quality`
		AirQuality        int    `json:"air-quality`
		Noise             int    `json:"noise"`

		//AuthMode      bool    `json:"auth-mode"`
		//AssociateAP APInfo `json:"ap"`
	} `json:"wireless"`
}

func LookupUserAddr(ipaddr string) (map[string]interface{}, int) {
	var result = make(map[string]interface{})
	var resulttype int = 0
	//TODO: add static wired client info
	if _mac, ok := Service["ARP"].Load(ipaddr); ok {
		if mac, ok := _mac.(string); ok && len(mac) > 0 {
			if clientinfo, valid := Service["WlanClient"].Load(mac); valid {

				clientinfo, _ := clientinfo.(WirelessClientInfo)
				result["HostInfo"] = clientinfo
				//result = kjson.JsonFlatten(structs.Map(clientinfo.(WirelessClientInfo)))
				result["hostType"] = "Client"

				result["Location"] = clientinfo.Location
				result["Owner"] = clientinfo.UserName
				result["HostName"] = clientinfo.Device.HostName
				result["OSInfo"] = clientinfo.Device.OS + ":" + clientinfo.Device.Version
				resulttype = 64
				return result, resulttype
			}
		}
	}
	return result, 0
}
