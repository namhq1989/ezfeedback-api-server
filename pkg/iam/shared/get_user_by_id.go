package shared

import (
	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/domain"
	"github.com/namhq1989/go-utilities/appcontext"
)

func (s Service) GetUserByID(ctx *appcontext.AppContext, id string) (*domain.User, error) {
	ctx.Logger().Info("[service] get user by id", appcontext.Fields{"id": id})

	ctx.Logger().Text("find user in caching")
	user, err := s.cachingRepository.GetUserByID(ctx, id)
	if user != nil {
		ctx.Logger().Text("user found in caching, return")
		return user, nil
	}
	if err != nil {
		ctx.Logger().Error("failed to find user in caching", err, appcontext.Fields{})
	}

	ctx.Logger().Text("user not found in caching, find in db")
	user, err = s.userRepository.FindByID(ctx, id)
	if err != nil {
		ctx.Logger().Error("failed to find user in db", err, appcontext.Fields{})
		return nil, err
	}
	if user == nil {
		ctx.Logger().ErrorText("user not found")
		return nil, apperrors.User.UserNotFound
	}

	ctx.Logger().Text("set user in caching")
	if err = s.cachingRepository.SetUserByID(ctx, id, *user); err != nil {
		ctx.Logger().Error("failed to set user in caching", err, appcontext.Fields{})
	}
	return user, nil
}
