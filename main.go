package main

import (
	"net/http"

	"github.com/chimaera/prototype/agents"
	"github.com/chimaera/prototype/agents/scraper"
	"github.com/chimaera/prototype/core"
	"github.com/chimaera/prototype/db"
)

var (
	olympus = (*core.Orchestrator)(nil)
)

func main() {
	olympus = core.NewOrchestrator(32)

	go http.ListenAndServe(":8181", core.NewHTTPProxy())
	defer core.GetHTTPResponseCacheInstance().DestroyCacheStore()

	olympus.Register(agents.NewDNSEnum())
	olympus.Register(agents.NewIPChecker())
	olympus.Register(agents.NewWhoisChecker())
	olympus.Register(agents.NewTCPPortscanner())
	// olympus.Register(agents.NewUDPPortscanner())
	olympus.Register(agents.NewConfigChecker())
	olympus.Register(agents.NewTakeoverChecker())
	olympus.Register(scraper.NewSimpleScraper())
	agents.RegisterPassiveDNSAgents(olympus)

	inEvent := "new:hostname"
	inType := db.NodeTypeHostname
	inValue := "www.freelancer.com"

	olympus.Start(inEvent, inType, inValue)

	// dbase.Root().Print(log.Printf, 0)

	olympus.Wait()
}
