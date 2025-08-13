package repositories

import (
	"sfit-platform-web-backend/dtos"

	"gorm.io/gorm"
)

type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{
		db: db,
	}
}

func (r *TaskRepository) GetListTasks(pageNum, pageSize int) ([]dtos.TaskGeneralInformationResponse, int, error) {
	var tasks []dtos.TaskGeneralInformationResponse
	var totalTasks int

	taskQuery := `SELECT id, event_id, name, description, start_date, dead_line, assignee, percent_complete, create_at, update_at FROM tasks LIMIT ? OFFSET ?`
	err := r.db.Raw(taskQuery, pageSize, (pageNum-1)*pageSize).Scan(&tasks).Error
	if err != nil {
		return nil, 0, err
	}

	countQuery := `SELECT COUNT(*) FROM tasks`
	err = r.db.Raw(countQuery).Scan(&totalTasks).Error
	if err != nil {
		return nil, 0, err
	}

	return tasks, totalTasks, nil
}