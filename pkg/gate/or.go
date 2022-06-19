package gate

import (
	"fmt"

	"github.com/ncatelli/gates/pkg/models"
)

type Or struct{}

func (or *Or) Inputs() uint {
	return 2
}

func (or *Or) Compute(tick uint, inputs []models.IO) (models.IO, error) {
	if len(inputs) != int(or.Inputs()) {
		return false, fmt.Errorf("input does not match expected length of %d", or.Inputs())
	}

	a := inputs[0]
	b := inputs[1]
	output := a || b

	return output, nil
}
