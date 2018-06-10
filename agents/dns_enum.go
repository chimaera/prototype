package agents

import (
	"fmt"
	// "log"
	"net"

	"github.com/chimaera/prototype/core"

	"github.com/bobesa/go-domain-util/domainutil"
)

type DNSEnum struct {
	state        *core.State
	orchestrator *core.Orchestrator
}

func NewDNSEnum() *DNSEnum {
	return &DNSEnum{
		state: core.NewState(),
	}
}

func (d *DNSEnum) ID() string {
	return "dns:enum"
}

func (d *DNSEnum) Register(o *core.Orchestrator) error {
	// what is this agent interested into?
	o.Subscribe(core.NewHostname, d.onNewHostname)
	// we'll need it to publish results and run tasks
	d.orchestrator = o

	// log.Printf("subscribed %s to `new:hostname` event", d.ID())

	return nil
}

func (d *DNSEnum) onNewHostname(hostname string) {
	domainName := domainutil.Domain(hostname)
	if d.state.DidProcess(domainName) {
		return
	}

	d.state.Add(domainName)

	// log.Printf("got new domain to scan for subdomains: %s", domainName)

	// TODO: load this from a file :P
	wordlist := []string{
		"www",
		"www2",
		"dev",
		"app",
		"beta",
	}

	for _, word := range wordlist {
		// we need this to capture `word`
		func(sub string) {
			d.orchestrator.RunTask(func() {
				hostname := fmt.Sprintf("%s.%s", sub, domainName)
				if _, err := net.LookupHost(hostname); err == nil {
					d.orchestrator.Publish(core.NewSubdomain, hostname)
				}
			})
		}(word)
	}
}
