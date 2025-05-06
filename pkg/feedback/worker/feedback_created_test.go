package worker_test

import (
	"context"
	"testing"

	mockfeedback "github.com/namhq1989/ezfeedback-api-server/internal/mock/feedback"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/worker"
	"github.com/namhq1989/go-utilities/appcontext"
	"github.com/namhq1989/go-utilities/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type feedbackCreatedTestSuite struct {
	suite.Suite
	handler             worker.FeedbackCreatedHandler
	mockCtrl            *gomock.Controller
	mockProjectHub      *mockfeedback.MockProjectHub
	mockNotificationHub *mockfeedback.MockNotificationHub
}

func (s *feedbackCreatedTestSuite) SetupSuite() {
	s.setupApplication()
}

func (s *feedbackCreatedTestSuite) setupApplication() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockProjectHub = mockfeedback.NewMockProjectHub(s.mockCtrl)
	s.mockNotificationHub = mockfeedback.NewMockNotificationHub(s.mockCtrl)

	s.handler = worker.NewFeedbackCreatedHandler(s.mockProjectHub, s.mockNotificationHub)
}

func (s *feedbackCreatedTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

//
// CASES
//

func (s *feedbackCreatedTestSuite) Test_1_Success() {
	// mock
	s.mockProjectHub.EXPECT().
		GetProjectByID(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&domain.Project{
			ID:     uuid.New(),
			UserID: uuid.New(),
			Title:  "ProjectTitle",
		}, nil)

	s.mockNotificationHub.EXPECT().
		CreateNewFeedbackNotificationDocument(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil)

	s.mockNotificationHub.EXPECT().
		CreateNotificationReminder(gomock.Any(), gomock.Any()).
		Return(nil)

	// call
	ctx := appcontext.NewGRPC(context.Background())
	err := s.handler.FeedbackCreated(ctx, domain.QueueFeedbackCreatedPayload{
		Feedback: domain.Feedback{
			ID: uuid.New(),
		},
	})

	assert.Nil(s.T(), err)
}

//
// END OF CASES
//

func TestFeedbackCreatedTestSuite(t *testing.T) {
	suite.Run(t, new(feedbackCreatedTestSuite))
}
