package referral

import (
	"UserService/internal/model"
	"github.com/google/uuid"
	"time"
)

type CreateReferral struct {
	ReferrerID uuid.UUID `json:"referrer_id" binding:"required"`
	RefereeID  uuid.UUID `json:"referee_id" binding:"required"`
}

func (r *CreateReferral) ToReferral() *model.Referral {
	return &model.Referral{
		ReferrerID: &r.ReferrerID,
		RefereeID:  &r.RefereeID,
	}
}

type PatchReferral struct {
	ReferrerID *uuid.UUID `json:"referrer_id"`
	RefereeID  *uuid.UUID `json:"referee_id"`
	Confirmed  *bool      `json:"confirmed"`
	CreatedAt  *time.Time `json:"created_at"`
}

func (r *PatchReferral) ToReferral(ID int64) *model.Referral {
	return &model.Referral{
		ID:         &ID,
		ReferrerID: r.ReferrerID,
		RefereeID:  r.RefereeID,
		Confirmed:  r.Confirmed,
		CreatedAt:  r.CreatedAt,
	}
}

type ReferralResponse struct {
	ID         int64      `json:"id" binding:"required"`
	ReferrerID uuid.UUID  `json:"referrer_id" binding:"required"`
	RefereeID  uuid.UUID  `json:"referee_id" binding:"required"`
	Confirmed  bool       `json:"confirmed" binding:"required"`
	CreatedAt  time.Time  `json:"created_at" binding:"required"`
	UpdatedAt  *time.Time `json:"updated_at"`
}
