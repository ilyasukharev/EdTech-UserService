package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"time"
)

const (
	registrationIDPrefix = "user-registration-id-"
	registrationIDTTL    = time.Minute * 60

	verificationCodePrefix = "user-verification-code-"
	verificationCodeTTL    = time.Minute
)

type UserRedisRepository struct {
	Client *redis.Client
}

func NewRedisUserRepository(client *redis.Client) *UserRedisRepository {
	return &UserRedisRepository{Client: client}
}

func (r *UserRedisRepository) SaveVerificationCode(ctx context.Context, email string, code string) error {
	return r.Client.Set(ctx, verificationCodePrefix+email, code, verificationCodeTTL).Err()
}

func (r *UserRedisRepository) GetVerificationCode(ctx context.Context, email string) (string, error) {
	return r.Client.GetDel(ctx, verificationCodePrefix+email).Result()
}

func (r *UserRedisRepository) GetEmailByRegistrationID(ctx context.Context, regID uuid.UUID) (string, error) {
	return r.Client.GetDel(ctx, registrationIDPrefix+regID.String()).Result()
}

func (r *UserRedisRepository) SaveRegistrationID(ctx context.Context, regID uuid.UUID, email string) error {
	return r.Client.Set(ctx, registrationIDPrefix+regID.String(), email, registrationIDTTL).Err()
}
