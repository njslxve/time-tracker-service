package service

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/njslxve/time-tracker-service/internal/config"
	"github.com/njslxve/time-tracker-service/internal/model/dto"
	"github.com/njslxve/time-tracker-service/internal/model/entity"
)

type StrorageInterface interface {
	AddUser(entity.User) error
	GetUser(int) (entity.User, error)
	GetUsers(entity.FilterOptions) ([]entity.User, error)
	UpdateUser(entity.User) error
	DeleteUser(int) error
	AddTask(entity.Task) error
	GetTask(string, int) (entity.Task, error)
	GetTasks(int, int) ([]entity.Task, error)
	UpdateTask(entity.Task) error
	TokenData(string) (int, entity.TokenData, error)
	AddToken(string, int, []byte) error
}

type APIInterface interface {
	Info(string) (dto.UserInfoResponse, error)
}

type Service struct {
	cfg    *config.Config
	logger *slog.Logger
	db     StrorageInterface
	api    APIInterface
}

func New(cfg *config.Config, logger *slog.Logger, db StrorageInterface, api APIInterface) *Service {
	return &Service{
		cfg:    cfg,
		logger: logger,
		db:     db,
		api:    api,
	}
}

func (s *Service) AddUser(req dto.AddUserRequest) error {
	const op = "service.Service.AddUser"

	userInfo, err := s.api.Info(req.Passport)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	user := entity.User{
		Passport:   req.Passport,
		Name:       userInfo.Name,
		Surmame:    userInfo.Surname,
		Patronymic: userInfo.Patronymic,
		Adress:     userInfo.Adress,
	}

	err = s.db.AddUser(user)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) AddTask(req dto.TaskRequest) error {
	task := entity.Task{
		TaskID:    req.TaskID,
		UserID:    req.UserID,
		StartTime: time.Now(),
	}

	return s.db.AddTask(task)
}

func (s *Service) EndTask(req dto.TaskRequest) error {
	const op = "service.Service.EndTask"

	task, err := s.db.GetTask(req.TaskID, req.UserID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	task.EndTime = time.Now()

	task.Duration = int(task.EndTime.Sub(task.StartTime).Minutes())

	err = s.db.UpdateTask(task)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) UpdateUser(userID string, req dto.UpdateUserRequest) error {
	const op = "service.Service.UpdateUser"

	id, _ := strconv.Atoi(userID)

	user, err := s.db.GetUser(id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if req.Name != "" {
		user.Name = req.Name
	}

	if req.Surname != "" {
		user.Surmame = req.Surname
	}

	if req.Patronymic != "" {
		user.Patronymic = req.Patronymic
	}

	if req.Passport != "" {
		user.Passport = req.Passport
	}

	if req.Adress != "" {
		user.Adress = req.Adress
	}

	err = s.db.UpdateUser(user)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) DeleteUser(userID string) error {
	const op = "service.Service.DeleteUser"

	id, _ := strconv.Atoi(userID)

	err := s.db.DeleteUser(id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) GetTasks(userID string, interval string) ([]dto.TaskResponse, error) {
	const op = "service.Service.GetTasks"

	id, _ := strconv.Atoi(userID)
	intrval, _ := strconv.Atoi(interval)

	tasks, err := s.db.GetTasks(id, intrval)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	tasksRes := make([]dto.TaskResponse, 0)

	for _, task := range tasks {
		t := dto.TaskResponse{
			TaskID:    task.TaskID,
			UserID:    task.UserID,
			StartTime: task.StartTime,
			EndTime:   task.EndTime,
			Duration:  convertDuration(task.Duration),
		}

		tasksRes = append(tasksRes, t)
	}

	return tasksRes, nil
}

func (s *Service) GetUsers(filterOpts entity.FilterOptions, paginationOpts entity.PaginationOptions) (dto.GetUsersResponse, error) {
	const op = "service.Service.GetUsers"

	usersData, err := s.db.GetUsers(filterOpts)
	if err != nil {
		return dto.GetUsersResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	var oldLimit int

	if paginationOpts.Next != "" {
		old, tokenData, err := s.db.TokenData(paginationOpts.Next)
		if err != nil {
			return dto.GetUsersResponse{}, fmt.Errorf("%s: %w", op, err)
		}

		oldLimit = old

		if !s.validateToken(paginationOpts.Next, tokenData, filterOpts) {
			paginationOpts.Next = ""
		}
	}

	users, nextToken := s.paginate(oldLimit, usersData, paginationOpts, filterOpts)

	usersRes := make([]dto.User, 0)

	for _, user := range users {
		u := dto.User{
			UserID:     user.UserID,
			Surname:    user.Surmame,
			Name:       user.Name,
			Patronymic: user.Patronymic,
			Passport:   user.Passport,
			Adress:     user.Adress,
		}

		usersRes = append(usersRes, u)
	}

	return dto.GetUsersResponse{Users: usersRes, Next: nextToken}, nil
}

func convertDuration(duration int) string {
	return fmt.Sprintf("%dh%dm", duration/60, duration%60)
}

func (s *Service) validateToken(token string, tokenData entity.TokenData, filterOpts entity.FilterOptions) bool {
	if token == "" {
		return true
	}

	if !tokenData.IsAlive {
		return false
	}

	dataOpts := entity.FilterOptions{}

	if err := json.Unmarshal(tokenData.Params, &dataOpts); err != nil {
		return false
	}

	return dataOpts == filterOpts
}

func (s *Service) paginate(oldLimit int, users []entity.User, paginationOpts entity.PaginationOptions, filterOpts entity.FilterOptions) ([]entity.User, string) {
	if paginationOpts.Limit == 0 {
		paginationOpts.Limit = 10
	}

	switch paginationOpts.Next {
	case "":
		if len(users) <= paginationOpts.Limit {
			return users, ""
		} else {
			return users[:paginationOpts.Limit], s.newToken(paginationOpts.Limit, filterOpts)
		}
	default:
		if len(users[oldLimit:]) <= paginationOpts.Limit {
			return users[oldLimit:], ""
		} else {
			return users[oldLimit : oldLimit+paginationOpts.Limit], s.newToken(oldLimit+paginationOpts.Limit, filterOpts)
		}
	}
}

func (s *Service) newToken(limit int, filterOpts entity.FilterOptions) string {
	params, _ := json.Marshal(filterOpts)

	str := uuid.NewString()
	new := strings.ReplaceAll(str, "-", "")
	token := new[:15]

	err := s.db.AddToken(string(token), limit, params)
	if err != nil {
		return ""
	}

	return token
}
