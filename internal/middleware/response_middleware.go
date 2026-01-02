package middleware

import (
	"UserService/internal/errors"
	"UserService/internal/utils"
	"net/http"
)

type HandlerFunc func(http.ResponseWriter, *http.Request) (any, error)

func WrapResponse(next HandlerFunc) http.HandlerFunc {
	return WrapResponseWithCode(next, http.StatusOK)
}

func WrapResponseWithCode(next HandlerFunc, successCode int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := next(w, r)

		ctx := r.Context()
		if err != nil {
			apiErr := errors.GetApiError(err)
			utils.SendErrorResponse(ctx, w, apiErr, err)
			return
		}

		utils.SendResponse(ctx, w, successCode, resp)
	}
}
