package rest

import (
	"fmt"
	"net/http"

	"github.com/018bf/example/internal/domain/interceptors"
	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/pkg/log"
	"github.com/gin-gonic/gin"
)

type MarkHandler struct {
	markInterceptor interceptors.MarkInterceptor
	logger          log.Logger
}

func NewMarkHandler(markInterceptor interceptors.MarkInterceptor, logger log.Logger) *MarkHandler {
	return &MarkHandler{markInterceptor: markInterceptor, logger: logger}
}

func (h *MarkHandler) Register(router *gin.Engine) {
	group := router.Group("/marks")
	group.POST("/", h.Create)
	group.GET("/", h.List)
	group.GET("/:id", h.Get)
	group.PATCH("/:id", h.Update)
	group.DELETE("/:id", h.Delete)
}

func (h *MarkHandler) Create(ctx *gin.Context) {
	requestUser := ctx.Request.Context().Value(UserContextKey).(*models.User)
	create := &models.MarkCreate{}
	if err := ctx.Bind(create); err != nil {
		return
	}
	marks, err := h.markInterceptor.Create(ctx.Request.Context(), create, requestUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusOK, marks)
}

func (h *MarkHandler) List(ctx *gin.Context) {
	requestUser := ctx.Request.Context().Value(UserContextKey).(*models.User)
	filter := &models.MarkFilter{}
	if err := ctx.Bind(filter); err != nil {
		return
	}
	marks, count, err := h.markInterceptor.List(ctx.Request.Context(), filter, requestUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	ctx.Header("count", fmt.Sprint(count))
	ctx.JSON(http.StatusOK, marks)
}

func (h *MarkHandler) Get(c *gin.Context) {
	requestUser := c.Request.Context().Value(UserContextKey).(*models.User)
	marks, err := h.markInterceptor.Get(c.Request.Context(), c.Param("id"), requestUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, marks)
}

func (h *MarkHandler) Update(c *gin.Context) {
	requestUser := c.Request.Context().Value(UserContextKey).(*models.User)
	update := &models.MarkUpdate{}
	if err := c.Bind(update); err != nil {
		return
	}
	update.ID = c.Param("id")
	marks, err := h.markInterceptor.Update(c.Request.Context(), update, requestUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, marks)
}

func (h *MarkHandler) Delete(c *gin.Context) {
	requestUser := c.Request.Context().Value(UserContextKey).(*models.User)
	err := h.markInterceptor.Delete(c.Request.Context(), c.Param("id"), requestUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.Status(http.StatusNoContent)
}
