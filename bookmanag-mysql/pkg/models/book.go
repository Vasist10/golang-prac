package models

import (
	"github.com/jinzhu/gorm"
	"bookmanag-mysql/pkg/config"
)
var db *gorm.DB

type Book struct {
	gorm.Model
	Name   string `json:"name"`
	Author string `json:"author"`
	Publication string `json:"publication"`
}

func init() {
	config.Connect()
	db = config.GetDB()
	db.AutoMigrate(&Book{})
}

func(b *Book) CreateBook() *Book {
	db.NewRecord(b)
	db.Create(&b)
	return b
}

func GetAllBooks() []Book {
	var Books []Book
	db.Find(&Books)
	return Books
}
func GetBookById(bookId int64) (*Book, *gorm.DB) {
	var getBook Book
	db := db.Where("ID=?", bookId).Find(&getBook)
	return &getBook, db

}

func DeleteBook(bookId int64) Book {
	var book Book
	db.Where("ID=?", bookId).Delete(book)
	return book
}