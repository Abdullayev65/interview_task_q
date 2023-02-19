package app

import (
	"fmt"
	"github.com/Abdullayev65/interview_task_q/internal/app/moduls"
	"github.com/Abdullayev65/interview_task_q/internal/app/utill"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/Abdullayev65/interview_task_q/internal/app/endpoint"
	"github.com/Abdullayev65/interview_task_q/internal/app/mw"
	"github.com/Abdullayev65/interview_task_q/internal/app/service"
)

type App struct {
	e    *endpoint.Endpoint
	s    *service.Service
	echo *echo.Echo
}

var create = true

func New() (*App, error) {
	a := &App{}

	dsn := "host=localhost user=postgres password=root123 dbname=postgres port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	tokenJWT := utill.NewToken("salt", 6*time.Hour)
	a.s = service.New(db, tokenJWT)

	a.e = endpoint.New(a.s)

	a.echo = echo.New()

	a.echo.Use(mw.ErrorHandler)

	if create {
		migrate(db)
	}
	a.initApis()

	return a, nil
}

func (a *App) Run() error {
	fmt.Println("server running")

	err := a.echo.Start(":8080")
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (a *App) initApis() {
	a.echo.POST("/sign-up", a.e.SignUp)
	a.echo.POST("/log-in", a.e.LogIn)
	a.echo.POST("/post", a.e.AddPost, a.e.SetUserIdForMW)
	a.echo.GET("/post/list", a.e.PostsOfUser)
	a.echo.POST("/like", a.e.Like, a.e.SetUserIdForMW)
	a.echo.POST("/comment", a.e.AddComment, a.e.SetUserIdForMW)
	a.echo.GET("/comment", a.e.CommentsByPostId)
}

func migrate(db *gorm.DB) {
	db.AutoMigrate(&moduls.User{})
	db.AutoMigrate(&moduls.Post{})
	db.AutoMigrate(&moduls.Comment{})
	db.AutoMigrate(&moduls.Like{})

}
