package config

import (
	"io/ioutil"

	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	CollectorAddr                   string `yaml:"ip"`
	IpfixPort                       []int  `yaml:"ipfix_port"`
	Nfv9Port                        []int  `yaml:"nfv9_port"`
	CiscoIEPath                     string `yaml:"cisco_ie"`
	IanaIEPath                      string `yaml:"iana_ie"`
	KeyIEPath                       string `yaml:"key_ie"`
	GeoCityPath                     string `yaml:"geo_city"`
	GeoASNPath                      string `yaml:"geo_asn"`
	ElasticURI                      string `yaml:"elastic_uri"`
	FlowNumWorkers                  int    `yaml:"collector_worker_num"`
	StreamNumWorkers                int    `yaml:"stream_worker_num"`
	ElasticNumWorkers               int    `yaml:"elastic_worker_num"`
	ElasticFlushInverval            int    `yaml:"elastic_flush_interval"`
	ElasticPerRecordIndexPrefix     string `yaml:"elastic_per_record_index"`
	ElasticPerRecordOutput          bool   `yaml:"elastic_per_record"`
	ElasticPerRecordMappingFileName string `yaml:"elastic_per_record_mapping"`
	ElasticPerRecordMapping         string
	ElasticSumamryIndexPrefix       string `yaml:"elastic_summary_index"`
	ElasticSummaryMappingFileName   string `yaml:"elastic_summary_mapping"`
	ElasticSummaryMapping           string
	Verbose                         bool    `yaml:"verbose"`
	GeoDefaultLat                   float64 `yaml:"geo_default_lat"`
	GeoDefaultLong                  float64 `yaml:"geo_default_long"`
	HTTPSPort                       int     `yaml:"https_port"`
	SSLCert                         string  `yaml:"ssl_cert"`
	SSLKey                          string  `yaml:"ssl_key"`
	ServerInfoDBPath                string  `yaml:"server_info_db"`
	ReputationDBPath                string  `yaml:"ip_reputation_db"`
	BotNetDBPath                    string  `yaml:"botnet_db"`
	GolapStatsConfigFilePath        string  `yaml:"golap_stats"`
}

func (c *Config) GetConf(filename string) *Config {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		logrus.Fatalf("[Error]Config file fetch error, %v", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		logrus.Fatalf("Unmarshal: %v", err)
	}
	perRecordMapping, err := ioutil.ReadFile(c.ElasticPerRecordMappingFileName)
	if err != nil {
		logrus.Fatalf("[Error]ElasticSearch Per Record Mapping Config file fetch error, %v", err)
	} else {
		c.ElasticPerRecordMapping = string(perRecordMapping)
	}

	summaryMapping, err := ioutil.ReadFile(c.ElasticSummaryMappingFileName)
	if err != nil {
		logrus.Fatalf("[Error]ElasticSearch Summary Mapping Config file fetch error, %v", err)
	} else {
		c.ElasticSummaryMapping = string(summaryMapping)
	}

	return c
}
