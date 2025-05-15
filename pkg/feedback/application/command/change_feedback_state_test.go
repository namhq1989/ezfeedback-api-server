package command_test

import (
	"context"
	"testing"

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

type changeFeedbackStateTestSuite struct {
	suite.Suite
	handler                            command.ChangeFeedbackStateHandler
	mockCtrl                           *gomock.Controller
	mockFeedbackRepository             *mockfeedback.MockFeedbackRepository
	mockFeedbackStateHistoryRepository *mockfeedback.MockFeedbackStateHistoryRepository
	mockCachingRepository              *mockfeedback.MockCachingRepository
	mockService                        *mockfeedback.MockService
}

func (s *changeFeedbackStateTestSuite) SetupSuite() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockFeedbackRepository = mockfeedback.NewMockFeedbackRepository(s.mockCtrl)
	s.mockFeedbackStateHistoryRepository = mockfeedback.NewMockFeedbackStateHistoryRepository(s.mockCtrl)
	s.mockCachingRepository = mockfeedback.NewMockCachingRepository(s.mockCtrl)
	s.mockService = mockfeedback.NewMockService(s.mockCtrl)
	s.handler = command.NewChangeFeedbackStateHandler(s.mockFeedbackRepository, s.mockFeedbackStateHistoryRepository, s.mockCachingRepository, s.mockService)
}

func (s *changeFeedbackStateTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

//
// CASES
//

func (s *changeFeedbackStateTestSuite) Test_1_Success() {
	s.mockService.EXPECT().
		GetFeedbackByID(gomock.Any(), gomock.Any()).
		Return(&domain.Feedback{ID: uuid.New()}, nil)

	s.mockFeedbackRepository.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		Return(nil)

	s.mockCachingRepository.EXPECT().
		SetFeedbackByID(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil)

	s.mockFeedbackStateHistoryRepository.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.ChangeFeedbackState(ctx, uuid.New(), uuid.New(), dto.ChangeFeedbackStateRequest{
		State: domain.FeedbackStateCompleted.String(),
	})

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
}

//
// END OF CASES
//

func TestChangeFeedbackStateTestSuite(t *testing.T) {
	suite.Run(t, new(changeFeedbackStateTestSuite))
}
