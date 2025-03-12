package repository

import (
	"errors"
	"todolist/internal/model"

	"gorm.io/gorm"
)

type TaskRepository struct {
	DB *gorm.DB
}

// 构造函数依赖注入
func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{DB: db}
}

// 用gorm的框架creater一条task
func (r *TaskRepository) Create(task *model.Task) error {
	return r.DB.Create(task).Error
}

// 删
func (r *TaskRepository) DeleteTask(taskID uint) error {
	// 使用 GORM 删除任务
	result := r.DB.Delete(&model.Task{}, taskID)
	if result.Error != nil {
		return result.Error
	}
	// 检查是否实际删除了记录
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// 改
func (r *TaskRepository) UpdateTask(task *model.Task) error {
	var existingTask model.Task
	result := r.DB.First(&existingTask, task.ID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return gorm.ErrRecordNotFound
		}
		return result.Error
	}

	// 使用 GORM 的 Save 方法更新任务
	result = r.DB.Save(task)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// 用gorm框架，根据ID查询记录，如果用sql原生框架就要自己手动写sql语句和Scan映射
func (r *TaskRepository) GetTaskByID(ID uint) (*model.Task, error) {
	var task model.Task
	err := r.DB.First(&task, ID).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// 根据重要程度和优先级进行筛选
func (r *TaskRepository) GetTasksByPriority(importance, urgency int) ([]*model.Task, error) {
	var tasks []*model.Task
	query := r.DB.Model(&model.Task{})

	// 动态组合筛选条件
	if importance >= 0 {
		query = query.Where("importance = ?", importance)
	}
	if urgency >= 0 {
		query = query.Where("urgency = ?", urgency)
	}

	// 执行查询
	result := query.Find(&tasks)
	if result.Error != nil {
		return nil, result.Error
	}
	return tasks, nil
}

// 查询一个用户名下所有task
func (r *TaskRepository) GetTasksByUserID(userID uint) ([]*model.Task, error) {
	var tasks []*model.Task

	// 使用 GORM 查询 user_id 匹配的所有任务
	result := r.DB.Where("user_id = ?", userID).Find(&tasks)
	if result.Error != nil {
		return nil, result.Error
	}

	return tasks, nil
}


