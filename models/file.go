package models

import "time"

type File struct {
	Name   string `json:"name" validate:"required"`
	Url    string `json:"url" validate:"required,url"`
	UserID string `json:"user_id" validate:"required"`
}

type CreateFile struct {
	ID        string `json:"id"`
	Name      string `json:"name" validate:"required,min=3,max=100"`
	Type      string `json:"type" validate:"required,oneof=file folder"`
	Path      string `json:"path" validate:"required,min=3,max=1000"`
	Size      int    `json:"size" validate:"required,min=0,max=100000000"`
	UserID    string `json:"userID"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type FileList struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Path      string `json:"path"`
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
	Name      string `json:"title" validate:"required,min=3,max=100"`
	OldTitle  string `json:"oldTitle" validate:"required,min=3,max=100"`
	UpdatedAt string `json:"updatedAt"`
}

type DeleteFilePayload struct {
	Id  string `json:"id" validate:"required"`
	Key string `json:"key" validate:"required"`
}

type PreSignedBody struct {
	FileID      string `json:"fileID"`
	Filename    string `json:"filename" validate:"required,min=3,max=100"`
	Path        string `json:"path" validate:"required,min=3,max=1000"`
	ContentType string `json:"contentType"`
	Size        int    `json:"size" validate:"required,min=0,max=100000000"` // in bytes
}

type PreSignedResponse struct {
	Key string `json:"key"`
	Url string `json:"url"`
}

type MoveFilePayload struct {
	ID        string `json:"id" validate:"required"`
	Path      string `json:"path" validate:"required,min=0,max=1000"`
	UpdatedAt string `json:"updatedAt"`
}

type CreateFolderPayload struct {
	Name string `json:"folderName" validate:"required,min=3,max=100"`
	Path string `json:"path" validate:"required,min=0,max=1000"`
}

type FileEntity struct {
	ID        string `db:"id"`
	Name      string `db:"name"`
	Type      string `db:"type"`
	Size      string `db:"size"`
	Path      string `db:"path"`
	UserID    string `db:"userID"`
	CreatedAt string `db:"createdAt"`
	UpdatedAt string `db:"updatedAt"`
}

type StarFilePayload struct {
	Type string   `db:"type"`
	IDs  []string `db:"ids"`
}

// models/file.go
type StarredFile struct {
	ID        string    `json:"fileID"`
	Name      string    `json:"filename"`
	UpdatedAt time.Time `json:"lastModified"`
	Size      int64     `json:"size"`
	Type      string    `json:"contentType"`
}
