package models

import (
	"errors"
	"fmt"
	"html"
	"strings"

	"gorm.io/gorm"
)

// Profile model represents user's profile details
type Profile struct {
	gorm.Model
	UserID      uint        `gorm:"not null" json:"user_id"`
	Name        string      `gorm:"type:varchar(50);not null" json:"name" validate:"min=2,max=50"`
	Title       string      `gorm:"type:varchar(100);not null" json:"title" validate:"max=100"`
	Bio         string      `gorm:"type:text;not null" json:"bio" validate:"max=500"`
	ProfilePic  string      `gorm:"type:varchar(255)" json:"profile_pic"`
	SocialLinks *SocialLink `json:"social_links"`
	Username    string      `gorm:"type:varchar(50)" json:"username"`
	CoverPic    string      `gorm:"type:varchar(255)" json:"cover_pic"`
}

func (p *Profile) Prepare() {
	// Sanitize and trim strings
	p.Name = html.EscapeString(strings.TrimSpace(p.Name))
	p.Title = html.EscapeString(strings.TrimSpace(p.Title))
	p.Bio = html.EscapeString(strings.TrimSpace(p.Bio))

	// Ensure SocialLink is initialized if nil
	if p.SocialLinks == nil {
		p.SocialLinks = &SocialLink{}
	} else {
		// Sanitize and trim social links only if SocialLinks is not nil
		p.SocialLinks.Website = html.EscapeString(strings.TrimSpace(p.SocialLinks.Website))
		p.SocialLinks.Github = html.EscapeString(strings.TrimSpace(p.SocialLinks.Github))
		p.SocialLinks.Linkedin = html.EscapeString(strings.TrimSpace(p.SocialLinks.Linkedin))
		p.SocialLinks.Twitter = html.EscapeString(strings.TrimSpace(p.SocialLinks.Twitter))
	}
}

func (p *Profile) AfterFind() (err error) {
	if err != nil {
		return err
	}

	// userAvatarPath := p.User.AvatarPath

	// if userAvatarPath != "" {
	// 	p.ProfilePic = userAvatarPath
	// }

	return nil
}

func (p *Profile) Validate(action string) map[string]string {
	errorMessages := make(map[string]string)

	switch strings.ToLower(action) {
	case "create", "update", "":
		if p.Name == "" {
			errorMessages["Required_name"] = "Name is required"
		} else if len(p.Name) < 2 || len(p.Name) > 50 {
			errorMessages["Name_length"] = "Name must be between 2 and 50 characters"
		}

		if p.Title == "" {
			errorMessages["Required_title"] = "Title is required"
		} else if len(p.Title) > 100 {
			errorMessages["Title_length"] = "Title cannot exceed 100 characters"
		}

		if p.Bio == "" {
			errorMessages["Required_bio"] = "Bio is required"
		} else if len(p.Bio) > 500 {
			errorMessages["Bio_length"] = "Bio cannot exceed 500 characters"
		}

		if p.UserID == 0 {
			errorMessages["user_id"] = "Invalid User ID"
		}

	case "delete":
		// No additional validation needed for delete action
	default:
		errorMessages["action"] = "Invalid action specified"
	}

	return errorMessages
}
func (p *Profile) SaveUserProfile(db *gorm.DB) (*Profile, error) {
	var err error

	// Check if the profile already exists
	if p.ID != 0 {
		return nil, errors.New("user already has a profile")
	}

	// Create the profile
	err = db.Debug().Model(&Profile{}).Create(&p).Error
	if err != nil {
		return nil, err
	}

	// Update the User.ProfileID with the newly created profile's ID
	if p.UserID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.UserID).Update("profile_id", p.ID).Error
		if err != nil {
			return nil, err
		}
	}

	return p, nil
}

// THE ONLY PERSON THAT NEED TO DO THIS IS THE ADMIN, SO I HAVE COMMENTED THE ROUTES, SO SOMEONE ELSE DONT VIES THIS DEATAILS
func (p *Profile) FindAllUsersProfile(db *gorm.DB) (*[]Profile, error) {
	var err error
	profiles := []Profile{}
	err = db.Debug().Model(&Profile{}).Limit(100).Find(&profiles).Error
	if err != nil {
		return &[]Profile{}, err
	}
	return &profiles, nil
}

func (p *Profile) FindUserProfileByID(db *gorm.DB, pid uint32) (*Profile, error) {
	var err error
	err = db.Debug().Model(Profile{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &Profile{}, errors.New("profile not found")
		}
		return &Profile{}, err
	}

	return p, err
}

func (p *Profile) UpdateAUserProfile(db *gorm.DB, pid uint32) (*Profile, error) {
	db = db.Debug().Model(&Profile{}).Where("id = ?", pid).Take(&Profile{}).UpdateColumns(
		map[string]interface{}{
			"name":         p.Name,
			"title":        p.Title,
			"bio":          p.Bio,
			"social_links": p.SocialLinks,
		},
	)

	if db.Error != nil {
		return &Profile{}, db.Error
	}

	// This is the display the updated profile
	err := db.Debug().Model(&Profile{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		fmt.Println("Erorr Profile", err)
		return &Profile{}, err
	}
	return p, nil
}

func (p *Profile) UpdateAUserProfilePic(db *gorm.DB, pid uint32, imageType string) (*Profile, error) {
	// Create a map to store the column to be updated
	updateColumns := map[string]interface{}{}

	// Check the image type and update the corresponding column
	switch imageType {
	case "profile_pic":
		updateColumns["profile_pic"] = p.ProfilePic
	case "cover_pic":
		updateColumns["cover_pic"] = p.CoverPic
	default:
		return nil, errors.New("invalid image type")
	}

	// Update the specified column
	db = db.Debug().Model(&Profile{}).Where("id = ?", pid).UpdateColumns(updateColumns)

	if db.Error != nil {
		return nil, db.Error
	}

	// Retrieve and return the updated profile
	err := db.Debug().Model(&Profile{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Profile) DeleteAUserProfile(db *gorm.DB, pid uint32) (int64, error) {
	// Retrieve the profile to be deleted
	var profileToDelete Profile
	err := db.Debug().Model(&Profile{}).Where("id = ?", pid).Take(&profileToDelete).Error
	if err != nil {
		return 0, err
	}

	// If the profile has a linked user, delete the user as well
	if profileToDelete.UserID != 0 {
		// Retrieve the user associated with the profile
		var userToDelete User
		err := db.Debug().Model(&User{}).Where("id = ?", profileToDelete.UserID).Take(&userToDelete).Error
		if err != nil {
			return 0, err
		}

		// Delete the user
		err = db.Debug().Delete(&userToDelete).Error
		if err != nil {
			return 0, err
		}
	}

	// Delete the profile
	db = db.Debug().Delete(&profileToDelete)

	if db.Error != nil {
		return 0, db.Error
	}

	return db.RowsAffected, nil
}
