package models

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"

	"github.com/shuaibu222/go-bookstore/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserProfile struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Username string             `bson:"username" json:"username"`
	Email    string             `bson:"email" json:"email"`
	FullName string             `bson:"full_name" json:"full_name"`
	Password string             `bson:"password" json:"password"`
	Bio      string             `bson:"bio" json:"bio"`
}

const mongoURL = "mongodb://bookstore_db:27017"

var UserColl *mongo.Collection
var BookColl *mongo.Collection

// connect with MongoDB
func init() {
	cred, err := config.LoadConfig()
	if err != nil {
		log.Println("failed to load config file:", err)
	}

	credential := options.Credential{
		Username: cred.MongoUsername,
		Password: cred.MongoPassword,
	}
	clientOpts := options.Client().ApplyURI(mongoURL).SetAuth(credential)
	client, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		log.Println("Error connecting to MongoDB")
		return
	}

	UserColl = client.Database("bookstore").Collection("users")
	BookColl = client.Database("bookstore").Collection("books")

	// collection instance
	log.Println("Collections instance is ready")
}

func (u *UserProfile) Validate() error {
	// Check if the username is not empty
	if u.Username == "" {
		return errors.New("username cannot be empty")
	}

	// Check if the email is a valid email address using a simple regex pattern
	emailPattern := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailPattern.MatchString(u.Email) {
		return errors.New("invalid email address")
	}

	// Check if the Full Name is not empty
	if u.FullName == "" {
		return errors.New("full Name cannot be empty")
	}

	return nil
}

func (u *UserProfile) CreateUser() (*mongo.InsertOneResult, error) {
	inserted, err := UserColl.InsertOne(context.Background(), u)
	if err != nil {
		log.Println(err)
	}
	return inserted, nil
}

func GetAllUsers() []primitive.M {
	cur, err := UserColl.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Println(err)
	}

	var users []primitive.M
	for cur.Next(context.Background()) {
		var user bson.M
		err := cur.Decode(&user)
		if err != nil {
			log.Println(err)
		}
		users = append(users, user)
	}

	defer cur.Close(context.Background())
	return users
}

func GetUserById(id string) (UserProfile, error) {
	var user UserProfile

	Id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Failed to convert id", err)
	}

	filter := bson.M{"_id": Id}
	err = UserColl.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		log.Println("Failed to get user", err)
	}

	return user, nil
}

// for authentication purposes only
func GetUserByUsername(username string) UserProfile {
	var user UserProfile
	filter := bson.M{"username": username}
	UserColl.FindOne(context.Background(), filter).Decode(&user)

	return user
}

func DeleteUser(id string) (*mongo.DeleteResult, *mongo.DeleteResult, error) {
	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Failed to convert id", err)
		return nil, nil, err
	}

	str := fmt.Sprint(ID)

	userDeleted, err := UserColl.DeleteOne(context.Background(), bson.M{"_id": ID})
	if err != nil {
		log.Println("Error deleting user", err)
		return nil, nil, err
	}

	booksDeleted, err := BookColl.DeleteMany(context.Background(), bson.M{"user.user_id": str})
	if err != nil {
		log.Println("Error deleting book", err)
		return nil, nil, err
	}

	return userDeleted, booksDeleted, nil
}
