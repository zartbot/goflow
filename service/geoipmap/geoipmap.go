package geoipmap

import (
	"fmt"
	"log"
	"net"

	"github.com/olivere/elastic"
	"github.com/zartbot/goflow/lib/iputil"

	geoip2 "github.com/oschwald/geoip2-golang"
)

type GeoIPCollector struct {
	CityDB       *geoip2.Reader
	ASNDB        *geoip2.Reader
	Default_Lat  float64
	Default_Long float64
}

func NewGeoIPCollector(citymmdb string, asnmmdb string, defaultLat float64, defaultLong float64) (*GeoIPCollector, error) {
	var r = GeoIPCollector{}
	var err error
	r.CityDB, err = geoip2.Open(citymmdb)
	r.ASNDB, err = geoip2.Open(asnmmdb)
	r.Default_Lat = defaultLat
	r.Default_Long = defaultLong
	if err != nil {
		log.Fatal(err)
	}
	return &r, err
}

/* Use the following address test hongkong/taiwan/Macau
test_ip = "14.0.207.94"//hongkong :location_region_name == "Hong Kong"
test_ip = "140.112.110.1"//taiwan :location_region_name == "Taiwan"
test_ip = "122.100.160.253" //Macau :location_region_name == "Macau"
*/

type GeoLocation struct {
	City           string
	Region         string
	Country        string
	ASN            string
	AccuracyRadius uint16
	Location       *elastic.GeoPoint `json:"location"`
}

func (g *GeoIPCollector) Lookup(ip net.IP) GeoLocation {
	var r GeoLocation
	if ip.IsGlobalUnicast() {
		if !iputil.IsPrivateIP(ip) {
			c, _ := g.CityDB.City(ip)
			asn, _ := g.ASNDB.ASN(ip)
			if c.City.GeoNameID != 0 {
				r.City = c.City.Names["en"]

			}
			if len(c.Subdivisions) > 0 {
				if c.Subdivisions[0].GeoNameID != 0 {
					r.Region = c.Subdivisions[0].Names["en"]
				}
			}
			if c.Country.GeoNameID != 0 {
				r.Country = c.Country.Names["en"]
			}
			if r.Country == "Hong Kong" {
				r.Country = "China"
				r.Region = "Hong Kong"
			}
			if r.Country == "Macau" {
				r.Country = "China"
				r.Region = "Macau"
			}
			if r.Country == "Taiwan" {
				r.Country = "China"
				r.Region = "Taiwan"
			}

			if r.Country == "" {
				r.Country = "Unknown"
			}
			if r.Region == "" {
				r.Region = "Unknown"
			}
			if r.City == "" {
				r.City = "Unknown"
			}
			r.AccuracyRadius = c.Location.AccuracyRadius
			r.Location = elastic.GeoPointFromLatLon(c.Location.Latitude, c.Location.Longitude)

			r.ASN = fmt.Sprintf("%d::%s", asn.AutonomousSystemNumber, asn.AutonomousSystemOrganization)
		} else {
			r.City = "Local Private Network"
			r.Region = "Local Private Network"
			r.Country = "Local Private Network"
			r.ASN = "0::LAN"
			r.Location = elastic.GeoPointFromLatLon(g.Default_Lat, g.Default_Long)
		}
	} else {
		r.City = "Other"
		r.Region = "Other"
		r.Country = "Other"
		r.ASN = "0::Unknown"
		r.Location = elastic.GeoPointFromLatLon(g.Default_Lat, g.Default_Long)
	}
	return r
}

func UpdateGeoLocationInfo(d map[string]interface{}, g *GeoIPCollector) error {
	var err error
	f, valid := d["flowinfo"]
	if valid {
		ipinfo := f.(iputil.FlowInfo)
		d["Location_A"] = g.Lookup(ipinfo.IPAddress_A)
		d["Location_B"] = g.Lookup(ipinfo.IPAddress_B)
	}
	return err
}
