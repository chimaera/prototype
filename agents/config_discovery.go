package agents

import (
	"fmt"
	"strings"

	"github.com/chimaera/prototype/core"
)

type Config struct {
	state        *core.State
	orchestrator *core.Orchestrator
}

func NewConfigChecker() *Config {
	return &Config{
		state: core.NewState(),
	}
}

func (d *Config) ID() string {
	return "web:config"
}

func (d *Config) Register(o *core.Orchestrator) error {
	o.Subscribe(core.NewPortTCP, d.onNewWebHost)
	d.orchestrator = o

	return nil
}

func TestGitRepo(url string) (urlRequest string, statusFound bool) {
	RequestURL := fmt.Sprintf("%s%s", url, ".git/config")
	resp, body, _ := core.Get(RequestURL, 10)
	if strings.Contains(body, "[core]") && resp.StatusCode != 404 {
		return RequestURL, true
	}
	return RequestURL, false
}

func (d *Config) onNewWebHost(port int, hostname string) {
	url := fmt.Sprintf("http://%s:%d/", hostname, port)
	d.state.Add(url, d.ID())

	d.orchestrator.RunTask(func() {
		if url, status := TestGitRepo(url); status == true {
			d.orchestrator.Publish(core.NewContent, url)
		}
	})
}
