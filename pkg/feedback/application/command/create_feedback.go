package command

import (
	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/internal/utils/validation"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/dto"
	"github.com/namhq1989/go-utilities/appcontext"
)

type CreateFeedbackHandler struct {
	feedbackRepository domain.FeedbackRepository
	queueRepository    domain.QueueRepository
	billingHub         domain.BillingHub
	projectHub         domain.ProjectHub
	service            domain.Service
}

func NewCreateFeedbackHandler(
	feedbackRepository domain.FeedbackRepository,
	queueRepository domain.QueueRepository,
	billingHub domain.BillingHub,
	projectHub domain.ProjectHub,
	service domain.Service,
) CreateFeedbackHandler {
	return CreateFeedbackHandler{
		feedbackRepository: feedbackRepository,
		queueRepository:    queueRepository,
		billingHub:         billingHub,
		projectHub:         projectHub,
		service:            service,
	}
}

// CreateFeedback godoc
// @tags     Feedback
// @summary  Create feedback
// @id       feedback-create
// @accept   json
// @produce  json
// @param    payload body    dto.CreateFeedbackRequest true "Body"
// @success  200     {object} dto.CreateFeedbackResponse
// @router   /api/feedback [post]
func (h CreateFeedbackHandler) CreateFeedback(ctx *appcontext.AppContext, ip, domainName string, req dto.CreateFeedbackRequest) (*dto.CreateFeedbackResponse, error) {
	ctx.Logger().Info("new create feedback request", appcontext.Fields{
		"ip": ip, "domainName": domainName, "campaignID": req.CampaignID, "email": req.Email, "categoryID": req.CategoryID,
		"content": req.Content, "rating": req.Rating, "context.userID": req.Context.UserID,
	})

	ctx.Logger().Text("get campaign data via grpc")
	campaignData, err := h.projectHub.GetProjectCampaignByID(ctx, req.CampaignID)
	if err != nil {
		ctx.Logger().Error("failed to get campaign data via grpc", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("validate data")
	if err = h.validateData(ctx, domainName, req, *campaignData); err != nil {
		ctx.Logger().Text("invalid data, respond")
		return nil, err
	}

	ctx.Logger().Text("count total created today")
	totalCreated, err := h.feedbackRepository.CountProjectTotalCreatedTodayByIp(ctx, campaignData.Project.ID, ip)
	if err != nil {
		ctx.Logger().Error("failed to count total created today", err, appcontext.Fields{})
		return nil, err
	}
	if domain.HasExceededDailyFeedbackLimit(totalCreated) {
		ctx.Logger().ErrorText("daily limit exceeded, respond")
		return nil, apperrors.Feedback.DailyLimitExceeded
	}

	ctx.Logger().Text("count total feedbacks of this month")
	thisMonthUsage, err := h.feedbackRepository.CountMonthlyUsageForProject(ctx, campaignData.Project.ID)
	if err != nil {
		ctx.Logger().Error("failed to count total feedbacks of this month", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("call billing service to check if exceed limit")
	canAcceptFeedback, err := h.billingHub.CanAcceptFeedback(ctx, campaignData.Project.ID, thisMonthUsage)
	if err != nil {
		ctx.Logger().Error("failed to call billing service to check if exceed limit", err, appcontext.Fields{})
		return nil, err
	}
	if !canAcceptFeedback {
		ctx.Logger().ErrorText("feedback usage exceed limit, respond")
		return nil, apperrors.Billing.FeedbackLimitExceeded
	}

	ctx.Logger().Text("get ip location data")
	var (
		country = ""
	)
	ipLocationData, err := h.service.GetIpLocationData(ctx, ip)
	if err != nil {
		ctx.Logger().Error("failed to get ip location data, use empty string", err, appcontext.Fields{})
	} else if ipLocationData == nil {
		ctx.Logger().ErrorText("ip location not found, use empty string")
	} else {
		country = ipLocationData.Country
	}

	ctx.Logger().Text("create new feedback model")
	feedback, err := domain.NewFeedback(
		campaignData.Project.ID,
		campaignData.ProjectCampaign.ID,
		req.Context.UserID,
		req.Email,
		req.CategoryID,
		req.Content,
		req.Rating,
		campaignData.ProjectCampaign.CampaignType.String(),
		ip,
		country,
	)
	if err != nil {
		ctx.Logger().Error("failed to create new feedback model", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("persist feedback")
	if err = h.feedbackRepository.Create(ctx, *feedback); err != nil {
		ctx.Logger().Error("failed to persist feedback", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("add new feedback to queue")
	if err = h.queueRepository.FeedbackCreated(ctx, domain.QueueFeedbackCreatedPayload{
		Feedback: *feedback,
	}); err != nil {
		ctx.Logger().Error("failed to add new feedback to queue", err, appcontext.Fields{})
	}

	ctx.Logger().Text("done create feedback")
	return &dto.CreateFeedbackResponse{}, nil
}

func (h CreateFeedbackHandler) validateData(ctx *appcontext.AppContext, domainName string, req dto.CreateFeedbackRequest, campaignData domain.ProjectCampaignHubData) error {
	if campaignData.ProjectCampaign.Status.IsInactive() {
		ctx.Logger().ErrorText("campaign is inactive")
		return apperrors.Project.InvalidCampaign
	}

	if campaignData.Project.Status.IsInactive() {
		ctx.Logger().ErrorText("project is inactive")
		return apperrors.Project.ProjectNotFound
	}

	if campaignData.Project.Setting.Domain != domainName {
		ctx.Logger().ErrorText("invalid domain")
		return apperrors.Project.ProjectNotFound
	}

	if len(req.Email) > 0 && !validation.IsValidEmail(req.Email) {
		ctx.Logger().ErrorText("invalid email")
		return apperrors.Common.InvalidEmail
	}

	if campaignData.ProjectCampaign.CampaignType.IsFeedback() {
		if len(req.Content) <= 0 || len(req.Content) > 2000 {
			ctx.Logger().ErrorText("campaign 'feedback' but invalid content")
			return apperrors.Feedback.InvalidContent
		}

		if req.Rating < 1 || req.Rating > 5 {
			ctx.Logger().ErrorText("campaign 'feedback' but invalid rating")
			return apperrors.Feedback.InvalidRating
		}

		if req.CategoryID == "" {
			ctx.Logger().ErrorText("campaign 'feedback' but invalid category")
			return apperrors.Project.InvalidCategory
		}
	} else if campaignData.ProjectCampaign.CampaignType.IsNPS() {
		if req.Rating < 1 || req.Rating > 10 {
			ctx.Logger().ErrorText("campaign 'nps' but invalid rating")
			return apperrors.Feedback.InvalidRating
		}
	} else if campaignData.ProjectCampaign.CampaignType.IsCSAT() {
		if req.Rating < 1 || req.Rating > 5 {
			ctx.Logger().ErrorText("campaign 'csat' but invalid rating")
			return apperrors.Feedback.InvalidRating
		}
	} else {
		ctx.Logger().ErrorText("invalid campaign type")
		return apperrors.Project.InvalidCampaignType
	}

	return nil
}
