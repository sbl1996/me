package post

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/jinzhu/now"
	"github.com/sbl1996/me/database"
	"github.com/sbl1996/me/utils"
)

type Post struct {
	gorm.Model
	Title   string
	Content string
	Date    time.Time
}

func GetByTitle(title string) (Post, error) {
	db := database.DB
	var p Post
	err := db.Where("title = ?", title).First(&p).Error
	return p, err
}

func GetByID(id uint) (Post, error) {
	db := database.DB
	var p Post
	err := db.Where("id = ?", id).First(&p).Error
	return p, err
}

func GetAllPostTitles() []string {
	db := database.DB
	var posts []Post
	db.Order("date desc").Select("title").Find(&posts)

	var titles []string
	for _, p := range posts {
		titles = append(titles, p.Title)
	}
	return titles
}

func GetAllPostTitlesOfYear(year string) []string {
	db := database.DB
	var posts []Post
	t, err := time.Parse("2006", year)
	utils.Check(err, "Parse year")
	start := now.With(t).BeginningOfYear()
	end := now.With(t).EndOfYear()
	db.Where("date between ? AND ?", start, end).Order("date desc").Select("title").Find(&posts)

	var titles []string
	for _, p := range posts {
		titles = append(titles, p.Title)
	}
	return titles
}

func CreatePost(title string, content string, date time.Time) error {
	db := database.DB
	return db.Create(&Post{
		Title:   title,
		Content: content,
		Date:    date,
	}).Error
}

func UpdatePost(id uint, title string, content string, date time.Time) error {
	db := database.DB
	var p Post
	err := db.First(&p, id).Error
	if err != nil {
		return err
	} else {
		db.Model(&p).Update(Post{Title: title, Content: content, Date: date})
		return nil
	}
}

func SearchPosts(query string) []string {
	db := database.DB
	var posts []Post
	db.Select("title").Where("to_tsvector(title || ' ' || content) @@ to_tsquery(?)", query).Find(&posts)

	var titles []string
	for _, p := range posts {
		titles = append(titles, p.Title)
	}
	return titles
}
