package status

import (
	"log"
	"main/alerter"
	"time"

	"github.com/johnjones4/Jabba/core"
)

const (
	StatusOk         = "ok"
	StatusRecovering = "recovering"
	StatusAbnormal   = "abnormal"
)

type StatusEngine interface {
	Start()
	GetStatusForVendorType(string) (*core.Status, error)
	HandleNewEvent(event core.Event) (*core.Status, error)
	ProcessEventsForVendorType(string) (*core.Status, error)
	GetVendorName(string) string
	GetAlerters() []alerter.AlertSender
}

func GenerateStatus(e StatusEngine, lastEvent core.Event) (core.Status, error) {
	yesterday := time.Now().UTC().Add(-24 * time.Hour)

	status := core.Status{
		LastEvent:       lastEvent,
		EventVendorType: lastEvent.EventVendorType,
		EventVendorName: e.GetVendorName(lastEvent.EventVendorType),
	}

	secondLastStatus, _ := e.GetStatusForVendorType(lastEvent.EventVendorType)

	if secondLastStatus != nil && secondLastStatus.LastEvent.ID == lastEvent.ID && secondLastStatus.LastEvent.Created.After(yesterday) {
		log.Printf("Using previously established status for %s", lastEvent.EventVendorType)
		return *secondLastStatus, nil
	}

	if !lastEvent.IsNormal || lastEvent.Created.Before(time.Now().UTC().Add(OneWeekAgo)) {
		status.Status = StatusAbnormal
	} else if secondLastStatus != nil && lastEvent.IsNormal && !secondLastStatus.LastEvent.IsNormal && secondLastStatus.LastEvent.Created.After(yesterday) {
		status.Status = StatusRecovering
	} else {
		status.Status = StatusOk
	}

	log.Printf("Status is now %s", status.Status)

	if secondLastStatus == nil || secondLastStatus.Status != status.Status {
		log.Println("Sending update on status")
		for _, a := range e.GetAlerters() {
			err := a.SendAlert(status)
			if err != nil {
				return core.Status{}, err
			}
		}
	}

	return status, nil
}
