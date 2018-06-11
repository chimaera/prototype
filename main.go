package main

import (
	// "log"

	"github.com/chimaera/prototype/agents"
	"github.com/chimaera/prototype/core"
)

var (
	olympus = (*core.Orchestrator)(nil)
)

func main() {
	olympus = core.NewOrchestrator(32)

	olympus.Register(agents.NewDNSEnum())
	olympus.Register(agents.NewIPChecker())
	olympus.Register(agents.NewTCPPortscanner())
	olympus.Register(agents.NewUDPPortscanner())
	olympus.Register(agents.NewConfigChecker())
	agents.RegisterPassiveDNSAgents(olympus)

	olympus.Start()

	olympus.Publish("new:hostname", "www.freelancer.com")

	// TODO: Publish Such webhost event for each subdomain found
	// TODO: Also, implement custom HOST header addition
	olympus.Publish("new:webhost", "http://www.madridghosttour.com/")

	olympus.Wait()
}
