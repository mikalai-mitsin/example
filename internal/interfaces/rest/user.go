package rest

import (
	"fmt"
	"net/http"

	"github.com/018bf/example/internal/domain/interceptors"
	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/pkg/log"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userInterceptor interceptors.UserInterceptor
	logger          log.Logger
}

func NewUserHandler(userInterceptor interceptors.UserInterceptor, logger log.Logger) *UserHandler {
	return &UserHandler{userInterceptor: userInterceptor, logger: logger}
}

func (h *UserHandler) Register(router *gin.RouterGroup) {
	group := router.Group("/users")
	group.POST("/", h.Create)
	group.GET("/", h.List)
	group.GET("/:id", h.Get)
	group.PATCH("/:id", h.Update)
	group.DELETE("/:id", h.Delete)
}

// Create        godoc
// @Summary      Store a new User
// @Description  Takes a User JSON and store in DB. Return saved JSON.
// @Tags         User
// @Produce      json
// @Param        User  body   models.UserCreate  true  "User JSON"
// @Success      201   {object}  models.User
// @Failure        400   {object}  errs.Error
// @Failure        401   {object}  errs.Error
// @Failure        403   {object}  errs.Error
// @Failure        404   {object}  errs.Error
// @Failure        405   {object}  errs.Error
// @Failure        500   {object}  errs.Error
// @Failure        503   {object}  errs.Error
// @Router       /users [post]
func (h *UserHandler) Create(ctx *gin.Context) {
	requestUser := ctx.Request.Context().Value(UserContextKey).(*models.User)
	create := &models.UserCreate{}
	_ = ctx.Bind(create)
	users, err := h.userInterceptor.Create(ctx.Request.Context(), create, requestUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusCreated, users)
}

// List          godoc
// @Summary      List User array
// @Description  Responds with the list of all User as JSON.
// @Tags         User
// @Produce      json
// @Param        filter  query models.UserFilter false "User filter"
// @Success      200  {array}  models.User
// @Failure      400  {object}  errs.Error
// @Failure      401  {object}  errs.Error
// @Failure      402  {object}  errs.Error
// @Failure      404  {object}  errs.Error
// @Failure      500  {object}  errs.Error
// @Router       /users [get]
func (h *UserHandler) List(ctx *gin.Context) {
	requestUser := ctx.Request.Context().Value(UserContextKey).(*models.User)
	filter := &models.UserFilter{}
	_ = ctx.Bind(filter)
	users, count, err := h.userInterceptor.List(ctx.Request.Context(), filter, requestUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	ctx.Header("count", fmt.Sprint(count))
	ctx.JSON(http.StatusOK, users)
}

// Get           godoc
// @Summary      Get single User by UUID
// @Description  Returns the User whose UUID value matches the UUID.
// @Tags         User
// @Produce      json
// @Param        uuid  path      string  true  "search User by UUID"
// @Success      200  {object}  models.User
// @Failure        400   {object}  errs.Error
// @Failure        401   {object}  errs.Error
// @Failure        403   {object}  errs.Error
// @Failure        404   {object}  errs.Error
// @Failure        405   {object}  errs.Error
// @Failure        500   {object}  errs.Error
// @Failure        503   {object}  errs.Error
// @Router       /users/{uuid} [get]
func (h *UserHandler) Get(c *gin.Context) {
	requestUser := c.Request.Context().Value(UserContextKey).(*models.User)
	users, err := h.userInterceptor.Get(
		c.Request.Context(),
		models.UUID(c.Param("id")),
		requestUser,
	)
	if err != nil {
		decodeError(c, err)
		return
	}
	c.JSON(http.StatusOK, users)
}

// Update        godoc
// @Summary      Update User by UUID
// @Description  Returns the updated User.
// @Tags         User
// @Produce      json
// @Param        uuid  path      string  true  "update User by UUID"
// @Param        User  body   models.UserUpdate  true  "User JSON"
// @Success      200  {object}  models.User
// @Failure        400   {object}  errs.Error
// @Failure        401   {object}  errs.Error
// @Failure        403   {object}  errs.Error
// @Failure        404   {object}  errs.Error
// @Failure        405   {object}  errs.Error
// @Failure        500   {object}  errs.Error
// @Failure        503   {object}  errs.Error
// @Router       /users/{uuid} [PATCH]
func (h *UserHandler) Update(c *gin.Context) {
	requestUser := c.Request.Context().Value(UserContextKey).(*models.User)
	update := &models.UserUpdate{}
	_ = c.Bind(update)
	update.ID = models.UUID(c.Param("id"))
	users, err := h.userInterceptor.Update(c.Request.Context(), update, requestUser)
	if err != nil {
		decodeError(c, err)
		return
	}
	c.JSON(http.StatusOK, users)
}

// Delete        godoc
// @Summary      Delete single User by UUID
// @Description  Delete the User whose UUID value matches the UUID.
// @Tags         User
// @Param        uuid  path      string  true  "delete User by UUID"
// @Success      204
// @Failure        400   {object}  errs.Error
// @Failure        401   {object}  errs.Error
// @Failure        403   {object}  errs.Error
// @Failure        404   {object}  errs.Error
// @Failure        405   {object}  errs.Error
// @Failure        500   {object}  errs.Error
// @Failure        503   {object}  errs.Error
// @Router       /users/{uuid} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	requestUser := c.Request.Context().Value(UserContextKey).(*models.User)
	err := h.userInterceptor.Delete(c.Request.Context(), models.UUID(c.Param("id")), requestUser)
	if err != nil {
		decodeError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
