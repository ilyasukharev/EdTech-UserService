package user

import (
	"UserService/internal/model"
	"UserService/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"time"
)

type CreateUser struct {
	FirstName     string  `json:"first_name" validate:"required" binding:"required"`
	LastName      *string `json:"last_name"`
	MiddleName    *string `json:"middle_name"`
	Email         string  `json:"email" validate:"required,email" binding:"required"`
	Phone         *string `json:"phone"`
	Notifications bool    `json:"notifications" default:"true" binding:"required"`
	Type          string  `json:"type" validate:"user_type" binding:"required"`
}

func CheckTypeValid(fl validator.FieldLevel) bool {
	v := fl.Field().String()
	if model.IsValidType(v) {
		return true
	}
	return false
}

func (c *CreateUser) ToUser() *model.User {
	ID := uuid.New()
	return &model.User{
		ID:            &ID,
		FirstName:     &c.FirstName,
		LastName:      c.LastName,
		MiddleName:    c.MiddleName,
		Email:         &c.Email,
		Phone:         c.Phone,
		Notifications: &c.Notifications,
		Type:          &c.Type,
		CreatedAt:     utils.TimePtr(time.Now()),
	}
}

type UpdateUser struct {
	FirstName     string    `json:"first_name" binding:"required"`
	LastName      string    `json:"last_name" binding:"required"`
	MiddleName    string    `json:"middle_name" binding:"required"`
	Email         string    `json:"email" binding:"required"`
	Phone         string    `json:"phone" binding:"required"`
	Notifications bool      `json:"notifications" binding:"required"`
	Type          string    `json:"type" binding:"required"`
	CreatedAt     time.Time `json:"created_at" binding:"required"`
}

func (c *UpdateUser) ToUser(ID uuid.UUID) *model.User {
	return &model.User{
		ID:            &ID,
		FirstName:     &c.FirstName,
		LastName:      &c.LastName,
		MiddleName:    &c.MiddleName,
		Email:         &c.Email,
		Phone:         &c.Phone,
		Notifications: &c.Notifications,
		Type:          &c.Type,
		CreatedAt:     &c.CreatedAt,
		UpdatedAt:     utils.TimePtr(time.Now()),
	}
}

type PatchUser struct {
	FirstName     *string    `json:"first_name"`
	LastName      *string    `json:"last_name"`
	MiddleName    *string    `json:"middle_name"`
	Email         *string    `json:"email"`
	Phone         *string    `json:"phone"`
	Notifications *bool      `json:"notifications"`
	Type          *string    `json:"type"`
	CreatedAt     *time.Time `json:"created_at"`
}

func (c *PatchUser) ToUser(ID uuid.UUID) *model.User {
	return &model.User{
		ID:            &ID,
		FirstName:     c.FirstName,
		LastName:      c.LastName,
		MiddleName:    c.MiddleName,
		Email:         c.Email,
		Phone:         c.Phone,
		Notifications: c.Notifications,
		Type:          c.Type,
		CreatedAt:     c.CreatedAt,
		UpdatedAt:     utils.TimePtr(time.Now()),
	}
}

type UserResponse struct {
	ID            uuid.UUID  `json:"id" binding:"required"`
	FirstName     string     `json:"first_name" binding:"required"`
	LastName      *string    `json:"last_name"`
	MiddleName    *string    `json:"middle_name"`
	Email         string     `json:"email" binding:"required"`
	Phone         *string    `json:"phone"`
	Notifications bool       `json:"notifications" binding:"required"`
	Type          string     `json:"type" binding:"required"`
	CreatedAt     time.Time  `json:"created_at" binding:"required"`
	UpdatedAt     *time.Time `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at"`
}
