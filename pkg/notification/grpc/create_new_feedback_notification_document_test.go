package grpc_test

import (
	"context"
	"testing"

	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/internal/genproto/notificationpb"
	mocknotification "github.com/namhq1989/ezfeedback-api-server/internal/mock/notification"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/grpc"
	"github.com/namhq1989/go-utilities/appcontext"
	"github.com/namhq1989/go-utilities/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type createNotificationDocumentTestSuite struct {
	suite.Suite
	handler                    grpc.CreateNewFeedbackNotificationDocumentHandler
	mockCtrl                   *gomock.Controller
	mockNotificationRepository *mocknotification.MockNotificationRepository
}

func (s *createNotificationDocumentTestSuite) SetupSuite() {
	s.setupApplication()
}

func (s *createNotificationDocumentTestSuite) setupApplication() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockNotificationRepository = mocknotification.NewMockNotificationRepository(s.mockCtrl)

	s.handler = grpc.NewCreateNewFeedbackNotificationDocumentHandler(s.mockNotificationRepository)
}

func (s *createNotificationDocumentTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

//
// CASES
//

func (s *createNotificationDocumentTestSuite) Test_1_Success() {
	// mock
	s.mockNotificationRepository.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(nil)

	// call
	ctx := appcontext.NewGRPC(context.Background())
	resp, err := s.handler.CreateNewFeedbackNotificationDocument(ctx, &notificationpb.CreateNewFeedbackNotificationDocumentRequest{
		TraceId: "trace-id",
		UserId:  uuid.New(),
		Metadata: &notificationpb.CreateNewFeedbackNotificationDocumentMetadata{
			ProjectId:    uuid.New(),
			ProjectTitle: "Project Title",
		},
	})

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
}

func (s *createNotificationDocumentTestSuite) Test_2_Fail_InvalidMetadata() {
	// call
	ctx := appcontext.NewGRPC(context.Background())
	_, err := s.handler.CreateNewFeedbackNotificationDocument(ctx, &notificationpb.CreateNewFeedbackNotificationDocumentRequest{
		TraceId: "trace-id",
		UserId:  uuid.New(),
		Metadata: &notificationpb.CreateNewFeedbackNotificationDocumentMetadata{
			ProjectId:    "invalid_id",
			ProjectTitle: "Project Title",
		},
	})

	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), apperrors.Notification.InvalidMetadata, err)
}

//
// END OF CASES
//

func TestCreateNotificationDocumentTestSuite(t *testing.T) {
	suite.Run(t, new(createNotificationDocumentTestSuite))
}
