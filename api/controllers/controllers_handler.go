package controllers

import (
	"encoding/json"

	"github.com/Mdromi/exp-blog-backend/api/models"
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
