package api

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sbl1996/me/models/post"
	"github.com/sbl1996/me/utils"
)

func CreatePostHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var p post.Post
		err := c.ShouldBindJSON(&p)
		if err != nil {
			log.Printf("%s: %q", "Parse post", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		if strings.Contains(p.Title, "Fuck") {
			err := "be polite"
			log.Printf(err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
		} else {
			id, err := post.CreatePost(p.Title, p.Content, p.Date)
			if err != nil {
				log.Printf("%s: %q", "Create post", err)
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
			} else {
				log.Printf("Created #%d", id)
				c.JSON(http.StatusCreated, gin.H{
					"ID": id,
				})
			}
		}
	}
}

func UpdatePostHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var p post.Post
		err := c.ShouldBindJSON(&p)
		if err != nil {
			log.Printf("%s: %q", "Parse post", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		err = post.UpdatePost(p.ID, p.Title, p.Content, p.Date)
		if err != nil {
			log.Printf("%s: %q", "Update post", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		} else {
			log.Printf("Updated #%d", p.ID)
			c.Writer.WriteHeader(http.StatusOK)
		}
	}
}

func GetPostHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var p post.Post

		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		utils.Check(err, "Parse ID")
		p, err = post.GetByID(uint(id))
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
		var items []post.PostItem
		var err error

		year := c.Query("year")
		q := c.Query("q")
		if (year == "") && (q == "") {
			items, err = post.GetAllPostItems()
		} else if year != "" {
			items, err = post.GetPostItemsByYear(year)
		} else {
			items, err = post.SearchPosts(q)
		}
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Posts not found",
			})
		} else {
			c.JSON(http.StatusOK, items)
		}
	}
}

func DeletePostHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := post.ParseID(idStr)
		if err != nil {
			log.Printf("%s: %q", "parse ID", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		err = post.DeletePost(id)
		if err != nil {
			log.Printf("%s: %q", "delete post", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		} else {
			log.Printf("Deleted #%d", id)
			c.Writer.WriteHeader(http.StatusOK)
		}
	}
}
