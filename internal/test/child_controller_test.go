package transport

import (
	"UserService/internal/model"
	"UserService/internal/transport/child"
	"UserService/internal/utils"
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestShouldCreateChild(t *testing.T) {
	r, td := configure(t)
	user1, _ := createDatabaseUsers(t, td.db)
	createChild(t, r, createChildModel(*user1.ID), http.StatusCreated)
}

func TestShouldCreateChildWithoutNonRequiredFields(t *testing.T) {
	r, td := configure(t)
	user1, _ := createDatabaseUsers(t, td.db)
	childModel := createChildModel(*user1.ID)
	childModel.Gender = nil
	childModel.Birthday = nil
	createChild(t, r, childModel, http.StatusCreated)
}

func TestShouldNotCreateChildIfParentIDDoesNotExist(t *testing.T) {
	r, _ := configure(t)
	createChild(t, r, createChildModel(uuid.New()), http.StatusBadRequest)
}

func TestShouldNotCreateChildIfGenderIsNotCorrect(t *testing.T) {
	r, td := configure(t)
	user1, _ := createDatabaseUsers(t, td.db)
	childModel := createChildModel(*user1.ID)
	childModel.Gender = utils.StringPtr("asd")
	createChild(t, r, childModel, http.StatusBadRequest)
}

func TestShouldGetChildByID(t *testing.T) {
	r, td := configure(t)
	user1, _ := createDatabaseUsers(t, td.db)
	childModel := createChildModel(*user1.ID)
	createdChild := createDatabaseChild(t, td.db, childModel)
	getChild(t, r, childModel, *createdChild.ID, http.StatusOK)
}

func TestShouldNoGetChildByUnknownID(t *testing.T) {
	r, _ := configure(t)
	getChild(t, r, createChildModel(uuid.New()), uuid.New(), http.StatusNotFound)
}

func TestShouldGetChildByParentID(t *testing.T) {
	r, td := configure(t)
	user1, _ := createDatabaseUsers(t, td.db)
	childModel := createChildModel(*user1.ID)
	createDatabaseChild(t, td.db, childModel)
	getChildByParent(t, r, childModel, *user1.ID, http.StatusOK)
}

func TestShouldNotGetChildByUnknownParentID(t *testing.T) {
	r, td := configure(t)
	user1, _ := createDatabaseUsers(t, td.db)
	childModel := createChildModel(*user1.ID)
	createDatabaseChild(t, td.db, childModel)
	getChildByParent(t, r, childModel, uuid.New(), http.StatusNotFound)
}

func TestShouldUpdateChild(t *testing.T) {
	r, td := configure(t)
	user1, _ := createDatabaseUsers(t, td.db)
	childModel := createChildModel(*user1.ID)
	createdChild := createDatabaseChild(t, td.db, childModel)
	updateChild(t, r, updateChildModel(*user1.ID, *createdChild.CreatedAt), *createdChild.ID, http.StatusOK)
}

func TestShouldNotUpdateChildIfIDIsUnknown(t *testing.T) {
	r, td := configure(t)
	user1, _ := createDatabaseUsers(t, td.db)
	childModel := createChildModel(*user1.ID)
	createdChild := createDatabaseChild(t, td.db, childModel)
	updateChild(t, r, updateChildModel(*user1.ID, *createdChild.CreatedAt), uuid.New(), http.StatusNotFound)
}

func TestShouldNotUpdateChildIfGenderIsIncorrect(t *testing.T) {
	r, td := configure(t)
	user1, _ := createDatabaseUsers(t, td.db)
	childModel := createChildModel(*user1.ID)
	createdChild := createDatabaseChild(t, td.db, childModel)
	updateModel := updateChildModel(*user1.ID, *createdChild.CreatedAt)
	updateModel.Gender = "asd"
	updateChild(t, r, updateModel, *createdChild.ID, http.StatusBadRequest)
}

func TestShouldPatchChild(t *testing.T) {
	r, td := configure(t)
	user1, _ := createDatabaseUsers(t, td.db)
	childModel := createChildModel(*user1.ID)
	createdChild := createDatabaseChild(t, td.db, childModel)
	patchChild(t, r, patchChildModel(*user1.ID, *createdChild.CreatedAt), *createdChild.ID, http.StatusOK)
}

func TestShouldNotPatchChildIfIDIsUnknown(t *testing.T) {
	r, td := configure(t)
	user1, _ := createDatabaseUsers(t, td.db)
	childModel := createChildModel(*user1.ID)
	createdChild := createDatabaseChild(t, td.db, childModel)
	patchChild(t, r, patchChildModel(*user1.ID, *createdChild.CreatedAt), uuid.New(), http.StatusNotFound)
}

func TestShouldNotPatchChildIfGenderIsIncorrect(t *testing.T) {
	r, td := configure(t)
	user1, _ := createDatabaseUsers(t, td.db)
	childModel := createChildModel(*user1.ID)
	createdChild := createDatabaseChild(t, td.db, childModel)
	patchModel := patchChildModel(*user1.ID, *createdChild.CreatedAt)
	patchModel.Gender = utils.StringPtr("asd")
	patchChild(t, r, patchModel, *createdChild.ID, http.StatusBadRequest)
}

func TestShouldDeleteChild(t *testing.T) {
	r, td := configure(t)
	user1, _ := createDatabaseUsers(t, td.db)
	createdChild := createDatabaseChild(t, td.db, createChildModel(*user1.ID))
	deleteChild(t, r, *createdChild.ID, http.StatusOK)
}

func TestShouldNotDeleteChildIfIdIsUnknown(t *testing.T) {
	r, _ := configure(t)
	deleteChild(t, r, uuid.New(), http.StatusNotFound)
}

func createChild(t *testing.T, r chi.Router, expected *child.CreateChild, expectedCode int) {
	bodyBytes, _ := json.Marshal(expected)
	req := httptest.NewRequest(http.MethodPost,
		model.NewApiPath("/v1/children"), bytes.NewBuffer(bodyBytes))

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != expectedCode {
		t.Fatalf("expected %d, got %d", expectedCode, rr.Code)
	}

	if expectedCode != http.StatusCreated {
		return
	}

	var actual child.ChildResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &actual); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	assertPtrEqual(t, "parent_id", &actual.ParentID, &expected.ParentID)
	assertPtrEqual(t, "name", &actual.Name, &expected.Name)
	assertPtrEqual(t, "age", &actual.Age, &expected.Age)
	assertPtrEqual(t, "gender", actual.Gender, expected.Gender)
	assertPtrEqual(t, "birthday", actual.Birthday, dateToString(t, expected.Birthday))
	var nilTime time.Time
	if actual.CreatedAt == nilTime {
		t.Fatalf("created at time is nil")
	}
}

func getChild(t *testing.T, r chi.Router, expected *child.CreateChild, childID uuid.UUID, expectedCode int) {
	req := httptest.NewRequest(http.MethodGet,
		model.NewApiPath("/v1/children/")+childID.String(), nil)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != expectedCode {
		t.Fatalf("expected %d, got %d", expectedCode, rr.Code)
	}

	if expectedCode != http.StatusOK {
		return
	}

	var actual child.ChildResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &actual); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	assertPtrEqual(t, "parent_id", &actual.ParentID, &expected.ParentID)
	assertPtrEqual(t, "name", &actual.Name, &expected.Name)
	assertPtrEqual(t, "age", &actual.Age, &expected.Age)
	assertPtrEqual(t, "gender", actual.Gender, expected.Gender)
	assertPtrEqual(t, "birthday", actual.Birthday, dateToString(t, expected.Birthday))
}

func getChildByParent(t *testing.T, r chi.Router, expected *child.CreateChild, parentID uuid.UUID, expectedCode int) {
	req := httptest.NewRequest(http.MethodGet,
		model.NewApiPath("/v1/children/by-parent/")+parentID.String(), nil)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != expectedCode {
		t.Fatalf("expected %d, got %d", expectedCode, rr.Code)
	}

	if expectedCode != http.StatusOK {
		return
	}

	var actual child.ChildResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &actual); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	assertPtrEqual(t, "parent_id", &actual.ParentID, &expected.ParentID)
	assertPtrEqual(t, "name", &actual.Name, &expected.Name)
	assertPtrEqual(t, "age", &actual.Age, &expected.Age)
	assertPtrEqual(t, "gender", actual.Gender, expected.Gender)
	assertPtrEqual(t, "birthday", actual.Birthday, dateToString(t, expected.Birthday))
}

func updateChild(t *testing.T, r chi.Router, expected *child.UpdateChild, childID uuid.UUID, expectedCode int) {
	bodyBytes, _ := json.Marshal(expected)
	req := httptest.NewRequest(http.MethodPut,
		model.NewApiPath("/v1/children/"+childID.String()), bytes.NewBuffer(bodyBytes))

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != expectedCode {
		t.Fatalf("expected %d, got %d", expectedCode, rr.Code)
	}

	if expectedCode != http.StatusOK {
		return
	}

	var actual child.ChildResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &actual); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	assertPtrEqual(t, "parent_id", &actual.ParentID, &expected.ParentID)
	assertPtrEqual(t, "name", &actual.Name, &expected.Name)
	assertPtrEqual(t, "age", &actual.Age, &expected.Age)
	assertPtrEqual(t, "gender", actual.Gender, &expected.Gender)
	assertPtrEqual(t, "birthday", actual.Birthday, dateToString(t, &expected.Birthday))
	assertTimePtrEqual(t, "created_at", &actual.CreatedAt, &expected.CreatedAt)
	if actual.UpdatedAt == nil {
		t.Fatalf("updated at time is nil")
	}
}

func patchChild(t *testing.T, r chi.Router, expected *child.PatchChild, childID uuid.UUID, expectedCode int) {
	bodyBytes, _ := json.Marshal(expected)
	req := httptest.NewRequest(http.MethodPatch,
		model.NewApiPath("/v1/children/"+childID.String()), bytes.NewBuffer(bodyBytes))

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != expectedCode {
		t.Fatalf("expected %d, got %d", expectedCode, rr.Code)
	}

	if expectedCode != http.StatusOK {
		return
	}

	var actual child.ChildResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &actual); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	assertPtrEqual(t, "parent_id", &actual.ParentID, expected.ParentID)
	assertPtrEqual(t, "name", &actual.Name, expected.Name)
	assertPtrEqual(t, "age", &actual.Age, expected.Age)
	assertPtrEqual(t, "gender", actual.Gender, expected.Gender)
	assertPtrEqual(t, "birthday", actual.Birthday, dateToString(t, expected.Birthday))
	assertTimePtrEqual(t, "created_at", &actual.CreatedAt, expected.CreatedAt)
	if actual.UpdatedAt == nil {
		t.Fatalf("updated at time is nil")
	}
}

func deleteChild(t *testing.T, r chi.Router, childID uuid.UUID, expectedCode int) {
	req := httptest.NewRequest(http.MethodDelete,
		model.NewApiPath("/v1/children/"+childID.String()), nil)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != expectedCode {
		t.Fatalf("expected %d, got %d", expectedCode, rr.Code)
	}

	if expectedCode != http.StatusOK {
		return
	}

	var actual child.ChildResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &actual); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if actual.DeletedAt == nil {
		t.Fatalf("deleted at time is nil")
	}
	assertTimePtrEqual(t, "deleted_at/updated_at", actual.DeletedAt, actual.UpdatedAt)
}

func createChildModel(parentID uuid.UUID) *child.CreateChild {
	return &child.CreateChild{
		ParentID: parentID,
		Name:     "Евгений",
		Age:      7,
		Gender:   utils.StringPtr(model.Male),
		Birthday: getDefaultDate(),
	}
}

func updateChildModel(parentID uuid.UUID, createdAt time.Time) *child.UpdateChild {
	return &child.UpdateChild{
		ParentID:  parentID,
		Name:      "Константин",
		Age:       10,
		Gender:    model.Female,
		Birthday:  *getDefaultDate(),
		CreatedAt: createdAt,
	}
}

func patchChildModel(parentID uuid.UUID, createdAt time.Time) *child.PatchChild {
	return &child.PatchChild{
		ParentID:  &parentID,
		Name:      utils.StringPtr("Константин"),
		Age:       utils.IntPtr(10),
		Gender:    utils.StringPtr(model.Female),
		Birthday:  getDefaultDate(),
		CreatedAt: &createdAt,
	}
}

func createDatabaseChild(t *testing.T, db *sqlx.DB, childModel *child.CreateChild) *model.Child {
	t.Helper()

	query := `
	INSERT INTO children (PARENT_ID, NAME, AGE, GENDER, BIRTHDAY)
	VALUES ($1, $2, $3, $4, $5) 
	RETURNING *
	`

	var child model.Child
	err := db.Get(&child, query, childModel.ParentID, childModel.Name,
		childModel.Age, childModel.Gender, childModel.Birthday.AsTime())
	if err != nil {
		t.Fatalf("failed to insert child 1: %v", err)
	}

	return &child
}

func dateToString(t *testing.T, s *model.Date) *string {
	t.Helper()

	if s == nil {
		return nil
	}
	dateStr := s.AsTime().Format("2006-01-02")
	return &dateStr
}

func getDefaultDate() *model.Date {
	date := model.Date(time.Date(
		2001,
		time.January,
		1,
		0, 0, 0, 0,
		time.UTC,
	))
	return &date
}
