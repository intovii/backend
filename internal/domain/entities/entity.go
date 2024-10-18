package entities

import (
	// "database/sql"
	"time"
)


type ReviewDTO struct {
	ID uint64
	Text string
	Mark uint16
	Reviewer UserDTO
	Deal Deal
}
type Review struct {
	ID              uint64
	Text            string
	Mark            uint16
	Reviewer        User
	Deal Deal
}

func ConvertReviewToDTO(review *Review, dto *ReviewDTO) {
	dto.ID = review.ID
	dto.Text = review.Text
	dto.Mark = review.Mark
	ConvertUserToDTO(&review.Reviewer, &dto.Reviewer)
	dto.Deal = review.Deal
}
func ConvertDTOToReview(dto *ReviewDTO, review *Review) {
	review.ID = dto.ID
	review.Text = dto.Text
	review.Mark = dto.Mark
	ConvertDTOToUser(&dto.Reviewer, &review.Reviewer,)
	review.Deal = dto.Deal
}

type AdPhoto struct {
	ID              uint64
	Path            string
	AdvertisementID uint64
}

type Deal struct {
	ID              uint64
	AdvertisementID uint64
	BuyerID uint64
	DateDeal time.Time
}

type Statistic struct {
	AdID		uint64
	AdName		string
	AdPrice		float32
	AdReview	Review
	AdPhoto		AdPhoto
}

type StatisticDTO struct {
	AdID		int64
	AdName		string
	AdPrice		float32
	AdReview	Review
	AdPhoto		AdPhoto
}

func ConvertDTOToStatistic(dto *StatisticDTO, stat *Statistic) {
	stat.AdID = uint64(dto.AdID)
	stat.AdName = dto.AdName
	stat.AdPrice = dto.AdPrice
	stat.AdReview = dto.AdReview
	stat.AdPhoto = dto.AdPhoto
}