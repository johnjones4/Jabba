package alerter

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/smtp"
	"os"
	"strings"

	"github.com/johnjones4/Jabba/core"
)

type EmailAlerter struct {
	Recipient string `json:"recipient"`
	Sender    string `json:"sender"`
	Host      string `json:"host"`
	Port      string `json:"port"`
	Username  string `json:"username"`
	Password  string `json:"password"`
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
	host := fmt.Sprintf("%s:%s", a.Host, a.Port)

	conn, err := net.Dial("tcp", host)
	if err != nil {
		return err
	}

	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         a.Host,
	}

	client, err := smtp.NewClient(conn, a.Host)
	if err != nil {
		return err
	}

	err = client.Hello(a.Host)
	if err != nil {
		return err
	}

	err = client.StartTLS(tlsconfig)
	if err != nil {
		return err
	}

	auth := smtp.PlainAuth("", a.Username, a.Password, a.Host)

	err = client.Auth(auth)
	if err != nil {
		return err
	}

	err = client.Mail(a.Sender)
	if err != nil {
		return err
	}

	err = client.Rcpt(a.Recipient)
	if err != nil {
		return err
	}

	w, err := client.Data()
	if err != nil {
		return err
	}

	msg, err := a.formatEmailMessage(s)
	if err != nil {
		return err
	}

	_, err = w.Write(msg)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	client.Quit()

	return nil
}

func (a EmailAlerter) formatEmailMessage(s core.Status) ([]byte, error) {

	str := new(strings.Builder)
	str.WriteString(fmt.Sprintf("From: Jabba<%s>\r\n", a.Sender))
	str.WriteString(fmt.Sprintf("To: %s\r\n", a.Recipient))
	str.WriteString(fmt.Sprintf("Subject: %s has entered status \"%s\" at %s\r\n\r\n", s.EventVendorName, s.Status, s.LastEvent.Created.Local().String()))

	if valMap, ok := s.LastEvent.VendorInfo.(map[string]interface{}); ok {
		if log, ok := valMap["log"]; ok {
			str.WriteString(fmt.Sprint(log))
			str.WriteString("\r\n")
		} else if body, ok := valMap["body"]; ok {
			if statusCode, ok := valMap["statusCode"]; ok {
				str.WriteString(fmt.Sprint(statusCode))
				str.WriteString("\r\n")
			}
			str.WriteString(fmt.Sprint(body))
			str.WriteString("\r\n")
		} else {
			jsonBytes, err := json.MarshalIndent(s.LastEvent.VendorInfo, "", "  ")
			if err != nil {
				return nil, err
			}
			str.WriteString(string(jsonBytes))
		}
	}

	str.WriteString("\r\n")
	return []byte(str.String()), nil
}
