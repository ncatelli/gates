package outputter

import (
	"log"

	"github.com/ncatelli/gates/pkg/models"
)

type StdOutOutputter struct{}

func (so *StdOutOutputter) Output(tick uint, state models.IO) error {
	log.Printf("tick: %d, state: %v", tick, state)
	return nil
}
