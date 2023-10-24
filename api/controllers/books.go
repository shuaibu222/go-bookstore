package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shuaibu222/go-bookstore/models"
	"github.com/shuaibu222/go-bookstore/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateNewBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := utils.JwtUserIdUsername(w, r)

	CreateBook := &models.Book{}

	// anonymous field to insert user's Id immediately before creating a book instance
	CreateBook.User.UserId = id

	// parse the book instance
	json.NewDecoder(r.Body).Decode(&CreateBook)

	books := models.GetAllBooks(id)

	for _, book := range books {
		bookTitle, isBookExists := book["title"]
		bookAuthor, isAuthorExists := book["author_name"]

		if isBookExists && isAuthorExists && bookTitle == CreateBook.Title && bookAuthor == CreateBook.AuthorName {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode("This book already exists. No duplicate books")
			return
		}
	}

	// create a new book instance
	book, err := CreateBook.CreateBook()
	if err != nil {
		log.Println(err)
	}
	res, err := json.Marshal(book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(res)

}

func GetAllUserBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	id := utils.JwtUserIdUsername(w, r)

	books := models.GetAllBooks(id)

	res, err := json.Marshal(books)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(res)

}

func GetAllPublicBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	books := models.GetPublicBooks()

	res, err := json.Marshal(books)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(res)

}

func GetBookById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	params := mux.Vars(r)

	id := utils.JwtUserIdUsername(w, r)
	founded, err := models.GetBookById(params["id"])
	if err != nil {
		log.Println(err)
	}

	if founded.UserId == id {
		res, _ := json.Marshal(founded)
		w.Write(res)
	} else {
		json.NewEncoder(w).Encode("You are not authorized to view this book!")
	}
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var updateBook = &models.Book{} // initialize an empty book struct to hold the updated values
	utils.ParseBody(r, updateBook)  // parse the body taken from client request for golang to understand and forward it to the DB

	params := mux.Vars(r)

	// get that specific book for updating from URL params
	id := utils.JwtUserIdUsername(w, r)
	bookDetails, err := models.GetBookById(params["id"])
	if err != nil {
		log.Println(err)
	}

	// get all his books and check for duplicates even when editing
	books := models.GetAllBooks(id)

	for _, book := range books {
		bookTitle, isBookExists := book["title"]
		bookAuthor, isAuthorExists := book["author_name"]

		if isBookExists && isAuthorExists && bookTitle == updateBook.Title && bookAuthor == updateBook.AuthorName {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode("This book already exists. No duplicate books")
			return
		}
	}

	if bookDetails.UserId == id {

		if updateBook.Title != "" {
			bookDetails.Title = updateBook.Title
		}
		if updateBook.Description != "" {
			bookDetails.Description = updateBook.Description
		}
		if updateBook.AuthorName != "" {
			bookDetails.AuthorName = updateBook.AuthorName
		}
		if updateBook.AuthorBio != "" {
			bookDetails.AuthorBio = updateBook.AuthorBio
		}
		if updateBook.PublishDate != "" {
			bookDetails.PublishDate = updateBook.PublishDate
		}
		if updateBook.Genre != "" {
			bookDetails.Genre = updateBook.Genre
		}

		updateBook.User.UserId = id

		Id, err := primitive.ObjectIDFromHex(params["id"])
		if err != nil {
			log.Println(err)
		}

		result, err := models.BookColl.UpdateOne(
			context.Background(),
			bson.M{"_id": Id},
			bson.M{"$set": updateBook},
		)
		if err != nil {
			log.Println("Failed to update book")
		}

		res, _ := json.Marshal(result)
		w.Write(res)
	} else {
		json.NewEncoder(w).Encode("You are not authorized to edit this book!")
	}
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	params := mux.Vars(r)

	id := utils.JwtUserIdUsername(w, r)
	bookUserId, _ := models.GetBookById(params["id"])

	if bookUserId.UserId == id {
		book, err := models.DeleteBook(params["id"])
		if err != nil {
			log.Println(err)
		}

		res, _ := json.Marshal(book)
		w.Write(res)
	} else {
		json.NewEncoder(w).Encode("You are not authorized to delete this book!")
	}
}
