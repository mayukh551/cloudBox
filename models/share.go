package models

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
