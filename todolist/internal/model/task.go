package model

type Task struct {
	ID          uint    `json:"id" gorm:"primaryKey"`
	Title       string  `json:"title" gorm:"not null"`
	Description string  `json:"description"`
	DueDate     string  `json:"due_date"`
	Status      string  `json:"status" gorm:"default:pending"`
	UserID      float64 `json:"user_id" gorm:"not null"` // 关联用户
	Importance  int     `json:"importance"`              // 重要程度（0-1）
	Urgency     int     `json:"urgency"`                 // 紧急程度（0-1）
}

// 这里没有传入status，因为默认pending，直到用户标记为完成
type CreateTaskRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	DueDate     string `json:"due_date" binding:"required"`
	Importance  int    `json:"importance"` // 重要程度（0-1）
	Urgency     int    `json:"urgency"`    // 紧急程度（0-1）
}

type EditTaskRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description" binding:"required"`
	DueDate     string  `json:"due_date" binding:"required"`
	Status      string  `json:"status" binding:"required"`
	UserID      float64 `json:"user_id" binding:"required"`
	Importance  int     `json:"importance"` // 重要程度（0-1）
	Urgency     int     `json:"urgency"`    // 紧急程度（0-1）
}
