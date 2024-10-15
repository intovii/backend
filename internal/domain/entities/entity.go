package entities

import (
	"time"
)

type Workout struct {
	StartTime string
	EndTime   string
	FormatID  uint64
	UsersID   []uint64
	TrainerID uint64
	FilialID  uint64
	Status    string
	Date      string
}

type Day struct {
	Workouts []Workout
	Date     string
}

type SchedulerGetter struct {
	ID    uint64
	Start time.Time
	End   time.Time
}

type UserRole struct {
	ID   uint32
	Name string
}

type ProductCategory struct {
	ID   uint64
	Name string
}

type TypePromotion struct {
	ID       uint64
	Name     string
	Price    float32
	TimeLive time.Duration
}

type User struct {
	ID                 uint64
	PathAva            string
	Username           string
	Firstname          string
	Lastname           string
	NumberPhone        string
	Rating             float32
	VerificationStatus string
	Role               UserRole
}

type Advertisment struct {
	ID                  uint64
	User                User
	Name                string
	Description         string
	Price               float64
	DatePlacement       time.Time
	Location            string
	TypePromotion       TypePromotion
	ViewsCount          uint32
	DateExpirePromotion time.Time
	ProductCategory     ProductCategory
	Reviews             []Review
	Photos              []AdPhoto
}

type Review struct {
	ID              uint64
	Text            string
	Mark            uint16
	Reviewer        User
	AdvertisementID uint64
}

type AdPhoto struct {
	ID              uint64
	Path            string
	AdvertisementID uint64
}

type Deal struct {
	ID              uint64
	AdvertisementID uint64
	BuyerID         uint64
	DateDeal        uint64
}

type CreateUser struct {
	ID          uint64
	PathAva     string // Может быть null
	Username    string // Может быть null
	Firstname   string // Может быть null
	Lastname    string // Может быть null
	NumberPhone string // Может быть null
}
