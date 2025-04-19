package mapping

import (
	"github.com/namhq1989/ezfeedback-api-server/internal/database/gen/ezfeedback/public/model"
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/domain"
)

type VerificationCodeMapper struct{}

func (VerificationCodeMapper) FromModelToDomain(verificationCode model.VerificationCodes) (*domain.VerificationCode, error) {
	var result = &domain.VerificationCode{
		ID:        verificationCode.ID,
		Email:     verificationCode.Email,
		Code:      verificationCode.Code,
		Ip:        verificationCode.IP,
		ExpiresAt: verificationCode.ExpiresAt,
		IsUsed:    verificationCode.IsUsed,
		CreatedAt: verificationCode.CreatedAt,
	}

	return result, nil
}

func (VerificationCodeMapper) FromDomainToModel(verificationCode domain.VerificationCode) (*model.VerificationCodes, error) {
	var result = &model.VerificationCodes{
		ID:        verificationCode.ID,
		Email:     verificationCode.Email,
		Code:      verificationCode.Code,
		IP:        verificationCode.Ip,
		ExpiresAt: verificationCode.ExpiresAt,
		IsUsed:    verificationCode.IsUsed,
		CreatedAt: verificationCode.CreatedAt,
	}

	return result, nil
}
