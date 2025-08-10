package routes

import (
	"sfit-platform-web-backend/handlers"
	"sfit-platform-web-backend/middlewares"

	"github.com/gin-gonic/gin"
)

type EventRoutes struct {
	eventHandler *handlers.EventHandler
}

func NewEventRoute(eventHandler *handlers.EventHandler) *EventRoutes {
	return &EventRoutes{
		eventHandler: eventHandler,
	}
}

func (eventRou *EventRoutes) RegisterRoutes(router *gin.Engine) {
	eventHandler := eventRou.eventHandler

	task := router.Group("/events")
	task.GET("", eventHandler.GetEventList)
	task.GET("/:event_id", eventHandler.GetEventDetail)
	task.GET("/:event_id/users", eventHandler.GetUsersInEvent)

	taskAuth := router.Group("/events")
	taskAuth.Use(middlewares.EnforceAuthenticatedMiddleware())
	taskAuth.POST("", eventHandler.CreateEvent)
	taskAuth.PUT("", eventHandler.UpdateEvent)
	taskAuth.DELETE("/:event_id", eventHandler.DeleteEvent)

	attendEvent := router.Group("/users/:user_id/events/:event_id")
	attendEvent.Use(middlewares.EnforceAuthenticatedMiddleware())
	attendEvent.PUT("", eventHandler.UpdateStatusUserAttendance)
}
