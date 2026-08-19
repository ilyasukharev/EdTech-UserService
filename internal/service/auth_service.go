package service

import (
	"UserService/internal/config"
	"UserService/internal/errors"
	"UserService/internal/repository"
	"UserService/internal/utils"
	"context"
	"github.com/google/uuid"
)

const CodeNumbersCount = 6

type AuthService struct {
	Repo *repository.UserRedisRepository
}

func NewAuthService(repo *repository.UserRedisRepository) *AuthService {
	return &AuthService{
		Repo: repo,
	}
}

func (s *AuthService) SendCode(ctx context.Context, email string) error {
	var verifyCode string
	if config.IsDevStand() {
		verifyCode = "111111"
	} else {
		verifyCode = utils.GenerateRandNumbersAsString(CodeNumbersCount, 10)
	}

	err := s.Repo.SaveVerificationCode(ctx, email, verifyCode)
	if err != nil {
		return err
	}

	//send code to email

	return nil
}

func (s *AuthService) VerifyAuthCode(ctx context.Context, email string, code string) (uuid.UUID, error) {
	var registrationID uuid.UUID

	actualCode, err := s.Repo.GetVerificationCode(ctx, email)
	if err != nil {
		return registrationID, err
	} else if actualCode != code {
		return registrationID, errors.OTPCodeMismatchErr
	}

	registrationID = uuid.New()
	err = s.Repo.SaveRegistrationID(ctx, registrationID, email)
	return registrationID, err
}
