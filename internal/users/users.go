package users

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ReturnUserVal struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}