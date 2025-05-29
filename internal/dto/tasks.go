package dto

import (
	"skillsrock-test-task/internal/models"
	"time"
)

type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type CreateTaskResponse struct {
	ID        uint64    `json:"id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type UpdateTaskRequest struct {
	ID          uint64 `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type GetTaskByIDResponse struct {
	Task models.Task `json:"task"`
}

type GetTasksResponse struct {
	Tasks []models.Task `json:"tasks"`
}
