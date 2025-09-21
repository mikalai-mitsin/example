package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	httpServer "github.com/mikalai-mitsin/example/internal/pkg/http"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type PostHandler struct {
	postUseCase postUseCase
	logger      logger
}

func NewPostHandler(postUseCase postUseCase, logger logger) *PostHandler {
	return &PostHandler{postUseCase: postUseCase, logger: logger}
}

// Create
//
// @Summary Create post
// @Tags post
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param form body PostCreateDTO true "Create post request"
// @Success 201 {object} PostDTO "Created post"
// @Failure 400 {object} errs.Error "Invalid request body or validation error"
// @Failure 401 {object} errs.Error "Unauthorized"
// @Failure 404 {object} errs.Error "Not found"
// @Failure 500 {object} errs.Error "Internal server error"
// @Router /api/v1/posts/posts/ [POST]
func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	createDTO, err := NewPostCreateDTO(r)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	create, err := createDTO.toEntity()
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	post, err := h.postUseCase.Create(r.Context(), create)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	response, err := NewPostDTO(post)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, response)
}

// Get
//
// @Summary Get post by id
// @Tags post
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "UUID"
// @Success 200 {object} PostDTO "Requested post"
// @Failure 400 {object} errs.Error "Invalid request body or validation error"
// @Failure 401 {object} errs.Error "Unauthorized"
// @Failure 404 {object} errs.Error "Not found"
// @Failure 500 {object} errs.Error "Internal server error"
// @Router /api/v1/posts/posts/{id} [GET]
func (h *PostHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	post, err := h.postUseCase.Get(r.Context(), id)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	response, err := NewPostDTO(post)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}

// List
//
// @Summary List of posts
// @Tags post
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param filter query PostFilterDTO true "Filter of posts"
// @Success 200 {object} PostListDTO "Filtered list of posts"
// @Failure 400 {object} errs.Error "Invalid request body or validation error"
// @Failure 401 {object} errs.Error "Unauthorized"
// @Failure 404 {object} errs.Error "Not found"
// @Failure 500 {object} errs.Error "Internal server error"
// @Router /api/v1/posts/posts/ [GET]
func (h *PostHandler) List(w http.ResponseWriter, r *http.Request) {
	filterDTO, err := NewPostFilterDTO(r)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	filter, err := filterDTO.toEntity()
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	posts, count, err := h.postUseCase.List(r.Context(), filter)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	response, err := NewPostListDto(posts, count)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}

// Update
//
// @Summary Update post
// @Tags post
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "UUID"
// @Param form body PostUpdateDTO true "Update post request"
// @Success 200 {object} PostDTO "Updated post"
// @Failure 400 {object} errs.Error "Invalid request body or validation error"
// @Failure 401 {object} errs.Error "Unauthorized"
// @Failure 404 {object} errs.Error "Not found"
// @Failure 500 {object} errs.Error "Internal server error"
// @Router /api/v1/posts/posts/{id} [PATCH]
func (h *PostHandler) Update(w http.ResponseWriter, r *http.Request) {
	updateDTO, err := NewPostUpdateDTO(r)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	update, err := updateDTO.toEntity()
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	post, err := h.postUseCase.Update(r.Context(), update)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	response, err := NewPostDTO(post)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}

// Delete
//
// @Summary Delete post by id
// @Tags post
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "UUID"
// @Success 200 {object} PostDTO "Updated post"
// @Failure 400 {object} errs.Error "Invalid request body or validation error"
// @Failure 401 {object} errs.Error "Unauthorized"
// @Failure 404 {object} errs.Error "Not found"
// @Failure 500 {object} errs.Error "Internal server error"
// @Router /api/v1/posts/posts/{id} [DELETE]
func (h *PostHandler) Delete(w http.ResponseWriter, r *http.Request) {
	delDTO, err := NewPostDeleteDTO(r)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	del, err := delDTO.toEntity()
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	post, err := h.postUseCase.Delete(r.Context(), del)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	response, err := NewPostDTO(post)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}
func (h *PostHandler) router() chi.Router {
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
func (h *PostHandler) RegisterHTTP(httpServer *httpServer.Server) error {
	httpServer.Mount("/api/v1/posts/posts", h.router())
	return nil
}
