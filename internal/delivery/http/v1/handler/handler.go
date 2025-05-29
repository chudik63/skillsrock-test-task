package handler

import (
	"context"
	"errors"
	"skillsrock-test-task/internal/dto"
	"skillsrock-test-task/internal/models"
	"skillsrock-test-task/pkg/logger"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const requestTimeout = 1 * time.Second

type ErrorResponse struct {
	Error string `json:"error"`
}

type TaskService interface {
	CreateTask(ctx context.Context, task *dto.CreateTaskRequest) (*dto.CreateTaskResponse, error)
	GetTaskByID(ctx context.Context, id string) (*dto.GetTaskByIDResponse, error)
	GetTasks(ctx context.Context, page, limit string) (*dto.GetTasksResponse, error)
	DeleteTask(ctx context.Context, id string) error
	UpdateTask(ctx context.Context, id string, task *dto.UpdateTaskRequest) error
}

type Handler struct {
	service TaskService
	logger  logger.Logger
}

func NewHandler(serv TaskService, log logger.Logger) *Handler {
	return &Handler{
		service: serv,
		logger:  log,
	}
}

// CreateTask
// @Summary      Create a new task
// @Description  Creates a new task
// @Tags         tasks
// @Produce      json
// @Param book body dto.CreateTaskRequest true "Task"
// @Success      201  {object}  dto.CreateTaskResponse
// @Failure      400  {object}  ErrorResponse  "Invalid request body"
// @Failure      500  {object}  ErrorResponse  "Uknown error occured while creating the task"
// @Router /tasks [post]
func (h *Handler) CreateTask(ctx *fiber.Ctx) error {
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	var task dto.CreateTaskRequest

	if err := ctx.BodyParser(&task); err != nil {
		h.logger.Error(ctx.Context(), "Failed to parse request body", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "Invalid request body"})
	}

	res, err := h.service.CreateTask(ctxWithTimeout, &task)
	if err != nil {

		h.logger.Error(ctx.Context(), "Unknown error occured while creating the task", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "Uknown error occured while creating the task"})
	}

	return ctx.Status(fiber.StatusCreated).JSON(res)
}

// GetTask
// @Summary      Get a task by ID
// @Description  Retrieves a task by its ID
// @Tags         tasks
// @Produce      json
// @Param        id   path      string  true  "Task ID"
// @Success      200  {object}  dto.GetTaskByIDResponse  "Task details"
// @Failure      400  {object}  ErrorResponse  "Invalid task ID"
// @Failure      404  {object}  ErrorResponse  "Task not found"
// @Failure      500  {object}  ErrorResponse  "Unknown error occurred"
// @Router       /tasks/{id} [get]
func (h *Handler) GetTaskByID(c *fiber.Ctx) error {
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	taskID := c.Params("id")

	res, err := h.service.GetTaskByID(ctxWithTimeout, taskID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrFailedToParseID):
			return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: err.Error()})
		case errors.Is(err, models.ErrNotFound):
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{Error: err.Error()})
		default:
			h.logger.Error(c.Context(), "Unknown error occurred while getting the task", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: "Unknown error occurred while getting the task"})
		}
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// UpdateTask
// @Summary      Update a task by ID
// @Description  Updates a task's data
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id    path  string     true  "Task ID"
// @Param        task  body  dto.UpdateTaskRequest   true  "Task payload"
// @Success      200   {string}  string  "Updated successfully"
// @Failure      400   {object}  ErrorResponse  "Invalid input or task ID"
// @Failure      404   {object}  ErrorResponse  "Task not found"
// @Failure      500   {object}  ErrorResponse  "Unknown error occurred"
// @Router       /tasks/{id} [put]
func (h *Handler) UpdateTask(c *fiber.Ctx) error {
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	taskID := c.Params("id")

	var task dto.UpdateTaskRequest
	if err := c.BodyParser(&task); err != nil {
		h.logger.Error(c.Context(), "Invalid request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: "Invalid request body"})
	}

	if err := h.service.UpdateTask(ctxWithTimeout, taskID, &task); err != nil {
		switch {
		case errors.Is(err, models.ErrFailedToParseID):
			return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: err.Error()})
		case errors.Is(err, models.ErrNotFound):
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{Error: err.Error()})
		default:
			h.logger.Error(c.Context(), "Unknown error occurred while updating the task", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: "Unknown error occurred while updating the task"})
		}
	}

	return c.SendStatus(fiber.StatusOK)
}

// DeleteTask
// @Summary      Delete a task by ID
// @Description  Deletes a task by its ID
// @Tags         tasks
// @Produce      json
// @Param        id   path  string  true  "Task ID"
// @Success      200  {string}  string  "Deleted successfully"
// @Failure      400  {object}  ErrorResponse  "Invalid task ID"
// @Failure      404  {object}  ErrorResponse  "Task not found"
// @Failure      500  {object}  ErrorResponse  "Unknown error occurred"
// @Router       /tasks/{id} [delete]
func (h *Handler) DeleteTask(c *fiber.Ctx) error {
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	taskID := c.Params("id")

	if err := h.service.DeleteTask(ctxWithTimeout, taskID); err != nil {
		switch {
		case errors.Is(err, models.ErrFailedToParseID):
			return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: err.Error()})
		case errors.Is(err, models.ErrNotFound):
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{Error: err.Error()})
		default:
			h.logger.Error(c.Context(), "Unknown error occurred while deleting the task", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: "Unknown error occurred while deleting the task"})
		}
	}

	return c.SendStatus(fiber.StatusOK)
}

// GetTasks
// @Summary      Get tasks
// @Description  Retrieves a paginated list of tasks
// @Tags         tasks
// @Produce      json
// @Param        page   query     string  false  "Page number"
// @Param        limit  query     string  false  "Items per page"
// @Success      200    {object}  dto.GetTasksResponse  "List of tasks"
// @Failure      400    {object}  ErrorResponse  "Invalid pagination parameters"
// @Failure      404    {object}  ErrorResponse  "No tasks found"
// @Failure      500    {object}  ErrorResponse  "Unknown error occurred"
// @Router       /tasks [get]
func (h *Handler) GetTasks(c *fiber.Ctx) error {
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	page := c.Query("page")
	limit := c.Query("limit")

	res, err := h.service.GetTasks(ctxWithTimeout, page, limit)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrFailedToParsePage), errors.Is(err, models.ErrFailedToParseLimit):
			return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: err.Error()})
		case errors.Is(err, models.ErrNotFound):
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{Error: err.Error()})
		default:
			h.logger.Error(c.Context(), "Unknown error occurred while listing the tasks", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: "Unknown error occurred while listing the tasks"})
		}
	}

	return c.Status(fiber.StatusOK).JSON(res)
}
