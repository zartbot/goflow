package identity

import (
	"bytes"
	"net"
	"strings"
	"time"
)

func LookupDNSAddr(ipaddr string, ptrLookUp bool) (map[string]interface{}, int) {
	var result = make(map[string]interface{})
	result["hostType"] = "DNS/SSL"
	var resulttype int = 0
	vA, ok := Service["DNS"].Load(ipaddr)
	if ok {
		result["Name"] = vA.(string)
		resulttype = 1
	} else if ptrLookUp {
		rA, _ := net.LookupAddr(ipaddr)
		var buf bytes.Buffer
		for _, item := range rA {
			if len(item) > 0 {
				//some platform may add dot in suffix
				item = strings.TrimSuffix(item, ".")
				if !strings.HasSuffix(item, ".in-addr.arpa") {
					buf.WriteString(item)
				}
			}
		}
		dnsstr := buf.String()
		result["HostName"] = dnsstr
		resulttype = 1
		//store to cache
		Service["DNS"].Store(ipaddr, dnsstr, time.Now())
	}
	return result, resulttype
}
