package agents

import (
	"github.com/chimaera/prototype/agents/passivedns/crtsh"
	"github.com/chimaera/prototype/core"
)

func RegisterPassiveDNSAgents(o *core.Orchestrator) {
	o.Register(crtsh.NewCrtsh())

	return
}
