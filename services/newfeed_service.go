package services

type NewFeedService struct {
	courseSer *CourseService
	userSer   *UserService
	taskSer   *TaskService
	evenetSer *EventService
}

func NewNewFeedService(courseSer *CourseService, userSer *UserService, taskSer *TaskService, evenetSer *EventService) *NewFeedService {
	return &NewFeedService{
		courseSer: courseSer,
		userSer:   userSer,
		taskSer:   taskSer,
		evenetSer: evenetSer,
	}
}
