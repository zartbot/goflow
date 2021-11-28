package output

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/olivere/elastic"
	"github.com/zartbot/goflow/datarecord"
)

//NewElasticClient is a construct function for Client
func NewElasticClient(url string) (*elastic.Client, error) {
	var (
		err           error
		elasticClient *elastic.Client
	)
	for {
		elasticClient, err = elastic.NewClient(
			elastic.SetURL(url),
			elastic.SetSniff(false),
		)
		if err != nil {
			log.Println(err)
			time.Sleep(3 * time.Second)
		} else {
			break
		}
	}
	return elasticClient, err
}

//ElasticBulkProcessor elastic bulk import processor
type ElasticBulkProcessor struct {
	C             *elastic.Client
	P             *elastic.BulkProcessor
	Workers       int
	IndexPrefix   string
	FlushInterval int
	StopC         chan struct{} // stop channel for caller
	CreateIdxFlag chan struct{}
}

// Run starts the ElasticBulkProcessor.
func (b *ElasticBulkProcessor) Run() error {
	// Start bulk processor
	p, err := b.C.BulkProcessor().
		Workers(b.Workers). // # of workers
		BulkActions(5000).  // # of queued requests before committed//
		//BulkSize(1024000000).                                             // # of bytes in requests before committed
		FlushInterval(time.Duration(b.FlushInterval) * time.Millisecond). // autocommit every interval milliseconds
		Do(context.Background())
	if err != nil {
		return err
	}
	b.P = p
	// Start indexer that pushes data into bulk processor
	b.StopC = make(chan struct{})
	b.CreateIdxFlag = make(chan struct{})
	go b.ensureIndex()
	<-b.CreateIdxFlag
	return nil
}

// Close the bulker.
func (b *ElasticBulkProcessor) Close() error {
	b.StopC <- struct{}{}
	<-b.StopC
	close(b.StopC)
	return nil
}

//ElasticRecordProcessor is used bulkupload Record to ElasticSearch Server
func (b *ElasticBulkProcessor) ElasticRecordProcessor(dfchan chan *datarecord.DataFrame) {
	var (
		stop bool
		//rt    string
		//valid bool
	)

	for !stop {
		select {
		case <-b.StopC:
			stop = true
		case d := <-dfchan:
			exporttime := time.Unix(int64(d.ExportTime), 0).Format("2006.01.02")
			indexName := b.IndexPrefix + "-" + d.Type + "-" + exporttime
			//fmt.Println("DEBUG:::", indexName, exporttime)
			rl := d.RecordList()
			for _, item := range rl {
				//rt, valid = item["Type"].(string)
				//if valid {
				r := elastic.NewBulkIndexRequest().Index(indexName).Type("log").Doc(item)
				b.P.Add(r)
				//fmt.Printf("\n\n\nDEBUG[%s]::BULK-RECORD%+v\n\n\n", indexName, item)
				//}
			}
			// Sleep for a short time.
			time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
		}
	}
	b.StopC <- struct{}{} // ack stopping
}

//ElasticImport :Bulk import record to ElasticSearch
func ElasticBulkImportProcessor(url string, numWorkers int, indexPrefix string, FlushInterval int) (*ElasticBulkProcessor, error) {
	elasticClient, err := NewElasticClient(url)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("ElasticSearch Server connected.....")
	BulkProcessor := &ElasticBulkProcessor{
		C:             elasticClient,
		Workers:       numWorkers,
		IndexPrefix:   indexPrefix,
		FlushInterval: FlushInterval,
	}
	err = BulkProcessor.Run()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("ElasticSearch Bulk Processor is running...")
	return BulkProcessor, err
}

func (b *ElasticBulkProcessor) ensureIndex() {
	const date = "2006.01.02"
	var indexname, nextindexname string
	startflag := true

	mapping := `{
		"settings":{
            "index.refresh_interval":"30s",
			"number_of_shards":3,
			"number_of_replicas":0
		},
        "mappings":{
            "log":{
                "properties":{
                    "Location_A":{
                        "properties" : {
                            "location": {
                                "type":"geo_point"
                            }
                        }
                    },
                    "Location_B":{
                        "properties" : {
                            "location": {
                                "type":"geo_point"
                            }
                        }
					},
					"conn_delay_app_mean" : {
						"type":"float"
					},
					"conn_delay_client_to_server_mean" : {
						"type":"float"
					},
					"conn_delay_to_server_mean" : {
						"type":"float"
					},
					"conn_duration_mean" : {
						"type":"float"
					},
					"conn_delay_network_to_client_mean" : {
						"type":"float"
					},
					"conn_delay_network_to_server_mean" : {
						"type":"float"
					},
					"conn_delay_network_mean" : {
						"type":"float"
					},
					"conn_delay_network_client_to_server_mean" : {
						"type":"float"
					},
					"conn_transaction_duration_mean" : {
						"type":"float"
					},
					"conn_client_pkts_retransmit_rate" : {
						"type":"float"
					},
					"conn_server_pkts_retransmit_rate" : {
						"type":"float"
					},
					"conn_client_bytes_retransmit_rate" : {
						"type":"float"
					},
					"conn_server_bytes_retransmit_rate" : {
						"type":"float"
					}
                }
            }
		}
    }`

	if b.IndexPrefix == "" {
		log.Fatal("No indexprefix name.")
		return
	}
	for {
		currentTick := time.Now().Format(date)
		nextTick := time.Now().Add(time.Duration(3600*24) * time.Second).Format(date)
		for _, item := range datarecord.RecordTypeList {
			indexname = fmt.Sprintf("%s-%s-%s", b.IndexPrefix, item, currentTick)
			nextindexname = fmt.Sprintf("%s-%s-%s", b.IndexPrefix, item, nextTick)
			//fmt.Println(indexname, nextindexname)

			exists, _ := b.C.IndexExists(indexname).Do(context.Background())
			if !exists {
				log.Println("Creating index:", indexname)
				_, err := b.C.CreateIndex(indexname).BodyString(mapping).Do(context.Background())
				if err != nil {
					log.Println(err)
				}
			}

			nextExists, _ := b.C.IndexExists(nextindexname).Do(context.Background())
			if !nextExists {
				log.Println("Preparing nextday index:", nextindexname)
				_, err := b.C.CreateIndex(nextindexname).BodyString(mapping).Do(context.Background())
				if err != nil {
					log.Println(err)
				}
			}

		}

		if startflag {
			b.CreateIdxFlag <- struct{}{}
			startflag = false
			log.Println("Index Created... ")
		}

		time.Sleep(120 * time.Second)

	}

}
