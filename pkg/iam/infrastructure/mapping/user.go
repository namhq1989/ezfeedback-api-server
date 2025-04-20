package mapping

import (
	"github.com/namhq1989/ezfeedback-api-server/internal/database/gen/ezfeedback/public/model"
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/domain"
)

type UserMapper struct{}

func (UserMapper) FromModelToDomain(user model.Users) (*domain.User, error) {
	var result = &domain.User{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Status:    domain.ToUserStatus(user.Status.String()),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return result, nil
}

func (UserMapper) FromDomainToModel(user domain.User) (*model.Users, error) {
	var result = &model.Users{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Status:    model.UserStatus(user.Status.String()),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return result, nil
}
