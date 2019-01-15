package main

import (
	"contact-book-api/controller"
	"contact-book-api/model"
	"context"
	_ "fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// Conf holds the configuration struct from config.json
var Conf *Config
var DB *sqlx.DB

func init() {
	//Load config
	Conf = LoadConfiguration()
	dsn := Conf.DB.User + ":" + Conf.DB.Password + "@tcp(" + Conf.DB.Host + ":" + Conf.DB.Port + ")/" + Conf.DB.Name
	DB = model.InitDB(dsn)
}

func setupRouter() *gin.Engine {
	router := gin.New()
	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	return router
}

func main() {
	router := setupRouter()
	if Conf.App.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())

	// Auth Header Verification middleware
	router.Use(DatabaseMiddleware())
	router.Use(AuthHeaderMiddleware())

	// Enable gin logging if enabled in config
	if Conf.App.Logging {
		router.Use(gin.Logger())
	}

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Contact Book",
		})
	})

	v1 := router.Group("/v1")
	{
		v1.GET("/contacts/page/:page", controller.GetContacts)
		v1.DELETE("/contact/:id", controller.DeleteContactByID)
		v1.POST("/contact", controller.CreateContact)
		v1.PUT("/contact/:id", controller.EditContactByID)
		v1.PUT("/contact", controller.EditContactByEmail)
		v1.GET("/contacts/search", controller.SearchContact)

	}

	srv := &http.Server{
		Addr:    Conf.App.ListenPort,
		Handler: router,
	}
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting Down Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
