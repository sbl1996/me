package post

import (
	"strconv"
	"time"

	"github.com/jinzhu/now"
	"github.com/sbl1996/me/database"
	"github.com/sbl1996/me/utils"
)

type Model struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" sql:"index"`
}

type Post struct {
	Model
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Date    time.Time `json:"date"`
}

type PostItem struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
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

func GetAllPostItems() ([]PostItem, error) {
	db := database.DB
	var posts []Post
	err := db.Order("date desc").Select("id, title").Find(&posts).Error

	var items []PostItem
	for _, p := range posts {
		items = append(items, PostItem{ID: p.ID, Title: p.Title})
	}
	return items, err
}

func GetPostItemsByYear(year string) ([]PostItem, error) {
	db := database.DB
	var posts []Post
	t, err := time.Parse("2006", year)
	utils.Check(err, "Parse year")
	start := now.With(t).BeginningOfYear()
	end := now.With(t).EndOfYear()
	err = db.Where("date between ? AND ?", start, end).Order("date desc").Select([]string{"id", "title"}).Find(&posts).Error

	var items []PostItem
	for _, p := range posts {
		items = append(items, PostItem{ID: p.ID, Title: p.Title})
	}
	return items, err
}

func CreatePost(title string, content string, date time.Time) (uint, error) {
	db := database.DB
	post := Post{
		Title:   title,
		Content: content,
		Date:    date,
	}
	err := db.Create(&post).Error
	return post.ID, err
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

func SearchPosts(query string) ([]PostItem, error) {
	db := database.DB
	var posts []Post
	err := db.Select([]string{"id", "title"}).Where("to_tsvector(title || ' ' || content) @@ to_tsquery(?)", query).Find(&posts).Error

	var items []PostItem
	for _, p := range posts {
		items = append(items, PostItem{ID: p.ID, Title: p.Title})
	}
	return items, err
}

func DeletePost(id uint) error {
	db := database.DB
	p := Post{Model: Model{ID: id}}
	err := db.Delete(&p).Error
	return err
}

func ParseID(id string) (uint, error) {
	uid, err := strconv.ParseUint(id, 10, 64)
	return uint(uid), err
}
