package post

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/russross/blackfriday"

	"github.com/sbl1996/me/models/post"
	p "github.com/sbl1996/me/models/post"
	"github.com/sbl1996/me/utils"
)

func GetPostHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		title := c.Query("title")
		p, err := p.GetByTitle(title)
		if err != nil {
			c.HTML(http.StatusNotFound, "404.tmpl.html", gin.H{
				"title": "Not Found",
			})
		} else {
			postHTML := template.HTML(string(blackfriday.Run([]byte(p.Content))))
			c.HTML(http.StatusOK, "post.tmpl.html", gin.H{
				"title":   p.Title,
				"content": postHTML,
			})
		}
	}
}

func PostsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		year := c.Query("year")
		var titles []string
		if year == "" {
			titles = p.GetAllPostTitles()
		} else {
			titles = p.GetAllPostTitlesOfYear(year)
		}
		c.HTML(http.StatusOK, "posts.tmpl.html", gin.H{
			"titles": titles,
		})
	}
}

func SearchPostHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		q := c.Query("q")
		titles := post.SearchPosts(q)
		c.HTML(http.StatusOK, "posts.tmpl.html", gin.H{
			"titles": titles,
		})
	}
}

func EditPostHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		utils.Check(err, "Parsing ID")
		p, err := post.GetByID(uint(id))
		utils.Check(err, "Getting post")
		dateStr := p.Date.Format("2006-01-02T15:04")
		c.HTML(http.StatusOK, "edit.tmpl.html", gin.H{
			"id":      id,
			"title":   p.Title,
			"content": p.Content,
			"date":    dateStr,
		})
	}
}
