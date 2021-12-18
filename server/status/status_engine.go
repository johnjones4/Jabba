package status

import (
	"time"

	"github.com/johnjones4/Jabba/core"
)

const (
	StatusOk         = "ok"
	StatusRecovering = "recovering"
	StatusAbnormal   = "abnormal"
)

type Status struct {
	EventVendorType string     `json:"eventVendorType"`
	EventVendorName string     `json:"eventVendorName"`
	Status          string     `json:"status"`
	LastEvent       core.Event `json:"lastEvent"`
}

type StatusEngine interface {
	Start()
	GetStatusForVendorType(string) (*Status, error)
	HandleNewEvent(event core.Event) (*Status, error)
	ProcessEventsForVendorType(string) (*Status, error)
	GetVendorName(string) string
}

func GenerateStatus(e StatusEngine, lastEvent core.Event) (Status, error) {
	status := Status{
		LastEvent:       lastEvent,
		EventVendorType: lastEvent.EventVendorType,
		EventVendorName: e.GetVendorName(lastEvent.EventVendorType),
	}

	secondLastStatus, _ := e.GetStatusForVendorType(lastEvent.EventVendorType)

	if secondLastStatus != nil && secondLastStatus.LastEvent.ID == lastEvent.ID {
		return *secondLastStatus, nil
	}

	if !lastEvent.IsNormal || lastEvent.Created.Before(time.Now().UTC().Add(OneWeekAgo)) {
		status.Status = StatusAbnormal
	} else if secondLastStatus != nil && secondLastStatus.LastEvent.IsNormal && secondLastStatus.LastEvent.Created.After(time.Now().UTC().Add(-24*time.Hour)) {
		status.Status = StatusRecovering
	} else {
		status.Status = StatusOk
	}

	return status, nil
}
