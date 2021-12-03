package pollers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	jabbacore "github.com/johnjones4/Jabba/core"
)

type Site struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type SitesPoller struct {
	sourceFile string
	status     map[Site]bool
}

func NewSitesPoller(sourceFile string) *SitesPoller {
	return &SitesPoller{
		status:     make(map[Site]bool),
		sourceFile: sourceFile,
	}
}

func (p *SitesPoller) checkAndLogSite(site Site, u jabbacore.Upstream) error {
	res, err := http.Get(site.URL)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	status := res.StatusCode == http.StatusOK

	if status == p.status[site] {
		return nil
	}

	p.status[site] = status

	jEvent := jabbacore.Event{
		EventVendorType: site.ID,
		EventVendorID:   uuid.NewString(),
		VendorInfo: map[string]interface{}{
			"body":       string(body),
			"statusCode": res.StatusCode,
		},
		Created:  time.Now().UTC(),
		IsNormal: status,
		Alerts:   []jabbacore.Alert{},
	}
	err = u.LogEvent(&jEvent)
	if err != nil {
		return err
	}

	return nil
}

func (p *SitesPoller) runALoop(u jabbacore.Upstream) {
	for site := range p.status {
		err := p.checkAndLogSite(site, u)
		if err != nil {
			log.Println(err)
		}
	}
}

func (p *SitesPoller) Setup() error {
	contents, err := os.ReadFile(p.sourceFile)
	if err != nil {
		return err
	}

	var sites []Site
	err = json.Unmarshal(contents, &sites)
	if err != nil {
		return err
	}

	for _, site := range sites {
		p.status[site] = false
	}

	return nil
}

func (p *SitesPoller) Poll(u jabbacore.Upstream) {
	for {
		p.runALoop(u)
		time.Sleep(time.Minute)
	}
}
