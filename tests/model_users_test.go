package tests

import (
	"log"
	"testing"

	"github.com/Mdromi/exp-blog-backend/api/models"
	_ "github.com/jinzhu/gorm/dialects/mysql"    //mysql driver
	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres driver
	"github.com/stretchr/testify/assert"
)

func TestFindAllUsers(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	_, err = seedUsers()
	if err != nil {
		log.Fatal(err)
	}

	users, err := userInstance.FindAllUsers(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the users: %v\n", err)
		return
	}
	assert.Equal(t, len(*users), 2)
}

func TestSaveUser(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	newUser := models.User{
		Email:    "test@example.com",
		Username: "test",
		Password: "password",
	}

	savedUser, err := newUser.SaveUser(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the users: %v\n", err)
		return
	}

	assert.NotEqual(t, 0, savedUser.ID) // Check that the ID is not zero
	assert.Equal(t, newUser.Email, savedUser.Email)
	assert.Equal(t, newUser.Username, savedUser.Username)
	assert.Equal(t, newUser.ProfileID, uint32(0))
}

func TestFindUserByID(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("cannot seed users table: %v", err)
	}
	foundUser, err := userInstance.FindUserByID(server.DB, uint32(user.ID))
	if err != nil {
		t.Errorf("this is the error getting one user: %v\n", err)
		return
	}
	assert.Equal(t, foundUser.ID, user.ID)
	assert.Equal(t, foundUser.Email, user.Email)
	assert.Equal(t, foundUser.Username, user.Username)
}

func TestUpdateAuser(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Cannot seed user: %v\n", err)
	}

	userUpdate := models.User{
		Username: "modiUpdate",
		Email:    "modiupdate@example.com",
		Password: "password",
	}

	// Create a separate variable for ID
	userID := uint32(user.ID)

	updatedUser, err := userUpdate.UpdateAUser(server.DB, userID)
	if err != nil {
		t.Errorf("this is the error updating the user: %v\n", err)
		return
	}

	assert.Equal(t, uint32(updatedUser.ID), userID)
	assert.Equal(t, updatedUser.Email, userUpdate.Email)
	assert.Equal(t, updatedUser.Username, userUpdate.Username)
}

func TestDeleteAUser(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatalf("Cannot seed user: %v\n", err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Cannot seed user: %v\n", err)
	}
	isDeleted, err := userInstance.DeleteAUser(server.DB, uint32(user.ID))
	if err != nil {
		t.Errorf("this is the error updating the user: %v\n", err)
		return
	}
	assert.Equal(t, isDeleted, int64(1))
}
