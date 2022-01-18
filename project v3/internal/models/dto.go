package models

import "time"

type LogInDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type Session struct {
	RefreshToken string    `json:"refreshToken" db:"refreshToken"`
	ExpiresAt    time.Time `json:"expiresAt" db:"expiresAt"`
}

type OrderStatusDTO struct {
	Id     int    `json:"id" db:"id"`
	Status Status `json:"status" db:"status"`
}

type AuthorizedUserInfo struct {
	Id   int  `json:"id"`
	Role Role `json:"role"`
}
