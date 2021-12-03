package store

import "github.com/johnjones4/Jabba/core"

type Store interface {
	SaveEvent(event *core.Event) error
	GetEvents(limit int, offset int) ([]core.Event, error)
	GetEvent(id int) (core.Event, error)
	GetEventVendorTypes() ([]string, error)
	GetEventsForVendorType(t string, limit int, offset int) ([]core.Event, error)
}
