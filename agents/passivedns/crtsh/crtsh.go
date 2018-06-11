package crtsh

import (
	"encoding/json"
	"strings"

	"github.com/bobesa/go-domain-util/domainutil"
	"github.com/chimaera/prototype/core"
)

type Crtsh struct {
	state        *core.State
	orchestrator *core.Orchestrator
}

func NewCrtsh() *Crtsh {
	return &Crtsh{
		state: core.NewState(),
	}
}

func (d *Crtsh) ID() string {
	return "dns:crtsh"
}

func (d *Crtsh) Register(o *core.Orchestrator) error {
	o.Subscribe(core.NewHostname, d.onNewHostname)
	d.orchestrator = o
	return nil
}

type crtsh_object struct {
	Name_value string `json:"name_value"`
}

var crtsh_data []crtsh_object

func (d *Crtsh) onNewHostname(hostname string) {
	domainName := domainutil.Domain(hostname)
	if d.state.DidProcess(domainName, d.ID()) {
		return
	}

	d.state.Add(domainName, d.ID())

	d.orchestrator.RunTask(func() {
		_, resp_body, err := core.Get("https://crt.sh/?q=%25."+domainName+"&output=json", 120)

		if err == nil {
			if strings.Contains(string(resp_body), "The requested URL / was not found on this server.") == false {
				correct_format := strings.Replace(string(resp_body), "}{", "},{", -1)
				json_output := "[" + correct_format + "]"
				_ = json.Unmarshal([]byte(json_output), &crtsh_data)
				for _, subdomain := range crtsh_data {
					if strings.Contains(subdomain.Name_value, "*.") {
						subdomain.Name_value = strings.Split(subdomain.Name_value, "*.")[1]
					}

					d.orchestrator.Publish(core.NewSubdomain, subdomain.Name_value)
				}
			}
		}
	})
}
