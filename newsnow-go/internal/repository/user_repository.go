package repository

import (
	"newsnow-go/internal/model"
	"newsnow-go/internal/types"
	"time"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Init() error {
	return r.db.AutoMigrate(&model.User{})
}

func (r *UserRepository) AddUser(id string, email string, userType string) error {
	now := time.Now().UnixMilli()
	
	var existing model.User
	err := r.db.Where("id = ?", id).First(&existing).Error
	
	if err == gorm.ErrRecordNotFound {
		user := model.User{
			ID:      id,
			Email:   email,
			Data:    "",
			Type:    userType,
			Created: now,
			Updated: now,
		}
		return r.db.Create(&user).Error
	}
	
	if err != nil {
		return err
	}
	
	if existing.Email != email || existing.Type != userType {
		return r.db.Model(&existing).Updates(map[string]interface{}{
			"email":   email,
			"updated": now,
		}).Error
	}
	
	return nil
}

func (r *UserRepository) GetUser(id string) (*types.UserInfo, error) {
	var user model.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	return &types.UserInfo{
		ID:      user.ID,
		Email:   user.Email,
		Type:    user.Type,
		Data:    user.Data,
		Created: user.Created,
		Updated: user.Updated,
	}, nil
}

func (r *UserRepository) SetData(id string, data string, updatedTime int64) error {
	if updatedTime == 0 {
		updatedTime = time.Now().UnixMilli()
	}
	
	result := r.db.Model(&model.User{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"data":     data,
			"updated": updatedTime,
		})
	
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	
	return nil
}

func (r *UserRepository) GetData(id string) (*struct {
	Data     string
	Updated  int64
}, error) {
	var user model.User
	err := r.db.Select("data", "updated").Where("id = ?", id).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	return &struct {
		Data    string
		Updated int64
	}{
		Data:    user.Data,
		Updated: user.Updated,
	}, nil
}

func (r *UserRepository) DeleteUser(id string) error {
	result := r.db.Where("id = ?", id).Delete(&model.User{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *UserRepository) GetDB() *gorm.DB {
	return r.db
}
