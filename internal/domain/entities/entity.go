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
	AdvertisementID uint64
}
type Review struct {
	ID uint64
	Text string
	Mark uint16
	Reviewer User
	AdvertisementID uint64
}

func ConvertReviewToDTO(review *Review, dto *ReviewDTO) {
	dto.ID = review.ID
	dto.Text = review.Text
	dto.Mark = review.Mark
	ConvertUserToDTO(&review.Reviewer, &dto.Reviewer)
	dto.AdvertisementID = review.AdvertisementID
}
func ConvertDTOToReview(dto *ReviewDTO, review *Review) {
	review.ID = dto.ID
	review.Text = dto.Text
	review.Mark = dto.Mark
	ConvertDTOToUser(&dto.Reviewer, &review.Reviewer,)
	review.AdvertisementID = dto.AdvertisementID	
}

type AdPhoto struct {
	ID uint64
	Path string
	AdvertisementID uint64
}

type Deal struct {
	ID uint64
	AdvertisementID uint64
	BuyerID uint64
	DateDeal time.Time
}