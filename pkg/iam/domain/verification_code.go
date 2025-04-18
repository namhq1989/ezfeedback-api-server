package domain

import (
	"time"

	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/internal/utils/manipulation"
	"github.com/namhq1989/ezfeedback-api-server/internal/utils/validation"
	"github.com/namhq1989/go-utilities/uuid"
)

var (
	otpDigits           = 6
	verificationCodeTTL = 15 * time.Minute
)

type VerificationCode struct {
	ID        string
	Code      string
	Ip        string
	Email     string
	IsUsed    bool
	ExpiresAt time.Time
	CreatedAt time.Time
}

func NewVerificationCode(ip, email string) (*VerificationCode, error) {
	var (
		now = manipulation.NowUTC()
	)

	var c = &VerificationCode{
		ID:        uuid.New(),
		ExpiresAt: now.Add(verificationCodeTTL),
		IsUsed:    false,
		CreatedAt: now,
	}

	if err := c.SetIp(ip); err != nil {
		return nil, err
	}
	if err := c.SetEmail(email); err != nil {
		return nil, err
	}

	c.randomCode()

	return c, nil
}

func (c *VerificationCode) SetIp(ip string) error {
	if ip == "" {
		return apperrors.Common.InvalidIp
	}

	c.Ip = ip
	return nil
}

func (c *VerificationCode) SetEmail(email string) error {
	if !validation.IsValidEmail(email) {
		return apperrors.Common.InvalidEmail
	}

	c.Email = email
	return nil
}

func (c *VerificationCode) randomCode() {
	c.Code = manipulation.GenerateOTP(otpDigits)
}

func (c *VerificationCode) MarkAsUsed() {
	c.IsUsed = true
}

func (c *VerificationCode) IsValidCode(code string) bool {
	return c.Code == code
}

func (c *VerificationCode) IsExpired() bool {
	return c.ExpiresAt.Before(manipulation.NowUTC())
}
