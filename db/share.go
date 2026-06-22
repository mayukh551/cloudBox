package db

import (
	"context"
	"fmt"
	"time"

	"github.com/mayukh551/cloudbox/models"
	"github.com/mayukh551/cloudbox/utils"
)

func CreateShare(data models.ShareUser, ctxt context.Context) error {

	var id string = utils.GenerateUUID()

	queryCtxt, cancel := context.WithTimeout(ctxt, 30*time.Second)
	defer cancel()

	_, err := DB.ExecContext(queryCtxt,
		`INSERT INTO shares (id, sharedTo, sharedBy, fileID)
		 VALUES ($1, $2, $3, $4)`,
		id, data.SharedTo, data.SharedBy, data.FileID,
	)

	if err != nil {
		return fmt.Errorf("error creating new share entity: %w", err)
	}

	return nil
}

func ListShares(userID string, path string, page int, limit int, ctxt context.Context) ([]models.FileList, error) {

	var shares []models.FileList

	queryCtxt, cancel := context.WithTimeout(ctxt, 30*time.Second)
	defer cancel()

	skip := (page - 1) * limit

	rows, err := DB.QueryContext(queryCtxt,
		`SELECT * FROM files WHERE id in (
			SELECT fileID FROM shares WHERE sharedTo = $1 AND path = $2
		) ORDER BY updatedAt DESC LIMIT $3 OFFSET $4`,
		userID, path, limit, skip,
	)

	if err != nil {
		return nil, ErrShareFetch
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

		shares = append(shares, file)
	}

	if err := rows.Err(); err != nil {
		return nil, ErrShareFetch
	}

	return shares, nil
}

func RemoveSharedFile(id string, ctxt context.Context) error {
	queryCtxt, cancel := context.WithTimeout(ctxt, 30*time.Second)
	defer cancel()

	result, err := DB.ExecContext(queryCtxt,
		`DELETE FROM shares WHERE id = $1`,
		id,
	)

	if err != nil {
		return ErrShareNotfound
	}

	rowsAffectd, err := result.RowsAffected()

	if err != nil {
		return ErrInternal
	}

	if rowsAffectd == 0 {
		return ErrShareNotfound
	}

	return nil

}
