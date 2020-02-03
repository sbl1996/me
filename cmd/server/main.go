package main

import (
	"log"
	"os"

	"github.com/gin-gonic/contrib/static"
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
	// router.StaticFile("/manifest.json", "web/public/manifest.json")
	// router.StaticFile("/favicon.ico", "web/public/favicon.ico")
	// router.StaticFile("/", "web/public/index.html")
	router.Use(static.Serve("/", static.LocalFile("./web/public", true)))
	// router.Static("/", "web/public")

	apiRouter := router.Group("/api")
	{
		apiRouter.GET("/posts/:id", api.GetPostHandler())
		apiRouter.GET("/posts", api.GetPostsHandler())
		apiRouter.POST("/posts", api.CreatePostHandler())
		apiRouter.PATCH("/posts", api.UpdatePostHandler())
		apiRouter.DELETE("/posts/:id", api.DeletePostHandler())
	}

	router.Run(":" + port)
}
