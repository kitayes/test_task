package models

import "time"

type UpdateSubscriptionInput struct {
	ServiceName *string    `json:"service_name,omitempty" binding:"omitempty,min=1"`
	Price       *int       `json:"price,omitempty" binding:"omitempty,gte=0"`
	StartDate   *time.Time `json:"start_date,omitempty"`
	EndDate     *time.Time `json:"end_date,omitempty"`
}
