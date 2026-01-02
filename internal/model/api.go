package model

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"strings"
)

type LogLevel string

const (
	LogError LogLevel = "ERROR"
	LogWarn  LogLevel = "WARN"
	LogInfo  LogLevel = "INFO"
)

const (
	baseApiPath = "/api/user-service"

	XRequestIDHeader = "X-Request-ID"
)

const SourceApiErr = "API"

var Validator = validator.New()

type ApiError struct {
	Code    int           `json:"code"`
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
