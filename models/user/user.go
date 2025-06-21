package user

import (
	"github.com/google/uuid"
)

type User struct {
	Id       string `json:"id"`
	Token    string `json:"token"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Picture  string `json:"picture"`
}

// used to give id to player db
func GenUUID() uuid.UUID {
	id := uuid.New()
	return id
}
