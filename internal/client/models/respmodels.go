package models

import (
	"time"

	"github.com/google/uuid"
)

type ProjectResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedBy   uuid.UUID `json:"created_by"`
}

type SignUpResponse struct {
	UserID  string `json:"id"`
	Token   string `json:"access_token"`
	Message string `json:"message"`
}

type LoginResponse struct {
	AccessToken string `json:"token"`
}

type DBUserResponse struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	AccessToken string    `json:"access_token"`
}

type ProjectMemberAddResponse struct {
	ProjectName string `json:"project_name"`
	UserStatus  string `json:"user_status"`
}

type AssetResponse struct {
	AssetID   string `json:"asset_id"`
	UploadURL string `json:"upload_url"`
	S3Key     string `json:"s3_key"`
	ExpiresIn int    `json:"expires_in"`
}

type ProjectsLsResponse struct {
	UserName string `json:"user_name"`
	Projects map[string]string `json:"projects_map"`
}
