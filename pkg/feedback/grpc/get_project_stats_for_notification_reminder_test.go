package grpc_test

import (
	"context"
	"testing"

	"github.com/namhq1989/ezfeedback-api-server/internal/genproto/feedbackpb"
	mockfeedback "github.com/namhq1989/ezfeedback-api-server/internal/mock/feedback"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/grpc"
	"github.com/namhq1989/go-utilities/appcontext"
	"github.com/namhq1989/go-utilities/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type getProjectStatsForNotificationReminderTestSuite struct {
	suite.Suite
	handler         grpc.GetProjectStatsForNotificationReminderHandler
	mockCtrl        *gomock.Controller
	mockFeedbackHub *mockfeedback.MockFeedbackHub
	mockProjectHub  *mockfeedback.MockProjectHub
}

func (s *getProjectStatsForNotificationReminderTestSuite) SetupSuite() {
	s.setupApplication()
}

func (s *getProjectStatsForNotificationReminderTestSuite) setupApplication() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockFeedbackHub = mockfeedback.NewMockFeedbackHub(s.mockCtrl)
	s.mockProjectHub = mockfeedback.NewMockProjectHub(s.mockCtrl)

	s.handler = grpc.NewGetProjectStatsForNotificationReminderHandler(s.mockFeedbackHub, s.mockProjectHub)
}

func (s *getProjectStatsForNotificationReminderTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

//
// CASES
//

func (s *getProjectStatsForNotificationReminderTestSuite) Test_1_Success() {
	// mock
	s.mockFeedbackHub.EXPECT().
		CountFeedbackForProjectSinceTimestamp(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(int64(5), nil)

	s.mockFeedbackHub.EXPECT().
		FindFeedbackForProjectSinceTimestamp(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]domain.Feedback{
			{ID: uuid.New()},
		}, nil)

	s.mockProjectHub.EXPECT().
		GetProjectByID(gomock.Any(), gomock.Any()).
		Return(&domain.Project{ID: uuid.New()}, nil).
		AnyTimes()

	// call
	ctx := appcontext.NewGRPC(context.Background())
	resp, err := s.handler.GetProjectStatsForNotificationReminder(ctx, &feedbackpb.GetProjectStatsForNotificationReminderRequest{
		TraceId:   "trace-id",
		ProjectId: uuid.New(),
		Timestamp: timestamppb.Now(),
		Limit:     5,
	})

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
}

//
// END OF CASES
//

func TestGetProjectStatsForNotificationReminderTestSuite(t *testing.T) {
	suite.Run(t, new(getProjectStatsForNotificationReminderTestSuite))
}
