package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	httpServer "github.com/mikalai-mitsin/example/internal/pkg/http"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type LikeHandler struct {
	likeUseCase likeUseCase
	logger      logger
}

func NewLikeHandler(likeUseCase likeUseCase, logger logger) *LikeHandler {
	return &LikeHandler{likeUseCase: likeUseCase, logger: logger}
}

// Create
//
// @Summary Create like
// @Tags like
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param form body LikeCreateDTO true "Create like request"
// @Success 201 {object} LikeDTO "Created like"
// @Failure 400 {object} errs.Error "Invalid request body or validation error"
// @Failure 401 {object} errs.Error "Unauthorized"
// @Failure 404 {object} errs.Error "Not found"
// @Failure 500 {object} errs.Error "Internal server error"
// @Router /api/v1/posts/likes/ [POST]
func (h *LikeHandler) Create(w http.ResponseWriter, r *http.Request) {
	createDTO, err := NewLikeCreateDTO(r)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	create, err := createDTO.toEntity()
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	like, err := h.likeUseCase.Create(r.Context(), create)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	response, err := NewLikeDTO(like)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, response)
}

// Get
//
// @Summary Get like by id
// @Tags like
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "UUID"
// @Success 200 {object} LikeDTO "Requested like"
// @Failure 400 {object} errs.Error "Invalid request body or validation error"
// @Failure 401 {object} errs.Error "Unauthorized"
// @Failure 404 {object} errs.Error "Not found"
// @Failure 500 {object} errs.Error "Internal server error"
// @Router /api/v1/posts/likes/{id} [GET]
func (h *LikeHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	like, err := h.likeUseCase.Get(r.Context(), id)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	response, err := NewLikeDTO(like)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}

// List
//
// @Summary List of likes
// @Tags like
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param filter query LikeFilterDTO true "Filter of likes"
// @Success 200 {object} LikeListDTO "Filtered list of likes"
// @Failure 400 {object} errs.Error "Invalid request body or validation error"
// @Failure 401 {object} errs.Error "Unauthorized"
// @Failure 404 {object} errs.Error "Not found"
// @Failure 500 {object} errs.Error "Internal server error"
// @Router /api/v1/posts/likes/ [GET]
func (h *LikeHandler) List(w http.ResponseWriter, r *http.Request) {
	filterDTO, err := NewLikeFilterDTO(r)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	filter, err := filterDTO.toEntity()
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	likes, count, err := h.likeUseCase.List(r.Context(), filter)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	response, err := NewLikeListDto(likes, count)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}

// Update
//
// @Summary Update like
// @Tags like
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "UUID"
// @Param form body LikeUpdateDTO true "Update like request"
// @Success 200 {object} LikeDTO "Updated like"
// @Failure 400 {object} errs.Error "Invalid request body or validation error"
// @Failure 401 {object} errs.Error "Unauthorized"
// @Failure 404 {object} errs.Error "Not found"
// @Failure 500 {object} errs.Error "Internal server error"
// @Router /api/v1/posts/likes/{id} [PATCH]
func (h *LikeHandler) Update(w http.ResponseWriter, r *http.Request) {
	updateDTO, err := NewLikeUpdateDTO(r)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	update, err := updateDTO.toEntity()
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	like, err := h.likeUseCase.Update(r.Context(), update)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	response, err := NewLikeDTO(like)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}

// Delete
//
// @Summary Delete like by id
// @Tags like
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "UUID"
// @Success 200 {object} LikeDTO "Updated like"
// @Failure 400 {object} errs.Error "Invalid request body or validation error"
// @Failure 401 {object} errs.Error "Unauthorized"
// @Failure 404 {object} errs.Error "Not found"
// @Failure 500 {object} errs.Error "Internal server error"
// @Router /api/v1/posts/likes/{id} [DELETE]
func (h *LikeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	like, err := h.likeUseCase.Delete(r.Context(), id)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	response, err := NewLikeDTO(like)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}
func (h *LikeHandler) router() chi.Router {
	router := chi.NewRouter()
	router.Route("/", func(g chi.Router) {
		g.Post("/", h.Create)
		g.Get("/", h.List)
		g.Get("/{id}", h.Get)
		g.Patch("/{id}", h.Update)
		g.Delete("/{id}", h.Delete)
	})
	return router
}
func (h *LikeHandler) RegisterHTTP(httpServer *httpServer.Server) error {
	httpServer.Mount("/api/v1/posts/likes", h.router())
	return nil
}
