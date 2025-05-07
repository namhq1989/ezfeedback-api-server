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

type deleteFeedbackReplyTestSuite struct {
	suite.Suite
	handler                     command.DeleteFeedbackReplyHandler
	mockCtrl                    *gomock.Controller
	mockFeedbackRepository      *mockfeedback.MockFeedbackRepository
	mockFeedbackReplyRepository *mockfeedback.MockFeedbackReplyRepository
	mockService                 *mockfeedback.MockService
}

func (s *deleteFeedbackReplyTestSuite) SetupSuite() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockFeedbackRepository = mockfeedback.NewMockFeedbackRepository(s.mockCtrl)
	s.mockFeedbackReplyRepository = mockfeedback.NewMockFeedbackReplyRepository(s.mockCtrl)
	s.mockService = mockfeedback.NewMockService(s.mockCtrl)
	s.handler = command.NewDeleteFeedbackReplyHandler(s.mockFeedbackReplyRepository, s.mockFeedbackRepository, s.mockService)
}

func (s *deleteFeedbackReplyTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

//
// CASES
//

func (s *deleteFeedbackReplyTestSuite) Test_1_Success() {
	var feedbackID = uuid.New()

	s.mockService.EXPECT().
		GetFeedbackByID(gomock.Any(), gomock.Any()).
		Return(&domain.Feedback{ID: feedbackID}, nil)

	s.mockFeedbackReplyRepository.EXPECT().
		FindByID(gomock.Any(), gomock.Any()).
		Return(&domain.FeedbackReply{ID: uuid.New(), FeedbackID: feedbackID}, nil)

	s.mockFeedbackReplyRepository.EXPECT().
		Delete(gomock.Any(), gomock.Any()).
		Return(nil)

	s.mockFeedbackRepository.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		Return(nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.DeleteFeedbackReply(ctx, uuid.New(), feedbackID, uuid.New(), dto.DeleteFeedbackReplyRequest{})

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
}

func (s *deleteFeedbackReplyTestSuite) Test_2_Fail_NotBelongToFeedback() {
	s.mockService.EXPECT().
		GetFeedbackByID(gomock.Any(), gomock.Any()).
		Return(&domain.Feedback{ID: uuid.New()}, nil)

	s.mockFeedbackReplyRepository.EXPECT().
		FindByID(gomock.Any(), gomock.Any()).
		Return(&domain.FeedbackReply{ID: uuid.New()}, nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.DeleteFeedbackReply(ctx, uuid.New(), uuid.New(), uuid.New(), dto.DeleteFeedbackReplyRequest{})

	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), apperrors.Feedback.InvalidReply, err)
}

//
// END OF CASES
//

func TestDeleteFeedbackReplyTestSuite(t *testing.T) {
	suite.Run(t, new(deleteFeedbackReplyTestSuite))
}
