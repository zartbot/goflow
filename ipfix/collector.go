package ipfix

import (
	"log"
	"net"
	"sync"

	"github.com/zartbot/goflow/datarecord"
)

var TemplateMap sync.Map

type IPFIXCollector struct {
	addr       string
	port       int
	workerNum  int
	stop       bool
	outputChan chan *datarecord.DataFrame
}

type RawMessageUDP struct {
	remoteAddr *net.UDPAddr
	localport  int
	body       []byte
}

//IPFIXRawPacketChan is used to dispatch raw packet to multi goroutine workers.
var IPFIXRawPacketChan = make(chan *RawMessageUDP, 100)

func NewIPFIXCollector(addr string, port int, outputChan chan *datarecord.DataFrame, workerNum int) *IPFIXCollector {
	return &IPFIXCollector{
		addr:       addr,
		port:       port,
		workerNum:  workerNum,
		outputChan: outputChan,
	}
}

func IPFIXWorker(output chan *datarecord.DataFrame) {

	var msg *RawMessageUDP
	var ok bool
LOOP:
	for {
		select {
		case msg, ok = <-IPFIXRawPacketChan:
			if !ok {
				break LOOP
			}
			PacketParser(msg, output)
		}

	}
}

func (i *IPFIXCollector) Run() {
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP(i.addr), Port: i.port})
	if err != nil {
		log.Println(err)
		return
	}

	for idx := 0; idx < i.workerNum; idx++ {
		go IPFIXWorker(i.outputChan)
	}

	for !i.stop {
		data := make([]byte, 1500)
		n, remoteAddr, err := listener.ReadFromUDP(data)
		if err != nil {
			log.Printf("[IPFIX Collector]error during read: %s", err)
		}

		rcvpkt := &RawMessageUDP{
			remoteAddr: remoteAddr,
			localport:  i.port,
			body:       data[:n],
		}

		//logrus.Warn("Recive Packet: ", remoteAddr.IP.String())

		//if rcvpkt.remoteAddr.IP.String() == "192.168.99.32" {
		IPFIXRawPacketChan <- rcvpkt
		//}

	}
}
