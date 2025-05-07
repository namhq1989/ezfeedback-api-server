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

type getFeedbackRepliesTestSuite struct {
	suite.Suite
	handler                     query.GetFeedbackRepliesHandler
	mockCtrl                    *gomock.Controller
	mockFeedbackReplyRepository *mockfeedback.MockFeedbackReplyRepository
	mockIAMHub                  *mockfeedback.MockIAMHub
}

func (s *getFeedbackRepliesTestSuite) SetupSuite() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockFeedbackReplyRepository = mockfeedback.NewMockFeedbackReplyRepository(s.mockCtrl)
	s.mockIAMHub = mockfeedback.NewMockIAMHub(s.mockCtrl)
	s.handler = query.NewGetFeedbackRepliesHandler(s.mockFeedbackReplyRepository, s.mockIAMHub)
}

func (s *getFeedbackRepliesTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

//
// CASES
//

func (s *getFeedbackRepliesTestSuite) Test_1_Success() {
	s.mockFeedbackReplyRepository.EXPECT().
		FindWithFilter(gomock.Any(), gomock.Any()).
		Return([]domain.FeedbackReply{
			{ID: uuid.New()},
		}, nil)

	s.mockIAMHub.EXPECT().
		GetUserByID(gomock.Any(), gomock.Any()).
		Return(&domain.User{ID: uuid.New()}, nil).
		AnyTimes()

	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.GetFeedbackReplies(ctx, uuid.New(), uuid.New(), dto.GetFeedbackRepliesRequest{
		Page: 0,
	})

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.Equal(s.T(), 1, len(resp.Replies))
}

//
// END OF CASES
//

func TestGetFeedbackRepliesTestSuite(t *testing.T) {
	suite.Run(t, new(getFeedbackRepliesTestSuite))
}
