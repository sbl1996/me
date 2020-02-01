package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sbl1996/me/models/post"
	"github.com/sbl1996/me/utils"
)

func CreatePostHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		title := c.PostForm("title")
		content := c.PostForm("content")
		dateStr := c.PostForm("date")
		date, err := time.Parse("2006-01-02T15:04", dateStr)
		utils.Check(err, "Parse date")
		err = post.CreatePost(title, content, date)
		utils.Check(err, "Inserting post")

		c.Redirect(http.StatusMovedPermanently, "/post/"+title)
		// c.JSON(http.StatusCreated, gin.H{
		// 	title: title,
		// })
	}
}

func UpdatePostHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.PostForm("id"), 10, 64)
		utils.Check(err, "Parse ID")
		title := c.PostForm("title")
		content := c.PostForm("content")
		dateStr := c.PostForm("date")
		date, err := time.Parse("2006-01-02T15:04", dateStr)
		utils.Check(err, "Parse date")
		err = post.UpdatePost(uint(id), title, content, date)
		utils.Check(err, "Updating post")

		c.Redirect(http.StatusMovedPermanently, "/post/"+title)
	}
}

func GetPostHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var p post.Post
		var err error

		id := c.Query("id")
		if id != "" {
			id, err := strconv.ParseUint(id, 10, 64)
			utils.Check(err, "Parse ID")
			p, err = post.GetByID(uint(id))
		} else {
			title := c.Query("title")
			p, err = post.GetByTitle(title)
		}
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Post not found",
			})
		} else {
			c.JSON(http.StatusOK, p)
		}
	}
}

func GetPostsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var titles []string
		var err error

		year := c.Query("year")
		q := c.Query("q")
		if (year == "") && (q == "") {
			titles, err = post.GetAllPostTitles()
		} else if year != "" {
			titles, err = post.GetPostTitlesOfYear(year)
		} else {
			titles, err = post.SearchPosts(q)
		}
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Posts not found",
			})
		} else {
			c.JSON(http.StatusOK, titles)
		}
	}
}
