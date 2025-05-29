package repository

import (
	"context"
	"skillsrock-test-task/internal/database/postgres"
	"skillsrock-test-task/internal/models"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

type TaskRepository struct {
	db sq.StatementBuilderType
	pg *postgres.Database
}

func NewTaskRepository(pg *postgres.Database) *TaskRepository {
	return &TaskRepository{
		db: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		pg: pg,
	}
}

func (r *TaskRepository) CreateTask(ctx context.Context, task *models.Task) (uint64, error) {
	query := r.db.
		Insert("tasks").
		Columns("title", "description", "status", "created_at").
		Values(task.Title, task.Description, task.Status, task.CreatedAt).
		Suffix("RETURNING id")

	sql, args, err := query.ToSql()
	if err != nil {
		return 0, err
	}

	var id uint64
	err = r.pg.Pool.QueryRow(ctx, sql, args...).Scan(&id)
	return id, err
}

func (r *TaskRepository) GetTaskByID(ctx context.Context, id uint64) (*models.Task, error) {
	query := r.db.
		Select("id", "title", "description", "status", "created_at", "updated_at").
		From("tasks").
		Where(sq.Eq{"id": id}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var task models.Task
	err = r.pg.Pool.QueryRow(ctx, sql, args...).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, models.ErrNotFound
	}
	return &task, err
}

func (r *TaskRepository) DeleteTask(ctx context.Context, id uint64) error {
	query := r.db.
		Delete("tasks").
		Where(sq.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	cmdTag, err := r.pg.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return models.ErrNotFound
	}
	return nil
}

func (r *TaskRepository) GetTasks(ctx context.Context, limit, offset uint64) ([]*models.Task, error) {
	query := r.db.
		Select("id", "title", "description", "status", "created_at", "updated_at").
		From("tasks").
		Limit(limit).
		Offset(offset).
		OrderBy("id DESC")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.pg.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.CreatedAt,
			&task.UpdatedAt,
		); err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}

	if len(tasks) == 0 {
		return nil, models.ErrNotFound
	}

	return tasks, nil
}

func (r *TaskRepository) UpdateTask(ctx context.Context, id uint64, task *models.Task) error {
	query := r.db.
		Update("tasks").
		SetMap(map[string]interface{}{
			"title":       task.Title,
			"description": task.Description,
			"status":      task.Status,
			"updated_at":  task.UpdatedAt,
		}).
		Where(sq.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	cmdTag, err := r.pg.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return models.ErrNotFound
	}
	return nil
}
