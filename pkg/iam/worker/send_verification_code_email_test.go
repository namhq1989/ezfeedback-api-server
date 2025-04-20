package worker_test

import (
	"context"
	"testing"

	mockiam "github.com/namhq1989/ezfeedback-api-server/internal/mock/iam"
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/worker"
	"github.com/namhq1989/go-utilities/appcontext"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type sendVerificationCodeEmailTestSuite struct {
	suite.Suite
	handler          worker.SendSignInVerificationCodeEmailHandler
	mockCtrl         *gomock.Controller
	mockMailerRepository *mockiam.MockMailerRepository
}

func (s *sendVerificationCodeEmailTestSuite) SetupSuite() {
	s.setupApplication()
}

func (s *sendVerificationCodeEmailTestSuite) setupApplication() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockMailerRepository = mockiam.NewMockMailerRepository(s.mockCtrl)

	s.handler = worker.NewSendSignInVerificationCodeEmailHandler(s.mockMailerRepository)
}

func (s *sendVerificationCodeEmailTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

//
// CASES
//

func (s *sendVerificationCodeEmailTestSuite) Test_1_Success() {
	// mock
	payload := domain.QueueSendVerificationCodeEmailPayload{
		Email: "test@example.com",
		Code:  "123456",
	}
	s.mockMailerRepository.EXPECT().
		SendVerificationCodeEmail(gomock.Any(), payload.Email, payload.Code).
		Return(nil)

	// call
	ctx := appcontext.NewGRPC(context.Background())
	err := s.handler.SendSignInVerificationCodeEmail(ctx, payload)

	// assert
	assert.Nil(s.T(), err)
}

//
// END OF CASES
//

func TestSendVerificationCodeEmailTestSuite(t *testing.T) {
	suite.Run(t, new(sendVerificationCodeEmailTestSuite))
}
