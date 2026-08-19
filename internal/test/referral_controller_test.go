package transport

import (
	"UserService/internal/model"
	"UserService/internal/transport/referral"
	"UserService/internal/utils"
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func TestShouldCreateRef(t *testing.T) {
	r, td := configure(t)
	user1, user2 := createDatabaseUsers(t, td.db)
	createReferral(t, r, createReferralModel(*user1.ID, *user2.ID), http.StatusCreated)
}

func TestShouldNotCreateRefIfRefIDsDoNotExist(t *testing.T) {
	r, _ := configure(t)
	createReferral(t, r, createReferralModel(uuid.New(), uuid.New()), http.StatusBadRequest)
}

func TestShouldNotCreateRefIfReferrerIDIsEqualToRefereeID(t *testing.T) {
	r, td := configure(t)
	user1, _ := createDatabaseUsers(t, td.db)
	createReferral(t, r, createReferralModel(*user1.ID, *user1.ID), http.StatusInternalServerError)
}

func TestShouldGetByReferrerID(t *testing.T) {
	r, td := configure(t)
	user1, user2 := createDatabaseUsers(t, td.db)
	refModel := createReferralModel(*user1.ID, *user2.ID)
	createDatabaseReferrals(t, td.db, refModel.ReferrerID, refModel.RefereeID)
	getReferralByReferrer(t, r, refModel, refModel.ReferrerID, http.StatusOK)
}

func TestShouldNotGetByReferrerIDIfIDIsUnknown(t *testing.T) {
	r, _ := configure(t)
	getReferralByReferrer(t, r, nil, uuid.New(), http.StatusNotFound)
}

func TestShouldGetByRefereeID(t *testing.T) {
	r, td := configure(t)
	user1, user2 := createDatabaseUsers(t, td.db)
	refModel := createReferralModel(*user1.ID, *user2.ID)
	createDatabaseReferrals(t, td.db, refModel.ReferrerID, refModel.RefereeID)
	getReferralByReferee(t, r, refModel, refModel.RefereeID, http.StatusOK)
}

func TestShouldNotGetByRefereeIDIfIDIsUnknown(t *testing.T) {
	r, _ := configure(t)
	getReferralByReferee(t, r, nil, uuid.New(), http.StatusNotFound)
}

func TestShouldPatchReferral(t *testing.T) {
	r, td := configure(t)
	user1, user2 := createDatabaseUsers(t, td.db)
	refModel := createReferralModel(*user1.ID, *user2.ID)
	ref := createDatabaseReferrals(t, td.db, refModel.ReferrerID, refModel.RefereeID)
	patchModel := patchReferralModel(*user1.ID, *user2.ID)
	patchReferral(t, r, patchModel, *ref.ID, http.StatusOK)
}

func TestShouldNotPatchReferralIfRefIDDoesNotExist(t *testing.T) {
	r, td := configure(t)
	user1, user2 := createDatabaseUsers(t, td.db)
	refModel := createReferralModel(*user1.ID, *user2.ID)
	createDatabaseReferrals(t, td.db, refModel.ReferrerID, refModel.RefereeID)
	patchModel := patchReferralModel(*user1.ID, *user2.ID)
	patchReferral(t, r, patchModel, 400000, http.StatusNotFound)
}

func createReferral(t *testing.T, r chi.Router, expected *referral.CreateReferral, expectedCode int) *referral.ReferralResponse {
	bodyBytes, _ := json.Marshal(expected)
	req := httptest.NewRequest(http.MethodPost, model.NewApiPath("/v1/referrals"), bytes.NewBuffer(bodyBytes))

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != expectedCode {
		t.Fatalf("expected %d, got %d", expectedCode, rr.Code)
	}

	if expectedCode != http.StatusCreated {
		return nil
	}

	var actual referral.ReferralResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &actual); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	assertPtrEqual(t, "referrer_id", &actual.ReferrerID, &expected.ReferrerID)
	assertPtrEqual(t, "referee_id", &actual.RefereeID, &expected.RefereeID)
	assertPtrEqual(t, "confirmed", &actual.Confirmed, utils.BoolPtr(false))
	var nilTime time.Time
	if actual.CreatedAt == nilTime {
		t.Fatalf("created at time is nil")
	}

	return &actual
}

func getReferralByReferrer(
	t *testing.T,
	r chi.Router,
	expected *referral.CreateReferral,
	referrerID uuid.UUID,
	expectedCode int,
) *referral.ReferralResponse {
	req := httptest.NewRequest(http.MethodGet,
		model.NewApiPath("/v1/referrals/by-referrer/"+referrerID.String()), nil)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != expectedCode {
		t.Fatalf("expected %d, got %d", expectedCode, rr.Code)
	}

	if expectedCode != http.StatusOK {
		return nil
	}

	var actual referral.ReferralResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &actual); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	assertPtrEqual(t, "referrer_id", &actual.ReferrerID, &expected.ReferrerID)
	assertPtrEqual(t, "referee_id", &actual.RefereeID, &expected.RefereeID)

	return &actual
}

func getReferralByReferee(
	t *testing.T,
	r chi.Router,
	expected *referral.CreateReferral,
	referrerID uuid.UUID,
	expectedCode int,
) *referral.ReferralResponse {
	req := httptest.NewRequest(http.MethodGet,
		model.NewApiPath("/v1/referrals/by-referee/"+referrerID.String()), nil)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != expectedCode {
		t.Fatalf("expected %d, got %d", expectedCode, rr.Code)
	}

	if expectedCode != http.StatusOK {
		return nil
	}

	var actual referral.ReferralResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &actual); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	assertPtrEqual(t, "referrer_id", &actual.ReferrerID, &expected.ReferrerID)
	assertPtrEqual(t, "referee_id", &actual.RefereeID, &expected.RefereeID)

	return &actual
}

func patchReferral(t *testing.T, r chi.Router, expected *referral.PatchReferral, ID int64, expectedCode int) *referral.ReferralResponse {
	bodyBytes, _ := json.Marshal(expected)
	req := httptest.NewRequest(http.MethodPatch,
		model.NewApiPath("/v1/referrals/"+strconv.FormatInt(ID, 10)), bytes.NewBuffer(bodyBytes))

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != expectedCode {
		t.Fatalf("expected %d, got %d", expectedCode, rr.Code)
	}

	if expectedCode != http.StatusOK {
		return nil
	}

	var actual referral.ReferralResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &actual); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	assertPtrEqual(t, "referrer_id", &actual.ReferrerID, expected.ReferrerID)
	assertPtrEqual(t, "referee_id", &actual.RefereeID, expected.RefereeID)
	assertPtrEqual(t, "confirmed", &actual.Confirmed, expected.Confirmed)
	assertTimePtrEqual(t, "confirmed", &actual.CreatedAt, expected.CreatedAt)

	return &actual
}

func createReferralModel(referrerID uuid.UUID, refereeID uuid.UUID) *referral.CreateReferral {
	return &referral.CreateReferral{
		ReferrerID: referrerID,
		RefereeID:  refereeID,
	}
}

func patchReferralModel(referrerID uuid.UUID, refereeID uuid.UUID) *referral.PatchReferral {
	return &referral.PatchReferral{
		ReferrerID: &referrerID,
		RefereeID:  &refereeID,
		Confirmed:  utils.BoolPtr(true),
		CreatedAt:  utils.TimePtr(time.Now()),
	}
}

func createDatabaseReferrals(t *testing.T, db *sqlx.DB, referrerID uuid.UUID, refereeID uuid.UUID) *model.Referral {
	query := `
	INSERT INTO referrals (referrer_id, referee_id)
	VALUES ($1, $2) 
	RETURNING *
	`

	var ref model.Referral
	err := db.Get(&ref, query, referrerID, refereeID)
	if err != nil {
		t.Fatalf("failed to insert referrals 1: %v", err)
	}

	return &ref
}

func createDatabaseUsers(t *testing.T, db *sqlx.DB) (*model.User, *model.User) {
	uuid1 := uuid.New()
	uuid2 := uuid.New()

	user1 := createUserModel().ToUser()
	user1.ID = &uuid1

	user2 := createUserModel().ToUser()
	user2.ID = &uuid2

	query := `
	INSERT INTO users (id, first_name, email, notifications, type)
	VALUES ($1, $2, $3, $4, $5) 
	RETURNING *
	`

	_, err := db.Exec(query, user1.ID, user1.FirstName, user1.Email, user1.Notifications, user1.Type)
	if err != nil {
		t.Fatalf("failed to insert user 1: %v", err)
	}
	_, err = db.Exec(query, user2.ID, user2.FirstName, user2.Email, user2.Notifications, user2.Type)
	if err != nil {
		t.Fatalf("failed to insert user 2: %v", err)
	}

	return user1, user2
}
