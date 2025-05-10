package domain

type ProjectInvitationStatus string

const (
	ProjectInvitationStatusUnknown  ProjectInvitationStatus = ""
	ProjectInvitationStatusPending  ProjectInvitationStatus = "pending"
	ProjectInvitationStatusAccepted ProjectInvitationStatus = "accepted"
	ProjectInvitationStatusExpired  ProjectInvitationStatus = "expired"
	ProjectInvitationStatusDeclined ProjectInvitationStatus = "declined"
)

func (s ProjectInvitationStatus) String() string {
	return string(s)
}

func (s ProjectInvitationStatus) IsValid() bool {
	return s != ProjectInvitationStatusUnknown
}

func (s ProjectInvitationStatus) IsPending() bool {
	return s == ProjectInvitationStatusPending
}

func (s ProjectInvitationStatus) IsAccepted() bool {
	return s == ProjectInvitationStatusAccepted
}

func (s ProjectInvitationStatus) IsExpired() bool {
	return s == ProjectInvitationStatusExpired
}

func (s ProjectInvitationStatus) IsDeclined() bool {
	return s == ProjectInvitationStatusDeclined
}

func ToProjectInvitationStatus(s string) ProjectInvitationStatus {
	switch s {
	case ProjectInvitationStatusPending.String():
		return ProjectInvitationStatusPending
	case ProjectInvitationStatusAccepted.String():
		return ProjectInvitationStatusAccepted
	case ProjectInvitationStatusExpired.String():
		return ProjectInvitationStatusExpired
	case ProjectInvitationStatusDeclined.String():
		return ProjectInvitationStatusDeclined
	default:
		return ProjectInvitationStatusUnknown
	}
}
