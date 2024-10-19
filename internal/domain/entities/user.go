package entities

import (
	"database/sql"
)

type UserRole struct {
	ID   uint32
	Name string
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

type UserDTO struct {
	ID                 uint64         `json:"id" db:"id"`
	PathAva            sql.NullString `json:"path_ava" db:"path_ava"`
	Username           sql.NullString `json:"username" db:"username"`
	Firstname          sql.NullString `json:"firstname" db:"firstname"`
	Lastname           sql.NullString `json:"lastname" db:"lastname"`
	NumberPhone        sql.NullString `json:"number_phone" db:"number_phone"`
	Rating             float32        `json:"rating" db:"rating"`
	VerificationStatus string         `json:"verification_status" db:"verification_status"`
	Role               UserRole
}

func ConvertDTOToUser(dto *UserDTO, u *User) {
	u.ID = dto.ID
	u.PathAva = dto.PathAva.String
	u.Username = dto.Username.String
	u.Firstname = dto.Firstname.String
	u.Lastname = dto.Lastname.String
	u.NumberPhone = dto.NumberPhone.String
	u.Rating = dto.Rating
	u.VerificationStatus = dto.VerificationStatus
	u.Role = dto.Role
}

func ConvertUserToDTO(u *User, dto *UserDTO) {
	dto.ID = u.ID
	dto.PathAva = *NewNullString(u.PathAva)
	dto.Username = *NewNullString(u.Username)
	dto.Firstname = *NewNullString(u.Firstname)
	dto.Lastname = *NewNullString(u.Lastname)
	dto.NumberPhone = *NewNullString(u.NumberPhone)
	dto.Rating = u.Rating
	dto.VerificationStatus = u.VerificationStatus
	dto.Role = u.Role
}
