package models

type CreateUser struct {
	ID       string `json:"id" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type UpdateUser struct {
	ID       string `json:"id" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type ResponseUser struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type User struct {
	ID    string `json:"id" validate:"required"`
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

type File struct {
	Name   string `json:"name" validate:"required"`
	Url    string `json:"url" validate:"required,url"`
	UserID string `json:"user_id" validate:"required"`
}

type FileShare struct {
	Email  string `json:"email" validate:"required,email"`
	FileID string `json:"fileID" validate:"required"`
}

type RequestUser struct {
	ID    string `json:"id" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

type ShareUser struct {
	SharedTo string `json:"sharedTo" validate:"required"`
	SharedBy string `json:"sharedBy"`
	FileID   string `json:"fileID" validate:"required"`
}

type ShareList struct {
	ID         string `json:"id" validate:"required"`
	SharedTo   string `json:"sharedTo" validate:"required,email"`
	SharedBy   string `json:"sharedBy" validate:"required"`
	FileID     string `json:"fileID" validate:"required"`
	ModifiedAt string `json:"modifiedAt"`
}
