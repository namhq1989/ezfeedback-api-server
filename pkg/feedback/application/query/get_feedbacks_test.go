package query_test

import (
	"context"
	"testing"

	mockfeedback "github.com/namhq1989/ezfeedback-api-server/internal/mock/feedback"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/application/query"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/dto"
	"github.com/namhq1989/go-utilities/appcontext"
	"github.com/namhq1989/go-utilities/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type getFeedbacksTestSuite struct {
	suite.Suite
	handler                query.GetFeedbacksHandler
	mockCtrl               *gomock.Controller
	mockFeedbackRepository *mockfeedback.MockFeedbackRepository
	mockProjectHub         *mockfeedback.MockProjectHub
}

func (s *getFeedbacksTestSuite) SetupSuite() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockFeedbackRepository = mockfeedback.NewMockFeedbackRepository(s.mockCtrl)
	s.mockProjectHub = mockfeedback.NewMockProjectHub(s.mockCtrl)
	s.handler = query.NewGetFeedbacksHandler(s.mockFeedbackRepository, s.mockProjectHub)
}

func (s *getFeedbacksTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

//
// CASES
//

func (s *getFeedbacksTestSuite) Test_1_Success() {
	s.mockProjectHub.EXPECT().
		GetProjectByID(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&domain.Project{ID: uuid.New()}, nil)

	s.mockFeedbackRepository.EXPECT().
		FindWithFilter(gomock.Any(), gomock.Any()).
		Return([]domain.Feedback{
			{ID: uuid.New()},
		}, nil)

	s.mockProjectHub.EXPECT().
		GetProjectCategories(gomock.Any(), gomock.Any()).
		Return([]domain.ProjectCategory{
			{ID: uuid.New()},
		}, nil)

	s.mockProjectHub.EXPECT().
		GetProjectCampaigns(gomock.Any(), gomock.Any()).
		Return([]domain.ProjectCampaign{
			{ID: uuid.New()},
		}, nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.GetFeedbacks(ctx, uuid.New(), dto.GetFeedbacksRequest{
		ProjectID:    uuid.New(),
		CategoryID:   uuid.New(),
		CampaignType: domain.ProjectCampaignTypeUnknown.String(),
		Keyword:      "",
		State:        domain.FeedbackStateUnknown.String(),
		Rating:       0,
		Page:         0,
	})

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.Equal(s.T(), 1, len(resp.Feedbacks))
}

//
// END OF CASES
//

func TestGetFeedbacksTestSuite(t *testing.T) {
	suite.Run(t, new(getFeedbacksTestSuite))
}
