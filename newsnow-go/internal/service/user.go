package service

import (
	"encoding/json"
	"newsnow-go/internal/repository"
	"newsnow-go/internal/types"
	"newsnow-go/internal/utils"
	"time"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) AddUser(id string, email string, userType string) error {
	if s.repo == nil {
		return nil
	}
	return s.repo.AddUser(id, email, userType)
}

func (s *UserService) GetUser(id string) (*types.UserInfo, error) {
	if s.repo == nil {
		return nil, nil
	}
	return s.repo.GetUser(id)
}

func (s *UserService) GetUserData(id string) (*struct {
	Data    json.RawMessage
	Updated int64
}, error) {
	if s.repo == nil {
		return nil, nil
	}
	result, err := s.repo.GetData(id)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return &struct {
		Data    json.RawMessage
		Updated int64
	}{
		Data:    json.RawMessage(result.Data),
		Updated: result.Updated,
	}, nil
}

func (s *UserService) SetUserData(id string, data map[string]interface{}, updatedTime int64) error {
	if s.repo == nil {
		return nil
	}
	
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	
	return s.repo.SetData(id, string(jsonData), updatedTime)
}

func (s *UserService) ParseMetadata(data json.RawMessage) (map[string][]types.SourceID, error) {
	var metadata map[string][]types.SourceID
	if err := json.Unmarshal(data, &metadata); err != nil {
		utils.Warn("Failed to parse metadata: " + err.Error())
		return nil, err
	}
	return metadata, nil
}

func (s *UserService) UpdateMetadata(id string, metadata map[string]interface{}) error {
	if s.repo == nil {
		return nil
	}
	
	jsonData, err := json.Marshal(metadata)
	if err != nil {
		return err
	}
	
	return s.repo.SetData(id, string(jsonData), time.Now().UnixMilli())
}
