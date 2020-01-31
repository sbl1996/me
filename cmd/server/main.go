package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/sbl1996/me/controllers/api"
	"github.com/sbl1996/me/controllers/post"
	"github.com/sbl1996/me/database"
	"github.com/sbl1996/me/utils"
)

func helloHandler(c *gin.Context) {
	c.String(http.StatusOK, "Hello")
}

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
	router.LoadHTMLGlob("web/templates/*.tmpl.html")
	router.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.tmpl.html", gin.H{
			"title": "Not Found",
		})
	})
	router.StaticFile("/favicon.ico", "web/static/favicon.ico")
	router.Static("/static", "web/static")

	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/posts")
	})

	router.POST("/api/post", api.CreatePostHandler())
	router.POST("/api/post/update", api.UpdatePostHandler())

	router.GET("/posts", post.PostsHandler())
	router.GET("/post/new", func(c *gin.Context) {
		c.HTML(http.StatusOK, "create.tmpl.html", gin.H{
			"title": "New Post",
		})
	})
	router.GET("/post", post.GetPostHandler())
	router.GET("/post/edit/:id", post.EditPostHandler())
	router.GET("/posts/search", post.SearchPostHandler())
	router.GET("/hello", helloHandler)
	router.Run(":" + port)
}
