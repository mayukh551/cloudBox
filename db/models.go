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
	Path      string `db:"path"`
	UserID    string `db:"userID"`
	CreatedAt string `db:"createdAt"`
	UpdatedAt string `db:"updatedAt"`
}

type shareEntity struct {
	ID        string `db:"id"`
	SharedTo  string `db:"sharedTo"`
	SharedBy  string `db:"sharedBy"`
	FileID    string `db:"fileID"`
	CreatedAt string `db:"createdAt"`
	UpdatedAt string `db:"updatedAt"`
}

// why create a separate entity for marking stars to files or folders?
// - A file owned by the user or shared to another user can be starred by either user
// - if we create a separate column on fileEntity for star, then how can we track which is file is starred by which user,
// - since in fileEntity a file is already tied to a single user
// - so fileID and userID mapping must be unique for each row on starEntity
type starFileEntity struct {
	FileID    string `db:"fileID"`
	UserID    string `db:"userID"`
	CreatedAt string `db:"createdAt"`
}
