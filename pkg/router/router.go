package router

import (
	"github.com/gorilla/mux"
	"github.com/ncatelli/gates/pkg/gate"
)

// Generates routes for a given gate.
func New(g gate.Gate) (*mux.Router, error) {
	m := mux.NewRouter()

	/*
		inputs := g.Inputs()
			for i := uint(0); i < inputs; i++ {
				p, err := gate.OffsetToRune(i)
				if err != nil {
					return nil, err
				}
				path := fmt.Sprintf("/input/%s", p)

				route := m.Handle(path, r).Methods(http.MethodPost)
			}
	*/

	return m, nil
}
