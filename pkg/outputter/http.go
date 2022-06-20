package outputter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/ncatelli/gates/pkg/models"
)

type HttpOutputter struct {
	Endpoints []url.URL
}

func (ho *HttpOutputter) Output(tick uint, state models.IO) error {
	reqBody := models.ServicePostBody{
		Tick:  tick,
		State: bool(state),
	}

	buf, err := json.Marshal(&reqBody)
	if err != nil {
		return err
	}

	for _, u := range ho.Endpoints {
		resp, err := http.Post(u.String(), "application/json", bytes.NewReader(buf))
		if resp.StatusCode != 202 {
			return fmt.Errorf("invalid status code from outputter: expected 202 got %d", resp.StatusCode)
		} else if err != nil {
			return err
		}
	}

	return nil
}
