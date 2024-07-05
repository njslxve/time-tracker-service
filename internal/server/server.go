package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/njslxve/time-tracker-service/internal/config"
	"github.com/njslxve/time-tracker-service/internal/service"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Server struct {
	cfg     *config.Config
	logger  *slog.Logger
	service *service.Service
}

func New(cfg *config.Config, logger *slog.Logger, service *service.Service) *Server {
	return &Server{
		cfg:     cfg,
		logger:  logger,
		service: service,
	}
}

func (s *Server) Start() {
	r := chi.NewRouter()

	r.Use(middleware.Timeout(30 * time.Second))

	r.Route("/users", func(r chi.Router) {
		r.Get("/", s.getUsersHandler)
		r.Post("/add", s.addUserHandler)
		r.Patch("/{user}", s.updateUserHandler)
		r.Delete("/{user}", s.deleteUserHandler)
	})

	r.Route("/tasks", func(r chi.Router) {
		r.Get("/{user}", s.getTasksHandler)
		r.Post("/start", s.addTaskHandler)
		r.Post("/end", s.endTaskHandler)
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://%s/swagger/doc.json", s.cfg.Address)),
	))

	s.logger.Info("starting server",
		slog.String("address", s.cfg.Address))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         s.cfg.Address,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			s.logger.Debug("server error",
				slog.String("error", err.Error()),
			)

			if !errors.Is(err, http.ErrServerClosed) {
				s.logger.Error("failed to start server")
			}
		}
	}()

	s.logger.Info("server started")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	<-done
	s.logger.Info("shutting down server")

	if err := srv.Shutdown(ctx); err != nil {
		s.logger.Error("failed to shutdown server")
	}

	s.logger.Info("server shutdown")
}
