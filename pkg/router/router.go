package router

import (
	"github.com/gorilla/mux"
	"github.com/ncatelli/gates/pkg/models"
)

type PathGenerator interface {
	RegisterPath(*mux.Router, chan<- models.MessageInput) error
}

// Generates routes for a given gate.
func New(pg PathGenerator, inbound chan<- models.MessageInput) (*mux.Router, error) {
	m := mux.NewRouter()

	if err := pg.RegisterPath(m, inbound); err != nil {
		return nil, err
	}

	return m, nil
}
