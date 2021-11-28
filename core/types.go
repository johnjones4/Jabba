package core

import (
	"time"
)

type Alert struct {
	Info interface{} `json:"info"`
	Type string      `json:"type"`
}

type Event struct {
	ID              int         `json:"id"`
	EventVendorType string      `json:"eventVendorType"`
	EventVendorID   string      `json:"eventVendorID"`
	Created         time.Time   `json:"created"`
	VendorInfo      interface{} `json:"vendorInfo"`
	Alerts          []Alert     `json:"alerts"`
	IsNormal        bool        `json:"isNormal"`
}
