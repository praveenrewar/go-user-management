package users

//User represents a music album
type User struct {
	UserID   string `json:"user_id"`
	Password string `json:"password"`
}

//Users is an array of User
type Users []User
