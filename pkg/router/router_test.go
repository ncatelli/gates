package router

import (
	"errors"
	"testing"

	"github.com/gorilla/mux"
	"github.com/ncatelli/gates/pkg/models"
)

type mockPathGenerator struct {
	err error
}

func (mpg *mockPathGenerator) RegisterPath(mux *mux.Router, msg chan<- models.MessageInput) error {
	return mpg.err
}

func TestNewShouldInstantiateRouterWhenPassedValidPathGenerator(t *testing.T) {
	ibc := make(chan models.MessageInput)
	mux, err := New(&mockPathGenerator{err: nil}, ibc)
	if mux == nil {
		t.Fatal("mux should be non-nil on valid input.")
	} else if err != nil {
		t.Fatal(err)
	}
}

func TestNewShouldErrorWhenPassedAnInvalidPG(t *testing.T) {
	ibc := make(chan models.MessageInput)
	mux, err := New(&mockPathGenerator{err: errors.New("mock error")}, ibc)
	if mux != nil && err == nil {
		t.Fatal("router should fail to instantiate")
	}
}
