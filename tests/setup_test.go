package tests

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/Mdromi/exp-blog-backend/api/controllers"
	"github.com/Mdromi/exp-blog-backend/api/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var server = controllers.Server{}
var userInstance = models.User{}
var profileInstance = models.Profile{}
var linksInstance = models.SocialLink{}
var postInstance = models.Post{}
var likeInstance = models.Like{}
var commentInstance = models.Comment{}

func TestMain(m *testing.M) {
	//Since we add our .env in .gitignore, Circle CI cannot see it, so see the else statement
	if _, err := os.Stat("./../.env"); !os.IsNotExist(err) {
		var err error
		err = godotenv.Load(os.ExpandEnv("./../.env"))
		if err != nil {
			log.Fatalf("Error getting env %v\n", err)
		}
		Database()
		// AutoMigrateTables(server.DB)
	} else {
		CIBuild()
	}
	os.Exit(m.Run())
}

// When using CircleCI
func CIBuild() {
	var err error
	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", "127.0.0.1", "5432", "steven", "forum_db_test", "password")
	server.DB, err = gorm.Open(postgres.Open(DBURL), &gorm.Config{})
	if err != nil {
		fmt.Printf("Cannot connect to %s database\n", "postgres")
		log.Fatal("This is the error:", err)
	} else {
		fmt.Printf("We are connected to the %s database\n", "postgres")
	}
}

func Database() {

	var err error

	TestDbDriver := os.Getenv("TEST_DB_DRIVER")
	if TestDbDriver == "mysql" {
		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("TEST_DB_USER"), os.Getenv("TEST_DB_PASSWORD"), os.Getenv("TEST_DB_HOST"), os.Getenv("TEST_DB_PORT"), os.Getenv("TEST_DB_NAME"))
		server.DB, err = gorm.Open(mysql.Open(DBURL), &gorm.Config{})
		if err != nil {
			fmt.Printf("Cannot connect to %s database\n", TestDbDriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database\n", TestDbDriver)
		}
	}
	if TestDbDriver == "postgres" {
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", os.Getenv("TEST_DB_HOST"), os.Getenv("TEST_DB_PORT"), os.Getenv("TEST_DB_USER"), os.Getenv("TEST_DB_NAME"), os.Getenv("TEST_DB_PASSWORD"))
		server.DB, err = gorm.Open(postgres.Open(DBURL), &gorm.Config{})
		if err != nil {
			fmt.Printf("Cannot connect to %s database\n", TestDbDriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database\n", TestDbDriver)
		}
	}
}

func refreshAllTable() error {
	migrator := server.DB.Migrator()

	// Drop the Profile table if it exists
	err := migrator.DropTable(&models.User{}, &models.Profile{}, &models.SocialLink{}, &models.ResetPassword{}, &models.Post{}, &models.Like{}, &models.Dislike{}, &models.Comment{}, models.Replyes{})
	if err != nil {
		return err
	}

	// AutoMigrate to create the Profile table
	err = server.DB.AutoMigrate(&models.User{}, &models.Profile{}, &models.SocialLink{}, &models.ResetPassword{}, &models.Post{}, &models.Like{}, &models.Dislike{}, &models.Comment{}, models.Replyes{})
	if err != nil {
		return err
	}

	log.Printf("Successfully refreshed All table")
	return nil
}

func refreshUserTable() error {
	migrator := server.DB.Migrator()

	// Drop the User table if it exists
	err := migrator.DropTable(&models.User{})
	if err != nil {
		return err
	}

	// AutoMigrate to create the User table
	err = server.DB.AutoMigrate(&models.User{})
	if err != nil {
		return err
	}

	log.Printf("Successfully refreshed table")
	return nil
}

func refreshUserProfileTable() error {
	migrator := server.DB.Migrator()

	// Drop the Profile table if it exists
	err := migrator.DropTable(&models.User{}, &models.Profile{})
	if err != nil {
		return err
	}

	// AutoMigrate to create the Profile table
	err = server.DB.AutoMigrate(&models.User{}, &models.Profile{})
	if err != nil {
		return err
	}

	log.Printf("Successfully refreshed table")
	return nil
}

func seedOneUser() (models.User, error) {

	user := models.User{
		Username: "Pet",
		Email:    "pet@example.com",
		Password: "password",
	}

	err := server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func seedOneUserProfile() (models.Profile, error) {
	// Create a user
	user := models.User{
		Username:   "Pet",
		Email:      "pet@example.com",
		Password:   "password",
		AvatarPath: "image/pic",
	}

	// Save the user to the database
	err := server.DB.Create(&user).Error
	if err != nil {
		return models.Profile{}, err
	}

	// Create a profile and set the UserID to the user's ID
	profile := models.Profile{
		Name:       user.Username,
		Title:      "Profile Title for " + user.Username,
		Bio:        "Profile Bio for " + user.Username,
		UserID:     user.ID,
		ProfilePic: user.AvatarPath,
	}

	// Save the profile to the database
	err = server.DB.Create(&profile).Error
	if err != nil {
		return models.Profile{}, err
	}

	// // Create social links for the profile
	// links := models.SocialLink{
	// 	ProfileID: uint32(profile.ID),
	// 	Website:   "https://example.com",
	// 	Facebook:  "https://facebook.com/user",
	// 	Twitter:   "https://twitter.com/user",
	// 	Github:    "https://github.com/user",
	// }

	// // Save the links to the database using the Save method
	// err = server.DB.Save(&links).Error
	// if err != nil {
	// 	return models.Profile{}, err
	// }

	// // Associate the social links with the profile
	// profile.SocialLinks = &links
	// err = server.DB.Save(&profile).Error
	// if err != nil {
	// 	return models.Profile{}, err
	// }

	// Update the user's ProfileID with the ID of the newly created profile
	user.ProfileID = uint32(profile.ID)

	// Save the updated user to the database
	err = server.DB.Save(&user).Error
	if err != nil {
		return models.Profile{}, err
	}

	return profile, nil
}

func seedUsers() ([]models.User, error) {

	var err error
	if err != nil {
		return nil, err
	}
	users := []models.User{
		models.User{
			Username: "Steven",
			Email:    "steven@example.com",
			Password: "password",
		},
		models.User{
			Username: "Kenny",
			Email:    "kenny@example.com",
			Password: "password",
		},
	}

	for i, _ := range users {
		err := server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			return []models.User{}, err
		}
	}
	return users, nil
}

func seedUsersProfiles() ([]*models.Profile, error) {
	users, err := seedUsers()
	if err != nil {
		return nil, err
	}

	profiles := make([]*models.Profile, len(users))
	for i, user := range users {
		profile := &models.Profile{
			Name:       user.Username,
			Title:      "Profile Title for " + user.Username,
			Bio:        "Profile Bio for " + user.Username,
			UserID:     user.ID,
			ProfilePic: user.AvatarPath,
		}

		err := server.DB.Create(profile).Error
		if err != nil {
			return nil, err
		}

		profiles[i] = profile

		// Update the User model's Profile field with the created profile
		err = server.DB.Debug().Model(&models.User{}).Where("id = ?", user.ID).Take(&user.ProfileID).Error
		err = server.DB.Save(&user).Error
		if err != nil {
			return nil, err
		}
	}

	return profiles, nil
}

func refreshUserAndPostTable() error {
	migrator := server.DB.Migrator()

	// Drop the User and Post tables if they exist
	err := migrator.DropTable(&models.User{}, &models.Post{})
	if err != nil {
		return err
	}

	// AutoMigrate to create the User and Post tables
	err = server.DB.AutoMigrate(&models.User{}, &models.Post{})
	if err != nil {
		return err
	}

	log.Printf("Successfully refreshed tables")
	return nil
}

func refreshUserProfileAndPostTable() error {
	migrator := server.DB.Migrator()

	// Drop the User and Post tables if they exist
	err := migrator.DropTable(&models.User{}, &models.Post{}, &models.Profile{})
	if err != nil {
		return err
	}

	// AutoMigrate to create the User and Post tables
	err = server.DB.AutoMigrate(&models.User{}, &models.Post{}, &models.Profile{})
	if err != nil {
		return err
	}

	log.Printf("Successfully refreshed tables")
	return nil
}

func seedOneUserProfileAndOnePost() (models.Profile, models.Post, error) {

	profile, err := seedOneUserProfile()
	if err != nil {
		return models.Profile{}, models.Post{}, err
	}

	post := models.Post{
		Title:    "This is the title sam",
		Content:  "This is the content sam",
		AuthorID: uint(profile.ID),
	}
	err = server.DB.Model(&models.Post{}).Create(&post).Error
	if err != nil {
		return models.Profile{}, models.Post{}, err
	}
	return profile, post, nil
}

func seedUsersProfileAndPosts() ([]*models.Profile, []models.Post, error) {

	var err error

	profiles, err := seedUsersProfiles()

	if err != nil {
		return []*models.Profile{}, []models.Post{}, err
	}

	var posts = []models.Post{
		models.Post{
			Title:   "Title 1",
			Content: "Hello world 1",
		},
		models.Post{
			Title:   "Title 2",
			Content: "Hello world 2",
		},
	}

	for i, _ := range profiles {
		posts[i].AuthorID = profiles[i].ID

		err = server.DB.Model(&models.Post{}).Create(&posts[i]).Error
		if err != nil {
			log.Fatalf("cannot seed posts table: %v", err)
		}
	}
	return profiles, posts, nil
}

func refreshUserProfilePostAndLikeTable() error {
	migrator := server.DB.Migrator()

	// Delete records from the post_likes table first
	err := server.DB.Where("1 = 1").Delete(&models.Like{}).Error
	if err != nil {
		return err
	}

	// Drop the User, Profile, Post, and Like tables if they exist
	err = migrator.DropTable(&models.User{}, &models.Profile{}, &models.Post{}, &models.Like{})
	if err != nil {
		return err
	}

	// AutoMigrate to create the User, Profile, Post, and Like tables
	err = server.DB.AutoMigrate(&models.User{}, &models.Profile{}, &models.Post{}, &models.Like{})
	if err != nil {
		return err
	}

	log.Printf("Successfully refreshed user, profile, post, and like tables")
	return nil
}

func seedUsersPostsAndLikes() (*models.Post, []*models.Profile, []*models.Like, error) {
	// The idea here is: two users can like one post
	var err error

	profiles, err := seedUsersProfiles()

	if err != nil {
		return nil, nil, nil, err
	}

	post := models.Post{
		Title:   "This is the title",
		Content: "This is the content",
	}
	post.AuthorID = profiles[0].ID

	err = server.DB.Model(&models.Post{}).Create(&post).Error
	if err != nil {
		log.Fatalf("cannot seed post table: %v", err)
	}

	likes := []*models.Like{
		&models.Like{
			ProfileID: profiles[0].ID,
			PostID:    post.ID,
		},
		&models.Like{
			ProfileID: profiles[1].ID,
			PostID:    post.ID,
		},
	}
	for i := range profiles {
		err = server.DB.Model(&models.Like{}).Create(likes[i]).Error
		if err != nil {
			log.Fatalf("cannot seed likes table: %v", err)
		}
	}
	return &post, profiles, likes, nil
}

func seedUsersPostsAndDislikes() (*models.Post, []*models.Profile, []*models.Dislike, error) {
	// The idea here is: two users can dislike one post
	var err error

	profiles, err := seedUsersProfiles()

	if err != nil {
		return nil, nil, nil, err
	}

	post := models.Post{
		Title:   "This is the title",
		Content: "This is the content",
	}
	post.AuthorID = profiles[0].ID

	err = server.DB.Model(&models.Post{}).Create(&post).Error
	if err != nil {
		log.Fatalf("cannot seed post table: %v", err)
	}

	dislikes := []*models.Dislike{
		&models.Dislike{
			ProfileID: profiles[0].ID,
			PostID:    post.ID,
		},
		&models.Dislike{
			ProfileID: profiles[1].ID,
			PostID:    post.ID,
		},
	}
	for i := range profiles {
		err = server.DB.Model(&models.Dislike{}).Create(dislikes[i]).Error
		if err != nil {
			log.Fatalf("cannot seed dislikes table: %v", err)
		}
	}
	return &post, profiles, dislikes, nil
}

func refreshUserPostAndCommentTable() error {
	migrator := server.DB.Migrator()

	// Drop the User, Post, and Comment tables if they exist
	err := migrator.DropTable(&models.User{}, &models.Post{}, &models.Comment{})
	if err != nil {
		return err
	}

	// AutoMigrate to create the User, Post, and Comment tables
	err = server.DB.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{})
	if err != nil {
		return err
	}

	log.Printf("Successfully refreshed user, post, and comment tables")
	return nil
}

func seedUsersPostsAndComments() (models.Post, []models.User, []models.Comment, error) {
	// The idea here is: two users can comment one post
	var err error
	var users = []models.User{
		models.User{
			Username: "Steven",
			Email:    "steven@example.com",
			Password: "password",
		},
		models.User{
			Username: "Magu",
			Email:    "magu@example.com",
			Password: "password",
		},
	}
	post := models.Post{
		Title:   "This is the title",
		Content: "This is the content",
	}
	err = server.DB.Model(&models.Post{}).Create(&post).Error
	if err != nil {
		log.Fatalf("cannot seed post table: %v", err)
	}
	var comments = []models.Comment{
		models.Comment{
			Body:   "user 1 made this comment",
			UserID: 1,
			PostID: uint64(post.ID),
		},
		models.Comment{
			Body:   "user 2 made this comment",
			UserID: 2,
			PostID: uint64(post.ID),
		},
	}
	for i, _ := range users {
		err = server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		err = server.DB.Model(&models.Like{}).Create(&comments[i]).Error
		if err != nil {
			log.Fatalf("cannot seed comments table: %v", err)
		}
	}
	return post, users, comments, nil
}

func refreshUserAndResetPasswordTable() error {
	migrator := server.DB.Migrator()

	// Drop the User and ResetPassword tables if they exist
	err := migrator.DropTable(&models.User{}, &models.ResetPassword{})
	if err != nil {
		return err
	}

	// AutoMigrate to create the User and ResetPassword tables
	err = server.DB.AutoMigrate(&models.User{}, &models.ResetPassword{})
	if err != nil {
		return err
	}

	log.Printf("Successfully refreshed user and resetpassword tables")
	return nil
}

// Seed the reset password table with the token
func seedResetPassword() (models.ResetPassword, error) {

	resetDetails := models.ResetPassword{
		Token: "awesometoken",
		Email: "pet@example.com",
	}
	err := server.DB.Model(&models.ResetPassword{}).Create(&resetDetails).Error
	if err != nil {
		return models.ResetPassword{}, err
	}
	return resetDetails, nil
}
