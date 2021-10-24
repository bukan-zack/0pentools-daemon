package domain

type Domain struct {
	Name  string `json:"name"`
	ID    int32 `json:"id"`
	UUID  string `json:"uuid"`
	State string `json:"state"`
}

type DomainXML struct {
	UUID string
	Memory  uint64
	Port uint64
}

const (
	StateNone = iota
	StateRunning
	StatePaused
	StateBlocked
	StateCrashed
	StateSuspended
	StateShutdown
	StateKilled
	StateLast
)

func ConvDomainState(state uint8) string {
	switch state {
	case StateNone:
		return "None"
	case StateRunning:
		return "Running"
	case StatePaused:
		return "Paused"
	case StateBlocked:
		return "Blocked"
	case StateCrashed:
		return "Crashed"
	case StateSuspended:
		return "Suspended"
	case StateShutdown:
		return "Shutdown"
	case StateKilled:
		return "Killed"
	case StateLast:
		return "Unknown"
	default:
		return ""
	}
}