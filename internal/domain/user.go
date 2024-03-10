package domain

import "time"

type User struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	HashPass    string    `json:"hash_pass"`
	PhoneNumber string    `json:"phone_number"`
	BDay        time.Time `json:"b_day"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UserProfileUpdateReq struct {
	ID          uint
	Name        *string    `json:"name"`
	Email       *string    `json:"email"`
	PhoneNumber *string    `json:"phone_number"`
	BDay        *time.Time `json:"b_day" format:"RFC3339"`
}
