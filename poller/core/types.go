package core

import (
	"github.com/johnjones4/Jabba/core"
)

type Poller interface {
	Setup() error
	Poll(core.Upstream)
}
