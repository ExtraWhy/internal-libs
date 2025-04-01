package user

type User struct {
	Id       uint64 `json:"id"`
	Token    string `json:"token"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Picture  string `json:"picture"`
}
