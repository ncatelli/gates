package gate

import (
	"fmt"

	"github.com/ncatelli/gates/pkg/models"
)

type Nand struct{}

func (nand *Nand) Inputs() uint {
	return 2
}

func (nand *Nand) Compute(tick uint, inputs []models.IO) (models.IO, error) {
	if len(inputs) != int(nand.Inputs()) {
		return false, fmt.Errorf("input does not match expected length of %d", nand.Inputs())
	}

	a := inputs[0]
	b := inputs[1]
	output := !(a && b)

	return output, nil
}
