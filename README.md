# goFlow

goFlow is a golang based Netflow/IPFIX/Cisco HighSpeed Logging collector.
The flow record decoder is strictly followed by `RFC3954` and `RFC7011`

## Supported features & platforms 

### platforms
- Cisco IOS XE series all supported Flexible Netflow records(AVC-ART etc..)
- Cisco Viptela SDWAN
- Cisco IOS XE series ETTA Security feature
- Cisco IOS XE series Firewall and NAT HighSpeed Logging

### features 

- Server list mapping
- information Elemement extension 
- malicious host mapping
- spam and botnet mapping
- DNS ptr mapping
- Geolocation mapping

## Howto build and use it 

check out the `example` folder and use `go build` command to build the file.

if you need to enable geolocation mapping function , please download  ASN and City mmdb and put them in the `example/config/database/geolocation` folder

#### Config performance monitor on IOS XE Router

```bash
Router(config)#performance monitor context foo profile application-performance
Router(config-perf-mon)# mode optimized
Router(config-perf-mon)# exporter destination 192.168.99.101 source TenGigabitEthernet0/0/0.43 
Router(config-perf-mon)# traffic-monitor all

Router(config)# interface Gi0/1/0
Router(config-if)# performance monitor context foo
```

execute the goFlow example

```bash
zartbot@zartbotLinux:~$ ./example
```

#### example output

```bash
[IPFIX] | ==========================================================DataFrame Start========================================================
[IPFIX] | AgentIP: 192.168.99.199       | ExportTime: 2021-11-28 17:01:51 +0800 CST | DomainID/SetID: [  2560:309   ] 
[IPFIX] | Record:[   0]--------------------------------------------------------------------------------------------------------------------
[IPFIX] | Name: applicationAttribute                | Value: {Category:business-and-productivity-tools SubCategory:desktop-virtualization GroupName:other TrafficClass:multimedia-streaming BusinessRelevance:business-relevant p2pTech:no TunnelTech:no EncryptTech:yes SetName:desktop-virtualization-apps FamilyName:thin-client}
[IPFIX] | Name: flowinfo                            | Value: {IPAddress_A:10.1.1.1 Hash_A:942535861 Hash_B:3031818525 Hash_A_Mod:5861 Hash_B_Mod:8525 Port_A:0 IPAddress_B:101.25.1.1 Port_B:12346 FlowID:10.1.1.1-101.25.1.1 FlowIDWithPort:10.1.1.1:0-101.25.1.1:12346 Direction:true}
[IPFIX] | Name: hosta                               | Value: map[HostName: hostType:DNS/SSL risk:0]
[IPFIX] | Name: flowname                            | Value: <-->
[IPFIX] | Name: newConnectionDeltaCount             | Value: 0
[IPFIX] | Name: services_waas_segment               | Value: 16
[IPFIX] | Name: connectionSumDurationSeconds        | Value: 10120
[IPFIX] | Name: conn_client_cnt_bytes_network_long  | Value: 805
[IPFIX] | Name: applicationId                       | Value: 50335820
[IPFIX] | Name: flowDirection                       | Value: 1
[IPFIX] | Name: services_waas_passthrough_reason    | Value: 0
[IPFIX] | Name: flowStartSysUpTime                  | Value: 3527529164
[IPFIX] | Name: ipDiffServCodePoint                 | Value: 48
[IPFIX] | Name: observationInterface                | Value: GigabitEthernet0/1/0
[IPFIX] | Name: ingressVRFID                        | Value: 0
[IPFIX] | Name: AgentID                             | Value: 192.168.99.199
[IPFIX] | Name: Type                                | Value: conn
[IPFIX] | Name: hostb                               | Value: map[HostName: hostType:DNS/SSL risk:0]
[IPFIX] | Name: conn_client_ipv4_address            | Value: 10.1.1.1
[IPFIX] | Name: biflowDirection                     | Value: 1
[IPFIX] | Name: monitoringIntervalStartMilliSeconds | Value: 1638090000000
[IPFIX] | Name: flowtype                            | Value: 2
[IPFIX] | Name: flowEndSysUpTime                    | Value: 3527532852
[IPFIX] | Name: conn_server_cnt_bytes_network_long  | Value: 0
[IPFIX] | Name: applicationName                     | Value: {Name:pcoip Description:PC-over-IP - Virtual Desktop Infrastructure}
[IPFIX] | Name: conn_server_ipv4_address            | Value: 101.25.1.1
[IPFIX] | Name: conn_server_trans_port              | Value: 12346
[IPFIX] | Name: observationPointId                  | Value: 2
[IPFIX] | Name: responderPackets                    | Value: 0
[IPFIX] | Name: initiatorPackets                    | Value: 5
[IPFIX] | Name: CreateAt                            | Value: 2021-11-28 17:01:51 +0800 CST
[IPFIX] | Name: flowrisk                            | Value: 0
[IPFIX] | Name: protocolIdentifier                  | Value: 17
```


## Development

1. Load Information Elements database.

unlike other Netflow collector, goFlow use csv file defined all information elements,

```golang
	iedb.ReadIANA("./config/database/informationElement/iana_ie.csv")
	iedb.ReadCiscoIE("./config/database/informationElement/cisco_ie.csv")
	iedb.ReadCiscoIE("./config/database/informationElement/key_ie.csv")

	iedb.ShowIEDB()
```
csv format as below:

```
ElementID,EnterpriseNo,FieldType,DataType
1,0,octetDeltaCount,unsigned64
2,0,packetDeltaCount,unsigned64
3,0,deltaFlowCount,unsigned64
4,0,protocolIdentifier,unsigned8
5,0,ipClassOfService,unsigned8
...

```

2. Start Identity Services

goFlow could leverage the existing database to search the server identitiy

```golang
	logrus.Info("Starting Identification service....")
	identity.LoadServerInfo("./config/database/server/server.csv")
	identity.LoadReputationInfo("./config/database/security/reputation.csv")
	identity.LoadBotnetInfo("./config/database/security/spam_bot_net.csv")
	go identity.Run()
```

3. Create Flow collectors 

all flow record will be sent out to a streaming golang channel  `chan *datarecord.DataFrame`
goFlow support multiple thread decoding, you may change the woker number in `NewIPFIXCollector` and `NewNFV9Collector`

```golang
	//IPFIXDataFrameChan is dataframe output channel for future streaming process.
	var IPFIXDataFrameChan = make(chan *datarecord.DataFrame, 10000)
	var NFv9DataFrameChan = make(chan *datarecord.DataFrame, 10000)

	//input Collector
	var ipfixCollector []*ipfix.IPFIXCollector
	for _, port := range []int{4739, 4740} {
		c := ipfix.NewIPFIXCollector(ipaddr, port, IPFIXDataFrameChan, 1)
		ipfixCollector = append(ipfixCollector, c)
	}

	var nfv9Collector []*netflowv9.NetflowV9Collector
	for _, port := range []int{2055, 2056, 2057} {
		c := netflowv9.NewNFV9Collector(ipaddr, port, NFv9DataFrameChan, 1)
		nfv9Collector = append(nfv9Collector, c)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	for _, c := range ipfixCollector {
		go c.Run()
	}

	for _, c := range nfv9Collector {
		go c.Run()
	}

```    

4. Consuming the datastream

you may use the following method to start and comsuming the datarecord:

```golang
	for {
		select {
		case e1 := <-IPFIXDataFrameChan:
			RecordMap(e1)
			e1.Print("IPFIX")

		case e2 := <-NFv9DataFrameChan:
			RecordMap(e2)
			e2.Print("NFv9")
		}
	}
```

You can  update and correlate record with optional template and  update record with addtional services here.

```golang
func RecordMap(d *datarecord.DataFrame) {
	d.TypeAssertion()
	createAt := time.Unix(int64(d.ExportTime), 0)
	for _, value := range d.Record {
		value["CreateAt"] = createAt
		value["AgentID"] = d.AgentID
		value["Type"] = d.Type
		//Optional Template Correlation
		optiontemplatemap.UpdateInterfaceMap(value, d.AgentID)
		optiontemplatemap.UpdateAppMap(value, d.AgentID)
		optiontemplatemap.UpdateCiscoVarString(value, d.AgentID)
		optiontemplatemap.UpdateC3PLMap(value, d.AgentID)
		optiontemplatemap.UpdateDropCauseMap(value, d.AgentID)
		optiontemplatemap.UpdateViptelaTLOCMap(value, d.AgentID)
		optiontemplatemap.UpdateFWEventMap(value, d.AgentID)
		optiontemplatemap.UpdateFWZonePairMap(value, d.AgentID)
		optiontemplatemap.UpdateFWClassMapMap(value, d.AgentID)

		//Additional Service
		flowinfo.UpdateFlowInfo(value, true)
		ciscoavc.ExtendField(value)
		ciscoeta.ExtendField(value)
	}
}
```

## Export to external Database

goFlow export as a golang channel which can be easily consumed by kafka and other applications. the output module has a simple elastic-search exporter.
But we strongly recommend you to use `zartbot/golap` project for data pre-aggregation and in-line AI based performance scoring and reasoning.
