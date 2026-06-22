package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/mayukh551/cloudbox/models"
)

func CreateFile(data models.CreateFile, ctxt context.Context) error {

	queryCtxt, cancel := context.WithTimeout(ctxt, 30*time.Second)
	defer cancel()

	_, err := DB.ExecContext(queryCtxt,
		`INSERT INTO files (id, name, path, type, size, userID)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		data.ID,
		data.Name,
		data.Path,
		data.Type,
		data.Size,
		data.UserID,
	)

	if err != nil {
		return ErrFileNotCreated
	}

	return nil
}

func UpdatePath(data models.MoveFilePayload, ctxt context.Context) error {
	queryCtxt, cancel := context.WithTimeout(ctxt, 30*time.Second)
	defer cancel()

	result, err := DB.ExecContext(queryCtxt,
		`UPDATE files
		 SET path = $2, updatedAt = $3
		 WHERE id = $1`,
		data.ID,
		data.Path,
		data.UpdatedAt,
	)

	if err != nil {
		return ErrFilePathNotUpdated
	}

	rowsAffectd, err := result.RowsAffected()

	if err != nil {
		return ErrFilePathNotUpdated
	}

	if rowsAffectd == 0 {
		return ErrFileIsNotFound
	}

	return nil
}

func ListFiles(userID string, searchTerm string, path string, page int, limit int, ctxt context.Context) ([]models.FileList, error) {
	var files []models.FileList

	queryCtxt, cancel := context.WithTimeout(ctxt, 30*time.Second)
	defer cancel()

	var rows *sql.Rows
	var err error

	skip := (page - 1) * limit

	if searchTerm != "" {
		rows, err = DB.QueryContext(queryCtxt,
			`SELECT id, name, type, path, size, userID, createdAt, updatedAt FROM files WHERE userID = $1 AND name ~* $2 ORDER BY updatedAt DESC LIMIT $3`,
			userID, searchTerm, limit,
		)
	} else {
		rows, err = DB.QueryContext(queryCtxt,
			`SELECT id, name, type, path, size, userID, createdAt, updatedAt FROM files WHERE userID = $1 AND path = $2 ORDER BY updatedAt DESC LIMIT $3 OFFSET $4`,
			userID, path, limit, skip,
		)
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var file models.FileList
		if err := rows.Scan(
			&file.ID,
			&file.Name,
			&file.Type,
			&file.Path,
			&file.Size,
			&file.UserID,
			&file.CreatedAt,
			&file.UpdatedAt,
		); err != nil {
			return nil, ErrFileFetch
		}

		files = append(files, file)
	}

	if err := rows.Err(); err != nil {
		return nil, ErrFileFetch
	}

	return files, nil
}

func CheckIfFileNameExists(filename string, path string, ctxt context.Context) (map[string]bool, error) {
	queryCtxt, cancel := context.WithTimeout(ctxt, 10*time.Second)
	defer cancel()

	var samePathExists, diffPathExists, onlyPathExists bool

	err := DB.QueryRowContext(queryCtxt,
		`SELECT
			EXISTS(SELECT 1 FROM files WHERE name = $1 AND path = $2)  AS same_path_exists,
			EXISTS(SELECT 1 FROM files WHERE name = $1 AND path != $2) AS diff_path_exists,
			EXISTS(SELECT 1 FROM files WHERE path = $2) AS onlyPathExists
		`,
		filename, path,
	).Scan(&samePathExists, &diffPathExists, &onlyPathExists)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrFileOrPathNotFound
		}
		return nil, ErrInternal
	}

	return map[string]bool{
		"samePathExists": samePathExists,
		"diffPathExists": diffPathExists,
		"onlyPathExists": onlyPathExists,
	}, nil
}

func GetFileByID(fileID string, ctxt context.Context) (models.FileEntity, error) {
	var file models.FileEntity

	queryCtxt, cancel := context.WithTimeout(ctxt, 30*time.Second)
	defer cancel()

	err := DB.QueryRowContext(queryCtxt,
		`SELECT id, name, type, size, userID, createdAt, updatedAt
		 FROM files
		 WHERE id = $1`,
		fileID,
	).Scan(
		&file.ID,
		&file.Name,
		&file.Type,
		&file.Size,
		&file.UserID,
		&file.CreatedAt,
		&file.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return file, ErrFileIsNotFound
		}
		return file, ErrInternal
	}

	return file, nil
}

func UpdateFile(data models.CreateFile, ctxt context.Context) error {
	queryCtxt, cancel := context.WithTimeout(ctxt, 30*time.Second)
	defer cancel()

	_, err := DB.ExecContext(queryCtxt,
		`UPDATE files
		 SET name = $2, type = $3, size = $4, updatedAt = $5, path = $6
		 WHERE id = $1`,
		data.ID,
		data.Name,
		data.Type,
		data.Size,
		data.UpdatedAt,
		data.Path,
	)

	if err != nil {
		fmt.Printf("error updating file entity: %w", err)
		return ErrFileNotUpdated
	}

	return nil
}

func UpdateFileName(data models.UpdateFileNamePayload, ctxt context.Context) error {
	queryCtxt, cancel := context.WithTimeout(ctxt, 30*time.Second)
	defer cancel()

	_, err := DB.ExecContext(queryCtxt,
		`UPDATE files
		 SET name = $2, updatedAt = $3
		 WHERE id = $1`,
		data.Id,
		data.Name,
		data.UpdatedAt,
	)

	if err != nil {
		fmt.Printf("error updating file entity: %w", err)
		return ErrFileNameNotUpdated
	}

	return nil
}

func DeleteFile(fileID string, ctxt context.Context) error {
	queryCtxt, cancel := context.WithTimeout(ctxt, 30*time.Second)
	defer cancel()

	_, err := DB.ExecContext(queryCtxt,
		`DELETE FROM files WHERE id = $1`,
		fileID,
	)

	if err != nil {
		fmt.Printf("error deleting file entity: %w\n", err)
		return ErrDeleteFiles
	}

	return nil
}

func DeleteFolder(folderName string, id string, userID string, ctxt context.Context) error {

	queryCtxt, cancel := context.WithTimeout(ctxt, 30*time.Second)
	defer cancel()

	tx, err := DB.BeginTx(queryCtxt, nil)

	if err != nil {
		fmt.Println("error beginning transaction", err)
		return ErrInternal
	}

	defer tx.Rollback()

	_, err = tx.ExecContext(queryCtxt,
		`DELETE FROM files WHERE path ~ ('(^|/)' || $1 || '(/|$)') AND userID = $2`,
		folderName, userID,
	)

	if err != nil {
		fmt.Println("error deleting files that are stored inside the folder", err)
		return ErrDeleteFilesInsideFolder
	}

	_, err = tx.ExecContext(queryCtxt,
		`DELETE FROM files WHERE id = $1`,
		id,
	)

	if err != nil {
		fmt.Println("error deleting the selected folder", err)
		return ErrDeleteFolder
	}

	if err := tx.Commit(); err != nil {
		fmt.Println("commit folder delete", err)
		return ErrInternal
	}

	return nil
}

// ----------------------------------------------------

// Favourites API

func AddStar(fileID string, userID string, ctxt context.Context) error {

	queryCtxt, cancel := context.WithTimeout(ctxt, 30*time.Second)
	defer cancel()

	fmt.Println(fileID, userID)

	_, err := DB.ExecContext(queryCtxt, `
		INSERT INTO starFiles (fileID, userID)
		VALUES ($1, $2)`,
		fileID,
		userID,
	)

	if err != nil {
		fmt.Println("error inserting new starred file, %w", err)
		return ErrInternal
	}

	return nil

}

func RemoveStar(fileID string, userID string, ctxt context.Context) error {
	queryCtxt, cancel := context.WithTimeout(ctxt, 30*time.Second)
	defer cancel()

	_, err := DB.ExecContext(queryCtxt, `
		DELETE FROM starFiles WHERE fileID = $1 AND userID = $2`,
		fileID,
		userID,
	)

	if err != nil {
		fmt.Println("error inserting new starred file, %w", err)
		return ErrInternal
	}

	return nil
}

func GetStarredFiles(userID string, path string, page int, limit int, ctxt context.Context) ([]models.FileList, error) {

	queryCtxt, cancel := context.WithTimeout(ctxt, 30*time.Second)
	defer cancel()

	skip := (page - 1) * limit

	rows, err := DB.QueryContext(queryCtxt, `
		SELECT f.id, f.name, f.type, f.path, f.size, f.userID, f.createdAt, f.updatedAt
		FROM starFiles sf
		JOIN files f ON sf.fileID = f.id
		WHERE sf.userID = $1 AND f.path = $2
		ORDER BY f.updatedAt DESC
		LIMIT $3 OFFSET $4`,
		userID, path, limit, skip,
	)

	if err != nil {
		return nil, ErrInternal
	}

	defer rows.Close()

	var files []models.FileList

	for rows.Next() {
		var file models.FileList
		if err := rows.Scan(
			&file.ID,
			&file.Name,
			&file.Type,
			&file.Path,
			&file.Size,
			&file.UserID,
			&file.CreatedAt,
			&file.UpdatedAt,
		); err != nil {
			return nil, ErrFileFetch
		}

		files = append(files, file)
	}

	if err := rows.Err(); err != nil {
		return nil, ErrFileFetch
	}

	return files, nil
}

func SearchFilesByName(userID string, pattern string, page int, limit int, ctxt context.Context) ([]models.FileList, error) {
	queryCtxt, cancel := context.WithTimeout(ctxt, 30*time.Second)
	defer cancel()

	skip := (page - 1) * limit
	fmt.Println("Pattern", pattern, userID)

	rows, err := DB.QueryContext(queryCtxt, `
        SELECT id, name, type, path, size, userID, createdAt, updatedAt
        FROM files
        WHERE name ~* $1 AND userID = $2
        ORDER BY updatedAt DESC LIMIT $3 OFFSET $4`,
		pattern, userID, limit, skip,
	)

	fmt.Println(sql.ErrNoRows)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var files []models.FileList

	for rows.Next() {
		fmt.Println("yo")
		var file models.FileList
		if err := rows.Scan(
			&file.ID,
			&file.Name,
			&file.Type,
			&file.Path,
			&file.Size,
			&file.UserID,
			&file.CreatedAt,
			&file.UpdatedAt,
		); err != nil {
			return nil, ErrFileFetch
		}
		fmt.Println(file.Name)
		files = append(files, file)
	}

	if err := rows.Err(); err != nil {
		return nil, ErrFileFetch
	}

	return files, nil
}
