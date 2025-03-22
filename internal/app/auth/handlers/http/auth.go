package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
)

type AuthHandler struct {
	authUseCase authUseCase
	logger      logger
}

func NewAuthHandler(authUseCase authUseCase, logger logger) *AuthHandler {
	return &AuthHandler{authUseCase: authUseCase, logger: logger}
}

// ObtainTokenPair
//
// @Tags auth
// @Accept json
// @Produce json
// @Param form body ObtainTokenDTO true "Obtain token pair"
// @Success 200 {object} TokenPairDTO "Token pair"
// @Failure 400 {object} errs.Error "Invalid request body or validation error"
// @Failure 401 {object} errs.Error "Unauthorized"
// @Failure 404 {object} errs.Error "Not found"
// @Failure 500 {object} errs.Error "Internal server error"
// @Router /api/v1/auth/obtain [POST]
func (h *AuthHandler) ObtainTokenPair(w http.ResponseWriter, r *http.Request) {
	createDTO, err := NewObtainTokenDTO(r)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	create, err := createDTO.toEntity()
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	tokenPair, err := h.authUseCase.CreateToken(r.Context(), create)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	response, err := NewTokenPairDTO(tokenPair)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}

// RefreshTokenPair
//
// @Tags auth
// @Accept json
// @Produce json
// @Param form body RefreshTokenDTO true "Refresh token pair"
// @Success 200 {object} TokenPairDTO "Token pair"
// @Failure 400 {object} errs.Error "Invalid request body or validation error"
// @Failure 401 {object} errs.Error "Unauthorized"
// @Failure 404 {object} errs.Error "Not found"
// @Failure 500 {object} errs.Error "Internal server error"
// @Router /api/v1/auth/refresh [POST]
func (h *AuthHandler) RefreshTokenPair(w http.ResponseWriter, r *http.Request) {
	refreshTokenDTO, err := NewRefreshTokenDTO(r)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	refreshToken, err := refreshTokenDTO.toEntity()
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	tokenPair, err := h.authUseCase.RefreshToken(r.Context(), refreshToken)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	response, err := NewTokenPairDTO(tokenPair)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}
func (h *AuthHandler) ChiRouter() chi.Router {
	router := chi.NewRouter()
	router.Route("/", func(g chi.Router) {
		g.Post("/obtain", h.ObtainTokenPair)
		g.Post("/refresh", h.RefreshTokenPair)
	})
	return router
}
