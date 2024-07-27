package users

type User struct {
	ID                 int    `json:"id"`
	Email              string `json:"email"`
	Password           string `json:"password"`
	Expires_in_Seconds *int   `json:"expires_in_seconds"`
}

type ReturnUserVal struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}
