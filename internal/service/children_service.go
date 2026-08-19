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

type ChildrenService struct {
	TxManager *db.TxManager
	Repo      *repository.ChildrenRepository
}

func NewChildrenService(txManager *db.TxManager, repo *repository.ChildrenRepository) *ChildrenService {
	return &ChildrenService{
		TxManager: txManager,
		Repo:      repo,
	}
}

func (s *ChildrenService) Create(ctx context.Context, newChild *model.Child) (*model.Child, error) {
	var child *model.Child

	err := s.TxManager.RunReadCommited(ctx, func(tx *sqlx.Tx) error {
		var err error
		child, err = s.Repo.Create(ctx, tx, newChild)
		return err
	})
	if err != nil {
		return nil, err
	}

	return child, nil
}

func (s *ChildrenService) GetByID(ctx context.Context, ID uuid.UUID) (*model.Child, error) {
	var child *model.Child

	err := s.TxManager.RunReadCommited(ctx, func(tx *sqlx.Tx) error {
		var err error
		child, err = s.Repo.GetByID(ctx, tx, ID)
		return err
	})

	if err != nil && err == sql.ErrNoRows {
		return nil, fmt.Errorf("%w: %v", errors.ChildNotFoundErr, err)
	} else if err != nil {
		return nil, err
	}

	return child, nil
}

func (s *ChildrenService) GetByParentID(ctx context.Context, parentID uuid.UUID) (*model.Child, error) {
	var child *model.Child

	err := s.TxManager.RunReadCommited(ctx, func(tx *sqlx.Tx) error {
		var err error
		child, err = s.Repo.GetByParentID(ctx, tx, parentID)
		return err
	})

	if err != nil && err == sql.ErrNoRows {
		return nil, fmt.Errorf("%w: %v", errors.ChildNotFoundErr, err)
	} else if err != nil {
		return nil, err
	}

	return child, nil
}

func (s *ChildrenService) Update(ctx context.Context, childModel *model.Child) (*model.Child, error) {
	var child *model.Child

	err := s.TxManager.RunReadCommited(ctx, func(tx *sqlx.Tx) error {
		var err error
		child, err = s.Repo.Update(ctx, tx, childModel)
		return err
	})

	if err != nil && err == sql.ErrNoRows {
		return nil, fmt.Errorf("%w: %v", errors.ChildNotFoundErr, err)
	} else if err != nil {
		return nil, err
	}

	return child, nil
}

func (s *ChildrenService) Delete(ctx context.Context, ID uuid.UUID) (*model.Child, error) {
	var child *model.Child

	err := s.TxManager.RunReadCommited(ctx, func(tx *sqlx.Tx) error {
		var err error
		child, err = s.Repo.Delete(ctx, tx, ID)
		return err
	})

	if err != nil && err == sql.ErrNoRows {
		return nil, fmt.Errorf("%w: %v", errors.ChildNotFoundErr, err)
	} else if err != nil {
		return nil, err
	}

	return child, nil
}
