package model

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"strconv"
	"strings"
	"time"
)

type LogLevel string

const LogError LogLevel = "ERROR"

const (
	baseApiPath = "/api/user-service"

	XRequestIDHeader = "X-Request-ID"
)

var Validator = validator.New()

type ApiError struct {
	Code    int           `json:"-"`
	Message string        `json:"message"`
	Details *DebugDetails `json:"-"`
}

type DebugDetails struct {
	RequestId string
	Level     LogLevel
	Debug     []any
}

func NewApiPath(path string) string {
	return baseApiPath + path
}

func (e *ApiError) Error() string {
	if e.Code == 0 || e.Message == "" || e.Details.Level == "" {
		return "Api error must be filled correctly"
	}

	builder := &strings.Builder{}
	builder.WriteString(string(e.Details.Level) + ": ")
	builder.WriteString(fmt.Sprintf("with requestId: '%s' ", e.Details.RequestId))
	builder.WriteString(fmt.Sprintf("with code: '%d' ", e.Code))
	if e.Message != "" {
		builder.WriteString(fmt.Sprintf("with msg: '%s' ", e.Message))
	}

	if e.Details.Debug != nil {
		builder.WriteString(fmt.Sprintf("with info: '%v' ", e.Details.Debug))
	}

	return builder.String()
}

func NewApiError(code int, msg string) *ApiError {
	return &ApiError{
		Code:    code,
		Message: msg,
		Details: &DebugDetails{
			Level:     LogError,
			RequestId: uuid.New().String(),
		},
	}
}

func (e *ApiError) WithCode(code int) *ApiError {
	e.Code = code
	return e
}

func (e *ApiError) WithMessage(msg string) *ApiError {
	e.Message = msg
	return e
}

func (e *ApiError) WithDebug(debug ...any) *ApiError {
	e.Details.Debug = debug
	return e
}

func (e *ApiError) WithLevel(level LogLevel) *ApiError {
	e.Details.Level = level
	return e
}

func (e *ApiError) WithRequestId(requestId string) *ApiError {
	e.Details.RequestId = requestId
	return e
}

type PatchField struct {
	Value any
	Name  string
}

type Date time.Time

func (d *Date) UnmarshalJSON(data []byte) error {
	value, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}

	t, err := time.Parse("2006-01-02", value)
	if err != nil {
		return err
	}

	*d = Date(t)
	return nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(d).Format("2006-01-02"))
}

func (d *Date) AsTime() *time.Time {
	if d == nil {
		return nil
	}
	t := time.Time(*d)
	return &t
}
