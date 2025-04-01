package user

type User struct {
	Id       uint64 `json:"id"`
	Username uint64 `json:"username"`
	Email    string `json:"email"`
	Picture  string `json:"picture"`
}
