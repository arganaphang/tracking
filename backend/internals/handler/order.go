package handler

import (
	"application/internals/dto"
	"application/internals/model"
	"application/internals/service"
	"application/pkg/sse"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type IOrderHandler interface {
	FindAll(ctx *gin.Context)
	FindByID(ctx *gin.Context)
	WatchByID(ctx *gin.Context)
	Create(ctx *gin.Context)
	UpdateStatusByID(ctx *gin.Context)
}

type order struct {
	services  service.Services
	sentEvent *sse.SSEConn
}

func NewOrder(app *gin.Engine, services service.Services, sentEvent *sse.SSEConn, middlewares ...gin.HandlerFunc) IOrderHandler {
	o := order{services: services, sentEvent: sentEvent}
	g := app.Group("/order")
	g.Use(middlewares...)
	g.GET("/", o.FindAll)
	g.GET("/:id", o.FindByID)
	g.GET("/:id/watch", o.WatchByID)
	g.POST("/", o.Create)
	g.PUT("/:id/:status", o.UpdateStatusByID)

	return &o
}

func (o order) FindAll(ctx *gin.Context) {
	userToken, ok := ctx.Get(IDENTITY_KEY)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, dto.BadResponse{
			Message: "unauthorize",
		})
		return
	}
	var queries dto.FindAllRequest
	if err := ctx.ShouldBindQuery(&queries); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.BadResponse{
			Message: "failed to parse query params",
		})
		return
	}
	if queries.Page == nil {
		defaultPage := uint(1)
		queries.Page = &defaultPage
	}
	if queries.PerPage == nil {
		defaultPerPage := uint(10)
		queries.PerPage = &defaultPerPage
	}
	user, err := o.services.User.GetByID(ctx, userToken.(*model.User).ID)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dto.BadResponse{
			Message: "unauthorize",
		})
		return
	}
	var results []model.Order
	switch user.Role {
	case model.RoleApplicant:
		queries.CreatedBy = &user.ID
	}
	results, total, err := o.services.Order.FindAll(ctx, queries)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.BadResponse{
			Message: "failed to get all orders",
		})
		return
	}
	ctx.JSON(http.StatusOK, dto.FindAllResponse{
		Message: "get all orders",
		Data:    results,
		Meta: dto.Meta{
			Page:    *queries.Page,
			PerPage: *queries.PerPage,
			Total:   total,
		},
	})
}

func (o order) FindByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.BadResponse{
			Message: "id must be uuid",
		})
		return
	}
	result, err := o.services.Order.FindByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.BadResponse{
			Message: "failed to get order by id",
		})
		return
	}
	ctx.JSON(http.StatusOK, dto.FindByIDResponse{
		Message: "get order by id",
		Data:    result,
	})
}

func (o order) WatchByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.BadResponse{
			Message: "id must be uuid",
		})
		return
	}
	if _, err := o.services.Order.FindByID(ctx, id); err != nil {
		ctx.JSON(http.StatusNotFound, dto.BadResponse{
			Message: "order not found",
		})
		return
	}

	ctx.Writer.Header().Set("Content-Type", "text/event-stream")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")
	ctx.Writer.Header().Set("Transfer-Encoding", "chunked")
	ch := o.sentEvent.AddClient(id.String())
	go func() {
		o.sentEvent.Broadcast(id.String(), "initial")
	}()

	ctx.Stream(func(w io.Writer) bool {
		for {
			select {
			case message := <-*ch:
				result, err := o.services.Order.FindByID(ctx, id)
				if err != nil {
					return true
				}
				ctx.SSEvent("message", dto.FindByIDResponse{
					Message: fmt.Sprintf("notification: [%s]", message),
					Data:    result,
				})
				return true
			case <-ctx.Writer.CloseNotify(): // Act as defer function for http handler :man_shrugging: -> ref [https://github.com/gin-gonic/gin/issues/515#issuecomment-176434018]
				o.sentEvent.RemoveClient(id.String(), *ch)
				return false
			}
		}
	})
}

func (o order) Create(ctx *gin.Context) {
	var data dto.CreateOrderRequest
	if err := ctx.ShouldBind(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.BadResponse{
			Message: "failed to serialize request body",
		})
		return
	}
	user, ok := ctx.Get(IDENTITY_KEY)
	if ok {
		data.CreatedBy = user.(*model.User).ID
	}
	if err := o.services.Order.Create(ctx, data); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.BadResponse{
			Message: "failed to create order",
		})
		return
	}
	ctx.JSON(http.StatusCreated, dto.CreateOrderResponse{
		Message: "create order success",
	})
}

func (o order) UpdateStatusByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.BadResponse{
			Message: "id must be uuid",
		})
		return
	}
	status := strings.ToUpper(ctx.Param("status"))
	if !model.OrderStatusContains(status) {
		ctx.JSON(http.StatusBadRequest, dto.BadResponse{
			Message: "status not found",
		})
		return
	}
	var data dto.UpdateStatusOrderRequest
	if err := ctx.ShouldBind(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.BadResponse{
			Message: "failed to serialize request body",
		})
		return
	}
	data.ID = id
	data.Status = status
	user, ok := ctx.Get(IDENTITY_KEY)
	if ok {
		data.UpdatedBy = user.(*model.User).ID
	}
	if err := o.services.Order.UpdateStatusByID(ctx, data); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.BadResponse{
			Message: "failed to update order",
		})
		return
	}

	o.sentEvent.Broadcast(id.String(), "update status")
	ctx.JSON(http.StatusOK, dto.CreateOrderResponse{
		Message: "update order success",
	})
}
