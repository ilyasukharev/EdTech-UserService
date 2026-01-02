package transport

import (
	"UserService/internal/transport/user"
	"UserService/internal/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"math/rand/v2"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/lib/pq"
)

func TestShouldCreateUser(t *testing.T) {
	createUser(t, createUserModel(), http.StatusCreated, true)
}

func TestShouldReturnErrWhenParentDoesNotHaveChild(t *testing.T) {
	userModel := createUserModel()
	userModel.ChildName = nil
	userModel.ChildAge = nil

	createUser(t, userModel, http.StatusBadRequest, true)
}

func TestShouldReturnErrWhenChildAgeValueIsNotInBounds(t *testing.T) {
	userModel := createUserModel()
	if rand.IntN(2) == 0 {
		userModel.ChildAge = utils.IntPtr(0)
	} else {
		userModel.ChildAge = utils.IntPtr(19)
	}
	userModel.ChildAge = nil

	createUser(t, userModel, http.StatusBadRequest, true)
}

func TestShouldReturnErrWhenUserWithEmailAlreadyExists(t *testing.T) {
	userModel1 := createUserModel()
	createUser(t, userModel1, http.StatusCreated, true)

	userModel2 := createUserModel()
	userModel2.Email = userModel1.Email
	createUser(t, userModel2, http.StatusConflict, true)
}

func TestShouldReturnErrWhenRegistrationIDIsNotExists(t *testing.T) {
	createUser(t, createUserModel(), http.StatusBadRequest, false)
}

func TestShouldGetUserById(t *testing.T) {
	userModel := createUserModel()
	userResponse := createUser(t, userModel, http.StatusCreated, true)
	getUserById(t, userModel, userResponse.ID, http.StatusOK)
}

func TestShouldReturnErrWhenGetUserByUnExistId(t *testing.T) {
	getUserById(t, createUserModel(), uuid.New(), http.StatusNotFound)
}

func TestShouldGetUserByEmail(t *testing.T) {
	userModel := createUserModel()
	userResponse := createUser(t, userModel, http.StatusCreated, true)
	getUserByEmail(t, userModel, userResponse.Email, http.StatusOK)
}

func TestShouldReturnErrWhenGetUserByUnExistEmail(t *testing.T) {
	getUserByEmail(t, createUserModel(), uuid.NewString()+"@ya.ru", http.StatusNotFound)
}

func TestShouldUpdateUser(t *testing.T) {
	userModel := createUserModel()
	userResponse := createUser(t, userModel, http.StatusCreated, true)
	updateUser(t, updateUserModel(), userResponse.ID, http.StatusOK)
}

func TestShouldReturnErrWhenUpdateUserWithUnExistID(t *testing.T) {
	updateUser(t, updateUserModel(), uuid.New(), http.StatusNotFound)
}

func TestShouldPatchUser(t *testing.T) {
	userModel := createUserModel()
	userResponse := createUser(t, userModel, http.StatusCreated, true)
	patchUser(t, patchUserModel(), userResponse.ID, http.StatusOK)
}

func TestShouldReturnErrWhenPatchUserWithUnExistID(t *testing.T) {
	patchUser(t, patchUserModel(), uuid.New(), http.StatusNotFound)
}

func TestShouldDeleteUser(t *testing.T) {
	userModel := createUserModel()
	userResponse := createUser(t, userModel, http.StatusCreated, true)
	deleteUser(t, userModel, userResponse.ID, http.StatusOK)
}

func TestShouldReturnErrWhenDeleteUserWithUnExistID(t *testing.T) {
	deleteUser(t, createUserModel(), uuid.New(), http.StatusNotFound)
}

func createUser(t *testing.T, expected *user.CreateUser, expectedCode int, redisHasEmail bool) *user.UserResponse {
	r, e := configureEnvironment(t)

	regID := uuid.New()
	if redisHasEmail {
		_ = e.userController.Service.RedisRepo.SaveRegistrationID(t.Context(), regID, expected.Email)
	}

	bodyBytes, _ := json.Marshal(expected)
	req := httptest.NewRequest(http.MethodPost,
		fmt.Sprintf("/api/user-service/v1/users?reg_id=%s", regID), bytes.NewBuffer(bodyBytes))

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != expectedCode {
		t.Fatalf("expected %d, got %d", expectedCode, rr.Code)
	}

	if expectedCode != http.StatusCreated {
		return nil
	}

	var actual user.UserResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &actual); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	assertPtrEqual(t, "first_name", &actual.FirstName, &expected.FirstName)
	assertPtrEqual(t, "last_name", actual.LastName, expected.LastName)
	assertPtrEqual(t, "middle_name", actual.MiddleName, expected.MiddleName)
	assertPtrEqual(t, "email", &actual.Email, &expected.Email)
	assertPtrEqual(t, "phone", actual.Phone, expected.Phone)
	assertPtrEqual(t, "notifications", &actual.Notifications, &expected.Notifications)
	assertPtrEqual(t, "type", &actual.Type, &expected.Type)
	assertPtrEqual(t, "child_name", actual.ChildName, expected.ChildName)
	assertPtrEqual(t, "child_age", actual.ChildAge, expected.ChildAge)

	return &actual
}

func getUserById(t *testing.T, expected *user.CreateUser, ID uuid.UUID, expectedCode int) {
	getUser(t, expected, "/api/user-service/v1/users/"+ID.String(), expectedCode)
}

func getUserByEmail(t *testing.T, expected *user.CreateUser, email string, expectedCode int) {
	getUser(t, expected, fmt.Sprintf("/api/user-service/v1/users/by-email?email=%s", email), expectedCode)
}

func getUser(t *testing.T, expected *user.CreateUser, url string, expectedCode int) {
	r, _ := configureEnvironment(t)

	req := httptest.NewRequest(http.MethodGet, url, nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != expectedCode {
		t.Fatalf("expected %d, got %d", expectedCode, rr.Code)
	}

	if expectedCode != http.StatusOK {
		return
	}

	var actual user.UserResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &actual); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	assertPtrEqual(t, "first_name", &actual.FirstName, &expected.FirstName)
	assertPtrEqual(t, "last_name", actual.LastName, expected.LastName)
	assertPtrEqual(t, "middle_name", actual.MiddleName, expected.MiddleName)
	assertPtrEqual(t, "email", &actual.Email, &expected.Email)
	assertPtrEqual(t, "phone", actual.Phone, expected.Phone)
	assertPtrEqual(t, "notifications", &actual.Notifications, &expected.Notifications)
	assertPtrEqual(t, "type", &actual.Type, &expected.Type)
	assertPtrEqual(t, "child_name", actual.ChildName, expected.ChildName)
	assertPtrEqual(t, "child_age", actual.ChildAge, expected.ChildAge)
}

func updateUser(t *testing.T, expected *user.UpdateUser, ID uuid.UUID, expectedCode int) *user.UserResponse {
	r, _ := configureEnvironment(t)

	bodyBytes, _ := json.Marshal(expected)
	req := httptest.NewRequest(http.MethodPut, "/api/user-service/v1/users/"+ID.String(), bytes.NewBuffer(bodyBytes))

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != expectedCode {
		t.Fatalf("expected %d, got %d", expectedCode, rr.Code)
	}

	if expectedCode != http.StatusOK {
		return nil
	}

	var actual user.UserResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &actual); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	assertPtrEqual(t, "first_name", &actual.FirstName, &expected.FirstName)
	assertPtrEqual(t, "last_name", actual.LastName, &expected.LastName)
	assertPtrEqual(t, "middle_name", actual.MiddleName, &expected.MiddleName)
	assertPtrEqual(t, "email", &actual.Email, &expected.Email)
	assertPtrEqual(t, "phone", actual.Phone, &expected.Phone)
	assertPtrEqual(t, "notifications", &actual.Notifications, &expected.Notifications)
	assertPtrEqual(t, "type", &actual.Type, &expected.Type)
	assertPtrEqual(t, "child_name", actual.ChildName, &expected.ChildName)
	assertPtrEqual(t, "child_age", actual.ChildAge, &expected.ChildAge)

	return &actual
}

func patchUser(t *testing.T, expected *user.PatchUser, ID uuid.UUID, expectedCode int) *user.UserResponse {
	r, _ := configureEnvironment(t)

	bodyBytes, _ := json.Marshal(expected)
	req := httptest.NewRequest(http.MethodPatch, "/api/user-service/v1/users/"+ID.String(), bytes.NewBuffer(bodyBytes))

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != expectedCode {
		t.Fatalf("expected %d, got %d", expectedCode, rr.Code)
	}

	if expectedCode != http.StatusOK {
		return nil
	}

	var actual user.UserResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &actual); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	assertPtrEqual(t, "first_name", &actual.FirstName, expected.FirstName)
	assertPtrEqual(t, "last_name", actual.LastName, expected.LastName)
	assertPtrEqual(t, "middle_name", actual.MiddleName, expected.MiddleName)
	assertPtrEqual(t, "email", &actual.Email, expected.Email)
	assertPtrEqual(t, "phone", actual.Phone, expected.Phone)
	assertPtrEqual(t, "notifications", &actual.Notifications, expected.Notifications)
	assertPtrEqual(t, "type", &actual.Type, expected.Type)
	assertPtrEqual(t, "child_name", actual.ChildName, expected.ChildName)
	assertPtrEqual(t, "child_age", actual.ChildAge, expected.ChildAge)

	return &actual
}

func deleteUser(t *testing.T, expected *user.CreateUser, ID uuid.UUID, expectedCode int) *user.UserResponse {
	r, _ := configureEnvironment(t)

	bodyBytes, _ := json.Marshal(expected)
	req := httptest.NewRequest(http.MethodDelete, "/api/user-service/v1/users/"+ID.String(), bytes.NewBuffer(bodyBytes))

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != expectedCode {
		t.Fatalf("expected %d, got %d", expectedCode, rr.Code)
	}

	if expectedCode != http.StatusOK {
		return nil
	}

	var actual user.UserResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &actual); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	assertPtrEqual(t, "first_name", &actual.FirstName, &expected.FirstName)
	assertPtrEqual(t, "last_name", actual.LastName, expected.LastName)
	assertPtrEqual(t, "middle_name", actual.MiddleName, expected.MiddleName)
	assertPtrEqual(t, "email", &actual.Email, &expected.Email)
	assertPtrEqual(t, "phone", actual.Phone, expected.Phone)
	assertPtrEqual(t, "notifications", &actual.Notifications, &expected.Notifications)
	assertPtrEqual(t, "type", &actual.Type, &expected.Type)
	assertPtrEqual(t, "child_name", actual.ChildName, expected.ChildName)
	assertPtrEqual(t, "child_age", actual.ChildAge, expected.ChildAge)
	if actual.DeletedAt == nil {
		t.Fatal("expected DeletedAt to be non-nil")
	}

	return &actual
}
