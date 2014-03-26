package main

import (
	"github.com/abh/geodns/countries"
	"github.com/abh/geoip"
	"log"
	"net"
	"strings"
	"time"
)

type GeoIP struct {
	country         *geoip.GeoIP
	hasCountry      bool
	countryLastLoad time.Time

	city         *geoip.GeoIP
	cityLastLoad time.Time
	hasCity      bool
}

var geoIP = new(GeoIP)

func (g *GeoIP) GetCountry(ip net.IP) (country, continent string, netmask int) {
	if g.country == nil {
		return "", "", 0
	}

	country, netmask = geoIP.country.GetCountry(ip.String())
	if len(country) > 0 {
		country = strings.ToLower(country)
		continent = countries.CountryContinent[country]
	}
	return
}

func (g *GeoIP) GetCountryRegion(ip net.IP) (country, continent, regionGroup, region string, netmask int) {
	if g.city == nil {
		log.Println("No city database available")
		country, continent, netmask = g.GetCountry(ip)
		return
	}

	record := geoIP.city.GetRecord(ip.String())
	if record == nil {
		return
	}

	country = record.CountryCode
	region = record.Region
	if len(country) > 0 {
		country = strings.ToLower(country)
		continent = countries.CountryContinent[country]

		if len(region) > 0 {
			region = country + "-" + strings.ToLower(region)
			regionGroup = countries.CountryRegionGroup(country, region)
		}

	}
	return
}

func (g *GeoIP) setDirectory() {
	if len(Config.GeoIP.Directory) > 0 {
		geoip.SetCustomDirectory(Config.GeoIP.Directory)
	}
}

func (g *GeoIP) setupGeoIPCountry() {
	if g.country != nil {
		return
	}

	g.setDirectory()

	gi, err := geoip.OpenType(geoip.GEOIP_COUNTRY_EDITION)
	if gi == nil || err != nil {
		log.Printf("Could not open country GeoIP database: %s\n", err)
		return
	}
	g.countryLastLoad = time.Now()
	g.hasCity = true
	g.country = gi

}

func (g *GeoIP) setupGeoIPCity() {
	if g.city != nil {
		return
	}

	g.setDirectory()

	gi, err := geoip.OpenType(geoip.GEOIP_CITY_EDITION_REV0)
	if gi == nil || err != nil {
		log.Printf("Could not open city GeoIP database: %s\n", err)
		return
	}
	g.countryLastLoad = time.Now()
	g.hasCity = true
	g.city = gi

}
