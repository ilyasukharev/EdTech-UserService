package referral

import (
	"UserService/internal/middleware"
	"UserService/internal/model"
	"UserService/internal/service"
	"UserService/internal/utils"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type ReferralController struct {
	Service *service.ReferralService
}

func NewReferralController(service *service.ReferralService) *ReferralController {
	return &ReferralController{
		Service: service,
	}
}

func (c *ReferralController) RegisterRoutes(r chi.Router) {
	r.Post(model.NewApiPath("/v1/referrals"), middleware.WrapResponseWithCode(c.CreateReferral, http.StatusCreated))
	r.Get(model.NewApiPath("/v1/referrals/by-referrer/{id}"), middleware.WrapResponse(c.GetByReferrerID))
	r.Get(model.NewApiPath("/v1/referrals/by-referee/{id}"), middleware.WrapResponse(c.GetByRefereeID))
	r.Patch(model.NewApiPath("/v1/referrals/{id}"), middleware.WrapResponse(c.PatchReferral))
}

// CreateReferral godoc
// @Summary Создать запись с рефералкой
// @Tags Referral
// @Accept json
// @Produce json
// @Param request body CreateReferral true "Данные рефералки"
// @Success 201 {object} ReferralResponse
// @Failure 400 {object} model.ApiError "Некорректное тело запроса"
// @Failure 500 {object} model.ApiError
// @Router /api/user-service/v1/referrals [post]
func (c *ReferralController) CreateReferral(_ http.ResponseWriter, r *http.Request) (any, error) {
	var createReferral CreateReferral
	if err := utils.DecodeRequestAndValidate(r.Body, &createReferral); err != nil {
		return nil, err
	}
	referral, err := c.Service.CreateReferral(r.Context(), createReferral.ToReferral())
	return toReferralResponse(referral), err
}

// GetByReferrerID godoc
// @Summary Получить запись рефералки по ID referrer
// @Tags Referral
// @Produce json
// @Param id path string true "Идентификатор referrer" format(uuid)
// @Success 200 {object} ReferralResponse
// @Failure 400 {object} model.ApiError "Некорректное тело запроса"
// @Failure 404 {object} model.ApiError "Рефералка не найдена"
// @Failure 500 {object} model.ApiError
// @Router /api/user-service/v1/referrals/by-referrer/{id} [GET]
func (c *ReferralController) GetByReferrerID(_ http.ResponseWriter, r *http.Request) (any, error) {
	ID, err := utils.GetPathParam(r, "id", utils.ParseUUID)
	referral, err := c.Service.GetByReferrerID(r.Context(), ID)
	return toReferralResponse(referral), err
}

// GetByRefereeID godoc
// @Summary Получить запись рефералки по ID referee
// @Tags Referral
// @Produce json
// @Param id path string true "Идентификатор referee" format(uuid)
// @Success 200 {object} ReferralResponse
// @Failure 400 {object} model.ApiError "Некорректное тело запроса"
// @Failure 404 {object} model.ApiError "Рефералка не найдена"
// @Failure 500 {object} model.ApiError
// @Router /api/user-service/v1/referrals/by-referee/{id} [GET]
func (c *ReferralController) GetByRefereeID(_ http.ResponseWriter, r *http.Request) (any, error) {
	ID, err := utils.GetPathParam(r, "id", utils.ParseUUID)
	referral, err := c.Service.GetByRefereeID(r.Context(), ID)
	return toReferralResponse(referral), err
}

// PatchReferral godoc
// @Summary Обновить запись с рефералкой
// @Tags Referral
// @Accept json
// @Produce json
// @Param id path int true "Идентификатор рефералки"
// @Param request body PatchReferral true "Данные рефералки"
// @Success 200 {object} ReferralResponse
// @Failure 400 {object} model.ApiError "Некорректное тело запроса"
// @Failure 404 {object} model.ApiError "Рефералка не найдена"
// @Failure 500 {object} model.ApiError
// @Router /api/user-service/v1/referrals [patch]
func (c *ReferralController) PatchReferral(_ http.ResponseWriter, r *http.Request) (any, error) {
	ID, err := utils.GetPathParam(r, "id", utils.ParseInt64)

	var patchReferral PatchReferral
	if err = utils.DecodeRequestAndValidate(r.Body, &patchReferral); err != nil {
		return nil, err
	}

	referral, err := c.Service.PatchReferral(r.Context(), patchReferral.ToReferral(ID))
	return toReferralResponse(referral), err
}

func toReferralResponse(referral *model.Referral) *ReferralResponse {
	if referral == nil {
		return nil
	}

	return &ReferralResponse{
		*referral.ID,
		*referral.ReferrerID,
		*referral.RefereeID,
		*referral.Confirmed,
		*referral.CreatedAt,
		referral.UpdatedAt,
	}
}
