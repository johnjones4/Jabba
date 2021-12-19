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
	IsNormal        bool        `json:"isNormal"`
}

type Status struct {
	EventVendorType string `json:"eventVendorType"`
	EventVendorName string `json:"eventVendorName"`
	Status          string `json:"status"`
	LastEvent       Event  `json:"lastEvent"`
}
