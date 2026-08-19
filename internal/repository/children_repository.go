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

type ChildrenRepository struct {
	DB *sqlx.DB
}

func NewChildrenRepository(db *sqlx.DB) *ChildrenRepository {
	return &ChildrenRepository{
		DB: db,
	}
}

func (r *ChildrenRepository) Create(
	ctx context.Context,
	tx *sqlx.Tx,
	childModel *model.Child,
) (*model.Child, error) {
	query := `
			INSERT INTO children (parent_id, name, age, gender, birthday) 
			VALUES ($1, $2, $3, $4, $5)
			RETURNING *
	`

	return makeChildQueryRowQuery(ctx, tx, query, []any{
		childModel.ParentID, childModel.Name, childModel.Age, childModel.Gender, childModel.Birthday})
}

func (r *ChildrenRepository) GetByID(
	ctx context.Context,
	tx *sqlx.Tx,
	ID uuid.UUID,
) (*model.Child, error) {
	query := "SELECT * FROM children WHERE id = $1"

	return makeChildQueryRowQuery(ctx, tx, query, []any{ID})
}

func (r *ChildrenRepository) GetByParentID(
	ctx context.Context,
	tx *sqlx.Tx,
	parentID uuid.UUID,
) (*model.Child, error) {
	query := "SELECT * FROM children WHERE parent_id = $1"

	return makeChildQueryRowQuery(ctx, tx, query, []any{parentID})
}

func (r *ChildrenRepository) Update(
	ctx context.Context,
	tx *sqlx.Tx,
	childModel *model.Child,
) (*model.Child, error) {
	fields := []model.PatchField{
		{childModel.ParentID, "parent_id"},
		{childModel.Name, "name"},
		{childModel.Age, "age"},
		{childModel.Gender, "gender"},
		{childModel.Birthday, "birthday"},
		{childModel.CreatedAt, "created_at"},
		{childModel.UpdatedAt, "updated_at"},
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

	args = append(args, childModel.ID)
	queryFormatted := "UPDATE children SET " + strings.Join(query, ", ") +
		fmt.Sprintf(" WHERE id = $%v RETURNING *", argId)
	return makeChildQueryRowQuery(ctx, tx, queryFormatted, args)
}

func (r *ChildrenRepository) Delete(ctx context.Context, tx *sqlx.Tx, ID uuid.UUID) (*model.Child, error) {
	query := `
		UPDATE children SET deleted_at = CURRENT_TIMESTAMP, 
		                 updated_at = CURRENT_TIMESTAMP WHERE ID = $1
		RETURNING *
	`

	return makeChildQueryRowQuery(ctx, tx, query, []any{ID})
}

func makeChildQueryRowQuery(ctx context.Context, tx *sqlx.Tx, query string, args []any) (*model.Child, error) {
	var child model.Child
	row := tx.QueryRowxContext(ctx, query, args...)
	if err := row.StructScan(&child); err != nil {
		return nil, err
	}
	return &child, nil
}
