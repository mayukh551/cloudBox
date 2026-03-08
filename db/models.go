package db

type userEntity struct {
	ID        string `db:"id"`
	Name      string `db:"name"`
	Email     string `db:"email"`
	Password  string `db:"password"`
	CreatedAt string `db:"createdAt"`
	UpdatedAt string `db:"updatedAt"`
}

type fileEntity struct {
	ID        string `db:"id"`
	Name      string `db:"name"`
	Type      string `db:"type"`
	Size      string `db:"size"`
	Url       string `db:"url"`
	UserID    string `db:"userID"`
	CreatedAt string `db:"createdAt"`
	UpdatedAt string `db:"updatedAt"`
}

type shareEntity struct {
	ID        string `db:"id"`
	SharedTo  string `db:"sharedTo"`
	SharedBy  string `db:"sharedBy"`
	fileID    string `db:"fileID"`
	CreatedAt string `db:"createdAt"`
	UpdatedAt string `db:"updatedAt"`
}
