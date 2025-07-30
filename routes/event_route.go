package routes

import (
	"github.com/gin-gonic/gin"
	"sfit-platform-web-backend/handlers"
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

	task := router.Group("/event")
	task.GET("/list", eventHandler.GetEventList)
	task.GET("/:event_id", eventHandler.GetEventDetail)
	task.GET("/registed-events-list", eventHandler.GetRegistedEventList)
	task.POST("/subscribe", eventHandler.SubscribeEvent)
	task.POST("/unsubscribe", eventHandler.UnsubscribeEvent)
	task.DELETE("/delete/:event_id", eventHandler.DeleteEvent)
	task.PUT("/update", eventHandler.UpdateEvent)
}
