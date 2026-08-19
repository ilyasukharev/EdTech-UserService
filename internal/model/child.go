package model

import (
	"github.com/google/uuid"
	"time"
)

const (
	Male   = "MALE"
	Female = "FEMALE"
)

type Child struct {
	ID        *uuid.UUID `db:"id"`
	ParentID  *uuid.UUID `db:"parent_id"`
	Name      *string    `db:"name"`
	Age       *int       `db:"age"`
	Gender    *string    `db:"gender"`
	Birthday  *time.Time `db:"birthday"`
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

func IsValidGender(val string) bool {
	if val == Male || val == Female {
		return true
	}

	return false
}
