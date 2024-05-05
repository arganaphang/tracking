package handler

import (
	"application/internals/dto"
	"application/internals/model"
	"application/internals/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type IUserHandler interface {
	Register(ctx *gin.Context)
	ActiveUser(ctx *gin.Context)
}

type user struct {
	services service.Services
}

func NewUser(app *gin.Engine, services service.Services, middlewares ...gin.HandlerFunc) IUserHandler {
	u := user{services: services}

	app.POST("/register", u.Register)

	app.GET("/user/:id/active", middlewares[0], u.ActiveUser) // middleware[0] -> JWT

	return &u
}

func (u user) Register(ctx *gin.Context) {
	var body dto.RegisterRequest
	if err := ctx.ShouldBind(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.BadResponse{
			Message: "failed to serialize request body",
		})
		return
	}
	if err := u.services.User.Create(ctx, body); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.BadResponse{
			Message: "failed to register",
		})
		return
	}
	ctx.JSON(http.StatusBadRequest, dto.RegisterResponse{
		Message: "register success",
	})
}

func (u user) ActiveUser(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.BadResponse{
			Message: "id must be uuid",
		})
		return
	}

	userToken, ok := ctx.Get(IDENTITY_KEY)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, dto.BadResponse{
			Message: "you are not allowed to actived user",
		})
		return
	}
	user, err := u.services.User.GetByID(ctx, userToken.(*model.User).ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.BadResponse{
			Message: "you are not allowed to actived user",
		})
		return
	}
	if user.Role != model.RoleAdmin {
		ctx.JSON(http.StatusBadRequest, dto.BadResponse{
			Message: "you are not allowed to actived user",
		})
		return
	}
	if err := u.services.User.Active(ctx, id); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.BadResponse{
			Message: "failed to activate user",
		})
		return
	}
	ctx.JSON(http.StatusOK, dto.ActiveUserResponse{
		Message: "user activated",
	})
}
