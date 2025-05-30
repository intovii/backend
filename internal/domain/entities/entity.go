package entities

import (
	// "database/sql"
	"database/sql"
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

type ProfileStatistic struct {
	DealID 			uint64
	AdID			uint64
	DealReviewID	uint64
	AdName			string
	AdPrice			float32
	AdReviewMark	uint16
	AdPhotoPath		string
}

type MyAdvertisement struct {
	AdID					uint64
	AdPhotoPath				string
	AdName					string
	AdPrice					float64
	AdCountViews			uint32
	AdTypePromotionID		uint64
	AdTypePromotionName		string
	AdDateExpirePromotion	*time.Time
}

type ProfileReview struct {
	AdID uint64
	DealID uint64
	ReviewID uint64
	ReviewerID uint64
	ReviewText string
	ReviewMark uint16
	ReviewerPathAva string
	ReviewerUsername string
	ReviewerFirstname string
	ReviewerLastname string
}
type ProfileReviewDTO struct {
	AdID uint64
	DealID uint64
	ReviewID uint64
	ReviewerID uint64
	ReviewText string
	ReviewMark uint16
	ReviewerPathAva sql.NullString
	ReviewerUsername sql.NullString
	ReviewerFirstname sql.NullString
	ReviewerLastname sql.NullString
}

func ConvertDTOToProfileReview(dto *ProfileReviewDTO, review *ProfileReview) {
	review.AdID = dto.AdID
	review.DealID = dto.DealID
	review.ReviewID = dto.ReviewID
	review.ReviewerID = dto.ReviewerID
	review.ReviewText = dto.ReviewText
	review.ReviewMark = dto.ReviewMark
	review.ReviewerPathAva = dto.ReviewerPathAva.String
	review.ReviewerUsername = dto.ReviewerUsername.String
	review.ReviewerFirstname = dto.ReviewerFirstname.String
	review.ReviewerLastname = dto.ReviewerLastname.String
}