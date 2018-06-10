package agents

import (
	// "log"
	"fmt"
	"net"
	"strings"

	"github.com/chimaera/prototype/core"
)

type IPChecker struct {
	state        *core.State
	orchestrator *core.Orchestrator
}

func NewIPChecker() *IPChecker {
	return &IPChecker{
		state: core.NewState(),
	}
}

func (c *IPChecker) ID() string {
	return "ip:checker"
}

func (c *IPChecker) Register(o *core.Orchestrator) error {
	o.Subscribe(core.NewHostname, c.onEndpoint)
	o.Subscribe(core.NewSubdomain, c.onEndpoint)

	c.orchestrator = o

	// log.Printf("subscribed %s to `new:hostname` and `new:subdomain` events", c.ID())

	return nil
}

func (c *IPChecker) onEndpoint(hostname string) {
	if c.state.DidProcess(hostname, c.ID()) {
		return
	}

	c.state.Add(hostname, c.ID())

	// log.Printf("got new endpoint to scan for ip addresses: %s", hostname)

	c.orchestrator.RunTask(func() {
		if addrs, err := net.LookupHost(hostname); err == nil {
			for _, addr := range addrs {
				if strings.Contains(addr, ":") {
					addr = fmt.Sprintf("[%s]", addr)
				}
				c.orchestrator.Publish(core.NewIP, addr)
			}
		}
	})
}
