package staff

import "time"

type Staff struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Username     string    `json:"username" gorm:"uniqueIndex;size:50;not null"`
	PasswordHash string    `json:"-" gorm:"size:100;not null"`
	FullName     string    `json:"full_name" gorm:"size:100;not null"`
	Role         string    `json:"role" gorm:"size:20;not null"`
	IsActive     bool      `json:"is_active" gorm:"default:true"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
}
