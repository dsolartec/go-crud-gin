package services

import (
	"go-crud-gin/internal/apperror"
	"go-crud-gin/internal/models"
	"go-crud-gin/internal/platform/logger"
	"strings"
)

type UsersService interface {
	Create(username, password string) (int, error)
	GetByID(id int) *models.User
	GetByUsername(username string) *models.User
	GetUsers() []models.User
	DeleteUser(username string) error
}

type usersService struct {
	BaseService
	users []models.User
}

func (service *usersService) Create(username, password string) (int, error) {
	lastID := 0
	for _, user := range service.users {
		if strings.EqualFold(user.Username, username) {
			return 0, apperror.NewErrUserAlreadyExists()
		}

		if user.ID > lastID {
			lastID = user.ID
		}
	}

	lastID++

	service.users = append(service.users, models.User{
		ID:       lastID,
		Username: username,
		Password: password,
	})

	service.logger.Infof("[UsersService] New user created %s!", username)

	return lastID, nil
}

func (service *usersService) GetByID(id int) *models.User {
	for _, user := range service.users {
		if user.ID == id {
			return &user
		}
	}

	return nil
}

func (service *usersService) GetByUsername(username string) *models.User {
	for _, user := range service.users {
		if strings.EqualFold(user.Username, username) {
			return &user
		}
	}

	return nil
}

func (service *usersService) GetUsers() []models.User {
	return service.users
}

func (service *usersService) DeleteUser(username string) error {
	newUsers := []models.User{}
	for _, user := range service.users {
		if strings.EqualFold(user.Username, username) {
			continue
		}

		newUsers = append(newUsers, user)
	}

	if len(service.users) == len(newUsers) {
		return apperror.NewErrUserNotFound()
	}

	service.users = newUsers

	return nil
}

func NewUsersService(
	logger logger.Logger,
) UsersService {
	return &usersService{
		BaseService: BaseService{
			logger: logger,
		},
		users: []models.User{
			{
				ID:       1,
				Username: "admin",
				Password: "admin",
			},
			{
				ID:       2,
				Username: "dsolarte",
				Password: "1234",
			},
		},
	}
}
