package users

//User represents a user
type User struct {
	UserID   string `json:"user_id"`
	Password string `json:"password,omitempty"`
	Role     string `json:"role"`
}

//Message is used to pass messages among model, view and controller
type Message struct {
	UserID  string `json:"user_id,omitempty"`
	Status  int32  `json:"status"`
	Message string `json:"message,omitempty"`
	User    *User  `json:"user,omitempty"`
	Users   *Users `json:"users,omitempty"`
	JWT     string `json:"JWT,omitempty"`
}

//UpdateUser is used to represent a user while updating
type UpdateUser struct {
	UserID      string `json:"user_id"`
	OldPassword string `json:"old_password"`
	Password    string `json:"new_password"`
}

//Users is an array of User
type Users []User
