package domain

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

func GetDomainState(state uint8) string {
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
