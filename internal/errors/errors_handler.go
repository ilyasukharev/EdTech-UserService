package errors

import (
	"UserService/internal/model"
	"errors"
	"github.com/lib/pq"
	"net/http"
)

const (
	internalServerError = "internal server error"
)

func GetApiError(err error) *model.ApiError {
	switch {
	case errors.Is(err, NothingToUpdateErr):
		return model.NewApiError(http.StatusNoContent, NothingToUpdateErr.Error())

	case errors.Is(err, BodyInvalidFormatErr):
		return model.NewApiError(http.StatusBadRequest, BodyInvalidFormatErr.Error())

	case errors.Is(err, BodyInvalidContentErr):
		return model.NewApiError(http.StatusBadRequest, BodyInvalidContentErr.Error())

	case errors.Is(err, RegistrationIDNotFoundErr):
		return model.NewApiError(http.StatusBadRequest, RegistrationIDNotFoundErr.Error())

	case errors.Is(err, UserNotFoundErr):
		return model.NewApiError(http.StatusNotFound, UserNotFoundErr.Error())

	case isPqError(err):
		return handlePqError(err)
	}

	return model.NewApiError(http.StatusInternalServerError, internalServerError)
}

func isPqError(err error) bool {
	var pqErr *pq.Error
	return errors.As(err, &pqErr)
}

const (
	pqUniqueViolationErrCode     = "23505"
	pqForeignKeyViolationErrCode = "23503"
	pqNotNullViolationErrCode    = "23502"
)

func handlePqError(err error) *model.ApiError {
	var pqErr *pq.Error
	errors.As(err, &pqErr)

	switch pqErr.Code {
	case pqUniqueViolationErrCode:
		return model.NewApiError(http.StatusConflict, DuplicateValueErr.Error())

	case pqForeignKeyViolationErrCode:
		return model.NewApiError(http.StatusBadRequest, "invalid reference")

	case pqNotNullViolationErrCode:
		return model.NewApiError(http.StatusBadRequest, "required field missing")
	}

	return model.NewApiError(
		http.StatusInternalServerError,
		internalServerError,
	)
}
