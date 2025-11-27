package model

import "time"

type Details struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Phone     string    `json:"phoneNumber"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
}
