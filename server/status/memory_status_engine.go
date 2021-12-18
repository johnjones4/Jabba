package status

import (
	"errors"
	"fmt"
	"log"
	"main/store"
	"time"

	"github.com/johnjones4/Jabba/core"
)

type MemoryStatusEngine struct {
	statuses    map[string]Status
	vendorNames map[string]string
	eventsStore store.Store
}

const OneWeekAgo = time.Hour * 24 * 7 * -1

func NewMemoryStatusEngine(vendorNames map[string]string, eventsStore store.Store) *MemoryStatusEngine {
	return &MemoryStatusEngine{
		statuses:    make(map[string]Status),
		vendorNames: vendorNames,
		eventsStore: eventsStore,
	}
}

func (e *MemoryStatusEngine) Start() {
	for {
		log.Println("Updating statuses ...")

		types, err := e.eventsStore.GetEventVendorTypes()
		if err != nil {
			log.Println(err)
			continue
		}

		for _, t := range types {
			log.Printf("Updating status for %s ... ", t)
			s, err := e.ProcessEventsForVendorType(t)
			if err != nil {
				log.Println(err)
			} else {
				log.Println(s.Status)
			}
		}

		time.Sleep(time.Hour)
	}
}

func (e *MemoryStatusEngine) ProcessEventsForVendorType(eventVendorType string) (*Status, error) {
	events, err := e.eventsStore.GetEventsForVendorType(eventVendorType, 1, 0)
	if err != nil {
		return nil, err
	}

	if len(events) == 0 {
		return nil, errors.New("not enough context for status")
	}

	return e.HandleNewEvent(events[0])
}

func (e *MemoryStatusEngine) GetStatusForVendorType(t string) (*Status, error) {
	if status, ok := e.statuses[t]; ok {
		return &status, nil
	}
	return nil, fmt.Errorf("no status for %s", t)
}

func (e *MemoryStatusEngine) HandleNewEvent(lastEvent core.Event) (*Status, error) {
	s, err := GenerateStatus(e, lastEvent)
	if err != nil {
		return nil, err
	}
	e.statuses[lastEvent.EventVendorType] = s
	return &s, nil
}

func (e *MemoryStatusEngine) GetVendorName(t string) string {
	if name, ok := e.vendorNames[t]; ok {
		return name
	}
	return t
}
