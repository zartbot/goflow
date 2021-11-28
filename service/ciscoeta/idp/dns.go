package idp

import (
	"encoding/binary"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/zartbot/goflow/lib/ksyncmap"

	"github.com/zartbot/goflow/lib/iputil"
	"github.com/zartbot/goflow/service/identity"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/zartbot/go_utils/map2value"
)

var DNSRTT ksyncmap.Map

func init() {
	DNSRTT.Name = "IDP_DNS_RTT"
	DNSRTT.Timeout = 5
	DNSRTT.CheckFreq = 5
	DNSRTT.Verbose = false
	go DNSRTT.Run()
}

type DnsMsg struct {
	TrasactionID uint16
	QueryName    string
	AnswersIP    []string
	AnswersName  []string
	AnswersCount int
	ResponseType string
	ResponseCode int
	OpCode       int
}

func (d *DnsMsg) String() string {
	var r string
	r += fmt.Sprintf("------------------------\n")
	r += fmt.Sprintf("    DNS Transaction: %d\n", d.TrasactionID)
	r += fmt.Sprintf("    DNS OpCode: %d\n", d.OpCode)
	r += fmt.Sprintf("    DNS ResponseCode: %s[%d]\n", d.ResponseType, d.ResponseCode)
	r += fmt.Sprintf("    DNS # Answers: %d\n", d.AnswersCount)
	r += fmt.Sprintf("    DNS Question: %s\n", d.QueryName)
	for k, item := range d.AnswersIP {
		r += fmt.Sprintf("    DNS Answer: %s : %s\n", d.AnswersName[k], item)
	}
	return r
}

func IsDNSPacket(d map[string]interface{}) bool {
	ProtocolID, _ := map2value.MapToUInt8(d, "protocolIdentifier")
	srcPort := iputil.FetchPort(d, "src")
	dstPort := iputil.FetchPort(d, "dst")

	if (ProtocolID == 17) && ((srcPort == 53) || (dstPort == 53)) {
		return true
	}
	return false
}

func DecodeDNS(packet gopacket.Packet) []*DnsMsg {
	var result []*DnsMsg
	dnsLayer := packet.Layer(layers.LayerTypeDNS)
	if dnsLayer != nil {
		dns, _ := dnsLayer.(*layers.DNS)
		payload := dnsLayer.LayerContents()
		var transid uint16
		if len(payload) > 4 {
			transid = binary.BigEndian.Uint16(payload[0:2])
		}
		dnsOpCode := int(dns.OpCode)
		dnsResponseCode := int(dns.ResponseCode)
		dnsANCount := int(dns.ANCount)

		for _, dnsQuestion := range dns.Questions {
			d := DnsMsg{TrasactionID: transid,
				QueryName:    string(dnsQuestion.Name),
				OpCode:       dnsOpCode,
				ResponseCode: dnsResponseCode,
				ResponseType: dns.ResponseCode.String(),
				AnswersCount: dnsANCount}

			if dnsANCount > 0 {
				for _, dnsAnswer := range dns.Answers {
					//d.DnsAnswerTTL = append(d.DnsAnswerTTL, fmt.Sprint(dnsAnswer.TTL))
					/*
						//Debug log
							recordString, error := json.MarshalIndent(dnsAnswer, "", "\t")
							if error == nil {
								logrus.Warn(fmt.Sprintf("RecordMap: %s", recordString))
							} else {
								logrus.Warn(fmt.Sprintf("RecordMap: %+v", dnsAnswer))
							}
					*/
					if dnsAnswer.IP != nil {
						d.AnswersIP = append(d.AnswersIP, dnsAnswer.IP.String())
						d.AnswersName = append(d.AnswersName, string(dnsAnswer.Name[:]))
					}
				}
			}
			result = append(result, &d)
		}
	}
	return result
}

func CalculateDNSRTT(d map[string]interface{}, r *DnsMsg) float64 {
	var key string
	var rtt float64
	flowinfo, valid := d["flowinfo"]
	timestamp, _ := map2value.MapToUInt64(d, "flowStartMilliseconds")
	if valid {
		v, _ := flowinfo.(iputil.FlowInfo)
		key = fmt.Sprintf("%s__%d__%s", v.FlowIDWithPort, r.TrasactionID, r.QueryName)
		cachetime, cvalid := DNSRTT.Load(key)
		if cvalid {
			t1 := cachetime.(uint64)
			rtt = math.Abs(float64(t1) - float64(timestamp))
			//DNSRTT.Delete(key) //system has timeout mechanism
		} else {
			DNSRTT.Store(key, timestamp, time.Unix(0, int64(timestamp)*1e6))
			rtt = -1 // set rtt as invalid value, filter later
		}
	}
	/*recordTime := time.Unix(0, int64(timestamp)*1e6)
	logrus.WithFields(logrus.Fields{
		"TimeStamp": recordTime,
		"Key":       key,
		"RTT":       rtt,
	}).Warn(r)*/
	return rtt
}

func ParseDNSPacket(d map[string]interface{}, packet gopacket.Packet) {
	if IsDNSPacket(d) {
		rlist := DecodeDNS(packet)
		for idx, r := range rlist {
			if !strings.HasSuffix(r.QueryName, ".in-addr.arpa") {
				if idx == 0 {
					//d["IDP_DNS"] = r
					d["dns_domain-name"] = r.QueryName
					d["dns_response_code"] = r.ResponseType
					rtt := CalculateDNSRTT(d, r)
					if rtt >= 0 {
						d["dns_rtt"] = rtt
					}
				}
				if r.AnswersCount > 0 {
					for _, item := range r.AnswersIP {
						//add to identity service database
						identity.Service["DNS"].Store(item, r.QueryName, time.Now())
					}
				}
			}
		}
	}
}
