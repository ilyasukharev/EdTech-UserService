package transport

import (
	"UserService/internal/model"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	uuid2 "github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"testing"
)

const defaultOTPCode = "111111"

func TestShouldSendOTPCode(t *testing.T) {
	r, td := configure(t)
	sendCode(t, r, td, "test@mail.ru", http.StatusNoContent)
}

func TestShouldVerifyOTPCode(t *testing.T) {
	r, td := configure(t)
	email := "test1@mail.ru"
	err := td.authController.Service.Repo.SaveVerificationCode(t.Context(), email, defaultOTPCode)
	if err != nil {
		t.Errorf("Error saving verification code: %v", err)
	}
	verifyCode(t, r, td, email, http.StatusOK)
}

func TestShouldNotVerifyOTPCodeIfEmailIsUnknown(t *testing.T) {
	r, td := configure(t)
	verifyCode(t, r, td, "unknown-email", http.StatusInternalServerError)
}

func sendCode(t *testing.T, r chi.Router, td *TestData, email string, expectedCode int) {
	req := httptest.NewRequest(http.MethodGet,
		model.NewApiPath(fmt.Sprintf("/v1/auth/code/send?email=%s", email)), nil)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != expectedCode {
		t.Fatalf("expected %d, got %d", expectedCode, rr.Code)
	}

	code, err := td.authController.Service.Repo.GetVerificationCode(t.Context(), email)
	if err != nil {
		t.Fatalf("code must be found %v", err)
	}
	if code != defaultOTPCode {
		t.Fatalf("expected %s, got %s", defaultOTPCode, code)
	}
}

func verifyCode(t *testing.T, r chi.Router, td *TestData, email string, expectedCode int) {
	req := httptest.NewRequest(http.MethodGet,
		model.NewApiPath(fmt.Sprintf(
			"/v1/auth/code/verify?email=%s&code=%s", email, defaultOTPCode)), nil)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != expectedCode {
		t.Fatalf("expected %d, got %d", expectedCode, rr.Code)
	}

	if expectedCode != http.StatusOK {
		return
	}

	var uuid uuid2.UUID
	if err := json.Unmarshal(rr.Body.Bytes(), &uuid); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	expectedEmail, err := td.authController.Service.Repo.GetEmailByRegistrationID(t.Context(), uuid)
	if err != nil {
		t.Fatalf("email must be found %v", err)
	}
	if expectedEmail != email {
		t.Fatalf("expected %s, got %s", expectedEmail, email)
	}
}
