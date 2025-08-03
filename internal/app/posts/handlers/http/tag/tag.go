package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type TagHandler struct {
	tagUseCase tagUseCase
	logger     logger
}

func NewTagHandler(tagUseCase tagUseCase, logger logger) *TagHandler {
	return &TagHandler{tagUseCase: tagUseCase, logger: logger}
}

// Create
//
// @Summary Create tag
// @Tags tag
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param form body TagCreateDTO true "Create tag request"
// @Success 201 {object} TagDTO "Created tag"
// @Failure 400 {object} errs.Error "Invalid request body or validation error"
// @Failure 401 {object} errs.Error "Unauthorized"
// @Failure 404 {object} errs.Error "Not found"
// @Failure 500 {object} errs.Error "Internal server error"
// @Router /api/v1/posts/tags/ [POST]
func (h *TagHandler) Create(w http.ResponseWriter, r *http.Request) {
	createDTO, err := NewTagCreateDTO(r)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	create, err := createDTO.toEntity()
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	tag, err := h.tagUseCase.Create(r.Context(), create)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	response, err := NewTagDTO(tag)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, response)
}

// Get
//
// @Summary Get tag by id
// @Tags tag
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "UUID"
// @Success 200 {object} TagDTO "Requested tag"
// @Failure 400 {object} errs.Error "Invalid request body or validation error"
// @Failure 401 {object} errs.Error "Unauthorized"
// @Failure 404 {object} errs.Error "Not found"
// @Failure 500 {object} errs.Error "Internal server error"
// @Router /api/v1/posts/tags/{id} [GET]
func (h *TagHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	tag, err := h.tagUseCase.Get(r.Context(), id)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	response, err := NewTagDTO(tag)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}

// List
//
// @Summary List of tags
// @Tags tag
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param filter query TagFilterDTO true "Filter of tags"
// @Success 200 {object} TagListDTO "Filtered list of tags"
// @Failure 400 {object} errs.Error "Invalid request body or validation error"
// @Failure 401 {object} errs.Error "Unauthorized"
// @Failure 404 {object} errs.Error "Not found"
// @Failure 500 {object} errs.Error "Internal server error"
// @Router /api/v1/posts/tags/ [GET]
func (h *TagHandler) List(w http.ResponseWriter, r *http.Request) {
	filterDTO, err := NewTagFilterDTO(r)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	filter, err := filterDTO.toEntity()
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	tags, count, err := h.tagUseCase.List(r.Context(), filter)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	response, err := NewTagListDTO(tags, count)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}

// Update
//
// @Summary Update tag
// @Tags tag
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "UUID"
// @Param form body TagUpdateDTO true "Update tag request"
// @Success 200 {object} TagDTO "Updated tag"
// @Failure 400 {object} errs.Error "Invalid request body or validation error"
// @Failure 401 {object} errs.Error "Unauthorized"
// @Failure 404 {object} errs.Error "Not found"
// @Failure 500 {object} errs.Error "Internal server error"
// @Router /api/v1/posts/tags/{id} [PATCH]
func (h *TagHandler) Update(w http.ResponseWriter, r *http.Request) {
	updateDTO, err := NewTagUpdateDTO(r)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	update, err := updateDTO.toEntity()
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	tag, err := h.tagUseCase.Update(r.Context(), update)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	response, err := NewTagDTO(tag)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}

// Delete
//
// @Summary Delete tag by id
// @Tags tag
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "UUID"
// @Success 204 "No content"
// @Failure 400 {object} errs.Error "Invalid request body or validation error"
// @Failure 401 {object} errs.Error "Unauthorized"
// @Failure 404 {object} errs.Error "Not found"
// @Failure 500 {object} errs.Error "Internal server error"
// @Router /api/v1/posts/tags/{id} [DELETE]
func (h *TagHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	if err := h.tagUseCase.Delete(r.Context(), id); err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	render.Status(r, http.StatusNoContent)
	render.NoContent(w, r)
}
func (h *TagHandler) ChiRouter() chi.Router {
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
