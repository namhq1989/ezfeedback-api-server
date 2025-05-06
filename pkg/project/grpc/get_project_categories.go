package grpc

import (
	"github.com/namhq1989/ezfeedback-api-server/internal/genproto/projectpb"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
	"github.com/namhq1989/go-utilities/appcontext"
)

type GetProjectCategoriesHandler struct {
	service domain.Service
}

func NewGetProjectCategoriesHandler(service domain.Service) GetProjectCategoriesHandler {
	return GetProjectCategoriesHandler{
		service: service,
	}
}

func (h GetProjectCategoriesHandler) GetProjectCategories(ctx *appcontext.AppContext, req *projectpb.GetProjectCategoriesRequest) (*projectpb.GetProjectCategoriesResponse, error) {
	ctx.SetTraceID(req.GetTraceId())
	ctx.Logger().Info("new get project categories request", appcontext.Fields{"projectID": req.GetProjectId()})

	ctx.Logger().Text("find categories in db")
	categories, err := h.service.GetProjectCategoriesByProjectID(ctx, req.GetProjectId(), domain.StatusUnknown)
	if err != nil {
		ctx.Logger().Error("failed to find categories in db", err, appcontext.Fields{})
		return nil, err
	}

	var result = make([]*projectpb.ProjectCategory, 0)
	for _, category := range categories {
		result = append(result, &projectpb.ProjectCategory{
			Id:   category.ID,
			Name: category.Name,
		})
	}

	ctx.Logger().Text("done get project categories request")
	return &projectpb.GetProjectCategoriesResponse{
		Categories: result,
	}, nil
}
