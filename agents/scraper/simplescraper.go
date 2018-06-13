package scraper

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/bobesa/go-domain-util/domainutil"
	"github.com/chimaera/prototype/core"
)

type SimpleScraper struct {
	state        *core.State
	orchestrator *core.Orchestrator
}

func NewSimpleScraper() *SimpleScraper {
	return &SimpleScraper{
		state: core.NewState(),
	}
}

func (d *SimpleScraper) ID() string {
	return "dns:simplescraper"
}

func (d *SimpleScraper) Register(o *core.Orchestrator) error {
	o.Subscribe(core.NewHostname, d.onNewHostname)
	d.orchestrator = o
	return nil
}

func (d *SimpleScraper) onNewHostname(hostname string) {
	domainName := domainutil.Domain(hostname)
	if d.state.DidProcess(domainName, d.ID()) {
		return
	}

	d.state.Add(domainName, d.ID())

	d.orchestrator.RunTask(func() {
		searchQuery := fmt.Sprintf("site:%[1]s -site:www.%[1]s", domainName)
		re := regexp.MustCompile(`([a-z0-9]+\.)+` + domainName)
		for currentPage := 0; currentPage < 10; currentPage++ {
			_, responseBody, err := core.Get("https://www.baidu.com/s?rn=100&pn="+strconv.Itoa(currentPage)+"&wd=site:"+searchQuery+"&oq="+searchQuery, 120)
			if err != nil {
				break
			}

			matches := re.FindAllString(string(responseBody), -1)

			for _, subdomain := range matches {
				d.orchestrator.Publish(core.NewSubdomain, subdomain)
			}
		}

	})
}
