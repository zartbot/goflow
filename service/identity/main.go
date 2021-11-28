package identity

import (
	"math/rand"
	"sync"
	"time"

	"github.com/zartbot/goflow/lib/ksyncmap"
)

type DeviceInfo struct {
	Location string
	Owner    string
	HostName string
	HostType string
	OSInfo   string
}

/*
TODO : Add Server Indetify service and other RESTAPI support
       Reduce the previously definition...
*/
var Service = make(map[string]*ksyncmap.Map)

func init() {
	Service["WlanClient"] = &ksyncmap.Map{
		Timeout:   86400,
		CheckFreq: 720,
		Verbose:   false,
	}

	/*
		Service["Server"] = &ksyncmap.Map{
			Timeout:   86400,
			CheckFreq: 720,
			Verbose:   false,
		}*/

	Service["AP"] = &ksyncmap.Map{
		Timeout:   86400,
		CheckFreq: 720,
		Verbose:   false,
	}
	Service["ARP"] = &ksyncmap.Map{
		Timeout:   86400,
		CheckFreq: 720,
		Verbose:   false,
	}
	Service["DNS"] = &ksyncmap.Map{
		Timeout:   86400,
		CheckFreq: 720,
		Verbose:   false,
	}
	Service["NAT"] = &ksyncmap.Map{
		Timeout:   7200,
		CheckFreq: 120,
		Verbose:   false,
	}
}

func Set(Timeout int64, CheckFreq int64, verbose bool) {
	for key, v := range Service {
		v.CheckFreq = CheckFreq
		v.Timeout = Timeout
		v.Name = key
		v.Verbose = verbose
	}
}

func Run() {
	rand.Seed(time.Now().UnixNano())
	var wg sync.WaitGroup
	wg.Add(3)
	for _, v := range Service {
		go v.Run()
		time.Sleep(time.Millisecond * time.Duration(rand.Int63n(1000)))
	}
	wg.Wait()
}
