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
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Date    time.Time `json:"date"`
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

func GetAllPostTitles() ([]string, error) {
	db := database.DB
	var posts []Post
	err := db.Order("date desc").Select("title").Find(&posts).Error

	var titles []string
	for _, p := range posts {
		titles = append(titles, p.Title)
	}
	return titles, err
}

func GetPostTitlesOfYear(year string) ([]string, error) {
	db := database.DB
	var posts []Post
	t, err := time.Parse("2006", year)
	utils.Check(err, "Parse year")
	start := now.With(t).BeginningOfYear()
	end := now.With(t).EndOfYear()
	err = db.Where("date between ? AND ?", start, end).Order("date desc").Select("title").Find(&posts).Error

	var titles []string
	for _, p := range posts {
		titles = append(titles, p.Title)
	}
	return titles, err
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

func SearchPosts(query string) ([]string, error) {
	db := database.DB
	var posts []Post
	err := db.Select("title").Where("to_tsvector(title || ' ' || content) @@ to_tsquery(?)", query).Find(&posts).Error

	var titles []string
	for _, p := range posts {
		titles = append(titles, p.Title)
	}
	return titles, err
}
