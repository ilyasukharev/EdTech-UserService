package repository

import (
	"UserService/internal/errors"
	"UserService/internal/model"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"reflect"
	"strings"
)

type ReferralsRepository struct {
	DB *sqlx.DB
}

func NewReferralsRepository(db *sqlx.DB) *ReferralsRepository {
	return &ReferralsRepository{
		DB: db,
	}
}

func (r ReferralsRepository) Create(ctx context.Context, tx *sqlx.Tx, referralModel *model.Referral) (*model.Referral, error) {
	query := `
		INSERT INTO referrals (referrer_id, referee_id)
		values ($1, $2)
		RETURNING *
	`

	return makeRefQueryxRowContext(ctx, tx, query, []any{referralModel.ReferrerID, referralModel.RefereeID})
}

func (r ReferralsRepository) GetByID(ctx context.Context, tx *sqlx.Tx, ID int64) (*model.Referral, error) {
	query := `
		SELECT * FROM REFERRALS WHERE id=$1
	`
	return makeRefQueryxRowContext(ctx, tx, query, []any{ID})
}

func (r ReferralsRepository) GetByReferrerID(ctx context.Context, tx *sqlx.Tx, referrerID uuid.UUID) (*model.Referral, error) {
	query := `
		SELECT * FROM REFERRALS WHERE referrer_id=$1
	`
	return makeRefQueryxRowContext(ctx, tx, query, []any{referrerID})
}

func (r ReferralsRepository) GetByRefereeID(ctx context.Context, tx *sqlx.Tx, refereeID uuid.UUID) (*model.Referral, error) {
	query := `
		SELECT * FROM REFERRALS WHERE referee_id=$1
	`
	return makeRefQueryxRowContext(ctx, tx, query, []any{refereeID})
}

func (r ReferralsRepository) Update(ctx context.Context, tx *sqlx.Tx, referralModel *model.Referral) (*model.Referral, error) {
	fields := []model.PatchField{
		{referralModel.ReferrerID, "referrer_id"},
		{referralModel.RefereeID, "referee_id"},
		{referralModel.Confirmed, "confirmed"},
		{referralModel.CreatedAt, "created_at"},
		{referralModel.UpdatedAt, "updated_at"},
	}

	var query []string
	var args []any
	argId := 1
	for _, f := range fields {
		v := reflect.ValueOf(f.Value)
		if v.Kind() == reflect.Ptr && !v.IsNil() {
			query = append(query, fmt.Sprintf("%s=$%d", f.Name, argId))
			args = append(args, f.Value)
			argId++
		}
	}

	if len(query) == 0 {
		return nil, errors.NothingToUpdateErr
	}

	args = append(args, referralModel.ID)
	queryFormatted := "UPDATE referrals SET " + strings.Join(query, ", ") +
		fmt.Sprintf(" WHERE id = $%v RETURNING *", argId)
	return makeRefQueryxRowContext(ctx, tx, queryFormatted, args)
}

func makeRefQueryxRowContext(ctx context.Context, tx *sqlx.Tx, query string, args []any) (*model.Referral, error) {
	var referral model.Referral
	row := tx.QueryRowxContext(ctx, query, args...)
	if err := row.StructScan(&referral); err != nil {
		return nil, err
	}
	return &referral, nil
}
