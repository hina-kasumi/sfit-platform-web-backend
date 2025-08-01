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
	task.GET("/registed-event-list", eventHandler.GetRegistedEventList)
	task.GET("/:event_id", eventHandler.GetEventDetail)
	task.POST("/user-attendance", eventHandler.EventAttendance)
	task.POST("/subscribe", eventHandler.SubscribeEvent)
	task.POST("", eventHandler.CreateEvent)
	task.PUT("", eventHandler.UpdateEvent)
	task.DELETE("/delete/:event_id", eventHandler.DeleteEvent)
	task.POST("/unsubscribe", eventHandler.UnsubscribeEvent)
	task.GET("/list-register", eventHandler.GetEventRegistedList)
	task.GET("/list-attendance", eventHandler.GetEventAttendanceList)
}
