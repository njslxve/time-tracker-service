package storage

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/njslxve/time-tracker-service/internal/model/entity"
)

type Storage struct {
	logger *slog.Logger
	db     *pgx.Conn
}

func New(logger *slog.Logger, client *pgx.Conn) *Storage {
	return &Storage{
		logger: logger,
		db:     client,
	}
}

var (
	qb = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
)

func (s *Storage) AddUser(user entity.User) error {
	const op = "transport.storage.AddUser"

	uuid := uuid.NewString()

	querry := qb.Insert("users").
		Columns("id", "passport", "first_name", "last_name", "patronymic", "adress").
		Values(uuid, user.Passport, user.Name, user.Surmame, user.Patronymic, user.Adress)

	sql, args, err := querry.ToSql()
	if err != nil {
		s.logger.Debug("could not convert query to sql",
			slog.String("description", op),
			slog.String("error", err.Error()),
		)

		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.db.Exec(context.Background(), sql, args...)
	if err != nil {
		s.logger.Debug("sql error",
			slog.String("description", op),
			slog.String("sql", sql),
			slog.Any("args", args),
			slog.String("error", err.Error()),
		)

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetUser(userID int) (entity.User, error) {
	const op = "transport.storage.GetUser"

	querry := qb.Select("user_id", "passport", "first_name", "last_name", "patronymic", "adress").
		From("users").
		Where(sq.Eq{"user_id": userID})

	sql, args, err := querry.ToSql()
	if err != nil {
		s.logger.Debug("could not convert query to sql",
			slog.String("description", op),
			slog.String("error", err.Error()),
		)

		return entity.User{}, fmt.Errorf("%s: %w", op, err)
	}

	row := s.db.QueryRow(context.Background(), sql, args...)

	var user entity.User

	err = row.Scan(&user.UserID, &user.Passport, &user.Name, &user.Surmame, &user.Patronymic, &user.Adress)
	if err != nil {
		s.logger.Debug("could not scan row",
			slog.String("description", op),
			slog.String("error", err.Error()),
		)

		return entity.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Storage) GetUsers(opts entity.FilterOptions) ([]entity.User, error) {
	const op = "transport.storage.GetUsers"

	querry := qb.Select("user_id", "passport", "first_name", "last_name", "patronymic", "adress").
		From("users")

	if opts.Name != "" {
		querry = querry.Where(sq.Eq{"first_name": opts.Name})
	}

	if opts.Surname != "" {
		querry = querry.Where(sq.Eq{"last_name": opts.Surname})
	}

	if opts.Patronymic != "" {
		querry = querry.Where(sq.Eq{"patronymic": opts.Patronymic})
	}

	if opts.Adress != "" {
		querry = querry.Where(sq.ILike{"adress": fmt.Sprintf("%%%s%%", opts.Adress)})
	}

	sql, args, err := querry.ToSql()
	if err != nil {
		s.logger.Debug("could not convert query to sql",
			slog.String("description", op),
			slog.String("error", err.Error()),
		)

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := s.db.Query(context.Background(), sql, args...)
	if err != nil {
		s.logger.Debug("sql error",
			slog.String("description", op),
			slog.String("sql", sql),
			slog.Any("args", args),
			slog.String("error", err.Error()),
		)

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer rows.Close()

	var users []entity.User

	for rows.Next() {
		var user entity.User

		err = rows.Scan(&user.UserID, &user.Passport, &user.Name, &user.Surmame, &user.Patronymic, &user.Adress)
		if err != nil {
			s.logger.Debug("could not scan row",
				slog.String("description", op),
				slog.String("error", err.Error()),
			)

			return nil, fmt.Errorf("%s: %w", op, err)
		}

		users = append(users, user)
	}

	return users, nil
}

func (s *Storage) UpdateUser(user entity.User) error {
	const op = "transport.storage.UpdateUser"

	querry := qb.Update("users").
		SetMap(map[string]interface{}{
			"passport":   user.Passport,
			"first_name": user.Name,
			"last_name":  user.Surmame,
			"patronymic": user.Patronymic,
			"adress":     user.Adress,
		}).
		Where(sq.Eq{"user_id": user.UserID})

	sql, args, err := querry.ToSql()
	if err != nil {
		s.logger.Debug("could not convert query to sql",
			slog.String("description", op),
			slog.String("error", err.Error()),
		)

		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.db.Exec(context.Background(), sql, args...)
	if err != nil {
		s.logger.Debug("sql error",
			slog.String("description", op),
			slog.String("sql", sql),
			slog.Any("args", args),
			slog.String("error", err.Error()),
		)

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteUser(userID int) error {
	const op = "transport.storage.DeleteUser"

	querry := qb.Delete("users").
		Where(sq.Eq{"user_id": userID})

	sql, args, err := querry.ToSql()
	if err != nil {
		s.logger.Debug("could not convert query to sql",
			slog.String("description", op),
			slog.String("error", err.Error()),
		)

		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.db.Exec(context.Background(), sql, args...)
	if err != nil {
		s.logger.Debug("sql error",
			slog.String("description", op),
			slog.String("sql", sql),
			slog.Any("args", args),
			slog.String("error", err.Error()),
		)

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) AddTask(task entity.Task) error {
	const op = "transport.storage.AddTask"

	uuid := uuid.NewString()

	querry := qb.Insert("tasks").
		Columns("id", "user_id", "task_id", "start_time").
		Values(uuid, task.UserID, task.TaskID, task.StartTime)

	sql, args, err := querry.ToSql()
	if err != nil {
		s.logger.Debug("could not convert query to sql",
			slog.String("description", op),
			slog.String("error", err.Error()),
		)

		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.db.Exec(context.Background(), sql, args...)
	if err != nil {
		s.logger.Debug("sql error",
			slog.String("description", op),
			slog.String("sql", sql),
			slog.Any("args", args),
			slog.String("error", err.Error()),
		)

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetTask(taskID string, userID int) (entity.Task, error) {
	const op = "transport.storage.GetTask"

	querry := qb.Select("id", "user_id", "task_id", "start_time").
		From("tasks").
		Where(sq.And{
			sq.Eq{"user_id": userID},
			sq.Eq{"task_id": taskID},
		})

	sql, args, err := querry.ToSql()
	if err != nil {
		s.logger.Debug("could not convert query to sql",
			slog.String("description", op),
			slog.String("error", err.Error()),
		)

		return entity.Task{}, fmt.Errorf("%s: %w", op, err)
	}

	var task entity.Task

	err = s.db.QueryRow(context.Background(), sql, args...).Scan(&task.ID, &task.UserID, &task.TaskID, &task.StartTime)
	if err != nil {
		s.logger.Debug("sql error",
			slog.String("description", op),
			slog.String("sql", sql),
			slog.Any("args", args),
			slog.String("error", err.Error()),
		)

		return entity.Task{}, fmt.Errorf("%s: %w", op, err)
	}

	return task, nil
}

func (s *Storage) GetTasks(userID int, interval int) ([]entity.Task, error) {
	const op = "transport.storage.GetTasks"

	querry := qb.Select("user_id", "task_id", "start_time", "end_time", "duration").
		From("tasks").
		Where(sq.And{
			sq.Eq{"user_id": userID},
			sq.NotEq{"duration": nil},
		}).
		OrderBy("duration DESC")

	if interval != 0 {
		querry = querry.
			Where(fmt.Sprintf("start_time >= current_date - interval '%d days'", interval))
	}

	sql, args, err := querry.ToSql()
	if err != nil {
		s.logger.Debug("could not convert query to sql",
			slog.String("description", op),
			slog.String("error", err.Error()),
		)

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := s.db.Query(context.Background(), sql, args...)
	if err != nil {
		s.logger.Debug("sql error",
			slog.String("description", op),
			slog.String("sql", sql),
			slog.Any("args", args),
			slog.String("error", err.Error()),
		)

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var tasks []entity.Task
	for rows.Next() {
		var task entity.Task
		err = rows.Scan(&task.UserID, &task.TaskID, &task.StartTime, &task.EndTime, &task.Duration)
		if err != nil {
			s.logger.Debug("sql error",
				slog.String("description", op),
				slog.String("sql", sql),
				slog.Any("args", args),
				slog.String("error", err.Error()),
			)

			return nil, fmt.Errorf("%s: %w", op, err)
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (s *Storage) UpdateTask(task entity.Task) error {
	const op = "transport.storage.UpdateTask"

	querry := qb.Update("tasks").
		SetMap(sq.Eq{"end_time": task.EndTime, "duration": task.Duration}).
		Where(sq.And{
			sq.Eq{"id": task.ID},
			sq.Eq{"user_id": task.UserID},
		})

	sql, args, err := querry.ToSql()
	if err != nil {
		s.logger.Debug("could not convert query to sql",
			slog.String("description", op),
			slog.String("error", err.Error()),
		)

		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.db.Exec(context.Background(), sql, args...)
	if err != nil {
		s.logger.Debug("sql error",
			slog.String("description", op),
			slog.String("sql", sql),
			slog.Any("args", args),
			slog.String("error", err.Error()),
		)

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) TokenData(token string) (int, entity.TokenData, error) {
	const op = "transport.storage.TokenData"

	querry := qb.Select("old_limit", "filter_params", "is_alive").
		From("pagination_tokens").
		Where(sq.Eq{"token": token})

	sql, args, err := querry.ToSql()
	if err != nil {
		s.logger.Debug("could not convert query to sql",
			slog.String("description", op),
			slog.String("error", err.Error()),
		)

		return 0, entity.TokenData{}, fmt.Errorf("%s: %w", op, err)
	}

	var tokenData entity.TokenData

	err = s.db.QueryRow(context.Background(), sql, args...).Scan(&tokenData.Old, &tokenData.Params, &tokenData.IsAlive)
	if err != nil {
		s.logger.Debug("sql error",
			slog.String("description", op),
			slog.String("sql", sql),
			slog.Any("args", args),
			slog.String("error", err.Error()),
		)

		return 0, entity.TokenData{}, fmt.Errorf("%s: %w", op, err)
	}

	return tokenData.Old, tokenData, nil
}

func (s *Storage) AddToken(token string, limit int, params []byte) error {
	const op = "transport.storage.AddToken"

	querry := qb.Insert("pagination_tokens").
		Columns("token", "old_limit", "filter_params", "created_at").
		Values(token, limit, params, time.Now())

	sql, args, err := querry.ToSql()
	if err != nil {
		s.logger.Debug("could not convert query to sql",
			slog.String("description", op),
			slog.String("error", err.Error()),
		)

		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.db.Exec(context.Background(), sql, args...)
	if err != nil {
		s.logger.Debug("sql error",
			slog.String("description", op),
			slog.String("sql", sql),
			slog.Any("args", args),
			slog.String("error", err.Error()),
		)

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
