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

func ListShares(userID string, ctxt context.Context) []models.ShareList {

	var shares []models.ShareList

	queryCtxt, cancel := context.WithTimeout(ctxt, 30*time.Second)
	defer cancel()

	rows, err := DB.QueryContext(queryCtxt,
		`SELECT id, sharedTo, sharedBy, fileID, updatedAt FROM shares WHERE sharedBy = $1`,
		userID,
	)

	if err != nil {
		return nil
	}

	defer rows.Close()

	for rows.Next() {
		var share models.ShareList
		if err := rows.Scan(&share.ID, &share.SharedTo, &share.SharedBy, &share.FileID, &share.ModifiedAt); err != nil {
			return nil
		}
		shares = append(shares, share)
	}

	return shares
}

func ListSharedWithMe(userID string, ctxt context.Context) []models.ShareList {

	var shares []models.ShareList

	queryCtxt, cancel := context.WithTimeout(ctxt, 30*time.Second)
	defer cancel()

	rows, err := DB.QueryContext(queryCtxt,
		`SELECT id, sharedTo, sharedBy, fileID, updatedAt FROM shares WHERE sharedTo = $1`,
		userID,
	)

	if err != nil {
		return nil
	}

	defer rows.Close()

	for rows.Next() {
		var share models.ShareList
		if err := rows.Scan(&share.ID, &share.SharedTo, &share.SharedBy, &share.FileID, &share.ModifiedAt); err != nil {
			return nil
		}
		shares = append(shares, share)
	}

	return shares
}
