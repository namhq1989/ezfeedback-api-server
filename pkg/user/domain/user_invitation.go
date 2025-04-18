package domain

import "time"

type UserInvitation struct {
	ID        string
	Email     string
	InviterID string
	ProjectID string
	Role      ProjectRole
	Token     string
	Code      string
	Status    UserInvitationStatus
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
