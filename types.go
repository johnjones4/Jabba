package main

import (
	"regexp"
	"time"
)

type Regexp struct {
	*regexp.Regexp
}

type JobDefinition struct {
	Name    string   `json:"name"`
	Regexes []Regexp `json:"regexes"`
}

type Alert struct {
	Id          int    `json:"id"`
	Line        int    `json:"line"`
	Rule        string `json:"rule"`
	Description string `json:"description"`
}

type JobRun struct {
	Id      int       `json:"id"`
	Created time.Time `json:"created"`
	Job     string    `json:"job"`
	Log     string    `json:"log"`
	Alerts  []Alert   `json:"alerts"`
}

type Message struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

type HTTPReqInfo struct {
	method    string
	uri       string
	referer   string
	ipaddr    string
	code      int
	size      int64
	duration  time.Duration
	userAgent string
}

type ListPageData struct {
	Page         int
	NextPage     int
	PreviousPage int
	JobRuns      []JobRun
}

type DetailPageData struct {
	JobRun JobRun
}

type ErrorPageData struct {
	Error string
}
