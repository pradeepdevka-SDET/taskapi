package repository

import (
	"database/sql"
	"taskapi/internal/model"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

// Create inserts a task owned by userId and returns the full row
func (r *TaskRepository) Create(userId int, title string) (model.Task, error) {
	var t model.Task
	err := r.db.QueryRow(
		`INSERT INTO tasks (user_id,title) VALUES ($1,$2)
		RETURNING id, user_id,title,done,created_at`,
		userId, title,
	).Scan(&t.ID, &t.UserId, &t.Title, &t.Done, &t.CreatedAt)
	return t, err
}

// GetAllByUser returns only the tasks belonging to userId
func (r *TaskRepository) GetAllByUser(userId int) ([]model.Task, error) {
	rows, err := r.db.Query(
		`SELECT id, user_id,title,done,created_at 
		FROM tasks WHERE user_id=$1 ORDER BY id`,
		userId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tasks := []model.Task{}
	for rows.Next() {
		var t model.Task
		if err := rows.Scan(&t.ID, &t.UserId, &t.Title, &t.Done, &t.CreatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	if err := rows.Err(); err != nil {
		return tasks, err
	}
	return tasks, nil
}

// GetById returns one task only if it belongs to userId
func (r *TaskRepository) GetByID(userId, id int) (model.Task, error) {
	var t model.Task
	err := r.db.QueryRow(
		`SELECT id, user_id,title,done,created_at 
		FROM tasks WHERE user_id=$1 and id=$2`,
		userId, id,
	).Scan(&t.ID, &t.UserId, &t.Title, &t.Done, &t.CreatedAt)
	return t, err
}

// Update changes a task ONLY if it belongs to userID.
func (r *TaskRepository) Update(userId, id int, title string, done bool) (int64, error) {
	result, err := r.db.Exec(
		`UPDATE tasks SET title=$1, done =$2 WHERE id=$3 AND user_id=$3`,
		title, done, id, userId,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// DELTE removes a task only if it belongs to userId
func (r *TaskRepository) Delete(userId, id int) (int64, error) {
	result, err := r.db.Exec(
		`DELETE FROM tasks WHERE id=$1 AND user_id=$2`,
		id, userId,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
