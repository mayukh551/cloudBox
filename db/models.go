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
	Url       string `db:"url"`
	UserID    string `db:"user_id"`
	CreatedAt string `db:"createdAt"`
	UpdatedAt string `db:"updatedAt"`
}
