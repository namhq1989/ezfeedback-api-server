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

type createFeedbackReplyTestSuite struct {
	suite.Suite
	handler                     command.CreateFeedbackReplyHandler
	mockCtrl                    *gomock.Controller
	mockFeedbackRepository      *mockfeedback.MockFeedbackRepository
	mockFeedbackReplyRepository *mockfeedback.MockFeedbackReplyRepository
	mockService                 *mockfeedback.MockService
}

func (s *createFeedbackReplyTestSuite) SetupSuite() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockFeedbackRepository = mockfeedback.NewMockFeedbackRepository(s.mockCtrl)
	s.mockFeedbackReplyRepository = mockfeedback.NewMockFeedbackReplyRepository(s.mockCtrl)
	s.mockService = mockfeedback.NewMockService(s.mockCtrl)
	s.handler = command.NewCreateFeedbackReplyHandler(s.mockFeedbackReplyRepository, s.mockFeedbackRepository, s.mockService)
}

func (s *createFeedbackReplyTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

//
// CASES
//

func (s *createFeedbackReplyTestSuite) Test_1_Success() {
	s.mockService.EXPECT().
		GetFeedbackByID(gomock.Any(), gomock.Any()).
		Return(&domain.Feedback{ID: uuid.New()}, nil)

	s.mockFeedbackReplyRepository.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(nil)

	s.mockFeedbackRepository.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		Return(nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.CreateFeedbackReply(ctx, uuid.New(), uuid.New(), dto.CreateFeedbackReplyRequest{
		Content: "This is a test reply",
	})

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
}

func (s *createFeedbackReplyTestSuite) Test_2_Fail_InvalidContent() {
	s.mockService.EXPECT().
		GetFeedbackByID(gomock.Any(), gomock.Any()).
		Return(&domain.Feedback{ID: uuid.New()}, nil)

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.CreateFeedbackReply(ctx, uuid.New(), uuid.New(), dto.CreateFeedbackReplyRequest{
		Content: "",
	})

	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), resp)
	assert.Equal(s.T(), apperrors.Feedback.InvalidContent, err)
}

//
// END OF CASES
//

func TestCreateFeedbackReplyTestSuite(t *testing.T) {
	suite.Run(t, new(createFeedbackReplyTestSuite))
}
