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

func (c *UserController) GetUserById(_ http.ResponseWriter, r *http.Request) (any, error) {
	ID, err := utils.GetPathParam(r, "id", utils.ParseUUID)
	if err != nil {
		return nil, err
	}

	usr, err := c.Service.GetUserById(r.Context(), ID)
	return toUserResponse(usr), err
}

func (c *UserController) GetUserByEmail(_ http.ResponseWriter, r *http.Request) (any, error) {
	emailParam, err := utils.GetQueryParam(r, "email", utils.ParseString)
	if err != nil {
		return nil, err
	}

	usr, err := c.Service.GetUserByEmail(r.Context(), emailParam)
	return toUserResponse(usr), err
}

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
		ChildName:     dbUser.ChildName,
		ChildAge:      dbUser.ChildAge,
		CreatedAt:     *dbUser.CreatedAt,
		UpdatedAt:     dbUser.UpdatedAt,
		DeletedAt:     dbUser.DeletedAt,
	}
}
