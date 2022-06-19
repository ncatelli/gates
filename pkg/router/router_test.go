package router

import (
	"testing"

	"github.com/gorilla/mux"
	"github.com/ncatelli/gates/pkg/models"
)

type mockPathGenerator struct{}

func (mpg *mockPathGenerator) RegisterPath(mux *mux.Router, msg chan<- models.MessageInput) error {
	return nil
}

func TestRouterShouldInstantiateMuxWhenPassedValidPathGenerator(t *testing.T) {
	ibc := make(chan models.MessageInput)
	mux, err := New(&mockPathGenerator{}, ibc)
	if mux == nil {
		t.Fatal("mux should be non-nil on valid input.")
	} else if err != nil {
		t.Fatal(err)
	}
}
