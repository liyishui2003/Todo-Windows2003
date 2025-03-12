package repository

import (
	"todolist/internal/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	if err := r.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Create(user *model.User) error {
	return r.DB.Create(user).Error
}

func (r *UserRepository) UpdateAvatar(userID uint, avatarPath string) error {
	var user model.User
	result := r.DB.Model(&user).Where("id = ?", userID).Update("avatar", avatarPath)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
