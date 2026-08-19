package model

import (
	"github.com/google/uuid"
	"time"
)

type Referral struct {
	ID         *int64     `db:"id"`
	ReferrerID *uuid.UUID `db:"referrer_id"`
	RefereeID  *uuid.UUID `db:"referee_id"`
	Confirmed  *bool      `db:"confirmed"`
	CreatedAt  *time.Time `db:"created_at"`
	UpdatedAt  *time.Time `db:"updated_at"`
}
