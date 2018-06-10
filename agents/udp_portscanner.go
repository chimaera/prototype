package agents

import (
  "fmt"
  "log"
  "net"
  "time"

  "github.com/chimaera/prototype/core"
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
  return "ip:portscanner"
}

func (c *UDPPortscanner) Register(o *core.Orchestrator) error {
  o.Subscribe(core.NewIP, c.onEndpoint)

  c.orchestrator = o

  // log.Printf("subscribed %s to `new:ip` events", c.ID())

  return nil
}

// TODO: Make this an actually functional UDP scanner...
func (c *UDPPortscanner) onEndpoint(ip string) {
  if c.state.DidProcess(ip) {
    return
  }

  c.state.Add(ip)

  log.Printf("got new IP to scan for UDP ports: %s", ip)

  c.orchestrator.RunTask(func() {
    for _, port := range UDPPorts {
      host := fmt.Sprintf("%s:%d", ip, port)
      conn, _ := net.DialTimeout("udp", host, 1000*time.Millisecond)
      if conn != nil {
        conn.Close()
        c.orchestrator.Publish(core.NewPortTCP, port, ip)
      }
    }
  })
}
