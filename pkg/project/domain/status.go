package domain

type Status string

const (
	StatusUnknown  Status = ""
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
)

func (s Status) String() string {
	return string(s)
}

func (s Status) IsValid() bool {
	return s != StatusUnknown
}

func (s Status) IsActive() bool {
	return s == StatusActive
}

func (s Status) IsInactive() bool {
	return s == StatusInactive
}

func (s Status) IsEqual(status string) bool {
	dStatus := ToStatus(status)
	return s == dStatus
}

func ToStatus(s string) Status {
	switch s {
	case StatusActive.String():
		return StatusActive
	case StatusInactive.String():
		return StatusInactive
	default:
		return StatusUnknown
	}
}
