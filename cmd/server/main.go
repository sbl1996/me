package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/sbl1996/me/controllers/api"
	"github.com/sbl1996/me/database"
	"github.com/sbl1996/me/utils"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	db, err := database.Init()
	utils.Check(err, "Connecting database")
	defer db.Close()

	router := gin.New()
	router.Use(gin.Logger())

	router.NoRoute(func(c *gin.Context) {
		c.File("web/public/index.html")
	})
	router.StaticFile("/manifest.json", "web/public/manifest.json")
	router.StaticFile("/favicon.ico", "web/public/favicon.ico")
	router.StaticFile("/", "web/public/index.html")
	router.Static("/static", "web/public")

	router.GET("/api/post", api.GetPostHandler())
	router.GET("/api/posts", api.GetPostsHandler())
	router.POST("/api/post", api.CreatePostHandler())
	router.POST("/api/post/update", api.UpdatePostHandler())

	router.Run(":" + port)
}
