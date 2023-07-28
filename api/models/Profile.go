package models

import (
	"errors"
	"html"
	"strings"

	"gorm.io/gorm"
)

type Profile struct {
	gorm.Model
	UserID        uint32  `gorm:"not null" json:"user_id"`
	User          *User   `gorm:"foreignKey:UserID" json:"user"`
	Name          string  `gorm:"type:varchar(50);not null" json:"name" validate:"min=2,max=50"`
	Title         string  `gorm:"type:varchar(100);not null" json:"title" validate:"max=100"`
	Bio           string  `gorm:"type:text;not null" json:"bio" validate:"max=500"`
	ProfilePic    string  `gorm:"type:varchar(255)" json:"profile_pic"`
	Links         Links   `json:"links"`
	Posts         []*Post `gorm:"many2many:profile_posts;" json:"posts"`
	Bookmarks     []*Post `gorm:"many2many:profile_bookmarks;" json:"bookmarks"`
	Flowing       []*User `gorm:"many2many:user_followers;association_foreignkey:FlowingID;" json:"flowing"`
	LikedPosts    []*Post `gorm:"many2many:user_liked_posts;" json:"liked_posts"`
	DislikedPosts []*Post `gorm:"many2many:user_disliked_posts;" json:"disliked_posts"`
}

func (u *Profile) BeforeSave() error {
	return nil
}

func (p *Profile) Prepare() {
	p.Name = html.EscapeString(strings.TrimSpace(p.Name))
	p.Title = html.EscapeString(strings.TrimSpace(p.Title))
	p.Bio = html.EscapeString(strings.TrimSpace(p.Bio))
	p.User = &User{}
	p.Links = Links{}
	p.Posts = []*Post{}
	p.Bookmarks = []*Post{}
	p.LikedPosts = []*Post{}
	p.DislikedPosts = []*Post{}
}

func (p *Profile) AfterFind() (err error) {
	if err != nil {
		return err
	}

	userAvatarPath := p.User.AvatarPath

	if userAvatarPath != "" {
		p.ProfilePic = userAvatarPath
	}

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

	// Check if the user already has a profile
	if p.UserID != 0 {
		// Retrieve the user from the database
		var user User
		err = db.Debug().Model(&User{}).Where("id = ?", p.UserID).Take(&user).Error
		if err != nil {
			return &Profile{}, err
		}

		// Check if the user already has a profile
		if user.ProfileID != 0 {
			// User already has a profile, return an error
			return &Profile{}, errors.New("user already has a profile")
		}
	}

	// Create the profile
	err = db.Debug().Model(&Profile{}).Create(&p).Error
	if err != nil {
		return &Profile{}, err
	}

	// Update the User.ProfileID
	if p.ID != 0 {
		// First, retrieve the user from the database
		var user User
		err = db.Debug().Model(&User{}).Where("id = ?", p.UserID).Take(&user).Error
		if err != nil {
			return &Profile{}, err
		}

		// Then, update the User.ProfileID
		user.ProfileID = uint32(p.ID)
		err = db.Debug().Model(&User{}).Where("id = ?", p.UserID).Update("profile_id", p.ID).Error
		if err != nil {
			return &Profile{}, err
		}
	}

	return p, nil
}

// THE ONLY PERSON THAT NEED TO DO THIS IS THE ADMIN, SO I HAVE COMMENTED THE ROUTES, SO SOMEONE ELSE DONT VIES THIS DEATAILS
func (p *Profile) FindAllUsersProfile(db *gorm.DB) (*[]Profile, error) {
	var err error
	profiles := []Profile{}
	// err = db.Debug().Model(&Profile{}).Limit(100).Find(&profiles).Error
	err = db.Debug().Preload("User").Limit(100).Find(&profiles).Error
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
	// Find the profile by ID and preload the User association
	err := db.Debug().Preload("User").Model(&Profile{}).Where("id = ?", pid).Take(p).Error
	if err != nil {
		return nil, err
	}

	// Update the profile fields
	p.Name = html.EscapeString(strings.TrimSpace(p.Name))
	p.Title = html.EscapeString(strings.TrimSpace(p.Title))
	p.Bio = html.EscapeString(strings.TrimSpace(p.Bio))

	// Save the updated profile to the database
	err = db.Debug().Save(p).Error
	if err != nil {
		return nil, err
	}

	// Update the user's ProfileID in the user model
	updatedUser := p.User
	updatedUser.ProfileID = uint32(p.ID)
	err = db.Debug().Save(updatedUser).Error
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Profile) UpdateAUserProfilePic(db *gorm.DB, pid uint32) (*Profile, error) {
	// Update the profile_pic field
	db = db.Debug().Model(&Profile{}).Where("id = ?", pid).UpdateColumns(
		map[string]interface{}{
			"profile_pic": p.ProfilePic,
		},
	)

	if db.Error != nil {
		return nil, db.Error
	}

	// Retrieve and return the updated user
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

	// Delete the profile
	db = db.Debug().Model(&Profile{}).Where("id = ?", pid).Delete(&Profile{})
	if db.Error != nil {
		return 0, db.Error
	}

	// If the profile has a linked user, delete the user as well
	if profileToDelete.UserID != 0 {
		db = db.Debug().Model(&User{}).Where("id = ?", profileToDelete.UserID).Delete(&User{})
		if db.Error != nil {
			return 0, db.Error
		}
	}

	return db.RowsAffected, nil
}