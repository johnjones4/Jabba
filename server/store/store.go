package store

import "github.com/johnjones4/Jabba/core"

type Store interface {
	SaveEvent(event *core.Event) error
	GetEvents(limit int, offset int) ([]core.Event, error)
	GetEvent(id int) (core.Event, error)
}
