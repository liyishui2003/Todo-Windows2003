package service

import (
	"errors"
	"time"
	"todolist/internal/model"
	"todolist/internal/repository"
)

type TaskService struct {
	taskRepo *repository.TaskRepository
}

func NewTaskService(taskRepo *repository.TaskRepository) *TaskService {
	return &TaskService{taskRepo: taskRepo}
}

// 增
func (s *TaskService) CreateTask(task *model.Task) error {
	return s.taskRepo.Create(task)
}

// 删
func (s *TaskService) DeleteTask(taskID uint) error {
	// 调用 Repository 删除任务
	err := s.taskRepo.DeleteTask(taskID)
	if err != nil {
		return err
	}
	return nil
}

// 改
func (s *TaskService) EditTask(ID uint, req *model.EditTaskRequest) (*model.Task, error) {
	//解析时间
	dueDate, err := time.Parse(time.RFC3339, req.DueDate)
	if err != nil {
		return nil, errors.New("invalid due data format")
	}

	dueDateStr := dueDate.Format(time.RFC3339)
	updateTask := &model.Task{
		ID:          ID,
		Title:       req.Title,
		Description: req.Description,
		DueDate:     dueDateStr,
		Status:      req.Status,
		UserID:      req.UserID,
		Importance:  req.Importance,
		Urgency:     req.Urgency,
	}

	err = s.taskRepo.UpdateTask(updateTask)
	if err != nil {
		return nil, err
	}

	return updateTask, nil
}

// 查
func (s *TaskService) GetTask(ID uint) (*model.Task, error) {
	task, err := s.taskRepo.GetTaskByID(ID)
	if err != nil {
		return nil, err
	}
	return task, nil
}

// 根据重要程度和优先级查
func (s *TaskService) GetTasksByPriority(importance, urgency int) ([]*model.Task, error) {
	// 调用 Repository 层获取任务
	tasks, err := s.taskRepo.GetTasksByPriority(importance, urgency)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

// 根据用户id查所有task
func (s *TaskService) GetTasksByUserID(userID uint) ([]*model.Task, error) {
	// 调用 Repository 层的方法
	tasks, err := s.taskRepo.GetTasksByUserID(userID)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}
