package db

import (
	"context"
	"fmt"
	"time"

	"github.com/mayukh551/cloudbox/models"
)

func CreateFile(data models.CreateFile, ctxt context.Context) error {

	queryCtxt, cancel := context.WithTimeout(ctxt, 30*time.Second)
	defer cancel()

	_, err := DB.ExecContext(queryCtxt,
		`INSERT INTO files (id, name, type, size, userID)
		 VALUES ($1, $2, $3, $4, $5)`,
		data.ID,
		data.Title,
		data.Type,
		data.Size,
		data.UserID,
	)

	if err != nil {
		return fmt.Errorf("error creating file entity: %w", err)
	}

	return nil
}

func ListFiles(userID string, ctxt context.Context) ([]models.FileList, error) {
	var files []models.FileList

	queryCtxt, cancel := context.WithTimeout(ctxt, 30*time.Second)
	defer cancel()

	rows, err := DB.QueryContext(queryCtxt,
		`SELECT id, name, type, size, userID, createdAt, updatedAt
		 FROM files
		 WHERE userID = $1
		 ORDER BY updatedAt DESC`,
		userID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var file models.FileList
		if err := rows.Scan(
			&file.ID,
			&file.Title,
			&file.Type,
			&file.Size,
			&file.UserID,
			&file.CreatedAt,
			&file.UpdatedAt,
		); err != nil {
			return nil, err
		}

		files = append(files, file)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return files, nil
}

func GetFileByID(fileID string, ctxt context.Context) fileEntity {
	var file fileEntity

	queryCtxt, cancel := context.WithTimeout(ctxt, 30*time.Second)
	defer cancel()

	err := DB.QueryRowContext(queryCtxt,
		`SELECT id, name, type, size, userID, createdAt, updatedAt
		 FROM files
		 WHERE id = $1`,
		fileID,
	).Scan(
		&file.ID,
		&file.Title,
		&file.Type,
		&file.Size,
		&file.UserID,
		&file.CreatedAt,
		&file.UpdatedAt,
	)

	if err != nil {
		return file
	}

	return file
}

func UpdateFile(fileID string, data models.CreateFile, ctxt context.Context) error {
	queryCtxt, cancel := context.WithTimeout(ctxt, 30*time.Second)
	defer cancel()

	_, err := DB.ExecContext(queryCtxt,
		`UPDATE files
		 SET name = $2, type = $3, size = $4, updatedAt = $5
		 WHERE id = $1`,
		data.ID,
		data.Title,
		data.Type,
		data.Size,
		data.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("error updating file entity: %w", err)
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
		data.Title,
		data.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("error updating file entity: %w", err)
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
		return fmt.Errorf("error deleting file entity: %w", err)
	}

	return nil
}
