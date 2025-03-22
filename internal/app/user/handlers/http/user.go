package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type UserHandler struct {
	userUseCase userUseCase
	logger      logger
}

func NewUserHandler(userUseCase userUseCase, logger logger) *UserHandler {
	return &UserHandler{userUseCase: userUseCase, logger: logger}
}

// Create
//
// @Summary Create user
// @Tags user
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param form body UserCreateDTO true "Create user request"
// @Success 201 {object} UserDTO "Created user"
// @Failure 400 {object} errs.Error "Invalid request body or validation error"
// @Failure 401 {object} errs.Error "Unauthorized"
// @Failure 404 {object} errs.Error "Not found"
// @Failure 500 {object} errs.Error "Internal server error"
// @Router /api/v1/users/ [POST]
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	createDTO, err := NewUserCreateDTO(r)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	create, err := createDTO.toEntity()
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	user, err := h.userUseCase.Create(r.Context(), create)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	response, err := NewUserDTO(user)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, response)
}

// Get
//
// @Summary Get user by id
// @Tags user
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "UUID"
// @Success 200 {object} UserDTO "Requested user"
// @Failure 400 {object} errs.Error "Invalid request body or validation error"
// @Failure 401 {object} errs.Error "Unauthorized"
// @Failure 404 {object} errs.Error "Not found"
// @Failure 500 {object} errs.Error "Internal server error"
// @Router /api/v1/users/{id} [GET]
func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := uuid.UUID(chi.URLParam(r, "id"))
	user, err := h.userUseCase.Get(r.Context(), id)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	response, err := NewUserDTO(user)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}

// List
//
// @Summary List of users
// @Tags user
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param filter query UserFilterDTO true "Filter of users"
// @Success 200 {array} UserListDTO "Filtered list of users"
// @Failure 400 {object} errs.Error "Invalid request body or validation error"
// @Failure 401 {object} errs.Error "Unauthorized"
// @Failure 404 {object} errs.Error "Not found"
// @Failure 500 {object} errs.Error "Internal server error"
// @Router /api/v1/users/ [GET]
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	filterDTO, err := NewUserFilterDTO(r)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	filter, err := filterDTO.toEntity()
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	users, count, err := h.userUseCase.List(r.Context(), filter)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	response, err := NewUserListDTO(users, count)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}

// Update
//
// @Summary Update user
// @Tags user
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "UUID"
// @Param form body UserUpdateDTO true "Update user request"
// @Success 200 {object} UserDTO "Updated user"
// @Failure 400 {object} errs.Error "Invalid request body or validation error"
// @Failure 401 {object} errs.Error "Unauthorized"
// @Failure 404 {object} errs.Error "Not found"
// @Failure 500 {object} errs.Error "Internal server error"
// @Router /api/v1/users/{id} [PATCH]
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	updateDTO, err := NewUserUpdateDTO(r)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	update, err := updateDTO.toEntity()
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	user, err := h.userUseCase.Update(r.Context(), update)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	response, err := NewUserDTO(user)
	if err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}

// Delete
//
// @Summary Delete user by id
// @Tags user
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "UUID"
// @Success 204 "No content"
// @Failure 400 {object} errs.Error "Invalid request body or validation error"
// @Failure 401 {object} errs.Error "Unauthorized"
// @Failure 404 {object} errs.Error "Not found"
// @Failure 500 {object} errs.Error "Internal server error"
// @Router /api/v1/users/{id} [DELETE]
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := uuid.UUID(chi.URLParam(r, "id"))
	if err := h.userUseCase.Delete(r.Context(), id); err != nil {
		errs.RenderToHTTPResponse(err, w, r)
		return
	}
	render.Status(r, http.StatusNoContent)
	render.NoContent(w, r)
}
func (h *UserHandler) ChiRouter() chi.Router {
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
