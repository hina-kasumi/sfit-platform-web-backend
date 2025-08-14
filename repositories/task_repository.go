package repositories

import (
	"sfit-platform-web-backend/entities"
	"time"

	"github.com/google/uuid"
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

func (tr *TaskRepository) GetTasks(page, pageSize int, name, eventID string) ([]*entities.Task, int64, error) {
	var tasks []*entities.Task
	query := tr.db.Model(&entities.Task{})

	if name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}
	if eventID != "" {
		query = query.Where("event_id = ?", eventID)
	}

	var totalCount int64
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&tasks).Error; err != nil {
		return nil, 0, err
	}
	return tasks, totalCount, nil
}

func (tr *TaskRepository) CreateTask(createdBy, name, description, eventID string, startDate, deadline time.Time) (*entities.Task, error) {
	createdByID, err := uuid.Parse(createdBy)
	if err != nil {
		return nil, err
	}
	task := &entities.Task{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		StartDate:   startDate,
		Deadline:    deadline,
		EventID:     nil,
		CreateBy:    createdByID,
	}
	if eventID != "" {
		parseEventID, err := uuid.Parse(eventID)
		if err != nil {
			return nil, err
		}
		task.EventID = &parseEventID
	}
	if err := tr.db.Create(task).Error; err != nil {
		return nil, err
	}
	return task, nil
}

func (tr *TaskRepository) GetTaskByID(id string) (*entities.Task, error) {
	var task entities.Task
	if err := tr.db.First(&task, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (tr *TaskRepository) UpdateTask(id, name, description string, percent_complete float32, startTime, deadline time.Time) error {
	return tr.db.Model(&entities.Task{}).Where("id = ?", id).Updates(entities.Task{
		Name:            name,
		Description:     description,
		PercentComplete: percent_complete,
		StartDate:       startTime,
		Deadline:        deadline,
	}).Error
}

func (tr *TaskRepository) DeleteTask(id string) error {
	return tr.db.Delete(&entities.Task{}, "id = ?", id).Error
}

func (tr *TaskRepository) AddUserTask(taskID, userID uuid.UUID) (*entities.TaskAssignments, error) {
	taskAssignment := &entities.TaskAssignments{
		TaskID: taskID,
		UserID: userID,
	}
	if err := tr.db.Create(taskAssignment).Error; err != nil {
		return nil, err
	}
	return taskAssignment, nil
}

func (tr *TaskRepository) ListTasksByUserID(userID uuid.UUID, page, pageSize int) ([]*entities.Task, int64, error) {
	var tasks []*entities.Task
	var totalCount int64
	query := tr.db.Model(&entities.Task{}).Joins("JOIN task_assignments ON task_assignments.task_id = tasks.id").
		Where("task_assignments.user_id = ?", userID)

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&tasks).Error; err != nil {
		return nil, 0, err
	}
	return tasks, totalCount, nil
}

func (tr *TaskRepository) DeleteUserTask(taskID, userID uuid.UUID) error {
	taskAssign := &entities.TaskAssignments{
		TaskID: taskID,
		UserID: userID,
	}
	return tr.db.Delete(taskAssign).Error
}
