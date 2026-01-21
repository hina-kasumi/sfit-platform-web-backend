package server

import (
	"log"
	"sfit-platform-web-backend/internal/config"
	"sfit-platform-web-backend/internal/infrastructure/database"
	"sfit-platform-web-backend/internal/infrastructure/redis"
	middlewares "sfit-platform-web-backend/internal/middleware"
	"sfit-platform-web-backend/internal/model"
	"sfit-platform-web-backend/internal/router"

	"github.com/gin-gonic/gin"
)

type Server struct {
	cfg *config.Config
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		cfg: cfg,
	}
}

func (s *Server) Run() error {
	db := database.OpenDbConnection(&s.cfg.Database)
	db = db.Debug()
	if db != nil {
		err := db.AutoMigrate(
			&model.Log{}, &model.Device{}, &model.Users{}, &model.UserProfile{},
			&model.Teams{}, &model.Event{}, &model.Course{}, &model.Module{},
			&model.Task{}, &model.TaskAssignments{}, &model.Lesson{}, &model.Tag{},
			&model.FavoriteCourse{}, &model.EventAttendance{}, &model.TagTemp{}, &model.TeamMembers{},
			&model.UserCourse{}, &model.LessonAttendance{}, &model.Newsfeed{}, &model.UserRate{},
			&model.Role{}, &model.UserRole{},
		)
		if err != nil {
			log.Fatalf("AutoMigrate failed: %v", err)
		}
	}
	defer database.CloseDbConnection(db)

	redisClient, redisCtx := redis.InitRedis(&s.cfg.Redis)
	defer redisClient.Close()

	r := gin.Default()
	r.Use(middlewares.Cors(&s.cfg.Cors))

	router.Register(r, s.cfg, db, redisClient, redisCtx)

	return r.Run(":" + s.cfg.App.Port)
}
