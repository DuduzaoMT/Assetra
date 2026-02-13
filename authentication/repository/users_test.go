package repository

import (
	"assetra/authentication/models"
	"assetra/db"
	"log"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func init() {
	// if os.Getenv("ENV") != "dev" {
	// 	TODO: make K8s especial deployment
	// }
	err := godotenv.Load("../../.env") // load .env file
	if err != nil {
		log.Panic("Failed to load .env file")
	}
}

func TestUserRepositoryLifeCycle(t *testing.T) {
	// Initialize database connection
	cfg := db.NewConfig()
	conn := db.NewConnection(cfg)
	defer conn.Close()

	userRepository := NewUserRepository(conn)

	// Create a new user
	user := &models.User{
		Username:  "testuser",
		Email:     "testuser@example.com",
		Password:  "hashedpassword",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Roles:     []string{"user"},
	}

	// Save the user
	_, err := userRepository.Save(user)
	assert.NoError(t, err, "Failed to save user")

	// assert all the user fields (should be equal to the original user)
	fetchedUser, err := userRepository.FindByID(user.ID.String())
	assert.NoError(t, err, "Failed to find user by ID")
	assert.NotNil(t, fetchedUser, "Fetched user should not be nil")
	assert.Equal(t, user.Username, fetchedUser.Username)
	assert.Equal(t, user.Email, fetchedUser.Email)
	assert.Equal(t, user.Password, fetchedUser.Password)
	assert.Equal(t, user.Roles, fetchedUser.Roles)
	assert.Equal(t, user.CreatedAt.Unix(), fetchedUser.CreatedAt.Unix())
	assert.Equal(t, user.UpdatedAt.Unix(), fetchedUser.UpdatedAt.Unix())

	// Update some fields of the user
	user.Username = "updateduser"
	user.UpdatedAt = time.Now()
	err = userRepository.Update(user)
	assert.NoError(t, err, "Failed to update user")

	// Fetch the updated user and verify the updates
	updatedUser, err := userRepository.FindByEmail(user.Email)
	assert.NoError(t, err, "Failed to find user by email")
	assert.NotNil(t, updatedUser, "Updated user should not be nil")
	assert.Equal(t, user.Username, updatedUser.Username)
	assert.Equal(t, user.UpdatedAt.Unix(), updatedUser.UpdatedAt.Unix())

	// Fetch all users and verify at least one user exists
	allUsers, err := userRepository.FindAll()
	assert.NoError(t, err, "Failed to find all users")
	assert.GreaterOrEqual(t, len(allUsers), 1, "There should be at least one user")

	// Delete the user
	err = userRepository.Delete(user.ID.String())
	assert.NoError(t, err, "Failed to delete user")

	// Verify the user has been deleted
	deletedUser, err := userRepository.FindByID(user.ID.String())
	assert.Error(t, err, "Expected error when finding deleted user")
	assert.Nil(t, deletedUser, "Deleted user should be nil")
}
