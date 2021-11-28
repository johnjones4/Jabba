package upstream

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/johnjones4/Jabba/core"
)

type Upstream interface {
	LogEvent(*core.Event) error
}

type UpstreamConcrete struct {
	host string
}

func NewUpstreamConcrete(host string) Upstream {
	return &UpstreamConcrete{host}
}

func (u *UpstreamConcrete) LogEvent(e *core.Event) error {
	jsonBytes, err := json.Marshal(e)
	if err != nil {
		return err
	}

	resp, err := http.Post(fmt.Sprintf("http://%s/api/event", u.host), "application/json", io.NopCloser(bytes.NewBuffer(jsonBytes)))
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		return fmt.Errorf("bad response: %s", string(body))
	}

	return nil
}
