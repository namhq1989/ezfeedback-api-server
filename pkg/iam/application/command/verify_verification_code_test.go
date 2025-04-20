package command_test

import (
	"context"
	"testing"
	"time"

	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/dto"

	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	mockiam "github.com/namhq1989/ezfeedback-api-server/internal/mock/iam"
	"github.com/namhq1989/ezfeedback-api-server/internal/utils/manipulation"
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/application/command"
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/domain"
	"github.com/namhq1989/go-utilities/appcontext"
	"github.com/namhq1989/go-utilities/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type verifyVerificationCodeTestSuite struct {
	suite.Suite
	handler                        command.VerifyVerificationCodeHandler
	mockCtrl                       *gomock.Controller
	mockUserRepository             *mockiam.MockUserRepository
	mockVerificationCodeRepository *mockiam.MockVerificationCodeRepository
	mockJwtRepository              *mockiam.MockJwtRepository
}

func (s *verifyVerificationCodeTestSuite) SetupSuite() {
	s.setupApplication()
}

func (s *verifyVerificationCodeTestSuite) setupApplication() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockUserRepository = mockiam.NewMockUserRepository(s.mockCtrl)
	s.mockVerificationCodeRepository = mockiam.NewMockVerificationCodeRepository(s.mockCtrl)
	s.mockJwtRepository = mockiam.NewMockJwtRepository(s.mockCtrl)

	s.handler = command.NewVerifyVerificationCodeHandler(s.mockUserRepository, s.mockVerificationCodeRepository, s.mockJwtRepository)
}

func (s *verifyVerificationCodeTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

//
// CASES
//

func (s *verifyVerificationCodeTestSuite) Test_1_Success_NewUser() {
	// mock data
	var (
		code = "code"
	)

	s.mockVerificationCodeRepository.EXPECT().
		Find(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&domain.VerificationCode{
			ID:        uuid.New(),
			Code:      code,
			ExpiresAt: manipulation.NowUTC().Add(1 * time.Hour),
		}, nil)

	s.mockUserRepository.EXPECT().
		FindByEmail(gomock.Any(), gomock.Any()).
		Return(nil, nil)

	s.mockUserRepository.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(nil)

	s.mockJwtRepository.EXPECT().
		GenerateAccessToken(gomock.Any(), gomock.Any()).
		Return("token", nil)

	s.mockVerificationCodeRepository.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		Return(nil)

	// call
	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.VerifyVerificationCode(ctx, "127.0.0.1", dto.VerifyVerificationCodeRequest{
		Email: "john@gmail.com",
		Code:  code,
	})

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.NotEmpty(s.T(), resp.Token)
	assert.Equal(s.T(), true, resp.IsNewUser)
}

func (s *verifyVerificationCodeTestSuite) Test_1_Success_ExistedUser() {
	// mock data
	var (
		code = "code"
	)

	s.mockVerificationCodeRepository.EXPECT().
		Find(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&domain.VerificationCode{
			ID:        uuid.New(),
			Code:      code,
			ExpiresAt: manipulation.NowUTC().Add(1 * time.Hour),
		}, nil)

	s.mockUserRepository.EXPECT().
		FindByEmail(gomock.Any(), gomock.Any()).
		Return(&domain.User{
			ID: uuid.New(),
		}, nil)

	s.mockJwtRepository.EXPECT().
		GenerateAccessToken(gomock.Any(), gomock.Any()).
		Return("token", nil)

	s.mockVerificationCodeRepository.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		Return(nil)

	// call
	ctx := appcontext.NewRest(context.Background())
	resp, err := s.handler.VerifyVerificationCode(ctx, "127.0.0.1", dto.VerifyVerificationCodeRequest{
		Email: "john@gmail.com",
		Code:  code,
	})

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.NotEmpty(s.T(), resp.Token)
	assert.Equal(s.T(), false, resp.IsNewUser)
}

func (s *verifyVerificationCodeTestSuite) Test_1_Success_InvalidCode() {
	// mock data
	var (
		code = "code"
	)

	s.mockVerificationCodeRepository.EXPECT().
		Find(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, nil)

	// call
	ctx := appcontext.NewRest(context.Background())
	_, err := s.handler.VerifyVerificationCode(ctx, "127.0.0.1", dto.VerifyVerificationCodeRequest{
		Email: "john@gmail.com",
		Code:  code,
	})

	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), apperrors.User.InvalidVerificationCode, err)
}

func (s *verifyVerificationCodeTestSuite) Test_2_Fail_TokenExpired() {
	// mock data
	var (
		code = "code"
	)

	s.mockVerificationCodeRepository.EXPECT().
		Find(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&domain.VerificationCode{
			ID:        uuid.New(),
			Code:      code,
			ExpiresAt: manipulation.NowUTC().Add(-1 * time.Hour),
		}, nil)

	// call
	ctx := appcontext.NewRest(context.Background())
	_, err := s.handler.VerifyVerificationCode(ctx, "127.0.0.1", dto.VerifyVerificationCodeRequest{
		Email: "john@gmail.com",
		Code:  code,
	})

	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), apperrors.User.InvalidVerificationCode, err)
}

//
// END OF CASES
//

func TestVerifyVerificationCodeTestSuite(t *testing.T) {
	suite.Run(t, new(verifyVerificationCodeTestSuite))
}
