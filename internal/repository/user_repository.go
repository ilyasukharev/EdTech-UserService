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

type UserRepository struct {
	DB *sqlx.DB
}

func NewUserRepository(database *sqlx.DB) *UserRepository {
	return &UserRepository{
		DB: database,
	}
}

func (r *UserRepository) Create(ctx context.Context, tx *sqlx.Tx, user *model.User) (*model.User, error) {
	query := `
	INSERT INTO users (id, first_name, last_name, middle_name, email, phone, 
	                   notifications, type, child_name, child_age, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) 
	RETURNING *
	`

	row := tx.QueryRowxContext(ctx, query,
		user.ID, user.FirstName, user.LastName, user.MiddleName,
		user.Email, user.Phone, user.Notifications, user.Type,
		user.ChildName, user.ChildAge, user.CreatedAt)

	var result model.User
	if err := row.StructScan(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *UserRepository) GetByID(ctx context.Context, tx *sqlx.Tx, ID uuid.UUID) (*model.User, error) {
	query := `
	SELECT * FROM users WHERE ID = $1
	`

	var user model.User
	if err := tx.GetContext(ctx, &user, query, ID); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, tx *sqlx.Tx, email string) (*model.User, error) {
	query := `
	SELECT * FROM users WHERE email = $1
	`

	var user model.User
	if err := tx.GetContext(ctx, &user, query, email); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, tx *sqlx.Tx, userModel *model.User) (*model.User, error) {
	fields := []model.PatchField{
		{userModel.ID, "id"},
		{userModel.FirstName, "first_name"},
		{userModel.LastName, "last_name"},
		{userModel.MiddleName, "middle_name"},
		{userModel.Email, "email"},
		{userModel.Phone, "phone"},
		{userModel.Notifications, "notifications"},
		{userModel.Type, "type"},
		{userModel.ChildName, "child_name"},
		{userModel.ChildAge, "child_age"},
		{userModel.CreatedAt, "created_at"},
		{userModel.UpdatedAt, "updated_at"},
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

	args = append(args, userModel.ID)
	row := tx.QueryRowxContext(ctx, "UPDATE users SET "+strings.Join(query, ", ")+
		fmt.Sprintf(" WHERE id = $%v RETURNING *", argId), args...)

	var user model.User
	if err := row.StructScan(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Delete(ctx context.Context, tx *sqlx.Tx, ID uuid.UUID) (*model.User, error) {
	query := `
		UPDATE users SET deleted_at = CURRENT_TIMESTAMP, 
		                 updated_at = CURRENT_TIMESTAMP WHERE ID = $1
		RETURNING *
	`
	row := tx.QueryRowxContext(ctx, query, ID)

	var user model.User
	if err := row.StructScan(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
