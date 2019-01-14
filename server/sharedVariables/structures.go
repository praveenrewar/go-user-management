package shared

//User represents a user
type User struct {
	UserID   string `json:"user_id"`
	Password string `json:"password,omitempty"`
	Role     string `json:"role"`
}

//Users is an array of User
type Users []User
