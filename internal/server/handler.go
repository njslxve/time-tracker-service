package server

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/njslxve/time-tracker-service/internal/model/dto"
	"github.com/njslxve/time-tracker-service/internal/model/entity"
)

const (
	InternalError   = "Internal server error please try again later"
	BadRequestError = "Bad request, please check your request body"
)

// @Summary add user
// @Tags users
// @Description add user
// @Accept json
// @Produce json
// @Param request body dto.AddUserRequest true "request body"
// @Success 201
// @Failure 400 {object} dto.Error
// @Failure 422 {object} dto.Error
// @Failure 500 {object} dto.Error
// @Router       /users/add [post]
func (s *Server) addUserHandler(w http.ResponseWriter, r *http.Request) {
	const op = "server.Server.addUserHandler"

	var req dto.AddUserRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		e := dto.Error{
			Message: BadRequestError,
		}

		s.logger.Error(op, slog.String("error", err.Error()))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(e)

		return
	}

	if !strings.Contains(req.Passport, " ") {
		e := dto.Error{
			Message: "passport must contain space",
		}

		s.logger.Debug(op, slog.String("error", errors.New("passport must contain space").Error()))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(e)

		return
	}

	if err := s.service.AddUser(req); err != nil {
		e := dto.Error{
			Message: InternalError,
		}

		s.logger.Error(op, slog.String("error", err.Error()))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(e)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
}

// @Summary add task
// @Tags tasks
// @Description add task
// @Accept json
// @Produce json
// @Param request body dto.TaskRequest true "request body"
// @Success 201
// @Failure 400 {object} dto.Error
// @Failure 500 {object} dto.Error
// @Router       /tasks/start [post]
func (s *Server) addTaskHandler(w http.ResponseWriter, r *http.Request) {
	const op = "server.Server.addTaskHandler"

	var req dto.TaskRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		e := dto.Error{
			Message: BadRequestError,
		}

		s.logger.Error(op, slog.String("error", err.Error()))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(e)
	}

	if err := s.service.AddTask(req); err != nil {
		e := dto.Error{
			Message: InternalError,
		}

		s.logger.Error(op, slog.String("error", err.Error()))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(e)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
}

// @Summary end task
// @Tags tasks
// @Description end task
// @Accept json
// @Produce json
// @Param request body dto.TaskRequest true "request body"
// @Success 200
// @Failure 400 {object} dto.Error
// @Failure 500 {object} dto.Error
// @Router       /tasks/end [post]
func (s *Server) endTaskHandler(w http.ResponseWriter, r *http.Request) {
	const op = "server.Server.endTaskHandler"

	var req dto.TaskRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		e := dto.Error{
			Message: BadRequestError,
		}

		s.logger.Error(op, slog.String("error", err.Error()))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(e)
	}

	if err := s.service.EndTask(req); err != nil {
		e := dto.Error{
			Message: InternalError,
		}

		s.logger.Error(op, slog.String("error", err.Error()))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(e)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

// @Summary update user
// @Tags users
// @Description update user
// @Accept json
// @Produce json
// @Param user path string true "user id"
// @Param request body dto.UpdateUserRequest true "request body"
// @Success 200
// @Failure 400 {object} dto.Error
// @Failure 500 {object} dto.Error
// @Router       /users/{user} [patch]
func (s *Server) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	const op = "server.Server.updateUserHandler"

	var req dto.UpdateUserRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		e := dto.Error{
			Message: BadRequestError,
		}

		s.logger.Error(op, slog.String("error", err.Error()))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(e)
	}

	userID := chi.URLParam(r, "user")

	if err := s.service.UpdateUser(userID, req); err != nil {
		e := dto.Error{
			Message: InternalError,
		}

		s.logger.Error(op, slog.String("error", err.Error()))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(e)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

// @Summary delete user
// @Tags users
// @Description delete user
// @Accept json
// @Produce json
// @Param user path string true "user id"
// @Success 200
// @Failure 500 {object} dto.Error
// @Router       /users/{user} [delete]
func (s *Server) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	const op = "server.Server.deleteUserHandler"

	userID := chi.URLParam(r, "user")

	if err := s.service.DeleteUser(userID); err != nil {
		e := dto.Error{
			Message: InternalError,
		}

		s.logger.Error(op, slog.String("error", err.Error()))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(e)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

// @Summary get tasks
// @Tags tasks
// @Description get tasks
// @Accept json
// @Produce json
// @Param user path string true "user id"
// @Param interval query string false "interval"
// @Success 200 {array} dto.TaskResponse
// @Failure 500 {object} dto.Error
// @Router       /tasks/{user} [get]
func (s *Server) getTasksHandler(w http.ResponseWriter, r *http.Request) {
	const op = "server.Server.getTasksHandler"

	userID := chi.URLParam(r, "user")
	interval := r.URL.Query().Get("interval")

	tasks, err := s.service.GetTasks(userID, interval)
	if err != nil {
		e := dto.Error{
			Message: InternalError,
		}

		s.logger.Error(op, slog.String("error", err.Error()))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(e)
	} else {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(tasks)
	}
}

// @Summary get users
// @Tags users
// @Description get users
// @Accept json
// @Produce json
// @Param name query string false "name"
// @Param surname query string false "surname"
// @Param adress query string false "adress"
// @Param limit query int false "limit"
// @Param next_page query string false "next_page"
// @Success 200 {object} dto.GetUsersResponse
// @Failure 500 {object} dto.Error
// @Router       /users [get]
func (s *Server) getUsersHandler(w http.ResponseWriter, r *http.Request) {
	const op = "server.Server.getUsersHandler"

	filterOps := entity.FilterOptions{
		Name:       r.URL.Query().Get("name"),
		Surname:    r.URL.Query().Get("surname"),
		Patronymic: r.URL.Query().Get("patronymic"),
		Adress:     r.URL.Query().Get("adress"),
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	paginationOpts := entity.PaginationOptions{
		Limit: limit,
		Next:  r.URL.Query().Get("next_page"),
	}

	users, err := s.service.GetUsers(filterOps, paginationOpts)
	if err != nil {
		e := dto.Error{
			Message: InternalError,
		}

		s.logger.Error(op, slog.String("error", err.Error()))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(e)
	} else {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(users)
	}
}
