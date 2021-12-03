package core

import (
	"sync"

	"github.com/johnjones4/Jabba/core"
)

type PollWatcher struct {
	contLock sync.Mutex
	stop     bool
}

func NewPollWatcher() PollWatcher {
	return PollWatcher{
		contLock: sync.Mutex{},
		stop:     false,
	}
}

func (pw *PollWatcher) Continue() bool {
	pw.contLock.Lock()
	v := !pw.stop
	pw.contLock.Unlock()
	return v
}

func (pw *PollWatcher) Stop() {
	pw.contLock.Lock()
	pw.stop = true
	pw.contLock.Unlock()
}

type Poller interface {
	Poll(w PollWatcher, e chan error, u core.Upstream)
}
