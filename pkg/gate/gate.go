package gate

import (
	"encoding/json"
	"fmt"
	"net/http"
	"unicode"

	"github.com/gorilla/mux"
	"github.com/ncatelli/gates/pkg/models"
	"github.com/ncatelli/gates/pkg/outputter"
)

type Gate interface {
	Inputs() uint
	Compute(tick uint, input []models.IO) (models.IO, error)
}

type PendingInput struct {
	received bool
	state    models.IO
}

type TickState struct {
	inputs []PendingInput
}

func NewTickState(expectedInputs uint) *TickState {
	inputs := make([]PendingInput, expectedInputs)

	for i := 0; i < int(expectedInputs); i++ {
		inputs[i] = PendingInput{
			received: false,
			state:    false,
		}
	}

	return &TickState{
		inputs: inputs,
	}
}

// AllInputsReceived returns true if every pending input has been marked received.
func (ts *TickState) AllInputsReceived() bool {
	// check if all inputs have been received, if not return early.
	for _, input := range ts.inputs {
		if !input.received {
			return false
		}
	}

	return true
}

// ReturnInputsIfReady returns a slice of all IO state in order if called after
// all have been received. Otherwise an error is returned.
func (ts *TickState) ReturnInputsIfReady() ([]models.IO, error) {
	if !ts.AllInputsReceived() {
		return nil, fmt.Errorf("input is still pending")
	}

	inputs := make([]models.IO, 0, len(ts.inputs))
	for _, pending := range ts.inputs {
		inputs = append(inputs, pending.state)
	}

	return inputs, nil
}

type GateService struct {
	ticks          map[uint]*TickState
	expectedInputs uint
	gate           Gate
	op             outputter.Outputter
}

func NewGenericGate(g Gate, op outputter.Outputter) *GateService {
	return &GateService{
		ticks:          make(map[uint]*TickState),
		expectedInputs: g.Inputs(),
		gate:           g,
		op:             op,
	}
}

func (gs *GateService) Inputs() uint {
	return gs.expectedInputs
}

func (gs *GateService) Compute(tick uint, inputs []models.IO) (models.IO, error) {
	if len(inputs) != int(gs.Inputs()) {
		return false, fmt.Errorf("input does not match expected length of %d", gs.Inputs())
	}

	return gs.gate.Compute(tick, inputs)

}

func (gs *GateService) ReceiveInput(tick uint, input rune, state models.IO) (*TickState, error) {
	// calculate the offset of the input from the rune (path)
	inputOffset, err := runeToNormalizedOffset(input)
	if err != nil {
		return nil, err
	}

	// verify that the input is within range
	if (inputOffset + 1) > gs.expectedInputs {
		return nil, fmt.Errorf("offset %c(%d) greater than max %d", input, inputOffset, gs.expectedInputs)
	}

	ts, prs := gs.ticks[tick]
	if !prs {
		ts = NewTickState(gs.expectedInputs)
		gs.ticks[tick] = ts
	}

	// error if inputs are clobbered.
	if ts.inputs[inputOffset].received {
		return nil, fmt.Errorf("input (%c) for tick %d already set", input, tick)
	}

	ts.inputs[inputOffset] = PendingInput{
		received: true,
		state:    state,
	}

	return ts, nil
}

func (gs *GateService) RegisterPath(r *mux.Router, outbound chan<- models.MessageInput) error {
	inputs := gs.Inputs()
	for i := uint(0); i < inputs; i++ {
		p, err := OffsetToRune(i)
		if err != nil {
			return err
		}
		path := fmt.Sprintf("/input/%c", p)

		r.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			var inputId rune = p
			in := models.ServicePostBody{}
			op := gs.op

			dec := json.NewDecoder(r.Body)
			err := dec.Decode(&in)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			tick := in.Tick

			resp := make(chan models.GateResponse)
			msg := models.MessageInput{
				Resp:  resp,
				Tick:  tick,
				Path:  inputId,
				Input: models.IO(in.State),
			}

			// send the change to the gate service
			outbound <- msg
			// wait for the response from the gate service
			gateResp := <-resp
			if gateResp.Err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if gateResp.OutputReady {
				err = op.Output(tick, gateResp.Output)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			}

			w.WriteHeader(http.StatusAccepted)
		}).Methods(http.MethodPost)
	}

	return nil
}

func runeToNormalizedOffset(r rune) (uint, error) {
	min := uint('a')
	max := uint('z')
	normalized := unicode.ToLower(r)
	runeAsInt := uint(normalized)

	if runeAsInt <= max || runeAsInt >= min {
		offset := runeAsInt - min

		return offset, nil
	}

	return 0, fmt.Errorf("value out of range: must be between a-z, got %c", r)
}

func OffsetToRune(offset uint) (rune, error) {
	var min uint = 0
	max := uint('z') - uint('a')

	if offset <= max || offset >= min {
		r := offset + uint('a')

		return rune(r), nil
	}

	return 0, fmt.Errorf("value out of range: must be between %d-%d, got %d", min, max, offset)
}
