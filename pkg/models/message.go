package models

type IO bool

type GateResponse struct {
	OutputReady bool
	Output      IO
	Err         error
}

type MessageInput struct {
	Resp  chan GateResponse
	Tick  uint
	Path  rune
	Input IO
}

type ServicePostBody struct {
	State bool `json:"state"`
	Tick  uint `json:"tick"`
}

type SignalEvent struct {
	State bool `json:"state"`
	Tick  uint `json:"tick"`
}
