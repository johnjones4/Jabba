package pollers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/google/uuid"
	jabbacore "github.com/johnjones4/Jabba/core"
)

type AbodePoller struct {
	username    string
	password    string
	accessToken string
	tokenType   string
	expiration  time.Time
	lastEvent   abodeTimelineEvent
}

type abodeTimelineEvent struct {
	Id         string `json:"id"`
	EventUTC   string `json:"event_utc"`
	Device     string `json:"device_name"`
	Event      string `json:"event_name"`
	Severity   string `json:"severity"`
	EventCode  string `json:"event_code"`
	DeviceID   string `json:"device_id"`
	LabelCode  string `json:"label_code"`
	IsAlarm    string `json:"is_alarm"`
	HasActions string `json:"has_actions"`
	HasFaults  string `json:"hasFaults"`
}

type abodeLoginResponse struct {
	Token string `json:"token"`
}

type abodeClaimResponse struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

func NewAbodePoller(username string, password string) *AbodePoller {
	return &AbodePoller{
		username: username,
		password: password,
	}
}

func (p *AbodePoller) authorize() error {
	log.Println("Authenticating with Abode")

	form := make(url.Values)
	form.Add("id", p.username)
	form.Add("password", p.password)
	form.Add("uuid", uuid.NewString())
	form.Add("locale_code", "en-US")

	loginHttpResp, err := http.PostForm("https://my.goabode.com/api/auth2/login", form)
	if err != nil {
		return err
	}

	loginRespBody, err := io.ReadAll(loginHttpResp.Body)
	if err != nil {
		return err
	}

	loginResp := abodeLoginResponse{}

	err = json.Unmarshal(loginRespBody, &loginResp)
	if err != nil {
		return err
	}

	log.Println("Received login token from Abode")

	claimReq, err := http.NewRequest("GET", "https://my.goabode.com/api/auth2/claims", nil)
	if err != nil {
		return err
	}

	claimReq.Header.Add("ABODE-API-KEY", loginResp.Token)

	claimHttpResp, err := http.DefaultClient.Do(claimReq)
	if err != nil {
		return err
	}

	claimRespBody, err := io.ReadAll(claimHttpResp.Body)
	if err != nil {
		return err
	}

	claimResp := abodeClaimResponse{}
	err = json.Unmarshal(claimRespBody, &claimResp)
	if err != nil {
		return err
	}

	p.accessToken = claimResp.AccessToken
	p.tokenType = claimResp.TokenType
	p.expiration = time.Now().UTC().Add(time.Second * time.Duration(claimResp.ExpiresIn/2))

	log.Printf("Received access token from Abode. Expires on %s\n", p.expiration.String())

	return nil
}

func (p *AbodePoller) NeedsAuthorization() bool {
	return p.accessToken == "" || p.expiration.Before(time.Now().UTC())
}

func (p *AbodePoller) runALoop(u jabbacore.Upstream) error {
	allEvents, err := p.getEvents()
	if err != nil {
		return err
	}

	if p.lastEvent.EventUTC != "" {
		tstamp, err := strconv.Atoi(p.lastEvent.EventUTC)
		if err != nil {
			return err
		}
		events := make([]abodeTimelineEvent, 0)
		for _, e := range allEvents {
			tstamp1, err := strconv.Atoi(e.EventUTC)
			if err != nil {
				return err
			}

			if tstamp1 > tstamp {
				events = append(events, e)
			}
		}

		p.emitEvents(u, events)
	} else {
		p.emitEvents(u, allEvents)
	}

	if len(allEvents) > 0 {
		p.lastEvent = allEvents[0]
	}

	return nil
}

func (p *AbodePoller) getEvents() ([]abodeTimelineEvent, error) {
	log.Println("Getting recent Abode events")

	if p.NeedsAuthorization() {
		err := p.authorize()
		if err != nil {
			return nil, err
		}
	}

	params := make(url.Values)
	params.Add("size", "1000")
	params.Add("dir", "next")

	req, err := http.NewRequest("GET", "https://my.goabode.com/api/v1/timeline?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("%s %s", p.tokenType, p.accessToken))

	httpResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	httpBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	var allEvents []abodeTimelineEvent
	err = json.Unmarshal(httpBody, &allEvents)
	if err != nil {
		log.Println(string(httpBody))
		return nil, err
	}

	log.Printf("Got back %d Abode events.\n", len(allEvents))

	return allEvents, nil
}

func (p *AbodePoller) emitEvents(u jabbacore.Upstream, events []abodeTimelineEvent) error {
	for _, e := range events {
		tstamp, err := strconv.Atoi(e.EventUTC)
		if err != nil {
			return err
		}

		jEvent := jabbacore.Event{
			EventVendorType: "abode",
			EventVendorID:   e.Id,
			VendorInfo: map[string]interface{}{
				"device":     e.Device,
				"event":      e.Event,
				"severity":   e.Severity,
				"deviceId":   e.DeviceID,
				"labelCode":  e.LabelCode,
				"isAlarm":    e.IsAlarm != "0",
				"hasActions": e.HasActions != "0",
				"hasFaults":  e.HasFaults != "0",
			},
			Created:  time.Unix(int64(tstamp), 0),
			IsNormal: true,
		}
		err = u.LogEvent(&jEvent)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *AbodePoller) Setup() error {
	return p.authorize()
}

func (p *AbodePoller) Poll(u jabbacore.Upstream) {
	for {
		err := p.runALoop(u)
		if err != nil {
			log.Println(err)
		}

		time.Sleep(time.Minute)
	}
}
