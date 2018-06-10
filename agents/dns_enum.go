package agents

import (
	"fmt"
	"log"
	"time"

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
	o.Subscribe("new:hostname", d.onNewHostname)
	// we'll need it to publish results and run tasks
	d.orchestrator = o

	log.Printf("subscribed %s to `new:hostname` event", d.ID())

	return nil
}

func (d *DNSEnum) onNewHostname(hostname string) {
	domainName := domainutil.Domain(hostname)
	if d.state.DidProcess(domainName) {
		return
	}

	d.state.Add(domainName)

	log.Printf("got new domain to scan for subdomains: %s", domainName)

	// TODO stuff, just emitting values to show the idea
	d.orchestrator.Publish("new:subdomain", fmt.Sprintf("www.%s", domainName))
	time.Sleep(1 * time.Second)
	d.orchestrator.Publish("new:subdomain", fmt.Sprintf("app.%s", domainName))
	time.Sleep(1 * time.Second)
	d.orchestrator.Publish("new:subdomain", fmt.Sprintf("beta.%s", domainName))
}
