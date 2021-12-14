package shared

import (
	"errors"
	"main/store"
	"time"

	"github.com/johnjones4/Jabba/core"
)

const (
	StatusOk         = "ok"
	StatusRecovering = "recovering"
	StatusAbnormal   = "abnormal"
)

func GetEventVendorTypeInfo(s store.Store, eventVendorType string) (core.Event, string, error) {
	events, err := s.GetEventsForVendorType(eventVendorType, 2, 0)
	if err != nil {
		return core.Event{}, "", err
	}

	if len(events) == 0 {
		return core.Event{}, "", errors.New("not enough context for status")
	}

	lastEvent := events[0]

	oneWeekAgo := time.Hour * 24 * 7 * -1

	if !lastEvent.IsNormal || lastEvent.Created.Before(time.Now().UTC().Add(oneWeekAgo)) {
		return lastEvent, StatusAbnormal, nil
	}

	if len(events) > 1 && !events[1].IsNormal {
		return lastEvent, StatusRecovering, nil
	}

	return lastEvent, StatusOk, nil
}
