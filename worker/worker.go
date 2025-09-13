package worker

import (
	"fmt"
	"os"
	"sfit-platform-web-backend/services"
	"strconv"
	"time"
)

type Worker struct {
	eventService *services.EventService
}

func NewWorker(eventService *services.EventService) *Worker {
	return &Worker{eventService: eventService}
}

func (w *Worker) Start() {
	envSecond, err := strconv.ParseInt(os.Getenv("TICK_TIME_GAP"), 10, 64)
	if err != nil {
		fmt.Println("TICK_TIME_GAP is not set or invalid")
		os.Exit(1)
	}
	tickTimeGap := time.Duration(envSecond) * time.Second
	ticker := time.NewTicker(tickTimeGap)
	defer ticker.Stop()

	for range ticker.C {
		err := w.eventService.AutoUpdateStatusEvent()
		if err != nil {
			_ = fmt.Errorf("auto update status event error: %v", err)
		}
	}
}
