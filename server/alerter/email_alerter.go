package alerter

import (
	"encoding/json"
	"fmt"
	"net/smtp"
	"os"

	"github.com/johnjones4/Jabba/core"
)

type EmailAlerter struct {
	Recipients []string `json:"recipients"`
	Sender     string   `json:"sender"`
	Host       string   `json:"host"`
	Port       string   `json:"port"`
	Username   string   `json:"username"`
	Password   string   `json:"password"`
}

func LoadEmailAlerters(path string) ([]AlertSender, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var emailAlerters []EmailAlerter
	err = json.Unmarshal(bytes, &emailAlerters)
	if err != nil {
		return nil, err
	}

	alerters := make([]AlertSender, len(emailAlerters))
	for i, a := range emailAlerters {
		alerters[i] = a
	}

	return alerters, nil
}

func (a EmailAlerter) SendAlert(s core.Status) error {
	auth := smtp.PlainAuth("", a.Username, a.Password, a.Host)
	msg, err := formatEmailMessage(s)
	if err != nil {
		return err
	}
	return smtp.SendMail(fmt.Sprintf("%s:%s", a.Host, a.Port), auth, a.Sender, a.Recipients, msg)
}

func formatEmailMessage(s core.Status) ([]byte, error) {
	jsonBytes, err := json.MarshalIndent(s.LastEvent.VendorInfo, "", "  ")
	if err != nil {
		return nil, err
	}
	str := fmt.Sprintf("%s has entered status \"%s\" at %s:\n\n%s", s.EventVendorName, s.Status, s.LastEvent.Created.Local().String(), string(jsonBytes))
	return []byte(str), nil
}
