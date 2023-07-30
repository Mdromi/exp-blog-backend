package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Mdromi/exp-blog-backend/api/middlewares"
	"github.com/Mdromi/exp-blog-backend/api/models"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"    //mysql database driver
	"gorm.io/driver/postgres" //postgres database driver
	"gorm.io/gorm"
)

type Server struct {
	DB     *gorm.DB
	Router *gin.Engine
}

var errList = make(map[string]string)

func (server *Server) Initialize(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {
	var err error

	// If you are using mysql, i added support for you here (dont forgot to edit the .env file)

	if Dbdriver == "mysql" {
		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)
		server.DB, err = gorm.Open(mysql.Open(DBURL), &gorm.Config{})
		if err != nil {
			fmt.Printf("Cannot connect to %s database", Dbdriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connect to the %s database", Dbdriver)
		}
	} else if Dbdriver == "postgres" {
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
		server.DB, err = gorm.Open(postgres.Open(DBURL), &gorm.Config{})
		if err != nil {
			fmt.Printf("Cannot connect to %s database", Dbdriver)
			log.Fatal("This is the error connecting to postgres:", err)
		} else {
			fmt.Printf("We are connect to the %s database", Dbdriver)
		}
	} else {
		fmt.Println("Unknown Driver")
	}

	// database migration
	server.DB.Debug().AutoMigrate(
		&models.User{},
		&models.Profile{},
		&models.SocialLink{},
		&models.Post{},
		&models.ResetPassword{},
		&models.Like{},
		&models.Dislike{},
		&models.Comment{},
		&models.Replyes{},
	)

	server.Router = gin.Default()
	server.Router.Use(middlewares.CORSMiddleware())

	server.initializeRoutes()
}

func (server *Server) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, server.Router))
}
