package agents

import (
	// "log"
	"bufio"
	"net"
	"strings"
	"time"

	"github.com/chimaera/prototype/core"
)

// Note: Using HTTP whois resources could be an option.

// TODO: these can obviously be configurable later
const WhoisHost = "whois.verisign-grs.com"
const WhoisPort = "43"

type WhoisChecker struct {
	state        *core.State
	orchestrator *core.Orchestrator
}

func NewWhoisChecker() *WhoisChecker {
	return &WhoisChecker{
		state: core.NewState(),
	}
}

func (c *WhoisChecker) ID() string {
	return "whois:checker"
}

func (c *WhoisChecker) Register(o *core.Orchestrator) error {
	o.Subscribe(core.NewHostname, c.onEndpoint)
	o.Subscribe(core.NewSubdomain, c.onEndpoint)

	c.orchestrator = o

	// log.Printf("subscribed %s to `new:hostname` and `new:subdomain` events", c.ID())

	return nil
}

type WhoisRecord struct {
	DomainName           string
	RegistryDomainID     string
	RegistrarWHOISServer string
	RegistrarURL         string
	UpdatedDate          string
	CreationDate         string
	RegistryExpiryDate   string
	Registrar            string
	RegistrarIANAID      string
	NameServers          []string
}

func (c *WhoisChecker) onEndpoint(hostname string) {
	if c.state.DidProcess(hostname, c.ID()) {
		return
	}

	c.state.Add(hostname, c.ID())

	// log.Printf("got new endpoint to scan for ip addresses: %s", hostname)

	c.orchestrator.RunTask(func() {
		// remove this part of string
		hostname = strings.Replace(hostname, "www.", "", 1)

		conn, err := net.DialTimeout("tcp", net.JoinHostPort(WhoisHost, WhoisPort), time.Second*30)
		if err != nil {
			panic(err)
		}
		defer conn.Close()

		conn.Write([]byte("domain " + hostname + "\r\n"))

		scanner := bufio.NewScanner(conn)

		wr := WhoisRecord{NameServers: []string{}}

		for scanner.Scan() {
			line := scanner.Text()
			// Note: this could probably be improved, might be too defensive...
			// maybe not defensive enough? I dunno. Prototyping! :D
			switch {
			case strings.Contains(line, "Domain Name:"):
				s := strings.Split(line, ": ")
				wr.DomainName = strings.Join(s[1:], "")
			case strings.Contains(line, "Registry Domain ID:"):
				s := strings.Split(line, ": ")
				wr.RegistryDomainID = strings.Join(s[1:], "")
			case strings.Contains(line, "Registrar WHOIS Server"):
				s := strings.Split(line, ": ")
				wr.RegistrarWHOISServer = strings.Join(s[1:], "")
			case strings.Contains(line, "Registrar URL"):
				s := strings.Split(line, ": ")
				wr.RegistrarURL = strings.Join(s[1:], "")
			case strings.Contains(line, "Updated Date"):
				s := strings.Split(line, ": ")
				wr.UpdatedDate = strings.Join(s[1:], "")
			case strings.Contains(line, "Creation Date"):
				s := strings.Split(line, ": ")
				wr.CreationDate = strings.Join(s[1:], "")
			case strings.Contains(line, "Registry Expiry Date"):
				s := strings.Split(line, ": ")
				wr.RegistryExpiryDate = strings.Join(s[1:], "")
			case strings.Contains(line, "Registrar"):
				s := strings.Split(line, ": ")
				wr.Registrar = strings.Join(s[1:], "")
			case strings.Contains(line, "Registrar IANA ID"):
				s := strings.Split(line, ": ")
				wr.RegistrarIANAID = strings.Join(s[1:], "")
			case strings.Contains(line, "Name Server"):
				s := strings.Split(line, ": ")
				d := strings.Join(s[1:], "")
				wr.NameServers = append(wr.NameServers, d)
				c.orchestrator.Publish(core.NewNameServer, d)
			default:
				// move along!
			}
		}

		// TODO: make usre this cleans up better... lol
		if err := scanner.Err(); err != nil {
			panic(err)
		}

		c.orchestrator.Publish(core.NewWhois, wr)
	})
}
