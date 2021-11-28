package identity

import (
	"encoding/csv"
	"log"
	"net"
	"os"
	"sync"

	"github.com/zartbot/go_utils/net/netradix"
	"github.com/zartbot/goflow/lib/iputil"
)

var ReputationInfoDatabase sync.Map

//CSV Format
//based on https://reputation.alienvault.com/reputation.snort
//142.93.136.62,Spamming

func LoadReputationInfo(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Println("[Reputation Database Parser] Error:", err)
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
		var address net.IP
		var hostcategory string
		address = net.ParseIP(r[0])
		hostcategory = r[1]

		/* sync.Map hash key need string type..
		   convert Raw input string to Net.IP format as
		   a validation method to avoid human input error
		   and also filter private address by iputil.IsPrivate
		   function
		*/

		if address != nil && !iputil.IsPrivateIP(address) {
			ReputationInfoDatabase.Store(address.String(), hostcategory)
		}
	}
}

func LookupIPReputation(ipaddr string) (string, int) {
	var result string
	if _repinfo, ok := ReputationInfoDatabase.Load(ipaddr); ok {
		result = _repinfo.(string)
		return result, 8
	}
	return result, 0
}

var BotNet *netradix.NetRadixTree

func LoadBotnetInfo(filename string) error {
	rtree, err := netradix.NewNetRadixTree()
	if err != nil {
		panic(err)
	}

	file, err := os.Open(filename)
	if err != nil {
		log.Println("[BotNet Database Parser] Error:", err)
		return err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	for idx := 1; idx < len(records); idx++ {
		r := records[idx]
		var prefix string
		var hostcategory string
		prefix = r[0]
		hostcategory = r[1]

		if err = rtree.Add(prefix, hostcategory); err != nil {
			log.Fatal(err)
		}
	}
	BotNet = rtree
	return nil
}
