package postgres

import (
	"backend/config"
	"backend/internal/domain/entities"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type Repository struct {
	ctx context.Context
	log *zap.Logger
	cfg *config.ConfigModel
	DB  *pgxpool.Pool
}

func NewRepository(log *zap.Logger, cfg *config.ConfigModel, ctx context.Context) (*Repository, error) {
	return &Repository{
		ctx: ctx,
		log: log,
		cfg: cfg,
	}, nil
}

func (r *Repository) OnStart(_ context.Context) error {
	connectionUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		r.cfg.Postgres.Host,
		r.cfg.Postgres.Port,
		r.cfg.Postgres.User,
		r.cfg.Postgres.Password,
		r.cfg.Postgres.DBName,
		r.cfg.Postgres.SSLMode)
	pool, err := pgxpool.Connect(r.ctx, connectionUrl)
	if err != nil {
		return err
	}
	r.DB = pool
	return nil
}

func (r *Repository) OnStop(_ context.Context) error {
	r.DB.Close()
	return nil
}


const qGetClientWorkouts = `
WITH workout_data AS (
    SELECT
        w.date,
        w.start_time,
        w.end_time,
        w.format_id,
        w.coach_id AS trainerID,
        wf.filial_id,
        w.status,
        ARRAY_AGG(ctw.client_id) AS usersID
    FROM
        workouts w
            LEFT JOIN
        clientsToWorkout ctw ON w.id = ctw.workout_id
            LEFT JOIN
        workoutFormat wf ON w.format_id = wf.id
    GROUP BY
        w.date, w.start_time, w.end_time, w.format_id, w.coach_id, wf.filial_id, w.status
)
SELECT
    wd.date,
    JSON_AGG(
            JSON_BUILD_OBJECT(
                    'startTime', wd.start_time,
                    'endTime', wd.end_time,
                    'formatID', wd.format_id,
                    'usersID', wd.usersID,
                    'trainerID', wd.trainerID,
                    'filialID', wd.filial_id,
                    'status', wd.status
            )
    ) AS workouts
FROM
    workout_data wd
WHERE
    wd.date BETWEEN $1 AND $2
  AND $3 = ANY(wd.usersID)
GROUP BY
    wd.date
ORDER BY
    wd.date`

func (r *Repository) GetClientWorkouts(ctx context.Context, request *entities.SchedulerGetter) ([]entities.Day, error) {
	rows, err := r.DB.Query(ctx, qGetClientWorkouts,
		request.Start.Format(time.DateOnly), request.End.Format(time.DateOnly),
		request.ID)
	if err != nil {
		r.log.Error("fail to get client schedule", zap.Error(err),
			zap.String("query", qGetClientWorkouts),
			zap.String("id", strconv.FormatUint(request.ID, 10)),
			zap.String("start", request.Start.String()),
			zap.String("start", request.End.String()))
		return nil, err
	}
	defer rows.Close()

	var days []entities.Day

	for rows.Next() {
		var date string
		var workoutsJSON []byte

		err := rows.Scan(&date, &workoutsJSON)
		if err != nil {
			r.log.Error("fail to scan row", zap.Error(err))
			return nil, err
		}

		var workouts []entities.Workout
		err = json.Unmarshal(workoutsJSON, &workouts)
		if err != nil {
			r.log.Error("fail to unmarshal workouts JSON", zap.Error(err))
			continue
		}

		days = append(days, entities.Day{
			Date:     date,
			Workouts: workouts,
		})
	}

	if rows.Err() != nil {
		r.log.Error("rows iteration error", zap.Error(rows.Err()))
		return nil, rows.Err()
	}

	return days, nil
}
