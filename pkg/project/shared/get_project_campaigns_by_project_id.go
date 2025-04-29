package shared

import (
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
	"github.com/namhq1989/go-utilities/appcontext"
)

func (s Service) GetProjectCampaignsByProjectID(ctx *appcontext.AppContext, projectID string) ([]domain.ProjectCampaign, error) {
	ctx.Logger().Info("[service] get project campaigns by project id", appcontext.Fields{"projectID": projectID})

	ctx.Logger().Text("find project campaigns in caching")
	campaigns, err := s.cachingRepository.GetProjectCampaignsByProjectID(ctx, projectID)
	if campaigns != nil {
		ctx.Logger().Text("project campaigns found in caching, return")
		return campaigns, nil
	}
	if err != nil {
		ctx.Logger().Error("failed to find project campaigns in caching", err, appcontext.Fields{})
	}

	ctx.Logger().Text("project campaigns not found in caching, find in db")
	campaigns, err = s.projectCampaignRepository.FindByProjectID(ctx, projectID)
	if err != nil {
		ctx.Logger().Error("failed to find project campaigns in db", err, appcontext.Fields{})
		return nil, err
	}
	if campaigns == nil || len(campaigns) == 0 {
		ctx.Logger().Text("project campaigns not found")
		return make([]domain.ProjectCampaign, 0), nil
	}

	ctx.Logger().Text("set project campaigns in caching")
	if err = s.cachingRepository.SetProjectCampaignsByProjectID(ctx, projectID, campaigns); err != nil {
		ctx.Logger().Error("failed to set project campaigns in caching", err, appcontext.Fields{})
	}
	return campaigns, nil
}
