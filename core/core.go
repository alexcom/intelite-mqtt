package core

// ReturnAddress contains info required to respond to message
type ReturnAddress struct {
	RespondTo string `json:"respond_to"`
}

type LampMode string

const (
	LampModeOn    LampMode = "on"
	LampModeOff   LampMode = "off"
	LampModeNight LampMode = "night"
	LampModeMax   LampMode = "max"
)

// LampUpdate comes from IncomingTransport with lamp parameters
type LampUpdate struct {
	ReturnAddress
	Mode       LampMode `json:"mode"`
	Color      int      `json:"color"`
	Brightness int      `json:"brightness"`
}
