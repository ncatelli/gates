package outputter

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestOutputterRequest(t *testing.T) {
	t.Run("should return nil on valid response", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusAccepted)
		}))
		defer ts.Close()

		endpointUrl := ts.URL
		u, err := url.Parse(endpointUrl)
		if err != nil {
			t.Fatal(err)
		}

		op := HttpOutputter{Endpoints: []url.URL{*u}}
		if err = op.Output(0, false); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("should err on response error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
		}))
		defer ts.Close()

		endpointUrl := ts.URL
		u, err := url.Parse(endpointUrl)
		if err != nil {
			t.Fatal(err)
		}

		op := HttpOutputter{Endpoints: []url.URL{*u}}
		if err = op.Output(0, false); err == nil {
			t.Fatal(errors.New("outputter should return an error on non-202 response code"))
		}
	})
}
