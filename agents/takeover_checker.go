package agents

import (
	"fmt"
	"net"
	"strings"

	"github.com/chimaera/prototype/core"
)

type TakeoverChecker struct {
	state        *core.State
	orchestrator *core.Orchestrator
}

func NewTakeoverChecker() *TakeoverChecker {
	return &TakeoverChecker{
		state: core.NewState(),
	}
}

func (c *TakeoverChecker) ID() string {
	return "dns:takeover"
}

func (c *TakeoverChecker) Register(o *core.Orchestrator) error {
	o.Subscribe(core.NewHostname, c.onEndpoint)
	o.Subscribe(core.NewSubdomain, c.onEndpoint)

	c.orchestrator = o

	return nil
}

func (c *TakeoverChecker) onEndpoint(hostname string) {
	if c.state.DidProcess(hostname, c.ID()) {
		return
	}

	c.state.Add(hostname, c.ID())

	providers := map[string]string{
		"herokuapp.com": "<title>No such app</title>",
		"github.io":     "There isn't a GitHub Pages site here.",
	}

	c.orchestrator.RunTask(func() {
		if cname, err := net.LookupCNAME(hostname); err == nil {
			for pcname, response := range providers {
				if strings.Contains(cname, pcname) {
					if _, body, errs := core.Get(fmt.Sprintf("http://%s/", hostname), 120); len(errs) <= 0 {
						if strings.Contains(body, response) {
							c.orchestrator.Publish(core.NewTakeover, fmt.Sprintf("%s:%s", hostname, cname))
						}
					}
				}
			}
		}
	})
}
