package child

import (
	"UserService/internal/model"
	"UserService/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"time"
)

type CreateChild struct {
	ParentID uuid.UUID   `json:"parent_id" binding:"required"`
	Name     string      `json:"name" binding:"required"`
	Age      int         `json:"age" binding:"required"`
	Gender   *string     `json:"gender" validate:"omitempty,gender"`
	Birthday *model.Date `json:"birthday"`
}

func CheckGenderValid(fl validator.FieldLevel) bool {
	return model.IsValidGender(fl.Field().String())
}

func (c *CreateChild) ToChild() *model.Child {
	return &model.Child{
		ParentID: &c.ParentID,
		Name:     &c.Name,
		Age:      &c.Age,
		Gender:   c.Gender,
		Birthday: c.Birthday.AsTime(),
	}
}

type UpdateChild struct {
	ParentID  uuid.UUID  `json:"parent_id" binding:"required"`
	Name      string     `json:"name" binding:"required"`
	Age       int        `json:"age" binding:"required"`
	Gender    string     `json:"gender" validate:"gender" binding:"required"`
	Birthday  model.Date `json:"birthday" binding:"required"`
	CreatedAt time.Time  `json:"created_at" binding:"required"`
}

func (c *UpdateChild) ToChild(ID uuid.UUID) *model.Child {
	return &model.Child{
		ID:        &ID,
		ParentID:  &c.ParentID,
		Name:      &c.Name,
		Age:       &c.Age,
		Gender:    &c.Gender,
		Birthday:  c.Birthday.AsTime(),
		CreatedAt: &c.CreatedAt,
		UpdatedAt: utils.TimePtr(time.Now()),
	}
}

type PatchChild struct {
	ParentID  *uuid.UUID  `json:"parent_id"`
	Name      *string     `json:"name"`
	Age       *int        `json:"age"`
	Gender    *string     `json:"gender"  validate:"omitempty,gender"`
	Birthday  *model.Date `json:"birthday"`
	CreatedAt *time.Time  `json:"created_at"`
}

func (c *PatchChild) ToChild(ID uuid.UUID) *model.Child {
	return &model.Child{
		ID:        &ID,
		ParentID:  c.ParentID,
		Name:      c.Name,
		Age:       c.Age,
		Gender:    c.Gender,
		Birthday:  c.Birthday.AsTime(),
		CreatedAt: c.CreatedAt,
		UpdatedAt: utils.TimePtr(time.Now()),
	}
}

type ChildResponse struct {
	ID        uuid.UUID  `json:"id" binding:"required"`
	ParentID  uuid.UUID  `json:"parent_id" binding:"required"`
	Name      string     `json:"name" binding:"required"`
	Age       int        `json:"age" binding:"required"`
	Gender    *string    `json:"gender"`
	Birthday  *string    `json:"birthday"`
	CreatedAt time.Time  `json:"created_at" binding:"required"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
