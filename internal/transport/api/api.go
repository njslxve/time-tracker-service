package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/njslxve/time-tracker-service/internal/config"
	"github.com/njslxve/time-tracker-service/internal/model/dto"
)

type API struct {
	logger *slog.Logger
	cfg    *config.Config
}

func New(logger *slog.Logger, cfg *config.Config) *API {
	return &API{
		logger: logger,
		cfg:    cfg,
	}
}

func (a *API) Info(passport string) (dto.UserInfoResponse, error) {
	const op = "api.API.Info"

	data := strings.Split(passport, " ")

	serie, number := data[0], data[1]

	url := fmt.Sprintf("%s?passportSerie=%s&passportNumber=%s", a.cfg.InfoAPIURL, serie, number)

	resp, err := http.Get(url)
	if err != nil {
		a.logger.Debug("external api error",
			slog.String("description", op),
			slog.String("error", err.Error()),
		)

		return dto.UserInfoResponse{}, fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()

	var userInfo dto.UserInfoResponse

	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		a.logger.Debug("could not decode response",
			slog.String("description", op),
			slog.String("error", err.Error()),
		)

		return dto.UserInfoResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	return userInfo, nil
}
