package entity

import "time"

type BaseEntity struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy *string
	UpdatedBy *string
	IsDeleted int
	DeletedAt *time.Time
}
