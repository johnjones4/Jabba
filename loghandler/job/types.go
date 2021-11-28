package job

import (
	"regexp"

	"github.com/johnjones4/Jabba/core"
)

type Regexp struct {
	*regexp.Regexp
}

type JobDefinition struct {
	Name    string   `json:"name"`
	Regexes []Regexp `json:"regexes"`
}

type AlertGenerator interface {
	GenerateAlerts(e *core.Event) error
}

type AlertGeneratorConcrete struct {
	jobDefinitions []JobDefinition
}

type LogAlert struct {
	Line        int    `json:"line"`
	Description string `json:"description"`
	Rule        string `json:"rule"`
}
