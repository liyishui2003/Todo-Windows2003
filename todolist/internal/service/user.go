package service

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
	"todolist/internal/config"
	"todolist/internal/model"
	"todolist/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepor *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepor,
	}
}

// 通过username找ID，存储在前端的localStorage里
func (s *UserService) FindUserID(username string) (uint, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return 0, errors.New("用户不存在")
	}
	return user.ID, err
}

// 登录逻辑
func (s *UserService) Login(username, password string) (string, error) {
	// 1. 查询用户
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return "", errors.New("用户不存在")
	}

	fmt.Printf("username=%s,password=%s\n", username, password)
	fmt.Printf("password in DB =%s\n", user.Password)
	// 2. 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("密码错误")
	}

	// 3. 签发 JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // 有效期 24 小时
	})

	return token.SignedString([]byte(config.Cfg.JWTSecret))
}

// 注册逻辑
func (s *UserService) Register(username, password string) error {
	// 1. 检查用户名是否已存在
	_, err := s.userRepo.FindByUsername(username)
	if err == nil {
		return errors.New("用户名已存在")
	}

	// 2. 对密码进行哈希处理
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("密码哈希失败")
	}

	// 3. 创建用户
	user := model.User{
		Username: username,
		Password: string(hashedPassword), // 存储哈希值
	}

	// 4. 保存用户到数据库
	if err := s.userRepo.Create(&user).Error; err != nil {
		return errors.New("用户注册失败")
	}

	return nil
}

// 获取用户头像
func (s *UserService) GetAvatar(username string) (string, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return "", err
	}
	return user.Avatar, nil
}


