package db

import (
	"fmt"

	"github.com/mayukh551/cloudbox/models"
	"github.com/mayukh551/cloudbox/utils"
)

// CreateUser creates a new user and returns the created user with ID
func CreateUser(data models.CreateUser) (*models.User, error) {
	var user userEntity

	err := DB.QueryRow(
		`INSERT INTO users (id, name, email, password)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, name, email, password`,
		data.ID, data.Name, data.Email, data.Password,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Password)

	if err != nil {
		return nil, fmt.Errorf("error creating new user: %w", err)
	}

	return &models.User{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

// UpdateUser updates an existing user's information
func UpdateUser(id string, data models.UpdateUser) error {
	result, err := DB.Exec(
		`UPDATE users 
		 SET name = $1, email = $2, password = $3
		 WHERE id = $4`,
		data.Name, data.Email, data.Password, id,
	)

	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with id %s not found", id)
	}

	return nil
}

// DeleteUser deletes a user by ID
func DeleteUser(id string) error {
	result, err := DB.Exec(
		`DELETE FROM users WHERE id = $1`,
		id,
	)

	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with id %s not found", id)
	}

	return nil
}

// GetUserByID retrieves a user by ID (bonus function)
func GetUserByID(id string) (*models.User, error) {
	var user userEntity

	err := DB.QueryRow(
		`SELECT id, name, email FROM users WHERE id = $1`,
		id,
	).Scan(&user.ID, &user.Name, &user.Email)

	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	return &models.User{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

// GetUserByEmail retrieves a user by email (bonus function for auth)
func GetUserByEmail(email string) (*models.User, error) {
	var user userEntity

	err := DB.QueryRow(
		`SELECT id, name, email FROM users WHERE email = $1`,
		email,
	).Scan(&user.ID, &user.Name, &user.Email)

	if err != nil {
		return nil, fmt.Errorf("error getting user by email: %w", err)
	}

	return &models.User{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func VerifyUser(email string, password string) (models.User, error) {

	var user userEntity

	fmt.Println("email", email)

	err := DB.QueryRow(
		`SELECT id, name, email, password FROM users WHERE email = $1`,
		email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Password)

	if err != nil {
		fmt.Println(err)
		return models.User{}, fmt.Errorf("error getting user by email: %w", err)
	}

	isValid := utils.ValidatePassword(password, user.Password)

	if !isValid {
		return models.User{}, fmt.Errorf("invalid password")
	}

	return models.User{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil

}
