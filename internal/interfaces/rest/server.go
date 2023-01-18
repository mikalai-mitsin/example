package rest

import (
	"net/http"

	"github.com/018bf/example/pkg/log"
	"github.com/gin-gonic/gin"
)

type ctxKey string

const (
	UserContextKey		ctxKey	= "user"
	RequestIDContextKey	ctxKey	= "request_id"
)

func NewRouter(
	logger log.Logger,
	authMiddleware *AuthMiddleware, userSessionHandler *UserSessionHandler, equipmentHandler *EquipmentHandler, sessionHandler *SessionHandler, approachHandler *ApproachHandler, markHandler *MarkHandler,
) *gin.Engine {
	router := gin.Default()
	router.Use(authMiddleware.Middleware())
	router.Use(RequestMiddleware)
	router.GET("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})
	userSessionHandler.Register(router)
	equipmentHandler.Register(router)
	sessionHandler.Register(router)
	approachHandler.Register(router)
	markHandler.Register(router)
	return router
}
