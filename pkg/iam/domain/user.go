package domain

import (
	"time"

	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/internal/utils/manipulation"
	"github.com/namhq1989/ezfeedback-api-server/internal/utils/validation"
	"github.com/namhq1989/go-utilities/appcontext"
	"github.com/namhq1989/go-utilities/uuid"
)

type UserRepository interface {
	Create(ctx *appcontext.AppContext, user User) error
	Update(ctx *appcontext.AppContext, user User) error
	FindByID(ctx *appcontext.AppContext, userID string) (*User, error)
	FindByEmail(ctx *appcontext.AppContext, email string) (*User, error)
}

type User struct {
	ID        string
	Email     string
	Name      string
	Status    UserStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(email string) (*User, error) {
	var u = &User{
		ID:        uuid.New(),
		Status:    UserStatusActive,
		CreatedAt: manipulation.NowUTC(),
		UpdatedAt: manipulation.NowUTC(),
	}

	if err := u.SetEmail(email); err != nil {
		return nil, err
	}
	if err := u.SetName(email); err != nil {
		return nil, err
	}

	return u, nil
}

func (u *User) SetEmail(email string) error {
	if !validation.IsValidEmail(email) {
		return apperrors.Common.InvalidEmail
	}

	u.Email = email
	u.SetUpdatedAt()
	return nil
}

func (u *User) SetName(name string) error {
	if name == "" || len(name) < 3 || len(name) > 255 {
		return apperrors.Common.InvalidName
	}

	u.Name = name
	if u.Name == "" {
		u.Name = u.Email
	}

	u.SetUpdatedAt()
	return nil
}

func (u *User) SetStatus(status string) error {
	dStatus := ToUserStatus(status)
	if !dStatus.IsValid() {
		return apperrors.Common.InvalidStatus
	}

	u.Status = dStatus
	u.SetUpdatedAt()
	return nil
}

func (u *User) SetUpdatedAt() {
	u.UpdatedAt = manipulation.NowUTC()
}
