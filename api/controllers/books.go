package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/shuaibu222/go-bookstore/models"
	"github.com/shuaibu222/go-bookstore/utils"
)

func CreateNewBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	id, userName := utils.JwtUserIdUsername(w, r)

	CreateBook := &models.Books{}

	// anonymous field to insert user's Id immediately before creating a book instance
	CreateBook.User.UserId = id
	CreateBook.User.Username = userName

	// parse the book instance
	json.NewDecoder(r.Body).Decode(&CreateBook)

	books := models.GetAllBooks(id)

	for _, book := range books {
		if CreateBook.Title == book.Title && CreateBook.AuthorName == book.AuthorName {
			json.NewEncoder(w).Encode("This book already exists. No duplicate books")
			return
		}
	}

	// create a new book instance
	book := CreateBook.CreateBook()
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

	id, _ := utils.JwtUserIdUsername(w, r)

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
	ID, err := strconv.ParseInt(params["id"], 0, 0) // convert to int
	if err != nil {
		log.Println("Error while parsing!")
	}

	id, _ := utils.JwtUserIdUsername(w, r)
	founded, _ := models.GetBookById(ID)

	if founded.Privacy && founded.UserId == id {
		res, _ := json.Marshal(founded)
		w.Write(res)
	} else {
		json.NewEncoder(w).Encode("You are not authorized to view this book!")
	}
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var updateBook = &models.Books{} // initialize a empty book struct to hold the updated values
	utils.ParseBody(r, updateBook)   // parse the body taken from client request for golang to understand and forward it to the DB
	params := mux.Vars(r)
	bookId := params["id"]
	ID, err := strconv.ParseInt(bookId, 0, 0) // convert to int
	if err != nil {
		log.Println("Error while parsing!")
	}

	id, _ := utils.JwtUserIdUsername(w, r)
	bookDetails, db := models.GetBookById(ID)

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

		if err := db.Save(&bookDetails).Error; err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatal("Error updating book: ", err)
		}
		res, _ := json.Marshal(bookDetails)
		w.Write(res)
	} else {
		json.NewEncoder(w).Encode("You are not authorized to edit this book!")
	}
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	params := mux.Vars(r)
	bookId := params["id"]
	ID, err := strconv.ParseInt(bookId, 0, 0)
	if err != nil {
		log.Println("Error while parsing!: ", err)
	}

	id, _ := utils.JwtUserIdUsername(w, r)
	bookUserId, _ := models.GetBookById(ID)

	if bookUserId.UserId == id {
		book := models.DeleteBook(ID)
		res, _ := json.Marshal(book)
		w.Write(res)
	} else {
		json.NewEncoder(w).Encode("You are not authorized to delete this book!")
	}
}
