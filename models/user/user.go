package user

type User struct {
	Id       uint64 `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Picture  string `json:"picture"`
}
