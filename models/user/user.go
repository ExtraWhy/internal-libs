package user

type User struct {
	Id      uint64 `json:"id"`
	Name    uint64 `json:"name"`
	Email   string `json:"email"`
	Picture string `json:"picture"`
}
