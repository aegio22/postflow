package models

type ProjectRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type AddUserRequest struct {
	ProjectName string `json:"project_name"`
	UserEmail   string `json:"user_email"`
	UserStatus  string    `json:"user_status"`
}

type UserInfo struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
