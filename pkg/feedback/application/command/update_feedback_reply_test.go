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

type updateFeedbackReplyTestSuite struct {
	suite.Suite
	handler                     command.UpdateFeedbackReplyHandler
	mockCtrl                    *gomock.Controller
	mockFeedbackReplyRepository *mockfeedback.MockFeedbackReplyRepository
}

func (s *updateFeedbackReplyTestSuite) SetupSuite() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockFeedbackReplyRepository = mockfeedback.NewMockFeedbackReplyRepository(s.mockCtrl)
	s.handler = command.NewUpdateFeedbackReplyHandler(s.mockFeedbackReplyRepository)
}

func (s *updateFeedbackReplyTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

//
// CASES
//

func (s *updateFeedbackReplyTestSuite) Test_1_Success() {
	var feedbackID = uuid.New()

	s.mockFeedbackReplyRepository.EXPECT().
		FindByID(gomock.Any(), gomock.Any()).
		Return(&domain.FeedbackReply{ID: uuid.New(), FeedbackID: feedbackID}, nil)

	s.mockFeedbackReplyRepository.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		Return(nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.UpdateFeedbackReply(ctx, uuid.New(), feedbackID, uuid.New(), dto.UpdateFeedbackReplyRequest{
		Content: "Updated content",
	})

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
}

func (s *updateFeedbackReplyTestSuite) Test_2_Fail_NotBelongToFeedback() {
	s.mockFeedbackReplyRepository.EXPECT().
		FindByID(gomock.Any(), gomock.Any()).
		Return(&domain.FeedbackReply{ID: uuid.New()}, nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.UpdateFeedbackReply(ctx, uuid.New(), uuid.New(), uuid.New(), dto.UpdateFeedbackReplyRequest{
		Content: "",
	})

	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), apperrors.Feedback.InvalidReply, err)
}

func (s *updateFeedbackReplyTestSuite) Test_2_Fail_InvalidContent() {
	var feedbackID = uuid.New()

	s.mockFeedbackReplyRepository.EXPECT().
		FindByID(gomock.Any(), gomock.Any()).
		Return(&domain.FeedbackReply{ID: uuid.New(), FeedbackID: feedbackID}, nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.UpdateFeedbackReply(ctx, uuid.New(), feedbackID, uuid.New(), dto.UpdateFeedbackReplyRequest{
		Content: "",
	})

	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), apperrors.Feedback.InvalidContent, err)
}

//
// END OF CASES
//

func TestUpdateFeedbackReplyTestSuite(t *testing.T) {
	suite.Run(t, new(updateFeedbackReplyTestSuite))
}
