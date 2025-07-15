package models

import "time"

type SubscriptionFilter struct {
	UserID      string    `json:"user_id" binding:"required"`
	ServiceName string    `json:"service_name" binding:"required"`
	From        time.Time `json:"from" binding:"required"`
	To          time.Time `json:"to" binding:"required"`
}
