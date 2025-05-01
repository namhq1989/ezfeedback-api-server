package command_test

import (
	"context"
	"testing"

	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	mockiam "github.com/namhq1989/ezfeedback-api-server/internal/mock/iam"
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/application/command"
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/dto"
	"github.com/namhq1989/go-utilities/appcontext"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type requestVerificationCodeTestSuite struct {
	suite.Suite
	handler                        command.RequestVerificationCodeHandler
	mockCtrl                       *gomock.Controller
	mockVerificationCodeRepository *mockiam.MockVerificationCodeRepository
	mockQueueRepository            *mockiam.MockQueueRepository
}

func (s *requestVerificationCodeTestSuite) SetupSuite() {
	s.setupApplication()
}

func (s *requestVerificationCodeTestSuite) setupApplication() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockVerificationCodeRepository = mockiam.NewMockVerificationCodeRepository(s.mockCtrl)
	s.mockQueueRepository = mockiam.NewMockQueueRepository(s.mockCtrl)

	s.handler = command.NewRequestVerificationCodeHandler(s.mockVerificationCodeRepository, s.mockQueueRepository)
}

func (s *requestVerificationCodeTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

//
// CASES
//

func (s *requestVerificationCodeTestSuite) Test_1_Success() {
	// mock data
	s.mockVerificationCodeRepository.EXPECT().
		CountTotalSentTodayByIp(gomock.Any(), gomock.Any()).
		Return(int64(0), nil)

	s.mockVerificationCodeRepository.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(nil)

	s.mockQueueRepository.EXPECT().
		SendVerificationCodeEmail(gomock.Any(), gomock.Any()).
		Return(nil)

	// call
	ctx := appcontext.NewRest(context.Background())
	_, err := s.handler.RequestVerificationCode(ctx, "127.0.0.1", dto.RequestVerificationCodeRequest{
		Email: "john@gmail.com",
	})

	assert.Nil(s.T(), err)
}

func (s *requestVerificationCodeTestSuite) Test_2_Fail_ExceededDailyLimit() {
	// mock data
	s.mockVerificationCodeRepository.EXPECT().
		CountTotalSentTodayByIp(gomock.Any(), gomock.Any()).
		Return(int64(100), nil)

	// call
	ctx := appcontext.NewRest(context.Background())
	_, err := s.handler.RequestVerificationCode(ctx, "127.0.0.1", dto.RequestVerificationCodeRequest{
		Email: "john@gmail.com",
	})

	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), apperrors.User.DailyIpOtpLimitExceeded, err)
}

func (s *requestVerificationCodeTestSuite) Test_2_Fail_InvalidEmail() {
	// mock data
	s.mockVerificationCodeRepository.EXPECT().
		CountTotalSentTodayByIp(gomock.Any(), gomock.Any()).
		Return(int64(0), nil)

	// call
	ctx := appcontext.NewRest(context.Background())
	_, err := s.handler.RequestVerificationCode(ctx, "127.0.0.1", dto.RequestVerificationCodeRequest{
		Email: "invalid-email",
	})

	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), apperrors.Common.InvalidEmail, err)
}

//
// END OF CASES
//

func TestRequestVerificationCodeTestSuite(t *testing.T) {
	suite.Run(t, new(requestVerificationCodeTestSuite))
}
