package user

import (
	"UserService/internal/model"
	"UserService/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"time"
)

type CreateUser struct {
	FirstName     string  `json:"first_name" validate:"required"`
	LastName      *string `json:"last_name"`
	MiddleName    *string `json:"middle_name"`
	Email         string  `json:"email" validate:"required,email"`
	Phone         *string `json:"phone"`
	Notifications bool    `json:"notifications" default:"true"`
	Type          string  `json:"type" validate:"user_type"`
	ChildName     *string `json:"child_name" validate:"required_if=Type PARENT"`
	ChildAge      *int    `json:"child_age" validate:"required_if=Type PARENT,gte=1,lte=18"`
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
		ChildName:     c.ChildName,
		ChildAge:      c.ChildAge,
		CreatedAt:     utils.TimePtr(time.Now()),
	}
}

type UpdateUser struct {
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	MiddleName    string    `json:"middle_name"`
	Email         string    `json:"email"`
	Phone         string    `json:"phone"`
	Notifications bool      `json:"notifications"`
	Type          string    `json:"type"`
	ChildName     string    `json:"child_name"`
	ChildAge      int       `json:"child_age"`
	CreatedAt     time.Time `json:"created_at"`
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
		ChildName:     &c.ChildName,
		ChildAge:      &c.ChildAge,
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
	ChildName     *string    `json:"child_name"`
	ChildAge      *int       `json:"child_age"`
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
		ChildName:     c.ChildName,
		ChildAge:      c.ChildAge,
		CreatedAt:     c.CreatedAt,
		UpdatedAt:     utils.TimePtr(time.Now()),
	}
}

type UserResponse struct {
	ID            uuid.UUID  `json:"id"`
	FirstName     string     `json:"first_name"`
	LastName      *string    `json:"last_name"`
	MiddleName    *string    `json:"middle_name"`
	Email         string     `json:"email"`
	Phone         *string    `json:"phone"`
	Notifications bool       `json:"notifications"`
	Type          string     `json:"type"`
	ChildName     *string    `json:"child_name"`
	ChildAge      *int       `json:"child_age"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at"`
}
