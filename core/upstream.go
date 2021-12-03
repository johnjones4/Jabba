package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Upstream interface {
	LogEvent(*Event) error
}

type UpstreamConcrete struct {
	host string
}

func NewUpstreamConcrete(host string) Upstream {
	return &UpstreamConcrete{host}
}

func (u *UpstreamConcrete) LogEvent(e *Event) error {
	jsonBytes, err := json.Marshal(e)
	if err != nil {
		return err
	}

	resp, err := http.Post(fmt.Sprintf("http://%s/api/event", u.host), "application/json", io.NopCloser(bytes.NewBuffer(jsonBytes)))
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("bad response: %s", string(body))
	}

	err = json.Unmarshal(body, e)
	if err != nil {
		return err
	}

	return nil
}
