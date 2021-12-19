package alerter

import "github.com/johnjones4/Jabba/core"

type AlertSender interface {
	SendAlert(s core.Status) error
}
