package gate

import (
	"fmt"

	"github.com/ncatelli/gates/pkg/models"
)

type And struct{}

func (and *And) Inputs() uint {
	return 2
}

func (and *And) Compute(tick uint, inputs []models.IO) (models.IO, error) {
	if len(inputs) != int(and.Inputs()) {
		return false, fmt.Errorf("input does not match expected length of %d", and.Inputs())
	}

	a := inputs[0]
	b := inputs[1]
	output := a && b

	return output, nil
}
