package postgres

import (
	"backend/config"
	"backend/internal/domain/entities"
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	// "github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	// "log"
	// "strconv"
	// "time"
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

const queryGetAd = `
SELECT EXISTS (SELECT id
FROM advertisements
WHERE id = $1);
`

func (r *Repository) IsAdExist(ctx context.Context, advertisment *entities.Advertisment) (bool, error) {
	var res bool
	err := r.DB.QueryRow(ctx, queryGetAd, advertisment.ID).Scan(&res)
	if err != nil{
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

func (r *Repository) GetAdvertismentAllInfo(ctx context.Context, advertisment *entities.Advertisment) error {
	adto := &entities.AdvertismentDTO{
		ID: advertisment.ID,
	}
	if err := r.DB.QueryRow(
		ctx,
		queryGetAdInfo,
		adto.ID,
	).Scan(
		&adto.User.ID,
		&adto.Name,
		&adto.Description,
		&adto.Price,
		&adto.DatePlacement,
		&adto.Location,
		&adto.ViewsCount,
		&adto.DateExpirePromotion,
		&adto.AdvertismentCategory.ID,
		&adto.TypePromotion.ID,
		&adto.TypePromotion.Name,
		&adto.TypePromotion.Price,
		&adto.TypePromotion.TimeLive,
		&adto.AdvertismentCategory.Name,
	); err != nil{
		r.log.Error("GetAdvertismentAllInfo: error with SELECT FROM", zap.Error(err))
		return err
	}
	entities.ConvertDTOToAdvertisment(adto, advertisment)
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
	udto := &entities.UserDTO{
		ID: user.ID,
	}
	err := r.DB.QueryRow(
		ctx,
		queryGetUserInfo,
		user.ID,
	).Scan(
		&udto.PathAva,
		&udto.Username,
		&udto.Firstname,
		&udto.Lastname,
		&udto.NumberPhone,
		&udto.Rating,
		&udto.VerificationStatus,
		&udto.Role.ID,
		&udto.Role.Name,
		
	)
	entities.ConvertDTOToUser(udto, user)
	if err != nil{
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
	if err != nil{
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
    d.buyer_id,
	u.username,
	u.path_ava,
	u.firstname,
	u.lastname,
	u.number_phone,
	u.rating,
	u.verification_status,
	u.role_id,
	ur.name
FROM reviews r
JOIN deals d ON d.id = r.deal_id
JOIN users u ON d.buyer_id = u.id
JOIN user_roles ur ON u.role_id = ur.id
WHERE d.advertisement_id = $1;
`

func (r *Repository) GetAdvertismentReviews(ctx context.Context, advertisment *entities.Advertisment) error {
	rows, err := r.DB.Query(
		ctx,
		queryGetReviews,
		advertisment.ID,
	)
	if err != nil{
		r.log.Error("GetReviews: error with SELECT FROM", zap.Error(err))
		return err
	}
	defer rows.Close()

	// Сканируем результаты в срез
	for rows.Next() {
		var rdto entities.ReviewDTO
		var review entities.Review
		if err := rows.Scan(
			&rdto.ID,
			&rdto.Text,
			&rdto.Mark,
			&rdto.Reviewer.ID,
			&rdto.Reviewer.Username,
			&rdto.Reviewer.PathAva,
			&rdto.Reviewer.Firstname,
			&rdto.Reviewer.Lastname,
			&rdto.Reviewer.NumberPhone,
			&rdto.Reviewer.Rating,
			&rdto.Reviewer.VerificationStatus,
			&rdto.Reviewer.Role.ID,
			&rdto.Reviewer.Role.Name,
		); err != nil {
			r.log.Error("GetReviews: error with scan row", zap.Error(err))
			return err		
		}
		entities.ConvertDTOToReview(&rdto, &review)
		advertisment.Reviews = append(advertisment.Reviews, review)
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

func (r *Repository) GetAdvertismentPhotos(ctx context.Context, advertisment *entities.Advertisment) error {
	rows, err := r.DB.Query(
		ctx,
		queryGetPhotos,
		advertisment.ID,
	)
	if err != nil{
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
		advertisment.Photos = append(advertisment.Photos, adPhoto)
	}

	if err := rows.Err(); err != nil {
		r.log.Error("GetPhotos: error iterating through rows", zap.Error(err))
		return err
	}
	return nil
}

func (r *Repository) GetMainAdPhotoByAdID(ctx context.Context, adID uint64) (string, error) {
	var adPhotoPath string
	if err := r.DB.QueryRow(ctx, queryGetPhotos, adID,).Scan(&adPhotoPath,); err != nil{
		r.log.Error("GetStatisticAdPhoto: error with SELECT FROM", zap.Error(err))
		return "", err
	}
	
	return adPhotoPath, nil
}

const  queryGetReviewByDealID = `
SELECT EXIST	(SELECT 
					r.id
				FROM reviews r
				WHERE r.deal_id = $1;)
`

func (r *Repository) IsReviewExistByDealID(ctx context.Context, dealID uint64) (bool, error) {
	var res bool
	err := r.DB.QueryRow(ctx, queryGetReviewByDealID, dealID).Scan(&res)
	if err != nil{
		r.log.Error("IsAdExist: error with QueryRow", zap.Error(err))
		return false, err
	}
	return res, nil
}

const queryGetGetStatisticAdReviewMark = `
SELECT 
	r.id,
	r.mark
FROM reviews r
WHERE r.deal_id = $1;
`

func (r *Repository) GetStatisticAdReviewMark(ctx context.Context, stat *entities.ProfileStatistic) error {
	if err := r.DB.QueryRow(ctx, queryGetGetStatisticAdReviewMark, stat.AdID,).Scan(
		&stat.DealReviewID,
		&stat.AdReviewMark,
		); err != nil{
		r.log.Error("GetStatisticAdPhoto: error with SELECT FROM", zap.Error(err))
		return err
	}
	
	return nil
}


// const queryGetPhone = `
// SELECT 
//     *
// FROM users
// WHERE number_phone = $1
// `

// func (r *Repository) IsPhoneExist(ctx context.Context, user *entities.SqlUser) (bool, error) {
// 	var res bool
// 	err := r.DB.QueryRow(ctx, queryGetPhone, user.NumberPhone).Scan(&res)
// 	log.Println(err)
// 	if err != nil {
// 		r.log.Error("IsPhoneExist: error with QueryRow", zap.Error(err))
// 		return false, err
// 	}
// 	return res, nil

// }

// const queryGetUsername = `
// SELECT 
//     *
// FROM users
// WHERE username = $1
// `

// func (r *Repository) IsUsernameExist(ctx context.Context, user *entities.SqlUser) (bool, error) {
// 	var res bool
// 	err := r.DB.QueryRow(ctx, queryGetUsername, user.Username).Scan(&res)
// 	log.Println(err)
// 	if err != nil {
// 		r.log.Error("IsPhoneExist: error with QueryRow", zap.Error(err))
// 		return false, err
// 	}
// 	return res, nil

// }

// const queryCreateUser = `
// INSERT INTO users 
//     (id, path_ava, username, firstname, lastname, number_phone) 
// VALUES
//     ($1, $2, $3, $4, $5, $6)
// `

// func (r *Repository) CreateUser(ctx context.Context, user *entities.SqlUser) error {
// 	// Выполняем команду INSERT
// 	result, err := r.DB.Exec(
// 		ctx,
// 		queryCreateUser,
// 		user.ID,
// 		user.PathAva,
// 		user.Username,
// 		user.Firstname,
// 		user.Lastname,
// 		user.NumberPhone,
// 	)
// 	if err != nil {
// 		r.log.Error("CreateUser: error with INSERT INTO", zap.Error(err))
// 		return err
// 	}

// 	rowsAffected := result.RowsAffected()
// 	if rowsAffected == 0 {
// 		return fmt.Errorf("no rows affected, user may not be created")
// 	}

// 	return nil
// }


const queryGetAdsByBuyerID = `
	SELECT 
		d.id,
		d.advertisement_id,
		a.name,
		a.price
		FROM deals d
		JOIN advertisements a ON d.advertisement_id = a.id
	WHERE d.buyer_id = $1;
`
		// r.mark
	// JOIN reviews r ON d.id = r.deal_id

func (r *Repository) GetProfileUserStatistics(ctx context.Context, uID uint64, stats *[]*entities.ProfileStatistic) error {
	rows, err := r.DB.Query(ctx, queryGetAdsByBuyerID, uID)
	if err != nil{
		r.log.Error("GetAdvertismentsByBuyerID: error with SELECT FROM", zap.Error(err))
		return err
	}
	defer rows.Close()

	// Сканируем результаты в срез
	for rows.Next() {
		stat := &entities.ProfileStatistic{}
		// dto := &entities.StatisticDTO{
		// 	AdReview: entities.Review{},
		// }
		if err := rows.Scan(
			&stat.DealID,
			&stat.AdID,
			&stat.AdName,
			&stat.AdPrice,
			// &stat.AdReviewMark,
			); err != nil {
			r.log.Error("GetAdvertismentsByBuyerID: error with scan row", zap.Error(err))
			return err		
		}
		// entities.ConvertDTOToStatistic(dto, stat)
		*stats = append(*stats, stat)
	}

	if err := rows.Err(); err != nil {
		r.log.Error("GetAdvertismentsByBuyerID: error iterating through rows", zap.Error(err))
		return err
	}
	return nil
}


const queryGetProfileMyAdvertisments = `
SELECT
	a.id,
	a.name,
	a.price,
	a.views_count,
	a.type_id,
	tp.name,
	a.date_expire_promotion
FROM advertisements a
	JOIN types_promotion tp ON a.type_id = tp.id
WHERE 
	a.user_id = $1;
`


func (r *Repository) GetProfileMyAdvertisments(ctx context.Context, uID uint64, advertisements *[]*entities.MyAdvertisement) error {
	rows, err := r.DB.Query(ctx, queryGetProfileMyAdvertisments, uID)
	if err != nil{
		r.log.Error("GetProfileMyAdvertisments: error with SELECT FROM", zap.Error(err))
		return err
	}
	defer rows.Close()

	// Сканируем результаты в срез
	for rows.Next() {
		ad := entities.MyAdvertisement{}
		if err := rows.Scan(
			&ad.AdID,
			&ad.AdName,
			&ad.AdPrice,
			&ad.AdCountViews,
			&ad.AdTypePromotionID,
			&ad.AdTypePromotionName,
			&ad.AdDateExpirePromotion,
			); err != nil {
			r.log.Error("GetProfileMyAdvertisments: error with scan row", zap.Error(err))
			return err		
		}
		*advertisements = append(*advertisements, &ad)
	}

	if err := rows.Err(); err != nil {
		r.log.Error("GetProfileMyAdvertisments: error iterating through rows", zap.Error(err))
		return err
	}
	return nil
}

const queryGetProfileReviews = `
SELECT
	a.id,
	d.id,
	r.id,
	u.id,
	r.text,
	r.mark,
	u.path_ava,
	u.username,
	u.firstname,
	u.lastname
FROM advertisements a
	JOIN deals d ON a.id = d.advertisement_id
	JOIN reviews r ON d.id = r.deal_id
	JOIN users u ON d.buyer_id = u.id
WHERE 
	a.user_id = $1;
`

func (r *Repository) GetProfileReviews(ctx context.Context, uID uint64, reviews *[]*entities.ProfileReview) error {
	rows, err := r.DB.Query(ctx, queryGetProfileReviews, uID)
	if err != nil{
		r.log.Error("GetProfileReviews: error with SELECT FROM", zap.Error(err))
		return err
	}
	defer rows.Close()

	// Сканируем результаты в срез
	for rows.Next() {
		review := entities.ProfileReview{}
		dto := entities.ProfileReviewDTO{}
		if err := rows.Scan(
			&dto.AdID,
			&dto.DealID,
			&dto.ReviewID,
			&dto.ReviewerID,
			&dto.ReviewText,
			&dto.ReviewMark,
			&dto.ReviewerPathAva,
			&dto.ReviewerUsername,
			&dto.ReviewerFirstname,
			&dto.ReviewerLastname,
			); err != nil {
			r.log.Error("GetProfileReviews: error with scan row", zap.Error(err))
			return err		
		}
		entities.ConvertDTOToProfileReview(&dto, &review)
		*reviews = append(*reviews, &review)
	}

	if err := rows.Err(); err != nil {
		r.log.Error("GetProfileReviews: error iterating through rows", zap.Error(err))
		return err
	}
	return nil
}