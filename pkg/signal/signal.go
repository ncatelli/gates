package signal

import (
	"github.com/ncatelli/gates/pkg/models"
	"github.com/ncatelli/gates/pkg/outputter"
)

type Signal interface {
	Emitter() <-chan models.SignalEvent
}

type SignalService struct {
	signal Signal
	op     outputter.Outputter
}

func NewSignalService(signal Signal, op outputter.Outputter) *SignalService {
	return &SignalService{
		signal: signal,
		op:     op,
	}
}
