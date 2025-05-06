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

type countFeedbacksTestSuite struct {
	suite.Suite
	handler                query.CountFeedbacksHandler
	mockCtrl               *gomock.Controller
	mockFeedbackRepository *mockfeedback.MockFeedbackRepository
}

func (s *countFeedbacksTestSuite) SetupSuite() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockFeedbackRepository = mockfeedback.NewMockFeedbackRepository(s.mockCtrl)
	s.handler = query.NewCountFeedbacksHandler(s.mockFeedbackRepository)
}

func (s *countFeedbacksTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

//
// CASES
//

func (s *countFeedbacksTestSuite) Test_1_Success() {
	s.mockFeedbackRepository.EXPECT().
		CountWithFilter(gomock.Any(), gomock.Any()).
		Return(int64(10), nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.CountFeedbacks(ctx, uuid.New(), dto.CountFeedbacksRequest{
		ProjectID:    uuid.New(),
		CategoryID:   uuid.New(),
		CampaignType: domain.ProjectCampaignTypeUnknown.String(),
		Keyword:      "",
		State:        domain.FeedbackStateUnknown.String(),
		Rating:       0,
	})

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.Equal(s.T(), int64(10), resp.Total)
}

//
// END OF CASES
//

func TestCountFeedbacksTestSuite(t *testing.T) {
	suite.Run(t, new(countFeedbacksTestSuite))
}
