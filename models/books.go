package models

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Book struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	AuthorName  string             `bson:"author_name" json:"author_name"`
	AuthorBio   string             `bson:"author_bio" json:"author_bio"`
	PublishDate string             `bson:"publish_date" json:"publish_date"`
	Genre       string             `bson:"genre" json:"genre"`
	Privacy     bool               `bson:"privacy" json:"privacy"`
	User                           // anonymous user field
}

type User struct { // anonymous struct
	UserId   string `bson:"user_id" json:"user_id"`
	UserName string `bson:"user_name" json:"user_name"`
}

func (b *Book) CreateBook() (*mongo.InsertOneResult, error) {
	inserted, err := BookColl.InsertOne(context.Background(), b)
	if err != nil {
		log.Println("Failed creating book", err)
	}
	return inserted, nil
}

func GetAllBooks(id string) []primitive.M {

	cur, err := BookColl.Find(context.Background(), bson.M{"user.user_id": id})
	if err != nil {
		log.Println(err)
	}

	var books []primitive.M
	for cur.Next(context.Background()) {
		var book bson.M
		err := cur.Decode(&book)
		if err != nil {
			log.Println(err)
		}
		books = append(books, book)
	}

	defer cur.Close(context.Background())
	return books
}

func GetPublicBooks() []primitive.M {
	cur, err := BookColl.Find(context.Background(), bson.M{"privacy": false})
	if err != nil {
		log.Println(err)
	}

	var books []primitive.M
	for cur.Next(context.Background()) {
		var book bson.M
		err := cur.Decode(&book)
		if err != nil {
			log.Println(err)
		}
		books = append(books, book)
	}

	defer cur.Close(context.Background())
	return books
}

func GetBookById(id string) (Book, error) {
	var book Book

	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Failed to convert id", err)
	}

	filter := bson.M{"_id": ID}
	err = BookColl.FindOne(context.Background(), filter).Decode(&book)
	if err != nil {
		log.Println("Failed to get book", err)
	}

	return book, nil
}

func DeleteBook(id string) (*mongo.DeleteResult, error) {
	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Failed to convert id", err)
	}

	bookFilter := bson.M{"_id": ID}
	deleted, err := BookColl.DeleteOne(context.Background(), bookFilter)
	if err != nil {
		log.Println("Error deleting book", err)
	}

	return deleted, nil
}
