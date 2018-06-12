package agents

import (
	"fmt"
	// "log"
	"net"
	"time"

	"github.com/chimaera/prototype/core"
	"github.com/chimaera/prototype/db"
)

var (
	UDPPorts = []int{53, 69, 137, 138, 161, 162}
)

type UDPPortscanner struct {
	state        *core.State
	orchestrator *core.Orchestrator
}

func NewUDPPortscanner() *UDPPortscanner {
	return &UDPPortscanner{
		state: core.NewState(),
	}
}

func (c *UDPPortscanner) ID() string {
	return "portscanner:udp"
}

func (c *UDPPortscanner) Register(o *core.Orchestrator) error {
	o.Subscribe(core.NewIP, c.onEndpoint)

	c.orchestrator = o

	// log.Printf("subscribed %s to `new:ip` events", c.ID())

	return nil
}

// TODO: Make this an actual functional UDP scanner...
func (c *UDPPortscanner) onEndpoint(ip string) {
	if c.state.DidProcess(ip, c.ID()) {
		return
	}

	c.state.Add(ip, c.ID())

	// log.Printf("got new IP to scan for UDP ports: %s", ip)
	parent := db.Current.Search(db.NodeTypeIP, ip)

	c.orchestrator.RunTask(func() {
		for _, port := range UDPPorts {
			host := fmt.Sprintf("%s:%d", ip, port)
			conn, _ := net.DialTimeout("udp", host, 1000*time.Millisecond)
			if conn != nil {
				conn.Close()
				parent.Add(db.NodeTypePort, port)
				c.orchestrator.Publish(core.NewPortUDP, port, ip)
			}
		}
	})
}
