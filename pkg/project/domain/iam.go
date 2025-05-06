package domain

import "github.com/namhq1989/go-utilities/appcontext"

type IAMHub interface {
	GetUserByID(ctx *appcontext.AppContext, userID string) (*User, error)
}

type User struct {
	ID     string
	Name   string
	Email  string
	Status UserStatus
}

type UserStatus string

const (
	UserStatusUnknown  UserStatus = ""
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusDeleted  UserStatus = "deleted"
)

func (s UserStatus) String() string {
	return string(s)
}

func (s UserStatus) IsValid() bool {
	return s != UserStatusUnknown
}

func (s UserStatus) IsDeleted() bool {
	return s == UserStatusDeleted
}

func (s UserStatus) IsInactive() bool {
	return s == UserStatusInactive
}

func (s UserStatus) IsActive() bool {
	return s == UserStatusActive
}

func ToUserStatus(s string) UserStatus {
	switch s {
	case UserStatusActive.String():
		return UserStatusActive
	case UserStatusInactive.String():
		return UserStatusInactive
	case UserStatusDeleted.String():
		return UserStatusDeleted
	default:
		return UserStatusUnknown
	}
}
