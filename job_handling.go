package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
)

var jobDefinitions []JobDefinition

func preloadJobDefinitions() error {
	data, err := ioutil.ReadFile(os.Getenv("JOB_DEFINITIONS_FILE"))
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &jobDefinitions)
	if err != nil {
		return err
	}
	return nil
}

func getJobDefinition(jobRun JobRun) *JobDefinition {
	for _, jd := range jobDefinitions {
		if jobRun.Job == jd.Name {
			return &jd
		}
	}
	return nil
}

func generateAlerts(jobRun *JobRun) error {
	jd := getJobDefinition(*jobRun)
	if jd == nil {
		return fmt.Errorf("bad job: %s", jobRun.Job)
	}
	alerts := make([]Alert, 0)
	lines := strings.Split(jobRun.Log, "\n")
	for lineNo, line := range lines {
		for _, regex := range jd.Regexes {
			if regex.Match([]byte(line)) {
				alerts = append(alerts, Alert{
					Line:        lineNo,
					Rule:        regex.String(),
					Description: line,
				})
			}
		}
	}
	jobRun.Alerts = alerts
	return nil
}

func transmitStatus(jobRun JobRun) error {
	message := fmt.Sprintf("%s (#%d) executed %s. (%d alerts, %d bytes)", jobRun.Job, jobRun.Id, jobRun.Created.String(), len(jobRun.Alerts), len(jobRun.Log))

	accountSid := os.Getenv("TWILIO_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")

	urlStr := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", accountSid)

	msgData := url.Values{}
	msgData.Set("To", os.Getenv("TWILIO_NUMBER_TO"))
	msgData.Set("From", os.Getenv("TWILIO_NUMBER_FROM"))
	msgData.Set("Body", message)
	msgDataReader := *strings.NewReader(msgData.Encode())

	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlStr, &msgDataReader)
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	_, err := client.Do(req)
	return err
}

func jobLogPath(jobRun JobRun) string {
	return path.Join(os.Getenv("JOB_RUN_LOG_PATH"), fmt.Sprintf("%s_%d_%d.log", jobRun.Job, jobRun.Id, jobRun.Created.Unix()))
}

func loadJobRunLog(jobRun *JobRun) error {
	data, err := ioutil.ReadFile(jobLogPath(*jobRun))
	if err != nil {
		return err
	}
	jobRun.Log = string(data)
	return nil
}

func saveJobRunLog(jobRun JobRun) error {
	return ioutil.WriteFile(jobLogPath(jobRun), []byte(jobRun.Log), 0660)
}

func (r *Regexp) UnmarshalText(b []byte) error {
	regex, err := regexp.Compile(string(b))
	if err != nil {
		return err
	}

	r.Regexp = regex

	return nil
}

func (r *Regexp) MarshalText() ([]byte, error) {
	if r.Regexp != nil {
		return []byte(r.Regexp.String()), nil
	}

	return nil, nil
}
