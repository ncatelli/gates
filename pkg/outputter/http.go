package outputter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/ncatelli/gates/pkg/gate"
)

type Outputter interface {
	Output(uint, gate.IO) error
}

type HttpOutputter struct {
	endpoint *url.URL
}

func (ho *HttpOutputter) Output(tick uint, state gate.IO) error {
	reqBody := gate.ServicePostBody{
		Tick:  tick,
		State: bool(state),
	}

	buf, err := json.Marshal(&reqBody)
	if err != nil {
		return err
	}

	resp, err := http.Post(ho.endpoint.String(), "application/json", bytes.NewReader(buf))
	if resp.StatusCode != 202 {
		return fmt.Errorf("invalid status code from outputter: expected 202 got %d", resp.StatusCode)
	} else if err != nil {
		return err
	}

	return nil
}
