package service

import (
	"UserService/internal/errors"
	"UserService/internal/model"
	"UserService/internal/model/db"
	"UserService/internal/repository"
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserService struct {
	Repo      *repository.UserRepository
	RedisRepo *repository.UserRedisRepository
	TxManager *db.TxManager
}

func NewUserService(repo *repository.UserRepository, redisRepo *repository.UserRedisRepository, txManager *db.TxManager) *UserService {
	return &UserService{
		Repo:      repo,
		RedisRepo: redisRepo,
		TxManager: txManager,
	}
}

func (s *UserService) CreateUser(ctx context.Context, userModel *model.User, regID uuid.UUID) (*model.User, error) {
	var user *model.User

	var err error
	var email string
	email, err = s.RedisRepo.GetEmailByRegistrationID(ctx, regID)
	if err != nil {
		return nil, errors.RegistrationIDNotFoundErr
	}

	if email != *userModel.Email {
		return nil, errors.RegistrationEmailMismatchErr
	}

	err = s.TxManager.RunReadCommited(
		ctx, func(tx *sqlx.Tx) error {
			user, err = s.Repo.Create(ctx, tx, userModel)
			return err
		})
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetUserById(ctx context.Context, ID uuid.UUID) (*model.User, error) {
	var user *model.User

	err := s.TxManager.RunReadCommited(ctx, func(tx *sqlx.Tx) error {
		var err error
		user, err = s.Repo.GetByID(ctx, tx, ID)
		return err
	})

	if err != nil && err == sql.ErrNoRows {
		return nil, fmt.Errorf("%w: %v", errors.UserNotFoundErr, err)
	} else if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user *model.User

	err := s.TxManager.RunReadCommited(ctx, func(tx *sqlx.Tx) error {
		var err error
		user, err = s.Repo.GetByEmail(ctx, tx, email)
		return err
	})

	if err != nil && err == sql.ErrNoRows {
		return nil, fmt.Errorf("%w: %v", errors.UserNotFoundErr, err)
	} else if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, userModel *model.User) (*model.User, error) {
	var user *model.User

	err := s.TxManager.RunReadCommited(ctx, func(tx *sqlx.Tx) error {
		var err error
		user, err = s.Repo.Update(ctx, tx, userModel)
		return err
	})

	if err != nil && err == sql.ErrNoRows {
		return nil, fmt.Errorf("%w: %v", errors.UserNotFoundErr, err)
	} else if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) DeleteUser(ctx context.Context, ID uuid.UUID) (*model.User, error) {
	var user *model.User
	err := s.TxManager.RunReadCommited(ctx, func(tx *sqlx.Tx) error {
		var err error
		user, err = s.Repo.Delete(ctx, tx, ID)
		return err
	})

	if err != nil && err == sql.ErrNoRows {
		return nil, fmt.Errorf("%w: %v", errors.UserNotFoundErr, err)
	} else if err != nil {
		return nil, err
	}

	return user, nil
}
