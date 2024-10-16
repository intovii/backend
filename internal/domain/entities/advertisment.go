package entities

import (
	"database/sql"
	"time"
)


type TypePromotion struct {
	ID uint64
	Name string
	Price float32
	TimeLive time.Duration
}

type AdvertismentCategory struct {
	ID uint64
	Name string
}

type Advertisment struct {
	ID	uint64
	User User
	Name string
	Description string
	Price float64
	DatePlacement *time.Time
	Location string
	TypePromotion TypePromotion
	ViewsCount uint32
	DateExpirePromotion *time.Time
	AdvertismentCategory AdvertismentCategory 
	Reviews []Review
	Photos []AdPhoto
}

type AdvertismentDTO struct {
	ID	uint64 `json:"id" db:"id"`
	User UserDTO
	Name string `json:"name" db:"name"`
	Description sql.NullString `json:"description" db:"description"`
	Price float64 `json:"price" db:"price"`
	DatePlacement sql.NullTime `json:"date_placement" db:"date_placement"` //TODO 
	Location sql.NullString `json:"location" db:"location"`
	TypePromotion TypePromotion 
	ViewsCount uint32 `json:"views_count" db:"views_count"`
	DateExpirePromotion sql.NullTime `json:"date_expire_promotion" db:"date_expire_promotion"` //TODO 
	AdvertismentCategory AdvertismentCategory 
	Reviews []Review
	Photos []AdPhoto
}


func ConvertAdvertismentToDTO(a *Advertisment, dto *AdvertismentDTO) {
		dto.ID = a.ID
		ConvertUserToDTO(&a.User, &dto.User)
		dto.Name = a.Name
		dto.Description = *NewNullString(a.Description)
		dto.Price = a.Price
		dto.DatePlacement = *NewNullTime(a.DatePlacement)
		dto.Location = *NewNullString(a.Location)
		dto.TypePromotion = a.TypePromotion
		dto.ViewsCount = a.ViewsCount
		dto.DateExpirePromotion = *NewNullTime(a.DateExpirePromotion)
		dto.AdvertismentCategory = a.AdvertismentCategory
		dto.Reviews = a.Reviews
		dto.Photos = a.Photos
}

func ConvertDTOToAdvertisment(dto *AdvertismentDTO, a *Advertisment) {
	a.ID = dto.ID
	ConvertDTOToUser(&dto.User, &a.User)
	a.Name = dto.Name
	a.Description = dto.Description.String
	a.Price = dto.Price
	a.DatePlacement = &dto.DatePlacement.Time
	a.Location = dto.Location.String
	a.TypePromotion = dto.TypePromotion
	a.ViewsCount = dto.ViewsCount
	if dto.DateExpirePromotion.Valid {
		a.DateExpirePromotion = &dto.DateExpirePromotion.Time
	} else {
		a.DateExpirePromotion = nil //&dto.DateExpirePromotion.Time
	}
	a.AdvertismentCategory = dto.AdvertismentCategory
	a.Reviews = dto.Reviews
	a.Photos = dto.Photos
}