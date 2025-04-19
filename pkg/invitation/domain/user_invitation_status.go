package domain

type UserInvitationStatus string

const (
	UserInvitationStatusUnknown  UserInvitationStatus = ""
	UserInvitationStatusPending  UserInvitationStatus = "pending"
	UserInvitationStatusAccepted UserInvitationStatus = "accepted"
	UserInvitationStatusExpired  UserInvitationStatus = "expired"
	UserInvitationStatusDeclined UserInvitationStatus = "declined"
)

func (s UserInvitationStatus) String() string {
	return string(s)
}

func (s UserInvitationStatus) IsValid() bool {
	return s != UserInvitationStatusUnknown
}

func (s UserInvitationStatus) IsPending() bool {
	return s == UserInvitationStatusPending
}

func (s UserInvitationStatus) IsAccepted() bool {
	return s == UserInvitationStatusAccepted
}

func (s UserInvitationStatus) IsExpired() bool {
	return s == UserInvitationStatusExpired
}

func (s UserInvitationStatus) IsDeclined() bool {
	return s == UserInvitationStatusDeclined
}

func ToUserInvitationStatus(s string) UserInvitationStatus {
	switch s {
	case UserInvitationStatusPending.String():
		return UserInvitationStatusPending
	case UserInvitationStatusAccepted.String():
		return UserInvitationStatusAccepted
	case UserInvitationStatusExpired.String():
		return UserInvitationStatusExpired
	case UserInvitationStatusDeclined.String():
		return UserInvitationStatusDeclined
	default:
		return UserInvitationStatusUnknown
	}
}
