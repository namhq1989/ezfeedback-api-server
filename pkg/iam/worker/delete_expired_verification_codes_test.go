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

type deleteExpiredVerificationCodesTestSuite struct {
	suite.Suite
	handler                        worker.DeleteExpiredVerificationCodesHandler
	mockCtrl                       *gomock.Controller
	mockVerificationCodeRepository *mockiam.MockVerificationCodeRepository
}

func (s *deleteExpiredVerificationCodesTestSuite) SetupSuite() {
	s.setupApplication()
}

func (s *deleteExpiredVerificationCodesTestSuite) setupApplication() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockVerificationCodeRepository = mockiam.NewMockVerificationCodeRepository(s.mockCtrl)

	s.handler = worker.NewDeleteExpiredVerificationCodesHandler(s.mockVerificationCodeRepository)
}

func (s *deleteExpiredVerificationCodesTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

//
// CASES
//

func (s *deleteExpiredVerificationCodesTestSuite) Test_1_Success() {
	// mock
	s.mockVerificationCodeRepository.EXPECT().
		DeleteExpired(gomock.Any()).
		Return(nil)

	// call
	ctx := appcontext.NewGRPC(context.Background())
	err := s.handler.DeleteExpiredVerificationCodes(ctx, domain.QueueDeleteExpiredVerificationCodesPayload{})

	assert.Nil(s.T(), err)
}

//
// END OF CASES
//

func TestDeleteExpiredVerificationCodesTestSuite(t *testing.T) {
	suite.Run(t, new(deleteExpiredVerificationCodesTestSuite))
}
