package identity

import (
	"encoding/csv"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
)

var ServerInfoDatabase sync.Map

type ServerInfo struct {
	Name         string `json:"name"`
	AppGroup     string `json:"appgroup"`
	Owner        string `json:"owner"`
	HostCategory string `json:"host-category"`
	HostType     string `json:"host-type"`

	OS       string
	Version  string
	Location string `json:"location"`

	Ipv4 net.IP `json:"ipv4"`
	Ipv6 net.IP `json:"ipv6"`

	Speed int `json:"speed"`
}

//CSV Format
//name,appgroup,owner,hostcategory,hosttype,os,version,location,ipv4,ipv6,speed

func LoadServerInfo(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Println("[ServerInfo Parser] Error:", err)
		return
	}
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	for idx := 1; idx < len(records); idx++ {
		r := records[idx]
		var server ServerInfo

		server.Name = r[0]
		server.AppGroup = r[1]
		server.Owner = r[2]
		server.HostCategory = r[3]
		server.HostType = r[4]
		server.OS = r[5]
		server.Version = r[6]
		server.Location = r[7]
		server.Ipv4 = net.ParseIP(r[8])
		server.Ipv6 = net.ParseIP(r[9])
		linkspeed, _ := strconv.ParseUint(r[1], 0, 32)
		server.Speed = int(linkspeed)

		/* sync.Map hash key need string type..
		   convert Raw input string to Net.IP format as
		   a validation method to avoid human input error
		*/

		if server.Ipv4 != nil {
			ServerInfoDatabase.Store(server.Ipv4.String(), server)
		}
		if server.Ipv6 != nil {
			ServerInfoDatabase.Store(server.Ipv6.String(), server)
		}
	}
}

func LookupServerAddr(ipaddr string) (map[string]interface{}, int) {
	var result = make(map[string]interface{})
	var resulttype int = 0
	if srvinfoT, ok := ServerInfoDatabase.Load(ipaddr); ok {
		srvinfo, _ := srvinfoT.(ServerInfo)
		result["HostInfo"] = srvinfo
		result["HostName"] = srvinfo.Name
		result["Location"] = srvinfo.Location
		result["Owner"] = srvinfo.Owner
		result["HostType"] = "Server"
		result["OSInfo"] = srvinfo.OS + ":" + srvinfo.Version
		//result =kjson.JsonFlatten(structs.Map(_srvinfo.(ServerInfo)))
		resulttype = 8
		return result, resulttype
	}
	return result, 0
}
