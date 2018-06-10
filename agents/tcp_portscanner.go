package agents

import (
  "fmt"
  // "log"
  "net"
  "time"

  "github.com/chimaera/prototype/core"
)

var (
  TCPPorts = []int{80, 443, 8000, 8080, 8443}
)

type TCPPortscanner struct {
  state        *core.State
  orchestrator *core.Orchestrator
}

func NewTCPPortscanner() *TCPPortscanner {
  return &TCPPortscanner{
    state: core.NewState(),
  }
}

func (c *TCPPortscanner) ID() string {
  return "portscanner:tcp"
}

func (c *TCPPortscanner) Register(o *core.Orchestrator) error {
  o.Subscribe(core.NewIP, c.onEndpoint)

  c.orchestrator = o

  // log.Printf("subscribed %s to `new:ip` events", c.ID())

  return nil
}

func (c *TCPPortscanner) onEndpoint(ip string) {
  if c.state.DidProcess(ip, c.ID()) {
    return
  }

  c.state.Add(ip, c.ID())

  // log.Printf("got new IP to scan for ports: %s", ip)

  c.orchestrator.RunTask(func() {
    for _, port := range TCPPorts {
      host := fmt.Sprintf("%s:%d", ip, port)
      conn, _ := net.DialTimeout("tcp", host, 500*time.Millisecond)
      if conn != nil {
        conn.Close()
        c.orchestrator.Publish(core.NewPortTCP, port, ip)
      }
    }
  })
}
