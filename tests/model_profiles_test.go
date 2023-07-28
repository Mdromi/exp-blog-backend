package tests

import (
	"fmt"
	"log"
	"testing"

	"github.com/Mdromi/exp-blog-backend/api/models"
	"github.com/stretchr/testify/assert"
)

func TestFindAllUsersProfile(t *testing.T) {
	err := refreshUserProfileTable()
	if err != nil {
		log.Fatal(err)
	}

	_, err = seedUsersProfiles()
	if err != nil {
		log.Fatal(err)
	}

	profile, err := profileInstance.FindAllUsersProfile(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the profiles: %v\n", err)
		return
	}

	// for _, p := range *profile {
	// 	fmt.Printf("Profile ID: %d\n", p.ID)
	// 	fmt.Printf("Name: %s\n", p.Name)
	// 	fmt.Printf("Title: %s\n", p.Title)
	// 	fmt.Printf("ProfilePic: %s\n", p.ProfilePic)
	// 	fmt.Printf("UserID: %v\n", p.UserID)
	// 	fmt.Printf("User: %v\n", p.User)
	// 	fmt.Printf("UserEmail: %v\n", p.User.Email)
	// 	fmt.Println("---------------------")
	// }
	assert.Equal(t, len(*profile), 2)
}

func TestSaveUserProfile(t *testing.T) {
	err := refreshUserProfileTable()
	if err != nil {
		log.Fatal(err)
	}

	profile, err := seedOneUserProfile()
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println("profile.UserID", profile.UserID)
	// fmt.Println("profile.User.ID", profile.User.ID)
	assert.Equal(t, profile.UserID, uint32(profile.User.ID))
	assert.Equal(t, uint32(profile.ID), uint32(profile.User.ProfileID))
	assert.Equal(t, profile.Name, profile.User.Username)
	assert.Equal(t, profile.ProfilePic, profile.User.AvatarPath)
}

func TestFindUserProfileByID(t *testing.T) {
	err := refreshUserProfileTable()
	if err != nil {
		log.Fatal(err)
	}

	profile, err := seedOneUserProfile()
	if err != nil {
		log.Fatalf("cannot seed profile table: %v", err)
	}

	foundProfile, err := profileInstance.FindUserProfileByID(server.DB, uint32(profile.ID))
	if err != nil {
		t.Errorf("this is the error getting one profile: %v\n", err)
		return
	}

	assert.Equal(t, foundProfile.ID, profile.ID)
	assert.Equal(t, uint32(profile.ID), uint32(profile.User.ProfileID))
	assert.Equal(t, profile.Name, profile.User.Username)
	assert.Equal(t, profile.Title, "Profile Title for "+profile.Name)
	assert.Equal(t, foundProfile.Bio, "Profile Bio for "+profile.Name)
}

func TestUpdateUserProfile(t *testing.T) {
	err := refreshUserProfileTable()
	if err != nil {
		log.Fatal(err)
	}

	profile, err := seedOneUserProfile()
	if err != nil {
		log.Fatalf("cannot seed profile table: %v", err)
	}

	updateProfile := models.Profile{
		Name:  profile.Name + " - 2",
		Title: profile.Title + " - 2",
		Bio:   profile.Bio + " - 2",
	}

	profileID := uint32(profile.ID)

	updatedProfile, err := updateProfile.UpdateAUserProfile(server.DB, profileID)
	if err != nil {
		t.Errorf("this is the error updating the profile: %v\n", err)
		return
	}

	// Refresh the profile to get the updated data from the database
	err = server.DB.Model(&models.Profile{}).Where("id = ?", profileID).Take(&profile).Error
	if err != nil {
		t.Errorf("error refreshing profile: %v\n", err)
		return
	}

	fmt.Println("updatedProfile.User.ProfileID", updatedProfile.User.ProfileID)
	fmt.Println("profile.UserID", profile.UserID)

	assert.Equal(t, uint32(updatedProfile.ID), profileID)
	assert.Equal(t, updatedProfile.UserID, profile.UserID)
	assert.Equal(t, updatedProfile.User.ProfileID, uint32(profile.ID))

	fmt.Println("updatedProfile.Name", updatedProfile.Name)
	assert.Equal(t, updatedProfile.Name, profile.Name+" - 2")
	// assert.Equal(t, updatedProfile.Title, profile.Title+" - 2")
	// assert.Equal(t, updatedProfile.Bio, profile.Bio+" - 2")

	// Update the user's ProfileID in the user model
	updatedUser := profile.User
	updatedUser.ProfileID = uint32(updatedProfile.ID)
	err = server.DB.Save(&updatedUser).Error
	if err != nil {
		t.Errorf("error updating user's ProfileID: %v\n", err)
		return
	}
}

// func TestDeleteUserProfile(t *testing.T)
