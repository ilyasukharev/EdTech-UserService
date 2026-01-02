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

type RedisUserRepository struct {
	Client *redis.Client
}

func NewRedisUserRepository(client *redis.Client) *RedisUserRepository {
	return &RedisUserRepository{Client: client}
}

func (r *RedisUserRepository) SaveVerificationCode(ctx context.Context, email string, code string) error {
	return r.Client.Set(ctx, verificationCodePrefix+email, code, verificationCodeTTL).Err()
}

func (r *RedisUserRepository) GetVerificationCode(ctx context.Context, email string) (string, error) {
	return r.Client.Get(ctx, verificationCodePrefix+email).Result()
}

func (r *RedisUserRepository) GetEmailByRegistrationID(ctx context.Context, regID uuid.UUID) (string, error) {
	return r.Client.Get(ctx, registrationIDPrefix+regID.String()).Result()
}

func (r *RedisUserRepository) SaveRegistrationID(ctx context.Context, regID uuid.UUID, email string) error {
	return r.Client.Set(ctx, registrationIDPrefix+regID.String(), email, registrationIDTTL).Err()
}
