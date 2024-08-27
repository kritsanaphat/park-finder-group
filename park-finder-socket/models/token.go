package models

import "time"

type Token struct {
	UserID     string    `bson:"user_id"`
	Token      string    `bson:"token"`
	Role       string    `bson:"role"`
	Valid      bool      `bson:"valid"`
	RevokeDate time.Time `bson:"revoke_date"`
}
