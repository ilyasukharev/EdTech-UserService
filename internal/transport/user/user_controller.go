package user

import (
	"UserService/internal/middleware"
	"UserService/internal/model"
	"UserService/internal/service"
	"UserService/internal/utils"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type UserController struct {
	Service *service.UserService
}

func NewUserController(service *service.UserService) *UserController {
	return &UserController{
		Service: service,
	}
}

func (c *UserController) RegisterRoutes(r chi.Router) {
	r.Post(model.NewApiPath("/v1/users"), middleware.WrapResponseWithCode(c.CreateUser, http.StatusCreated))
	r.Get(model.NewApiPath("/v1/users/{id}"), middleware.WrapResponse(c.GetUserById))
	r.Get(model.NewApiPath("/v1/users/by-email"), middleware.WrapResponse(c.GetUserByEmail))
	r.Put(model.NewApiPath("/v1/users/{id}"), middleware.WrapResponse(c.UpdateUser))
	r.Patch(model.NewApiPath("/v1/users/{id}"), middleware.WrapResponse(c.PatchUser))
	r.Delete(model.NewApiPath("/v1/users/{id}"), middleware.WrapResponse(c.DeleteUser))
}

// CreateUser godoc
// @Summary Создать пользователя
// @Tags User
// @Accept json
// @Produce json
// @Param reg_id query string true "Идентификатор формы регистрации" format(uuid)
// @Param request body CreateUser true "Данные пользователя"
// @Success 201 {object} UserResponse
// @Failure 400 {object} model.ApiError "Некорректное тело запроса"
// @Failure 500 {object} model.ApiError
// @Router /api/user-service/v1/users [post]
func (c *UserController) CreateUser(_ http.ResponseWriter, r *http.Request) (any, error) {
	regID, err := utils.GetQueryParam(r, "reg_id", utils.ParseUUID)
	if err != nil {
		return nil, err
	}

	var createUser CreateUser
	if err := utils.DecodeRequestAndValidate(r.Body, &createUser); err != nil {
		return nil, err
	}

	usr, err := c.Service.CreateUser(r.Context(), createUser.ToUser(), regID)
	if err != nil {
		return nil, err
	}

	return toUserResponse(usr), nil
}

// GetUserById godoc
// @Summary Получить пользователя по ID
// @Tags User
// @Produce json
// @Param id path string true "Идентификатор пользователя" format(uuid)
// @Success 200 {object} UserResponse
// @Failure 400 {object} model.ApiError "Некорректное тело запроса"
// @Failure 404 {object} model.ApiError "Пользователь не найден"
// @Failure 500 {object} model.ApiError
// @Router /api/user-service/v1/users/{id} [get]
func (c *UserController) GetUserById(_ http.ResponseWriter, r *http.Request) (any, error) {
	ID, err := utils.GetPathParam(r, "id", utils.ParseUUID)
	if err != nil {
		return nil, err
	}

	usr, err := c.Service.GetUserById(r.Context(), ID)
	return toUserResponse(usr), err
}

// GetUserByEmail godoc
// @Summary Получить пользователя по email
// @Tags User
// @Produce json
// @Param email query string true "E-mail пользователя"
// @Success 200 {object} UserResponse
// @Failure 400 {object} model.ApiError "Некорректное тело запроса"
// @Failure 404 {object} model.ApiError "Пользователь не найден"
// @Failure 500 {object} model.ApiError
// @Router /api/user-service/v1/users/by-email [get]
func (c *UserController) GetUserByEmail(_ http.ResponseWriter, r *http.Request) (any, error) {
	emailParam, err := utils.GetQueryParam(r, "email", utils.ParseString)
	if err != nil {
		return nil, err
	}

	usr, err := c.Service.GetUserByEmail(r.Context(), emailParam)
	return toUserResponse(usr), err
}

// UpdateUser godoc
// @Summary Обновить пользователя
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "Идентификатор пользователя" format(uuid)
// @Param request body UpdateUser true "Данные пользователя"
// @Success 200 {object} UserResponse
// @Failure 400 {object} model.ApiError "Некорректное тело запроса"
// @Failure 404 {object} model.ApiError "Пользователь не найден"
// @Failure 500 {object} model.ApiError
// @Router /api/user-service/v1/users/{id} [put]
func (c *UserController) UpdateUser(_ http.ResponseWriter, r *http.Request) (any, error) {
	ID, err := utils.GetPathParam(r, "id", utils.ParseUUID)
	if err != nil {
		return nil, err
	}

	var updateUser UpdateUser
	if err := utils.DecodeRequestAndValidate(r.Body, &updateUser); err != nil {
		return nil, err
	}

	usr, err := c.Service.UpdateUser(r.Context(), updateUser.ToUser(ID))
	return toUserResponse(usr), err
}

// PatchUser godoc
// @Summary Обновить пользователя
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "Идентификатор пользователя" format(uuid)
// @Param request body PatchUser true "Данные пользователя"
// @Success 200 {object} UserResponse
// @Failure 400 {object} model.ApiError "Некорректное тело запроса"
// @Failure 404 {object} model.ApiError "Пользователь не найден"
// @Failure 500 {object} model.ApiError
// @Router /api/user-service/v1/users/{id} [patch]
func (c *UserController) PatchUser(_ http.ResponseWriter, r *http.Request) (any, error) {
	ID, err := utils.GetPathParam(r, "id", utils.ParseUUID)
	if err != nil {
		return nil, err
	}

	var patchUser PatchUser
	if err := utils.DecodeRequestAndValidate(r.Body, &patchUser); err != nil {
		return nil, err
	}

	usr, err := c.Service.UpdateUser(r.Context(), patchUser.ToUser(ID))
	return toUserResponse(usr), err
}

// DeleteUser godoc
// @Summary Удалить пользователя
// @Tags User
// @Produce json
// @Param id path string true "Идентификатор пользователя" format(uuid)
// @Success 200 {object} UserResponse
// @Failure 400 {object} model.ApiError "Некорректное тело запроса"
// @Failure 404 {object} model.ApiError "Пользователь не найден"
// @Failure 500 {object} model.ApiError
// @Router /api/user-service/v1/users/{id} [delete]
func (c *UserController) DeleteUser(_ http.ResponseWriter, r *http.Request) (any, error) {
	ID, err := utils.GetPathParam(r, "id", utils.ParseUUID)
	if err != nil {
		return nil, err
	}

	usr, err := c.Service.DeleteUser(r.Context(), ID)
	return toUserResponse(usr), err
}

func toUserResponse(dbUser *model.User) *UserResponse {
	if dbUser == nil {
		return nil
	}

	return &UserResponse{
		ID:            *dbUser.ID,
		FirstName:     *dbUser.FirstName,
		LastName:      dbUser.LastName,
		MiddleName:    dbUser.MiddleName,
		Email:         *dbUser.Email,
		Phone:         dbUser.Phone,
		Notifications: *dbUser.Notifications,
		Type:          *dbUser.Type,
		CreatedAt:     *dbUser.CreatedAt,
		UpdatedAt:     dbUser.UpdatedAt,
		DeletedAt:     dbUser.DeletedAt,
	}
}
