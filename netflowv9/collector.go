package netflowv9

import (
	"net"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/zartbot/goflow/datarecord"
)

var TemplateMap sync.Map

type RawMessageUDP struct {
	remoteAddr *net.UDPAddr
	localport  int
	body       []byte
}

//NFv9RawPacketChan is used to dispatch raw packet to multi goroutine workers.
var NFv9RawPacketChan = make(chan *RawMessageUDP, 100)

type NetflowV9Collector struct {
	addr       string
	port       int
	workerNum  int
	stop       bool
	outputChan chan *datarecord.DataFrame
}

func NewNFV9Collector(addr string, port int, outputChan chan *datarecord.DataFrame, workerNum int) *NetflowV9Collector {
	return &NetflowV9Collector{
		addr:       addr,
		port:       port,
		workerNum:  workerNum,
		outputChan: outputChan,
	}
}

func NFv9Worker(output chan *datarecord.DataFrame) {
	var msg *RawMessageUDP
	var ok bool
LOOP:
	for {
		select {
		case msg, ok = <-NFv9RawPacketChan:
			if !ok {
				break LOOP
			}
			PacketParser(msg, output)
		}

	}
}

func (i *NetflowV9Collector) Run() {
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP(i.addr), Port: i.port})
	if err != nil {
		logrus.Fatal(err)
		return
	}

	for idx := 0; idx < i.workerNum; idx++ {
		go NFv9Worker(i.outputChan)
	}

	for !i.stop {
		data := make([]byte, 1500)
		n, remoteAddr, err := listener.ReadFromUDP(data)
		if err != nil {
			logrus.Fatal("[NFv9]error during read: %s", err)
		}

		rcvpkt := &RawMessageUDP{
			remoteAddr: remoteAddr,
			localport:  i.port,
			body:       data[:n],
		}
		NFv9RawPacketChan <- rcvpkt
	}
}
