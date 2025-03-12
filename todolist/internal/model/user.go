package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"` // 存储哈希后的密码
	Avatar   string // 头像路径
}

// 后续可添加其他字段（如头像、手机号等）

//管理登录请求的结构体
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

//相应登录请求的结构体
type LoginResponse struct {
	Token string `json:"token"`
}

//负责注册的结构体
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
