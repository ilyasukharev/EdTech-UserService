package utils

import (
	"UserService/internal/errors"
	"UserService/internal/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
)

func DecodeRequestAndValidate(body io.ReadCloser, v any) error {
	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(v); err != nil {
		return fmt.Errorf("%w:%v", errors.BodyInvalidFormatErr, err)
	}

	if err := model.Validator.Struct(v); err != nil {
		return fmt.Errorf("%w:%v", errors.BodyInvalidContentErr, err)
	}

	return nil
}

func SendResponse(ctx context.Context, w http.ResponseWriter, status int, response any) {
	SetRequestId(GetRequestId(ctx), w.Header()).Set(
		"Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response)
}

func SendErrorResponse(ctx context.Context, w http.ResponseWriter, err *model.ApiError, originalErr error) {
	reqId := GetRequestId(ctx)
	SetRequestId(reqId, w.Header()).Set(
		"Content-Type", "application/json")
	w.WriteHeader(err.Code)
	log.Println(err.WithRequestId(reqId).WithDebug(originalErr).Error())
	_ = json.NewEncoder(w).Encode(err)
}

func CheckHttpClientStatusCodeIsPositive(source string, response *http.Response) error {
	errorMsg := "%s error: status=%d body=%s"
	if response.StatusCode >= 300 {
		b, _ := io.ReadAll(response.Body)
		return fmt.Errorf(errorMsg, source, response.StatusCode, b)
	}
	return nil
}

func GetRequestId(ctx context.Context) string {
	requestId := ctx.Value(model.XRequestIDHeader)
	return requestId.(string)
}

func SetRequestId(reqId string, header http.Header) http.Header {
	header.Set(model.XRequestIDHeader, reqId)
	return header
}

type Parser[T any] func(string) (T, error)

func GetPathParam[T any](
	r *http.Request,
	name string,
	parse Parser[T],
) (T, error) {
	value := chi.URLParam(r, name)
	return parse(value)
}

func GetQueryParam[T any](
	r *http.Request,
	name string,
	parse Parser[T],
) (T, error) {
	value := r.URL.Query().Get(name)
	return parse(value)
}

func ParseUUID(v string) (uuid.UUID, error) {
	return uuid.Parse(v)
}

func ParseString(v string) (string, error) {
	return v, nil
}
