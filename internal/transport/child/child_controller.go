package child

import (
	"UserService/internal/middleware"
	"UserService/internal/model"
	"UserService/internal/service"
	"UserService/internal/utils"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type ChildrenController struct {
	Service *service.ChildrenService
}

func NewChildrenController(service *service.ChildrenService) *ChildrenController {
	return &ChildrenController{
		Service: service,
	}
}

func (c *ChildrenController) RegisterRoutes(r chi.Router) {
	r.Post(
		model.NewApiPath("/v1/children"),
		middleware.WrapResponseWithCode(c.CreateChild, http.StatusCreated),
	)
	r.Get(
		model.NewApiPath("/v1/children/{id}"),
		middleware.WrapResponse(c.GetChild),
	)
	r.Get(
		model.NewApiPath("/v1/children/by-parent/{id}"),
		middleware.WrapResponse(c.GetChildByParentID),
	)
	r.Put(
		model.NewApiPath("/v1/children/{id}"),
		middleware.WrapResponse(c.UpdateChild),
	)
	r.Patch(
		model.NewApiPath("/v1/children/{id}"),
		middleware.WrapResponse(c.PatchChild),
	)
	r.Delete(
		model.NewApiPath("/v1/children/{id}"),
		middleware.WrapResponse(c.DeleteChild),
	)
}

// CreateChild godoc
// @Summary Создать сущность ребенка
// @Tags Child
// @Accept json
// @Produce json
// @Param request body CreateChild true "Данные ребенка"
// @Success 201 {object} ChildResponse
// @Failure 400 {object} model.ApiError "Некорректное тело запроса"
// @Failure 500 {object} model.ApiError
// @Router /api/user-service/v1/children [post]
func (c *ChildrenController) CreateChild(_ http.ResponseWriter, r *http.Request) (any, error) {
	var createChild CreateChild
	if err := utils.DecodeRequestAndValidate(r.Body, &createChild); err != nil {
		return nil, err
	}

	child, err := c.Service.Create(r.Context(), createChild.ToChild())
	return toResponse(child), err
}

// GetChild godoc
// @Summary Получить сущность ребенка
// @Tags Child
// @Produce json
// @Param id path string true "Идентификатор ребенка" format(uuid)
// @Success 200 {object} ChildResponse
// @Failure 400 {object} model.ApiError "Некорректный идентификатор"
// @Failure 404 {object} model.ApiError "Ребенок не найден"
// @Failure 500 {object} model.ApiError
// @Router /api/user-service/v1/children/{id} [get]
func (c *ChildrenController) GetChild(_ http.ResponseWriter, r *http.Request) (any, error) {
	ID, err := utils.GetPathParam(r, "id", utils.ParseUUID)
	if err != nil {
		return nil, err
	}

	child, err := c.Service.GetByID(r.Context(), ID)
	return toResponse(child), err
}

// GetChildByParentID godoc
// @Summary Получить сущность ребенка по идентификатору родителя
// @Tags Child
// @Produce json
// @Param id path string true "Идентификатор родителя" format(uuid)
// @Success 200 {object} ChildResponse
// @Failure 400 {object} model.ApiError "Некорректный идентификатор"
// @Failure 404 {object} model.ApiError "Ребенок не найден"
// @Failure 500 {object} model.ApiError
// @Router /api/user-service/v1/children/by-parent/{id} [get]
func (c *ChildrenController) GetChildByParentID(_ http.ResponseWriter, r *http.Request) (any, error) {
	parentID, err := utils.GetPathParam(r, "id", utils.ParseUUID)
	if err != nil {
		return nil, err
	}
	child, err := c.Service.GetByParentID(r.Context(), parentID)
	return toResponse(child), err
}

// UpdateChild godoc
// @Summary Обновить сущность ребенка
// @Tags Child
// @Accept json
// @Produce json
// @Param id path string true "Идентификатор ребенка" format(uuid)
// @Success 200 {object} ChildResponse
// @Failure 400 {object} model.ApiError "Некорректный идентификатор"
// @Failure 404 {object} model.ApiError "Ребенок не найден"
// @Failure 500 {object} model.ApiError
// @Router /api/user-service/v1/children/{id} [put]
func (c *ChildrenController) UpdateChild(_ http.ResponseWriter, r *http.Request) (any, error) {
	ID, err := utils.GetPathParam(r, "id", utils.ParseUUID)
	if err != nil {
		return nil, err
	}

	var updateChild UpdateChild
	if err := utils.DecodeRequestAndValidate(r.Body, &updateChild); err != nil {
		return nil, err
	}

	child, err := c.Service.Update(r.Context(), updateChild.ToChild(ID))
	return toResponse(child), err
}

// PatchChild godoc
// @Summary Обновить сущность ребенка
// @Tags Child
// @Accept json
// @Produce json
// @Param id path string true "Идентификатор ребенка" format(uuid)
// @Success 200 {object} ChildResponse
// @Failure 400 {object} model.ApiError "Некорректный идентификатор"
// @Failure 404 {object} model.ApiError "Ребенок не найден"
// @Failure 500 {object} model.ApiError
// @Router /api/user-service/v1/children/{id} [patch]
func (c *ChildrenController) PatchChild(_ http.ResponseWriter, r *http.Request) (any, error) {
	ID, err := utils.GetPathParam(r, "id", utils.ParseUUID)
	if err != nil {
		return nil, err
	}

	var patchChild PatchChild
	if err := utils.DecodeRequestAndValidate(r.Body, &patchChild); err != nil {
		return nil, err
	}

	child, err := c.Service.Update(r.Context(), patchChild.ToChild(ID))
	return toResponse(child), err
}

// DeleteChild godoc
// @Summary Удалить сущность ребенка
// @Tags Child
// @Produce json
// @Param id path string true "Идентификатор ребенка" format(uuid)
// @Success 200 {object} ChildResponse
// @Failure 400 {object} model.ApiError "Некорректный идентификатор"
// @Failure 404 {object} model.ApiError "Ребенок не найден"
// @Failure 500 {object} model.ApiError
// @Router /api/user-service/v1/children/{id} [delete]
func (c *ChildrenController) DeleteChild(_ http.ResponseWriter, r *http.Request) (any, error) {
	ID, err := utils.GetPathParam(r, "id", utils.ParseUUID)
	if err != nil {
		return nil, err
	}

	child, err := c.Service.Delete(r.Context(), ID)
	return toResponse(child), err
}

func toResponse(child *model.Child) *ChildResponse {
	if child == nil {
		return nil
	}

	var birthday *string
	if child.Birthday != nil {
		actualBirthday := *child.Birthday
		birthday = utils.StringPtr(actualBirthday.Format("2006-01-02"))
	}

	return &ChildResponse{
		*child.ID,
		*child.ParentID,
		*child.Name,
		*child.Age,
		child.Gender,
		birthday,
		*child.CreatedAt,
		child.UpdatedAt,
		child.DeletedAt,
	}
}
