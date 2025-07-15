package models

import "time"

type Subscription struct {
	ID          int        `json:"id" db:"id"`
	ServiceName string     `json:"service_name" db:"service_name" binding:"required"`
	Price       int        `json:"price" db:"price" binding:"required"`
	UserID      string     `json:"user_id" db:"user_id" binding:"required"`
	StartDate   time.Time  `json:"start_date" db:"start_date" binding:"required"`
	EndDate     *time.Time `json:"end_date,omitempty" db:"end_date"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
}
