package service

import (
	"UserService/internal/errors"
	"UserService/internal/model"
	"UserService/internal/model/db"
	"UserService/internal/repository"
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ReferralService struct {
	Repo      *repository.ReferralsRepository
	TxManager *db.TxManager
}

func NewReferralService(txManager *db.TxManager, repo *repository.ReferralsRepository) *ReferralService {
	return &ReferralService{
		TxManager: txManager,
		Repo:      repo,
	}
}

func (s *ReferralService) CreateReferral(ctx context.Context, referralModel *model.Referral) (*model.Referral, error) {
	var referral *model.Referral

	err := s.TxManager.RunReadCommited(ctx, func(tx *sqlx.Tx) error {
		var err error

		referral, err = s.Repo.Create(ctx, tx, referralModel)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return referral, nil
}

func (s *ReferralService) GetByReferrerID(ctx context.Context, referrerID uuid.UUID) (*model.Referral, error) {
	var referral *model.Referral

	err := s.TxManager.RunReadCommited(ctx, func(tx *sqlx.Tx) error {
		var err error
		referral, err = s.Repo.GetByReferrerID(ctx, tx, referrerID)
		return err
	})

	if err != nil && err == sql.ErrNoRows {
		return nil, errors.ReferralNotFoundErr
	} else if err != nil {
		return nil, err
	}

	return referral, nil
}

func (s *ReferralService) GetByRefereeID(ctx context.Context, refereeID uuid.UUID) (*model.Referral, error) {
	var referral *model.Referral

	err := s.TxManager.RunReadCommited(ctx, func(tx *sqlx.Tx) error {
		var err error
		referral, err = s.Repo.GetByRefereeID(ctx, tx, refereeID)
		return err
	})

	if err != nil && err == sql.ErrNoRows {
		return nil, errors.ReferralNotFoundErr
	} else if err != nil {
		return nil, err
	}

	return referral, nil
}

func (s *ReferralService) PatchReferral(ctx context.Context, referralModel *model.Referral) (*model.Referral, error) {
	var referral *model.Referral

	err := s.TxManager.RunReadCommited(ctx, func(tx *sqlx.Tx) error {
		var err error
		referral, err = s.Repo.Update(ctx, tx, referralModel)
		return err
	})

	if err != nil && err == sql.ErrNoRows {
		return nil, errors.ReferralNotFoundErr
	} else if err != nil {
		return nil, err
	}

	return referral, nil
}
