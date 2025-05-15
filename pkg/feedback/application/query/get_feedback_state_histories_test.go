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

type getFeedbackStateHistoriesTestSuite struct {
	suite.Suite
	handler                            query.GetFeedbackStateHistoriesHandler
	mockCtrl                           *gomock.Controller
	mockFeedbackStateHistoryRepository *mockfeedback.MockFeedbackStateHistoryRepository
	mockIAMHub                         *mockfeedback.MockIAMHub
}

func (s *getFeedbackStateHistoriesTestSuite) SetupSuite() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockFeedbackStateHistoryRepository = mockfeedback.NewMockFeedbackStateHistoryRepository(s.mockCtrl)
	s.mockIAMHub = mockfeedback.NewMockIAMHub(s.mockCtrl)
	s.handler = query.NewGetFeedbackStateHistoriesHandler(s.mockFeedbackStateHistoryRepository, s.mockIAMHub)
}

func (s *getFeedbackStateHistoriesTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

//
// CASES
//

func (s *getFeedbackStateHistoriesTestSuite) Test_1_Success() {
	s.mockFeedbackStateHistoryRepository.EXPECT().
		FindWithFilter(gomock.Any(), gomock.Any()).
		Return([]domain.FeedbackStateHistory{
			{ID: uuid.New()},
		}, nil)

	s.mockIAMHub.EXPECT().
		GetUserByID(gomock.Any(), gomock.Any()).
		Return(&domain.User{ID: uuid.New()}, nil).
		AnyTimes()

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.GetFeedbackStateHistories(ctx, uuid.New(), uuid.New(), dto.GetFeedbackStateHistoriesRequest{
		Page: 0,
	})

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.Equal(s.T(), 1, len(resp.Histories))
}

//
// END OF CASES
//

func TestGetFeedbackStateHistoriesTestSuite(t *testing.T) {
	suite.Run(t, new(getFeedbackStateHistoriesTestSuite))
}
