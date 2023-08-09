package testdata

import (
	"fmt"
	"net/http"

	"github.com/Mdromi/exp-blog-backend/api/models"
)

// CreateProfileTestCase represents a test case for creating a user profile
type CreateProfileTestCase struct {
	InputJSON  string
	StatusCode int
	UserID     uint
	Name       string
	Title      string
	ProfilePic string
	SocialLink models.SocialLink
}

type UpdateProfileTestCase struct {
	ID         string
	UpdateJSON string
	UserID     uint
	StatusCode int
	Name       string
	Title      string
	ProfilePic string
	TokenGiven string
}

func CreateProfileSamples(loginUserID uint) []CreateProfileTestCase {
	createProfileSamples := []CreateProfileTestCase{
		{
			InputJSON: fmt.Sprintf(`{
				"user_id": %d,
				"name": "Pet",
				"title": "This is the title",
				"bio": "This is the Bio",
				"profile_pic": "/images/profile.jpg", "social_links": {
					"facebook": "www.facebook.com",
					"twitter": "www.twitter.com",
					"instagram": "www.instagram.com"
				}
			}`, loginUserID),
			StatusCode: 201,
			UserID:     loginUserID,
			Name:       "Pet",
			Title:      "This is the title",
			ProfilePic: "/images/profile.jpg",
		},
		{
			InputJSON: fmt.Sprintf(`{
				"user_id": %d, 
				"name": "Pet", 
				"title": "This is the title", 
				"bio": "This is the Bio", 
				"profile_pic": "/images/profile.jpg", "social_links": {
					"facebook": "www.facebook.com", "twitter": "www.twitter.com", "instagram": "www.instagram.com"
				}
			}`, 0),
			StatusCode: 404,
			UserID:     0,
			Name:       "Pet",
			Title:      "This is the title",
			ProfilePic: "/images/profile.jpg",
		},
		{
			InputJSON: fmt.Sprintf(`{
				"user_id": %d, 
				"name": "Pet", 
				"title": "This is the title", 
				"bio": "This is the Bio", 
				"profile_pic": "/images/profile.jpg", "social_links": {
					"facebook": "www.facebook.com", "twitter": "www.twitter.com", "instagram": "www.instagram.com"
				}
			}`, 342049902),
			StatusCode: 404,
			UserID:     0,
			Name:       "Pet",
			Title:      "This is the title",
			ProfilePic: "/images/profile.jpg",
		},
		{
			InputJSON: fmt.Sprintf(`{
				"user_id": %d, 
				"name": "Pet", 
				"title": "", 
				"bio": "This is the Bio", 
				"profile_pic": "/images/profile.jpg", "social_links": {
					"facebook": "www.facebook.com", "twitter": "www.twitter.com", "instagram": "www.instagram.com"
				}
			}`, loginUserID),
			StatusCode: 400,
			UserID:     loginUserID,
			Name:       "Pet",
			Title:      "",
			ProfilePic: "/images/profile.jpg",
		},
		{
			InputJSON: fmt.Sprintf(`{
				"user_id": %d, 
				"name": "", 
				"title": "This is the title", 
				"bio": "This is the Bio", 
				"profile_pic": "/images/profile.jpg", "social_links": {
					"facebook": "www.facebook.com", "twitter": "www.twitter.com", "instagram": "www.instagram.com"
				}
			}`, loginUserID),
			StatusCode: 400,
			UserID:     loginUserID,
			Name:       "",
			Title:      "This is the title",
			ProfilePic: "/images/profile.jpg",
		},
		{
			InputJSON: fmt.Sprintf(`{
				"user_id": %d, 
				"name": "Pet", 
				"title": "This is the title", 
				"bio": "", 
				"profile_pic": "/images/profile.jpg", "social_links": {
					"facebook": "www.facebook.com", "twitter": "www.twitter.com", "instagram": "www.instagram.com"
				}
			}`, loginUserID),
			StatusCode: 400,
			UserID:     loginUserID,
			Name:       "Pet",
			Title:      "This is the title",
			ProfilePic: "/images/profile.jpg",
		},
		{
			InputJSON: fmt.Sprintf(`{
				"user_id": %d, 
				"name": "", 
				"title": "", 
				"bio": "", 
				"profile_pic": "/images/profile.jpg", "social_links": {}
			}`, loginUserID),
			StatusCode: 400,
			UserID:     loginUserID,
			Name:       "",
			Title:      "",
			ProfilePic: "/images/profile.jpg",
		},
		{
			InputJSON:  fmt.Sprintf(`{ "user_id": %d, "name": "", "title": "This is the title", "bio": "This is the Bio", "profile_pic": "/images/profile.jpg", "social_links": {} }`, loginUserID),
			StatusCode: http.StatusBadRequest,
		},
		{
			InputJSON:  fmt.Sprintf(`{ "user_id": %d, "name": "Pet", "title": "", "bio": "This is the Bio", "profile_pic": "/images/profile.jpg", "social_links": {} }`, loginUserID),
			StatusCode: http.StatusBadRequest,
		},
		{
			InputJSON:  fmt.Sprintf(`{ "user_id": %d, "name": "Pet", "title": "This is the title", "bio": "", "profile_pic": "/images/profile.jpg", "social_links": {} }`, loginUserID),
			StatusCode: http.StatusBadRequest,
		},
		{
			InputJSON:  fmt.Sprintf(`{ "user_id": %d, "name": "", "title": "", "bio": "", "profile_pic": "/images/profile.jpg", "social_links": {} }`, loginUserID),
			StatusCode: http.StatusBadRequest,
		},
	}
	return createProfileSamples
}

func UpdateProfileSamples(profileID, tokenString string, loginUserID uint) []UpdateProfileTestCase {
	updateProfileTestCase := []UpdateProfileTestCase{

		{
			ID: profileID,
			UpdateJSON: fmt.Sprintf(`{
				"user_id": %d, 
				"name": "Pet 1", 
				"title": "This is the title - 1", 
				"bio": "This is the Bio - 1", 
				"profile_pic": "/images/profile.jpg", "social_links": {
					"facebook": "www.facebook.com/mdromi",
					"twitter": "www.twitter.com/mdromi",
					"instagram": "www.instagram.com/mdromi"
				}
			}`, loginUserID),
			StatusCode: 200,
			UserID:     loginUserID,
			Name:       "Pet",
			Title:      "This is the title",
			ProfilePic: "/images/profile.jpg",
			TokenGiven: tokenString,
		},
		{
			ID: profileID,
			UpdateJSON: fmt.Sprintf(`{
				"user_id": %d,
				"name": "Pet",
				"title": "This is the title",
				"bio": "This is the Bio",
				"profile_pic": "/images/profile.jpg", "social_links": {
					"facebook": "www.facebook.com", "twitter": "www.twitter.com", "instagram": "www.instagram.com"
				}
			}`, 0),
			StatusCode: 401,
			UserID:     0,
			Name:       "Pet",
			Title:      "This is the title",
			ProfilePic: "/images/profile.jpg",
			TokenGiven: tokenString,
		},
		{
			ID: "342049902",
			UpdateJSON: fmt.Sprintf(`{
				"user_id": %d,
				"name": "Pet",
				"title": "This is the title",
				"bio": "This is the Bio",
				"profile_pic": "/images/profile.jpg", "social_links": {
					"facebook": "www.facebook.com", "twitter": "www.twitter.com", "instagram": "www.instagram.com"
				}
			}`, loginUserID),
			StatusCode: 404,
			UserID:     loginUserID,
			Name:       "Pet",
			Title:      "This is the title",
			ProfilePic: "/images/profile.jpg",
			TokenGiven: tokenString,
		},
		{
			ID: profileID,
			UpdateJSON: fmt.Sprintf(`{
				"user_id": %d,
				"name": "Pet",
				"title": "",
				"bio": "This is the Bio",
				"profile_pic": "/images/profile.jpg", "social_links": {
					"facebook": "www.facebook.com", "twitter": "www.twitter.com", "instagram": "www.instagram.com"
				}
			}`, loginUserID),
			StatusCode: 422,
			UserID:     loginUserID,
			Name:       "Pet",
			Title:      "",
			ProfilePic: "/images/profile.jpg",
			TokenGiven: tokenString,
		},
		{
			ID: profileID,
			UpdateJSON: fmt.Sprintf(`{
				"user_id": %d,
				"name": "",
				"title": "This is the title",
				"bio": "This is the Bio",
				"profile_pic": "/images/profile.jpg", "social_links": {
					"facebook": "www.facebook.com", "twitter": "www.twitter.com", "instagram": "www.instagram.com"
				}
			}`, loginUserID),
			StatusCode: 422,
			UserID:     loginUserID,
			Name:       "",
			Title:      "This is the title",
			ProfilePic: "/images/profile.jpg",
			TokenGiven: tokenString,
		},
		{
			ID: profileID,
			UpdateJSON: fmt.Sprintf(`{
				"user_id": %d,
				"name": "Pet",
				"title": "This is the title",
				"bio": "",
				"profile_pic": "/images/profile.jpg", "social_links": {
					"facebook": "www.facebook.com", "twitter": "www.twitter.com", "instagram": "www.instagram.com"
				}
			}`, loginUserID),
			StatusCode: 422,
			UserID:     loginUserID,
			Name:       "Pet",
			Title:      "This is the title",
			ProfilePic: "/images/profile.jpg",
			TokenGiven: tokenString,
		},
		{
			ID: profileID,
			UpdateJSON: fmt.Sprintf(`{
				"user_id": %d,
				"name": "",
				"title": "",
				"bio": "",
				"profile_pic": "/images/profile.jpg", "social_links": {}
			}`, loginUserID),
			StatusCode: 422,
			UserID:     loginUserID,
			Name:       "",
			Title:      "",
			ProfilePic: "/images/profile.jpg",
			TokenGiven: tokenString,
		},
		{
			ID: profileID,
			UpdateJSON: fmt.Sprintf(`{
				"user_id": %d,
				"name": "",
				"title": "",
				"bio": "",
				"profile_pic": "/images/profile.jpg", "social_links": {}
			}`, loginUserID),
			StatusCode: 422,
			UserID:     loginUserID,
			Name:       "",
			Title:      "",
			TokenGiven: tokenString,
		},
	}
	return updateProfileTestCase
}

// func setupTestDependencies() (*controllers.Server, uint, error) {
// 	err := refreshAllTable()
// 	if err != nil {
// 		return nil, 0, err
// 	}

// 	user, err := seedOneUser()
// 	if err != nil {
// 		return nil, 0, err
// 	}

// 	loginUserID := user.ID

// 	// Set up your router and other dependencies
// 	router := gin.Default()
// 	server := &controllers.Server{
// 		DB:     models.DB,
// 		Router: router,
// 	}

// 	return server, loginUserID, nil
// }
