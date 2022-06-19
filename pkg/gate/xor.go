package gate

import (
	"fmt"

	"github.com/ncatelli/gates/pkg/models"
)

type Xor struct{}

func (xor *Xor) Inputs() uint {
	return 2
}

func (xor *Xor) Compute(tick uint, inputs []models.IO) (models.IO, error) {
	if len(inputs) != int(xor.Inputs()) {
		return false, fmt.Errorf("input does not match expected length of %d", xor.Inputs())
	}

	a := inputs[0]
	b := inputs[1]
	output := models.IO(a != b)

	return output, nil
}
