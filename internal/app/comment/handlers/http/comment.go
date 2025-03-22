package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type CommentHandler struct {
	commentUseCase commentUseCase
	logger         logger
}

func NewCommentHandler(commentUseCase commentUseCase, logger logger) *CommentHandler {
	return &CommentHandler{commentUseCase: commentUseCase, logger: logger}
}

// Create
//
// @Summary Create comment
// @Tags comment
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param form body CommentCreateDTO true "Create comment request"
// @Success 201 {object} CommentDTO "Created comment"
// @Failure 400 {object} errs.Error "Invalid request body or validation error"
// @Failure 401 {object} errs.Error "Unauthorized"
// @Failure 404 {object} errs.Error "Not found"
// @Failure 500 {object} errs.Error "Internal server error"
// @Router /api/v1/comments/ [POST]
func (h *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
	createDTO, err := NewCommentCreateDTO(r)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	create, err := createDTO.toEntity()
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	comment, err := h.commentUseCase.Create(r.Context(), create)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	response, err := NewCommentDTO(comment)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, response)
}

// Get
//
// @Summary Get comment by id
// @Tags comment
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "UUID"
// @Success 200 {object} CommentDTO "Requested comment"
// @Failure 400 {object} errs.Error "Invalid request body or validation error"
// @Failure 401 {object} errs.Error "Unauthorized"
// @Failure 404 {object} errs.Error "Not found"
// @Failure 500 {object} errs.Error "Internal server error"
// @Router /api/v1/comments/{id} [GET]
func (h *CommentHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := uuid.UUID(chi.URLParam(r, "id"))
	comment, err := h.commentUseCase.Get(r.Context(), id)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	response, err := NewCommentDTO(comment)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}

// List
//
// @Summary List of comments
// @Tags comment
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param filter query CommentFilterDTO true "Filter of comments"
// @Success 200 {array} CommentListDTO "Filtered list of comments"
// @Failure 400 {object} errs.Error "Invalid request body or validation error"
// @Failure 401 {object} errs.Error "Unauthorized"
// @Failure 404 {object} errs.Error "Not found"
// @Failure 500 {object} errs.Error "Internal server error"
// @Router /api/v1/comments/ [GET]
func (h *CommentHandler) List(w http.ResponseWriter, r *http.Request) {
	filterDTO, err := NewCommentFilterDTO(r)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	filter, err := filterDTO.toEntity()
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	comments, count, err := h.commentUseCase.List(r.Context(), filter)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	response, err := NewCommentListDTO(comments, count)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}

// Update
//
// @Summary Update comment
// @Tags comment
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "UUID"
// @Param form body CommentUpdateDTO true "Update comment request"
// @Success 200 {object} CommentDTO "Updated comment"
// @Failure 400 {object} errs.Error "Invalid request body or validation error"
// @Failure 401 {object} errs.Error "Unauthorized"
// @Failure 404 {object} errs.Error "Not found"
// @Failure 500 {object} errs.Error "Internal server error"
// @Router /api/v1/comments/{id} [PATCH]
func (h *CommentHandler) Update(w http.ResponseWriter, r *http.Request) {
	updateDTO, err := NewCommentUpdateDTO(r)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	update, err := updateDTO.toEntity()
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	comment, err := h.commentUseCase.Update(r.Context(), update)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	response, err := NewCommentDTO(comment)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}

// Delete
//
// @Summary Delete comment by id
// @Tags comment
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "UUID"
// @Success 204 "No content"
// @Failure 400 {object} errs.Error "Invalid request body or validation error"
// @Failure 401 {object} errs.Error "Unauthorized"
// @Failure 404 {object} errs.Error "Not found"
// @Failure 500 {object} errs.Error "Internal server error"
// @Router /api/v1/comments/{id} [DELETE]
func (h *CommentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := uuid.UUID(chi.URLParam(r, "id"))
	if err := h.commentUseCase.Delete(r.Context(), id); err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	render.Status(r, http.StatusNoContent)
	render.NoContent(w, r)
}
func (h *CommentHandler) ChiRouter() chi.Router {
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
