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

		c.Redirect(http.StatusMovedPermanently, "/post?title="+title)
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

		c.Redirect(http.StatusMovedPermanently, "/post?title="+title)
	}
}
