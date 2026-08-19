package auth

import (
	"UserService/internal/middleware"
	"UserService/internal/model"
	"UserService/internal/service"
	"UserService/internal/utils"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type AuthController struct {
	Service *service.AuthService
}

func NewAuthController(service *service.AuthService) *AuthController {
	return &AuthController{
		Service: service,
	}
}

func (c *AuthController) RegisterRoutes(r chi.Router) {
	r.Get(model.NewApiPath("/v1/auth/code/send"), middleware.WrapResponseWithCode(c.SendAuthCode, http.StatusNoContent))
	r.Get(model.NewApiPath("/v1/auth/code/verify"), middleware.WrapResponse(c.VerifyAuthCode))
}

// SendAuthCode godoc
// @Summary Сохранить OTP код по переданному email пользователя
// @Tags Auth
// @Produce json
// @Param email query string true "email"
// @Success 204
// @Router /api/user-service/v1/auth/code/send [get]
func (c *AuthController) SendAuthCode(_ http.ResponseWriter, r *http.Request) (any, error) {
	email, err := utils.GetQueryParam(r, "email", utils.ParseString)
	if err != nil {
		return nil, err
	}
	return "", c.Service.SendCode(r.Context(), email)
}

// VerifyAuthCode godoc
// @Summary Верифицировать OTP код по переданному email пользователя
// @Tags Auth
// @Produce json
// @Param email query string true "email"
// @Param code query string true "code"
// @Success 200 {string} string"
// @Failure 400 {object} model.ApiError "Переданный OTP код неверный"
// @Failure 500 {object} model.ApiError
// @Router /api/user-service/v1/auth/code/verify [get]
func (c *AuthController) VerifyAuthCode(_ http.ResponseWriter, r *http.Request) (any, error) {
	email, err := utils.GetQueryParam(r, "email", utils.ParseString)
	if err != nil {
		return nil, err
	}

	code, err := utils.GetQueryParam(r, "code", utils.ParseString)
	if err != nil {
		return nil, err
	}

	return c.Service.VerifyAuthCode(r.Context(), email, code)
}
