package command_test

import (
	"context"
	"testing"

	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	mockfeedback "github.com/namhq1989/ezfeedback-api-server/internal/mock/feedback"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/application/command"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/dto"
	"github.com/namhq1989/go-utilities/appcontext"
	"github.com/namhq1989/go-utilities/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type createFeedbackTestSuite struct {
	suite.Suite
	handler          command.CreateFeedbackHandler
	mockCtrl         *gomock.Controller
	mockFeedbackRepo *mockfeedback.MockFeedbackRepository
	mockQueueRepo    *mockfeedback.MockQueueRepository
	mockBillingHub   *mockfeedback.MockBillingHub
	mockProjectHub   *mockfeedback.MockProjectHub
	mockService      *mockfeedback.MockService
}

func (s *createFeedbackTestSuite) SetupSuite() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockFeedbackRepo = mockfeedback.NewMockFeedbackRepository(s.mockCtrl)
	s.mockQueueRepo = mockfeedback.NewMockQueueRepository(s.mockCtrl)
	s.mockBillingHub = mockfeedback.NewMockBillingHub(s.mockCtrl)
	s.mockProjectHub = mockfeedback.NewMockProjectHub(s.mockCtrl)
	s.mockService = mockfeedback.NewMockService(s.mockCtrl)
	s.handler = command.NewCreateFeedbackHandler(s.mockFeedbackRepo, s.mockQueueRepo, s.mockBillingHub, s.mockProjectHub, s.mockService)
}

func (s *createFeedbackTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

//
// CASES
//

func (s *createFeedbackTestSuite) Test_1_Success() {
	projectID := uuid.New()
	campaignID := uuid.New()
	campaignData := &domain.ProjectCampaignHubData{
		Project: domain.Project{
			ID:     projectID,
			Status: domain.StatusActive,
			Setting: domain.ProjectSetting{
				Domain: "test.com",
			},
		},
		ProjectCampaign: domain.ProjectCampaign{
			ID:           campaignID,
			CampaignType: domain.ProjectCampaignTypeFeedback,
			Status:       domain.StatusActive,
		},
	}

	s.mockProjectHub.EXPECT().
		GetProjectCampaignByID(gomock.Any(), gomock.Any()).
		Return(campaignData, nil)

	s.mockFeedbackRepo.EXPECT().
		CountMonthlyUsageForProject(gomock.Any(), gomock.Any()).
		Return(int64(10), nil)

	s.mockBillingHub.EXPECT().
		CanAcceptFeedback(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(true, nil)

	s.mockService.EXPECT().
		GetIpLocationData(gomock.Any(), gomock.Any()).
		Return(&domain.IpLocationData{Country: "us"}, nil)

	s.mockFeedbackRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(nil)

	s.mockQueueRepo.EXPECT().
		FeedbackCreated(gomock.Any(), gomock.Any()).
		Return(nil)

	ctx := appcontext.NewRest(context.Background())
	userID := uuid.New()
	resp, err := s.handler.CreateFeedback(ctx, "127.0.0.1", "test.com", dto.CreateFeedbackRequest{
		CampaignID: campaignID,
		Email:      "test@example.com",
		CategoryID: uuid.New(),
		Content:    "This is a test feedback",
		Rating:     5,
		Context: dto.CreateFeedbackRequestContext{
			UserID: userID,
		},
	})

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
}

func (s *createFeedbackTestSuite) Test_2_Fail_GetProjectCampaignByIDError() {
	s.mockProjectHub.EXPECT().
		GetProjectCampaignByID(gomock.Any(), gomock.Any()).
		Return(nil, apperrors.Project.InvalidCampaign)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.CreateFeedback(ctx, "127.0.0.1", "test.com", dto.CreateFeedbackRequest{
		CampaignID: uuid.New(),
	})

	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), apperrors.Project.InvalidCampaign, err)
}

func (s *createFeedbackTestSuite) Test_2_Fail_InactiveCampaign() {
	projectID := uuid.New()
	campaignID := uuid.New()
	campaignData := &domain.ProjectCampaignHubData{
		Project: domain.Project{
			ID:     projectID,
			Status: domain.StatusActive,
			Setting: domain.ProjectSetting{
				Domain: "test.com",
			},
		},
		ProjectCampaign: domain.ProjectCampaign{
			ID:           campaignID,
			CampaignType: domain.ProjectCampaignTypeFeedback,
			Status:       domain.StatusInactive,
		},
	}

	s.mockProjectHub.EXPECT().
		GetProjectCampaignByID(gomock.Any(), gomock.Any()).
		Return(campaignData, nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.CreateFeedback(ctx, "127.0.0.1", "test.com", dto.CreateFeedbackRequest{
		CampaignID: campaignID,
	})

	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), apperrors.Project.InvalidCampaign, err)
}

func (s *createFeedbackTestSuite) Test_2_Fail_InactiveProject() {
	projectID := uuid.New()
	campaignID := uuid.New()
	campaignData := &domain.ProjectCampaignHubData{
		Project: domain.Project{
			ID:     projectID,
			Status: domain.StatusInactive,
			Setting: domain.ProjectSetting{
				Domain: "test.com",
			},
		},
		ProjectCampaign: domain.ProjectCampaign{
			ID:           campaignID,
			CampaignType: domain.ProjectCampaignTypeFeedback,
			Status:       domain.StatusActive,
		},
	}

	s.mockProjectHub.EXPECT().
		GetProjectCampaignByID(gomock.Any(), gomock.Any()).
		Return(campaignData, nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.CreateFeedback(ctx, "127.0.0.1", "test.com", dto.CreateFeedbackRequest{
		CampaignID: campaignID,
	})

	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), apperrors.Project.ProjectNotFound, err)
}

func (s *createFeedbackTestSuite) Test_2_Fail_InvalidDomain() {
	projectID := uuid.New()
	campaignID := uuid.New()
	campaignData := &domain.ProjectCampaignHubData{
		Project: domain.Project{
			ID:     projectID,
			Status: domain.StatusActive,
			Setting: domain.ProjectSetting{
				Domain: "correct-domain.com",
			},
		},
		ProjectCampaign: domain.ProjectCampaign{
			ID:           campaignID,
			CampaignType: domain.ProjectCampaignTypeFeedback,
			Status:       domain.StatusActive,
		},
	}

	s.mockProjectHub.EXPECT().
		GetProjectCampaignByID(gomock.Any(), gomock.Any()).
		Return(campaignData, nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.CreateFeedback(ctx, "127.0.0.1", "wrong-domain.com", dto.CreateFeedbackRequest{
		CampaignID: campaignID,
	})

	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), apperrors.Project.ProjectNotFound, err)
}

func (s *createFeedbackTestSuite) Test_2_Fail_InvalidEmail() {
	projectID := uuid.New()
	campaignID := uuid.New()
	campaignData := &domain.ProjectCampaignHubData{
		Project: domain.Project{
			ID:     projectID,
			Status: domain.StatusActive,
			Setting: domain.ProjectSetting{
				Domain: "test.com",
			},
		},
		ProjectCampaign: domain.ProjectCampaign{
			ID:           campaignID,
			CampaignType: domain.ProjectCampaignTypeFeedback,
			Status:       domain.StatusActive,
		},
	}

	s.mockProjectHub.EXPECT().
		GetProjectCampaignByID(gomock.Any(), gomock.Any()).
		Return(campaignData, nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.CreateFeedback(ctx, "127.0.0.1", "test.com", dto.CreateFeedbackRequest{
		CampaignID: campaignID,
		Email:      "invalid-email",
	})

	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), apperrors.Common.InvalidEmail, err)
}

func (s *createFeedbackTestSuite) Test_2_Fail_FeedbackCampaignInvalidContent() {
	projectID := uuid.New()
	campaignID := uuid.New()
	campaignData := &domain.ProjectCampaignHubData{
		Project: domain.Project{
			ID:     projectID,
			Status: domain.StatusActive,
			Setting: domain.ProjectSetting{
				Domain: "test.com",
			},
		},
		ProjectCampaign: domain.ProjectCampaign{
			ID:           campaignID,
			CampaignType: domain.ProjectCampaignTypeFeedback,
			Status:       domain.StatusActive,
		},
	}

	s.mockProjectHub.EXPECT().
		GetProjectCampaignByID(gomock.Any(), gomock.Any()).
		Return(campaignData, nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.CreateFeedback(ctx, "127.0.0.1", "test.com", dto.CreateFeedbackRequest{
		CampaignID: campaignID,
		Content:    "",
		Rating:     5,
		CategoryID: uuid.New(),
	})

	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), apperrors.Feedback.InvalidContent, err)
}

func (s *createFeedbackTestSuite) Test_2_Fail_FeedbackCampaignInvalidRating() {
	projectID := uuid.New()
	campaignID := uuid.New()
	campaignData := &domain.ProjectCampaignHubData{
		Project: domain.Project{
			ID:     projectID,
			Status: domain.StatusActive,
			Setting: domain.ProjectSetting{
				Domain: "test.com",
			},
		},
		ProjectCampaign: domain.ProjectCampaign{
			ID:           campaignID,
			CampaignType: domain.ProjectCampaignTypeFeedback,
			Status:       domain.StatusActive,
		},
	}

	s.mockProjectHub.EXPECT().
		GetProjectCampaignByID(gomock.Any(), gomock.Any()).
		Return(campaignData, nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.CreateFeedback(ctx, "127.0.0.1", "test.com", dto.CreateFeedbackRequest{
		CampaignID: campaignID,
		Content:    "Test content",
		Rating:     6,
		CategoryID: uuid.New(),
	})

	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), apperrors.Feedback.InvalidRating, err)
}

func (s *createFeedbackTestSuite) Test_2_Fail_FeedbackCampaignInvalidCategory() {
	// Mock campaign data
	projectID := uuid.New()
	campaignID := uuid.New()
	campaignData := &domain.ProjectCampaignHubData{
		Project: domain.Project{
			ID:     projectID,
			Status: domain.StatusActive,
			Setting: domain.ProjectSetting{
				Domain: "test.com",
			},
		},
		ProjectCampaign: domain.ProjectCampaign{
			ID:           campaignID,
			CampaignType: domain.ProjectCampaignTypeFeedback,
			Status:       domain.StatusActive,
		},
	}

	s.mockProjectHub.EXPECT().
		GetProjectCampaignByID(gomock.Any(), gomock.Any()).
		Return(campaignData, nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.CreateFeedback(ctx, "127.0.0.1", "test.com", dto.CreateFeedbackRequest{
		CampaignID: campaignID,
		Content:    "Test content",
		Rating:     5,
		CategoryID: "", // Missing category for feedback campaign
	})

	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), apperrors.Project.InvalidCategory, err)
}

func (s *createFeedbackTestSuite) Test_2_Fail_NPSCampaignInvalidRating() {
	// Mock campaign data
	projectID := uuid.New()
	campaignID := uuid.New()
	campaignData := &domain.ProjectCampaignHubData{
		Project: domain.Project{
			ID:     projectID,
			Status: domain.StatusActive,
			Setting: domain.ProjectSetting{
				Domain: "test.com",
			},
		},
		ProjectCampaign: domain.ProjectCampaign{
			ID:           campaignID,
			CampaignType: domain.ProjectCampaignTypeNPS,
			Status:       domain.StatusActive,
		},
	}

	s.mockProjectHub.EXPECT().
		GetProjectCampaignByID(gomock.Any(), gomock.Any()).
		Return(campaignData, nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.CreateFeedback(ctx, "127.0.0.1", "test.com", dto.CreateFeedbackRequest{
		CampaignID: campaignID,
		Rating:     11, // Invalid rating for NPS campaign (should be 1-10)
	})

	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), apperrors.Feedback.InvalidRating, err)
}

func (s *createFeedbackTestSuite) Test_2_Fail_CSATCampaignInvalidRating() {
	// Mock campaign data
	projectID := uuid.New()
	campaignID := uuid.New()
	campaignData := &domain.ProjectCampaignHubData{
		Project: domain.Project{
			ID:     projectID,
			Status: domain.StatusActive,
			Setting: domain.ProjectSetting{
				Domain: "test.com",
			},
		},
		ProjectCampaign: domain.ProjectCampaign{
			ID:           campaignID,
			CampaignType: domain.ProjectCampaignTypeCSAT,
			Status:       domain.StatusActive,
		},
	}

	s.mockProjectHub.EXPECT().
		GetProjectCampaignByID(gomock.Any(), gomock.Any()).
		Return(campaignData, nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.CreateFeedback(ctx, "127.0.0.1", "test.com", dto.CreateFeedbackRequest{
		CampaignID: campaignID,
		Rating:     0, // Invalid rating for CSAT campaign (should be 1-10)
	})

	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), apperrors.Feedback.InvalidRating, err)
}

func (s *createFeedbackTestSuite) Test_2_Fail_InvalidCampaignType() {
	// Mock campaign data with invalid campaign type
	projectID := uuid.New()
	campaignID := uuid.New()
	campaignData := &domain.ProjectCampaignHubData{
		Project: domain.Project{
			ID:     projectID,
			Status: domain.StatusActive,
			Setting: domain.ProjectSetting{
				Domain: "test.com",
			},
		},
		ProjectCampaign: domain.ProjectCampaign{
			ID:           campaignID,
			CampaignType: "invalid-type", // Invalid campaign type
			Status:       domain.StatusActive,
		},
	}

	s.mockProjectHub.EXPECT().
		GetProjectCampaignByID(gomock.Any(), gomock.Any()).
		Return(campaignData, nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.CreateFeedback(ctx, "127.0.0.1", "test.com", dto.CreateFeedbackRequest{
		CampaignID: campaignID,
		Rating:     5,
	})

	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), apperrors.Project.InvalidCampaignType, err)
}

func (s *createFeedbackTestSuite) Test_2_Fail_CountMonthlyUsageError() {
	// Mock campaign data
	projectID := uuid.New()
	campaignID := uuid.New()
	campaignData := &domain.ProjectCampaignHubData{
		Project: domain.Project{
			ID:     projectID,
			Status: domain.StatusActive,
			Setting: domain.ProjectSetting{
				Domain: "test.com",
			},
		},
		ProjectCampaign: domain.ProjectCampaign{
			ID:           campaignID,
			CampaignType: domain.ProjectCampaignTypeFeedback,
			Status:       domain.StatusActive,
		},
	}

	s.mockProjectHub.EXPECT().
		GetProjectCampaignByID(gomock.Any(), gomock.Any()).
		Return(campaignData, nil)

	s.mockFeedbackRepo.EXPECT().
		CountMonthlyUsageForProject(gomock.Any(), gomock.Any()).
		Return(int64(0), apperrors.Billing.FeedbackLimitExceeded)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.CreateFeedback(ctx, "127.0.0.1", "test.com", dto.CreateFeedbackRequest{
		CampaignID: campaignID,
		Email:      "test@example.com",
		CategoryID: uuid.New(),
		Content:    "Test content",
		Rating:     5,
	})

	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), apperrors.Billing.FeedbackLimitExceeded, err)
}

func (s *createFeedbackTestSuite) Test_2_Fail_CanAcceptFeedbackError() {
	// Mock campaign data
	projectID := uuid.New()
	campaignID := uuid.New()
	campaignData := &domain.ProjectCampaignHubData{
		Project: domain.Project{
			ID:     projectID,
			Status: domain.StatusActive,
			Setting: domain.ProjectSetting{
				Domain: "test.com",
			},
		},
		ProjectCampaign: domain.ProjectCampaign{
			ID:           campaignID,
			CampaignType: domain.ProjectCampaignTypeFeedback,
			Status:       domain.StatusActive,
		},
	}

	s.mockProjectHub.EXPECT().
		GetProjectCampaignByID(gomock.Any(), gomock.Any()).
		Return(campaignData, nil)

	s.mockFeedbackRepo.EXPECT().
		CountMonthlyUsageForProject(gomock.Any(), gomock.Any()).
		Return(int64(10), nil)

	s.mockBillingHub.EXPECT().
		CanAcceptFeedback(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(false, apperrors.Billing.FeedbackLimitExceeded)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.CreateFeedback(ctx, "127.0.0.1", "test.com", dto.CreateFeedbackRequest{
		CampaignID: campaignID,
		Email:      "test@example.com",
		CategoryID: uuid.New(),
		Content:    "Test content",
		Rating:     5,
	})

	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), apperrors.Billing.FeedbackLimitExceeded, err)
}

func (s *createFeedbackTestSuite) Test_2_Fail_FeedbackLimitExceeded() {
	// Mock campaign data
	projectID := uuid.New()
	campaignID := uuid.New()
	campaignData := &domain.ProjectCampaignHubData{
		Project: domain.Project{
			ID:     projectID,
			Status: domain.StatusActive,
			Setting: domain.ProjectSetting{
				Domain: "test.com",
			},
		},
		ProjectCampaign: domain.ProjectCampaign{
			ID:           campaignID,
			CampaignType: domain.ProjectCampaignTypeFeedback,
			Status:       domain.StatusActive,
		},
	}

	s.mockProjectHub.EXPECT().
		GetProjectCampaignByID(gomock.Any(), gomock.Any()).
		Return(campaignData, nil)

	s.mockFeedbackRepo.EXPECT().
		CountMonthlyUsageForProject(gomock.Any(), gomock.Any()).
		Return(int64(100), nil)

	s.mockBillingHub.EXPECT().
		CanAcceptFeedback(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(false, nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.CreateFeedback(ctx, "127.0.0.1", "test.com", dto.CreateFeedbackRequest{
		CampaignID: campaignID,
		Email:      "test@example.com",
		CategoryID: uuid.New(),
		Content:    "Test content",
		Rating:     5,
	})

	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), apperrors.Billing.FeedbackLimitExceeded, err)
}

//
// END OF CASES
//

func TestCreateFeedbackTestSuite(t *testing.T) {
	suite.Run(t, new(createFeedbackTestSuite))
}
