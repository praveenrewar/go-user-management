package usersmiddleware

import (
	"encoding/json"
	"net/http"

	jwtauthenticate "golang-mvc-boilerplate/server/middlewares/jwtAuthenticate"

	"github.com/gorilla/context"
	"github.com/thedevsaddam/govalidator"
)

//UserFormData is used to validate form data
type UserFormData struct {
	UserID   string `json:"user_id"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

//UpdateProfileFormData is used to validate form data for update_profile
type UpdateProfileFormData struct {
	UserID      string `json:"user_id,omitempty"`
	OldPassword string `json:"old_password"`
	Password    string `json:"new_password"`
}

//SignupValidator is used to validate form data for sign up api
func SignupValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rules := govalidator.MapData{
			"user_id":  []string{"required", "between:3,15"},
			"password": []string{"required", "min:8", "max:20"},
			"role":     []string{"required", "in:admin,datascientist"},
		}
		var user UserFormData
		type ValidationMessage struct {
			Status  int32  `json:"status"`
			Message string `json:"message"`
		}
		opts := govalidator.Options{
			Request: r,
			Data:    &user,
			Rules:   rules,
		}
		v := govalidator.New(opts)
		e := v.ValidateJSON()
		if len(e) != 0 {
			validationErr := map[string]interface{}{"validationError": e}
			validationError, _ := json.Marshal(validationErr)
			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(400)
			message := ValidationMessage{
				Status:  400,
				Message: string(validationError)}
			userJSON, _ := json.Marshal(message)
			w.Write(userJSON)
			return
		}
		context.Set(r, "user", user)
		next.ServeHTTP(w, r)
	})
}

//LoginValidator is used to validate form data for login apis
func LoginValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rules := govalidator.MapData{
			"user_id":  []string{"required"},
			"password": []string{"required"},
		}
		var user UserFormData
		type ValidationMessage struct {
			Status  int32  `json:"status"`
			Message string `json:"message"`
		}
		opts := govalidator.Options{
			Request: r,
			Data:    &user,
			Rules:   rules,
		}
		v := govalidator.New(opts)
		e := v.ValidateJSON()
		if len(e) != 0 {
			validationErr := map[string]interface{}{"validationError": e}
			validationError, _ := json.Marshal(validationErr)
			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(400)
			message := ValidationMessage{
				Status:  400,
				Message: string(validationError)}
			userJSON, _ := json.Marshal(message)
			w.Write(userJSON)
			return
		}
		context.Set(r, "user", user)
		next.ServeHTTP(w, r)
	})
}

//UpdatePasswordValidator is used to validate form data for update_profile
func UpdatePasswordValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		rules := govalidator.MapData{
			"old_password": []string{"required"},
			"new_password": []string{"required", "min:8", "max:20"},
		}
		var user UpdateProfileFormData
		type ValidationMessage struct {
			Status  int32  `json:"status"`
			Message string `json:"message"`
		}
		opts := govalidator.Options{
			Request: r,
			Data:    &user,
			Rules:   rules,
		}
		v := govalidator.New(opts)
		e := v.ValidateJSON()
		if len(e) != 0 {
			validationErr := map[string]interface{}{"validationError": e}
			validationError, _ := json.Marshal(validationErr)
			w.WriteHeader(400)
			message := ValidationMessage{
				Status:  400,
				Message: string(validationError)}
			userJSON, _ := json.Marshal(message)
			w.Write(userJSON)
			return
		}
		if user.OldPassword == user.Password {
			w.WriteHeader(400)
			message := ValidationMessage{
				Status:  400,
				Message: "New password cannot be same as old password"}
			userJSON, _ := json.Marshal(message)
			w.Write(userJSON)
			return
		}
		jwtData := context.Get(r, "user_jwt")
		jwtUser := jwtData.(jwtauthenticate.UserJWTData)
		userID := jwtUser.UserID
		user.UserID = userID
		context.Set(r, "user", user)
		next.ServeHTTP(w, r)
	})
}
