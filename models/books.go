package models

import (
	"github.com/shuaibu222/go-bookstore/config"
	"gorm.io/gorm"
)

var Mydb *gorm.DB

type Books struct {
	gorm.Model
	Title              string `json:"title"`
	Description        string `json:"description"`
	AuthorName         string `json:"author_name"`
	AuthorBio          string `json:"author_bio"`
	PublishDate        string `json:"publish_date"`
	Genre              string `json:"genre"`
	Privacy            bool   `json:"privacy"`
	UploadedBook       string `json:"uploaded_book"`
	UploadedCoverImage string `json:"uploaded_cover_image"`
	User                      // anonymous user field
}

type User struct { // anonymous struct
	UserId   string `json:"user_id"`
	Username string `json:"username"`
}

func init() {
	config.Connect()
	Mydb = config.GetDb()
	config.GetDb().AutoMigrate(&Books{})
}

func (b *Books) CreateBook() *Books {
	Mydb.Create(&b)
	return b
}

func GetAllBooks(id string) []Books {
	var books []Books
	Mydb.Where("user_id=?", id).Find(&books)
	return books
}

func GetPublicBooks() []Books {
	var books []Books
	Mydb.Where("privacy=?", false).Find(&books)
	return books
}

func GetBookById(id int64) (*Books, *gorm.DB) {
	var getBook Books
	db := Mydb.Where("ID=?", id).Find(&getBook)
	return &getBook, db
}

func DeleteBook(Id int64) Books {
	var book Books
	Mydb.Where("ID=?", Id).Delete(&book)
	return book
}
