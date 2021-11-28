package job

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/johnjones4/Jabba/core"
)

func NewAlertGeneratorConcrete(path string) (AlertGenerator, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	generator := AlertGeneratorConcrete{}
	err = json.Unmarshal(data, &generator.jobDefinitions)
	if err != nil {
		return nil, err
	}
	return &generator, nil
}

func (g *AlertGeneratorConcrete) GenerateAlerts(e *core.Event) error {
	jd := g.getJobDefinition(e)
	if jd == nil {
		return fmt.Errorf("bad job: %s", e.EventVendorType)
	}
	alerts := make([]core.Alert, 0)
	lines := strings.Split(e.Log, "\n")
	for lineNo, line := range lines {
		for _, regex := range jd.Regexes {
			if regex.Match([]byte(line)) {
				alerts = append(alerts, core.Alert{
					Type: "log",
					Info: LogAlert{
						Line:        lineNo,
						Rule:        regex.String(),
						Description: line,
					},
				})
			}
		}
	}
	e.Alerts = alerts
	return nil
}

func (g *AlertGeneratorConcrete) getJobDefinition(event *core.Event) *JobDefinition {
	for _, jd := range g.jobDefinitions {
		if event.EventVendorType == jd.Name {
			return &jd
		}
	}
	return nil
}
