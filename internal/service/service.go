package service

import (
	"context"
	"skillsrock-test-task/internal/dto"
	"skillsrock-test-task/internal/models"
	"strconv"
	"time"
)

const (
	statusNew      = "new"
	statusProgress = "in_progress"
	statusDone     = "done"
)

type TaskRepository interface {
	CreateTask(ctx context.Context, task *models.Task) (uint64, error)
	GetTaskByID(ctx context.Context, id uint64) (*models.Task, error)
	DeleteTask(ctx context.Context, id uint64) error
	GetTasks(ctx context.Context, limit, offset uint64) ([]*models.Task, error)
	UpdateTask(ctx context.Context, id uint64, task *models.Task) error
}

type TaskService struct {
	repo TaskRepository
}

func NewTaskService(repo TaskRepository) *TaskService {
	return &TaskService{
		repo: repo,
	}
}

func (s *TaskService) CreateTask(ctx context.Context, task *dto.CreateTaskRequest) (*dto.CreateTaskResponse, error) {
	now := time.Now()

	id, err := s.repo.CreateTask(ctx, &models.Task{
		Title:       task.Title,
		Description: task.Description,
		Status:      statusNew,
		CreatedAt:   now,
	})

	return &dto.CreateTaskResponse{
		ID:        id,
		Status:    statusNew,
		CreatedAt: now,
	}, err
}

func (s *TaskService) DeleteTask(ctx context.Context, taskIDStr string) error {
	taskID, err := strconv.ParseUint(taskIDStr, 10, 64)
	if err != nil {
		return models.ErrFailedToParseID
	}

	return s.repo.DeleteTask(ctx, taskID)
}

func (s *TaskService) GetTaskByID(ctx context.Context, taskIDStr string) (*dto.GetTaskByIDResponse, error) {
	taskID, err := strconv.ParseUint(taskIDStr, 10, 64)
	if err != nil {
		return nil, models.ErrFailedToParseID
	}

	task, err := s.repo.GetTaskByID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	return &dto.GetTaskByIDResponse{
		Task: task,
	}, nil
}

func (s *TaskService) GetTasks(ctx context.Context, pageStr string, limitStr string) (*dto.GetTasksResponse, error) {
	page, err := strconv.ParseUint(pageStr, 10, 64)
	if err != nil && pageStr != "" {
		return nil, models.ErrFailedToParsePage
	}

	limit, err := strconv.ParseUint(limitStr, 10, 64)
	if err != nil && limitStr != "" {
		return nil, models.ErrFailedToParseLimit
	}

	if page < 1 {
		page = 1
	}
	if limit == 0 {
		limit = 10
	}

	offset := (page - 1) * limit

	tasks, err := s.repo.GetTasks(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	return &dto.GetTasksResponse{
		Tasks: tasks,
	}, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, taskIDStr string, task *dto.UpdateTaskRequest) error {
	taskID, err := strconv.ParseUint(taskIDStr, 10, 64)
	if err != nil {
		return models.ErrFailedToParseID
	}

	if task.Status != statusNew && task.Status != statusProgress && task.Status != statusDone {
		return models.ErrInvalidStatus
	}

	return s.repo.UpdateTask(ctx, taskID, &models.Task{
		ID:          taskID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		UpdatedAt:   time.Now(),
	})
}
