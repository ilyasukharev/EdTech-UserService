package model

import (
	"github.com/google/uuid"
	"time"
)

const (
	Admin          = "ADMIN"
	ContentCreator = "CONTENT_CREATOR"
	Parent         = "PARENT"
)

type User struct {
	ID            *uuid.UUID `db:"id"`
	FirstName     *string    `db:"first_name"`
	LastName      *string    `db:"last_name"`
	MiddleName    *string    `db:"middle_name"`
	Email         *string    `db:"email"`
	Phone         *string    `db:"phone"`
	Notifications *bool      `db:"notifications"`
	Type          *string    `db:"type"`
	CreatedAt     *time.Time `db:"created_at"`
	UpdatedAt     *time.Time `db:"updated_at"`
	DeletedAt     *time.Time `db:"deleted_at"`
}

func IsValidType(val string) bool {
	if val == "" {
		return false
	}

	if val == ContentCreator || val == Parent || val == Admin {
		return true
	}

	return false
}
