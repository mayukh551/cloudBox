package models

type File struct {
	Name   string `json:"name" validate:"required"`
	Url    string `json:"url" validate:"required,url"`
	UserID string `json:"user_id" validate:"required"`
}

type CreateFile struct {
	ID        string `json:"id"`
	Title     string `json:"name"`
	Type      string `json:"type"`
	Size      int    `json:"size"`
	UserID    string `json:"userID"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type FileList struct {
	ID        string `json:"id"`
	Title     string `json:"name"`
	Type      string `json:"type"`
	Size      int    `json:"size"`
	UserID    string `json:"userID"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type FileShare struct {
	Email  string `json:"email" validate:"required,email"`
	FileID string `json:"fileID" validate:"required"`
}

type UpdateFileNamePayload struct {
	Id        string `json:"id" validate:"required"`
	FileID    string `json:"fileID" validate:"required"`
	Title     string `json:"title" validate:"required"`
	OldTitle  string `json:"oldTitle" validate:"required"`
	UpdatedAt string `json:"updatedAt"`
}

type DeleteFilePayload struct {
	Id  string `json:"id" validate:"required"`
	Key string `json:"key" validate:"required"`
}

type PreSignedBody struct {
	Filename    string `json:"filename"`
	ContentType string `json:"contentType"`
	Size        int    `json:"size"`
}

type PreSignedResponse struct {
	Key string `json:"key"`
	Url string `json:"url"`
}
