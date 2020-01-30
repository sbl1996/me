package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/russross/blackfriday"

	"github.com/sbl1996/me/pkg/util"
)

type Post struct {
	gorm.Model
	Title   string
	Content string
	Date    time.Time
}

func postHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		title := c.Param("title")
		var p Post
		if err := db.Where("title = ?", title).First(&p).Error; err != nil {
			c.Writer.WriteHeader(http.StatusNotFound)
			c.Writer.WriteString("Not Found")
		} else {
			postHTML := template.HTML(string(blackfriday.Run([]byte(p.Content))))
			c.HTML(http.StatusOK, "post.tmpl.html", gin.H{
				"title":   p.Title,
				"content": postHTML,
			})
		}
	}
}

func postsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var posts []Post
		db.Order("date desc").Select("title").Find(&posts)

		var titles []string
		for _, p := range posts {
			titles = append(titles, p.Title)
		}
		c.HTML(http.StatusOK, "posts.tmpl.html", gin.H{
			"titles": titles,
		})
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	dbURL := os.Getenv("DATABASE_URL")
	db, err := gorm.Open("postgres", dbURL)
	util.Check(err, "Error connecting database")
	defer db.Close()

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("web/templates/*.tmpl.html")
	router.StaticFile("/favicon.ico", "web/static/favicon.ico")
	router.Static("/static", "web/static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", gin.H{
			"title": "Getting Started with Go on Heroku",
		})
	})
	router.GET("/posts", postsHandler(db))
	router.GET("/post/:title", postHandler(db))

	router.Run(":" + port)
}
