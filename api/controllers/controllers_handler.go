package controllers

import (
	"encoding/json"

	"github.com/Mdromi/exp-blog-backend/api/models"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func FindUserByID(db *gorm.DB, userID uint32) (*models.User, error) {
	userModel := models.User{}
	user, err := userModel.FindUserByID(db, userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func FindUserProfileByID(db *gorm.DB, userID uint32) (*models.Profile, error) {
	profileModel := models.Profile{}
	profile, err := profileModel.FindUserProfileByID(db, userID)
	if err != nil {
		return nil, err
	}
	return profile, nil
}

func GetSocialLinksFromBody(requestBody map[string]string) *models.SocialLink {
	socialLinksStr, ok := requestBody["social_links"]
	if ok && socialLinksStr != "" {
		var socialLinksMap map[string]interface{}
		if err := json.Unmarshal([]byte(socialLinksStr), &socialLinksMap); err == nil {
			return mapToSocialLink(socialLinksMap)
		}
	}
	return nil
}

// handlare function
func ValidateProfileFields(profile *models.Profile) map[string]string {
	errList := map[string]string{}

	validate := validator.New()
	if err := validate.Struct(profile); err != nil {
		// Handle validation errors
		if _, ok := err.(*validator.InvalidValidationError); ok {
			// Handle error from the validation library itself (e.g., invalid struct)
			errList["Validation_error"] = "Invalid input data"
		} else {
			// Handle specific validation errors for each field
			for _, fieldErr := range err.(validator.ValidationErrors) {
				fieldName := fieldErr.Field()
				switch fieldName {
				case "Name":
					errList["Profile_name"] = "Name is required and should be between 2 and 50 characters"
				case "Title":
					errList["Profile_tile"] = "Title should be less than or equal to 100 characters"
				case "Bio":
					errList["Profile_bio"] = "Bio should be less than or equal to 500 characters"
					// Add more cases for other fields if needed
				}
			}
		}
	}

	return errList
}

func mapToSocialLink(input map[string]interface{}) *models.SocialLink {
	socialLink := &models.SocialLink{
		Facebook:  getString(input, "facebook"),
		Twitter:   getString(input, "twitter"),
		Instagram: getString(input, "instagram"),
		// Add more fields as needed
	}

	return socialLink
}

func getString(input map[string]interface{}, key string) string {
	if val, ok := input[key]; ok {
		if strVal, ok := val.(string); ok {
			return strVal
		}
	}
	return ""
}
