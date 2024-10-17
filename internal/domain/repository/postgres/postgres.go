package postgres

import (
	"backend/config"
	"backend/internal/domain/entities"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
	"log"
	"strconv"
	"time"
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

const queryGetAd = `
SELECT EXISTS (SELECT id
FROM advertisements
WHERE id = $1);
`

func (r *Repository) IsAdExist(ctx context.Context, product *entities.Advertisment) (bool, error) {
	var res bool
	err := r.DB.QueryRow(ctx, queryGetAd, product.ID).Scan(&res)
	if err != nil {
		r.log.Error("IsAdExist: error with QueryRow", zap.Error(err))
		return false, err
	}
	return res, nil
}

const queryGetAdInfo = `
SELECT 
    a.user_id,
    a.name AS advertisement_name,
    a.description,
    a.price,
    a.date_placement,
    a.location,
    a.views_count,
    a.date_expire_promotion,
    a.category_id,
    a.type_id,
    tp.name AS type_promotion_name,
    tp.price AS type_promotion_name,
    tp.time_live AS type_promotion_name,
    cp.name AS category_name
FROM advertisements a
JOIN types_promotion tp ON a.type_id = tp.id
JOIN categories_product cp ON a.category_id = cp.id
WHERE a.id = $1;
`

// JOIN categories_product u ON a.user_id = u.id

func (r *Repository) GetProductAllInfo(ctx context.Context, product *entities.Advertisment) error {
	err := r.DB.QueryRow(
		ctx,
		queryGetAdInfo,
		product.ID,
	).Scan(
		&product.User.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.DatePlacement,
		&product.Location,
		&product.ViewsCount,
		&product.DateExpirePromotion,
		&product.ProductCategory.ID,
		&product.TypePromotion.ID,
		&product.TypePromotion.Name,
		&product.TypePromotion.Price,
		&product.TypePromotion.TimeLive,
		&product.ProductCategory.Name,
	)
	if err != nil {
		r.log.Error("GetUserProfile: error with SELECT FROM", zap.Error(err))
		return err
	}
	return nil
}

const queryGetUserInfo = `
SELECT 
    u.path_ava,
    u.username,
    u.firstname,
    u.lastname,
    u.number_phone,
    u.rating,
    u.verification_status,
    u.role_id,
    r.name AS role_name
FROM users u
JOIN user_roles r ON u.role_id = r.id
WHERE u.id = $1;

`

func (r *Repository) GetUserInfo(ctx context.Context, user *entities.User) error {
	err := r.DB.QueryRow(
		ctx,
		queryGetUserInfo,
		user.ID,
	).Scan(
		&user.PathAva,
		&user.Username,
		&user.Firstname,
		&user.Lastname,
		&user.NumberPhone,
		&user.Rating,
		&user.VerificationStatus,
		&user.Role.ID,
		&user.Role.Name,
	)
	if err != nil {
		r.log.Error("GetUserInfo: error with SELECT FROM", zap.Error(err))
		return err
	}
	return nil
}

const queryGetUser = `
SELECT EXISTS (SELECT id
FROM users
WHERE id = $1);
`

func (r *Repository) IsUserExist(ctx context.Context, user *entities.User) (bool, error) {
	var res bool

	err := r.DB.QueryRow(ctx, queryGetUser, user.ID).Scan(&res)
	if err != nil {
		r.log.Error("IsAdExist: error with QueryRow", zap.Error(err))
		return false, err
	}
	return res, nil
}

const queryGetReviews = `
SELECT 
    r.id,
    r.text,
    r.mark,
    r.reviewer_id,
	u.username,
	u.path_ava,
	u.firstname,
	u.lastname,
	u.number_phone,
	u.rating,
	u.verification_status,
	u.role_id
FROM reviews r
JOIN users u ON r.reviewer_id = u.id
WHERE r.advertisement_id = $1;
`

func (r *Repository) GetReviews(ctx context.Context, product *entities.Advertisment) error {
	rows, err := r.DB.Query(
		ctx,
		queryGetReviews,
		product.ID,
	)
	if err != nil {
		r.log.Error("GetReviews: error with SELECT FROM", zap.Error(err))
		return err
	}
	defer rows.Close()

	// Сканируем результаты в срез
	for rows.Next() {
		var review entities.Review
		if err := rows.Scan(
			&review.ID,
			&review.Text,
			&review.Mark,
			&review.Reviewer.ID,
			&review.Reviewer.Username,
			&review.Reviewer.PathAva,
			&review.Reviewer.Firstname,
			&review.Reviewer.Lastname,
			&review.Reviewer.NumberPhone,
			&review.Reviewer.Rating,
			&review.Reviewer.VerificationStatus,
			&review.Reviewer.Role.ID,
		); err != nil {
			r.log.Error("GetReviews: error with scan row", zap.Error(err))
			return err
		}
		product.Reviews = append(product.Reviews, review)
	}

	if err := rows.Err(); err != nil {
		r.log.Error("GetReviews: error iterating through rows", zap.Error(err))
		return err
	}
	return nil
}

const queryGetPhotos = `
SELECT 
	ph.path
FROM ad_photos ph
WHERE ph.advertisement_id = $1;
`

func (r *Repository) GetPhotos(ctx context.Context, product *entities.Advertisment) error {
	rows, err := r.DB.Query(
		ctx,
		queryGetPhotos,
		product.ID,
	)
	if err != nil {
		r.log.Error("GetPhotos: error with SELECT FROM", zap.Error(err))
		return err
	}
	defer rows.Close()

	// Сканируем результаты в срез
	for rows.Next() {
		var adPhoto entities.AdPhoto
		if err := rows.Scan(
			&adPhoto.Path,
		); err != nil {
			r.log.Error("GetPhotos: error with scan row", zap.Error(err))
			return err
		}
		product.Photos = append(product.Photos, adPhoto)
	}

	if err := rows.Err(); err != nil {
		r.log.Error("GetPhotos: error iterating through rows", zap.Error(err))
		return err
	}
	return nil
}

const queryGetPhone = `
SELECT 
    *
FROM users
WHERE number_phone = $1
`

func (r *Repository) IsPhoneExist(ctx context.Context, user *entities.SqlUser) (bool, error) {
	var res bool
	err := r.DB.QueryRow(ctx, queryGetPhone, user.NumberPhone).Scan(&res)
	log.Println(err)
	if err != nil {
		r.log.Error("IsPhoneExist: error with QueryRow", zap.Error(err))
		return false, err
	}
	return res, nil

}

const queryGetUsername = `
SELECT 
    *
FROM users
WHERE username = $1
`

func (r *Repository) IsUsernameExist(ctx context.Context, user *entities.SqlUser) (bool, error) {
	var res bool
	err := r.DB.QueryRow(ctx, queryGetUsername, user.Username).Scan(&res)
	log.Println(err)
	if err != nil {
		r.log.Error("IsPhoneExist: error with QueryRow", zap.Error(err))
		return false, err
	}
	return res, nil

}

const queryCreateUser = `
INSERT INTO users 
    (id, path_ava, username, firstname, lastname, number_phone) 
VALUES
    ($1, $2, $3, $4, $5, $6)
`

func (r *Repository) CreateUser(ctx context.Context, user *entities.SqlUser) error {
	// Выполняем команду INSERT
	result, err := r.DB.Exec(
		ctx,
		queryCreateUser,
		user.ID,
		user.PathAva,
		user.Username,
		user.Firstname,
		user.Lastname,
		user.NumberPhone,
	)
	if err != nil {
		r.log.Error("CreateUser: error with INSERT INTO", zap.Error(err))
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected, user may not be created")
	}

	return nil
}
