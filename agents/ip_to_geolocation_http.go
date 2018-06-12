package agents

import (
	"encoding/json"
	"github.com/chimaera/prototype/core"
	"io/ioutil"
)

// NOTE: module implemented with help of:
// https://www.devdungeon.com/content/ip-geolocation-go
type GeoIP struct {
	Ip          string  `json:"ip"`
	CountryCode string  `json:"country_code"`
	CountryName string  `json:"country_name""`
	RegionCode  string  `json:"region_code"`
	RegionName  string  `json:"region_name"`
	City        string  `json:"city"`
	Zipcode     string  `json:"zipcode"`
	Lat         float32 `json:"latitude"`
	Lon         float32 `json:"longitude"`
	MetroCode   int     `json:"metro_code"`
	AreaCode    int     `json:"area_code"`
}

// TODO: these can obviously be configurable later, and it would
// be nice to support multiple backends and spread the queries
// across them.
const IP2GeoHost = "https://freegeoip.net/json/"

type IP2Geo struct {
	state        *core.State
	orchestrator *core.Orchestrator
}

func NewIP2Geo() *IP2Geo {
	return &IP2Geo{
		state: core.NewState(),
	}
}

func (c *IP2Geo) ID() string {
	return "ip:geo:checker"
}

func (c *IP2Geo) Register(o *core.Orchestrator) error {
	o.Subscribe(core.NewIP, c.onEndpoint)

	c.orchestrator = o

	return nil
}

func (c *IP2Geo) onEndpoint(ip string) {
	if c.state.DidProcess(ip, c.ID()) {
		return
	}

	c.state.Add(ip, c.ID())

	c.orchestrator.RunTask(func() {
		// NOTE: this is sort of hacky, using the same netClient
		// from the whois http checking agent.
		resp, err := netClient.Get(IP2GeoHost + ip)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		geo := GeoIP{}

		err = json.Unmarshal(body, &geo)
		if err != nil {
			panic(err)
		}

		c.orchestrator.Publish(core.NewIPGeo, geo)
	})
}
