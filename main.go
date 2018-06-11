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

	agents.RegisterPassiveAgents(olympus)

	olympus.Start()

	olympus.Publish("new:hostname", "www.freelancer.com")

	olympus.Wait()
}
