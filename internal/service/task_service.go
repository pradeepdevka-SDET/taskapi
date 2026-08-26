package service

import (
	"database/sql"
	"errors"
	"strings"
	"taskapi/internal/model"
	"taskapi/internal/repository"
)

// Domain errors - the vocabulary of this layer. The handler maps these to http codes
var (
	ErrInvalidTitle = errors.New("title is required")
	ErrTaskNotFound = errors.New("task not found")
)

type TaskService struct {
	repo *repository.TaskRepository
}

func NewTaskService(repo *repository.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(userId int, title string) (model.Task, error) {
	if strings.TrimSpace(title) == "" {
		return model.Task{}, ErrInvalidTitle
	}
	return s.repo.Create(userId, title)
}
func (s *TaskService) GetTasks(userId int) ([]model.Task, error) {
	return s.repo.GetAllByUser(userId)
}

func (s *TaskService) GetTask(userId, id int) (model.Task, error) {
	task, err := s.repo.GetByID(userId, id)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Task{}, ErrTaskNotFound
	}
	if err != nil {
		return model.Task{}, err
	}
	return task, nil
}

func (s *TaskService) UpdateTask(userId, id int, title string, done bool) (model.Task, error) {
	if strings.TrimSpace(title) == "" {
		return model.Task{}, ErrInvalidTitle
	}
	count, err := s.repo.Update(userId, id, title, done)
	if err != nil {
		return model.Task{}, err
	}
	if count == 0 {
		return model.Task{}, ErrTaskNotFound
	}
	return model.Task{ID: id, UserId: userId, Title: title, Done: done}, nil
}

func (s *TaskService) DeleteTask(userId, id int) error {
	count, err := s.repo.Delete(userId, id)
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrTaskNotFound
	}
	return nil
}
