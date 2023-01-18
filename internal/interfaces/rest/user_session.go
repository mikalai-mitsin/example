package rest

import (
	"fmt"
	"net/http"

	"github.com/018bf/example/internal/domain/interceptors"
	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/pkg/log"
	"github.com/gin-gonic/gin"
)

type UserSessionHandler struct {
	userSessionInterceptor interceptors.UserSessionInterceptor
	logger          log.Logger
}

func NewUserSessionHandler(userSessionInterceptor interceptors.UserSessionInterceptor, logger log.Logger) *UserSessionHandler {
	return &UserSessionHandler{userSessionInterceptor: userSessionInterceptor, logger: logger}
}

func (h *UserSessionHandler) Register(router *gin.Engine) {
	group := router.Group("/user_sessions")
	group.POST("/", h.Create)
	group.GET("/", h.List)
	group.GET("/:id", h.Get)
	group.PATCH("/:id", h.Update)
	group.DELETE("/:id", h.Delete)
}

func (h *UserSessionHandler) Create(ctx *gin.Context) {
	requestUser := ctx.Request.Context().Value(UserContextKey).(*models.User)
	create := &models.UserSessionCreate{}
	if err := ctx.Bind(create); err != nil {
		return
	}
	marks, err := h.userSessionInterceptor.Create(ctx.Request.Context(), create, requestUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusOK, marks)
}

func (h *UserSessionHandler) List(ctx *gin.Context) {
	requestUser := ctx.Request.Context().Value(UserContextKey).(*models.User)
	filter := &models.UserSessionFilter{}
	if err := ctx.Bind(filter); err != nil {
		return
	}
	marks, count, err := h.userSessionInterceptor.List(ctx.Request.Context(), filter, requestUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	ctx.Header("count", fmt.Sprint(count))
	ctx.JSON(http.StatusOK, marks)
}

func (h *UserSessionHandler) Get(c *gin.Context) {
	requestUser := c.Request.Context().Value(UserContextKey).(*models.User)
	marks, err := h.userSessionInterceptor.Get(c.Request.Context(), c.Param("id"), requestUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, marks)
}

func (h *UserSessionHandler) Update(c *gin.Context) {
	requestUser := c.Request.Context().Value(UserContextKey).(*models.User)
	update := &models.UserSessionUpdate{}
	if err := c.Bind(update); err != nil {
		return
	}
	update.ID = c.Param("id")
	marks, err := h.userSessionInterceptor.Update(c.Request.Context(), update, requestUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, marks)
}

func (h *UserSessionHandler) Delete(c *gin.Context) {
	requestUser := c.Request.Context().Value(UserContextKey).(*models.User)
	err := h.userSessionInterceptor.Delete(c.Request.Context(), c.Param("id"), requestUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.Status(http.StatusNoContent)
}
