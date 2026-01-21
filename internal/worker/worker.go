package worker

import (
	"log"
	"sfit-platform-web-backend/internal/services"
	"time"
)

type Worker struct {
	eventService *services.EventService
}

func NewWorker(eventService *services.EventService) *Worker {
	return &Worker{eventService: eventService}
}

func (w *Worker) Start() {
	envSecond := 900 // default 15 minutes
	tickTimeGap := time.Duration(envSecond) * time.Second
	ticker := time.NewTicker(tickTimeGap)
	defer ticker.Stop()

	for range ticker.C {
		err := w.eventService.AutoUpdateStatusEvent()
		if err != nil {
			log.Printf("auto update status event error: %v", err)
		}
	}
}
