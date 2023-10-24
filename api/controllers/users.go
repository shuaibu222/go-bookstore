package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shuaibu222/go-bookstore/models"
	"github.com/shuaibu222/go-bookstore/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// TODO: still creates user in db even if user already exists

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userProfile := &models.UserProfile{}

	err := json.NewDecoder(r.Body).Decode(&userProfile)
	if err != nil {
		// Handle JSON decoding error
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid JSON data")
		return
	}

	users := models.GetAllUsers()

	// user can only be created once no duplicate
	for _, user := range users {
		userEmail, emailExists := user["email"]
		userName, usernameExists := user["username"]

		if emailExists && usernameExists && userEmail == userProfile.Email && userName == userProfile.Username {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode("This user already exists. Try changing your username and email")
			return
		}
	}

	// hash the user entered password for security purposes
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userProfile.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error generating password: ", err)
	}

	userProfile.Password = string(hashedPassword)

	if userProfile.Validate() != nil {
		w.WriteHeader(http.StatusBadRequest) // 400 Bad Request status code
		json.NewEncoder(w).Encode("Invalid email or username and fullname are empty!")
	} else {
		u, err := userProfile.CreateUser()
		if err != nil {
			log.Println(err)
		}
		json, _ := json.Marshal(u)
		w.Write(json)
	}
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) { // for development purposes
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	users := models.GetAllUsers()

	res, _ := json.Marshal(users)
	w.Write(res)
}

func GetUserByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	params := mux.Vars(r)

	id := utils.JwtUserIdUsername(w, r)

	founded, err := models.GetUserById(params["id"])
	if err != nil {
		log.Println(err)
	}

	if founded.ID.String() == id {
		user, err := json.Marshal(founded)
		if err != nil {
			log.Println("Failed to marshal: ", err)
		}
		w.Write(user)
	} else {
		json.NewEncoder(w).Encode("You are not authorized to view this profile!")
	}
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userUpdate := &models.UserProfile{}
	utils.ParseBody(r, userUpdate)

	params := mux.Vars(r)

	id := utils.JwtUserIdUsername(w, r)
	founded, err := models.GetUserById(params["id"])
	if err != nil {
		log.Println(err)
	}

	if founded.ID.String() == id {

		if userUpdate.Username != "" {
			founded.Username = userUpdate.Username
		}
		if userUpdate.Password != "" {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userUpdate.Password), bcrypt.DefaultCost)
			if err != nil {
				log.Println("Error generating hashed password: ", err)
			}

			userUpdate.Password = string(hashedPassword)
			founded.Password = userUpdate.Password
		}
		if userUpdate.Email != "" {
			founded.Email = userUpdate.Email
		}
		if userUpdate.Bio != "" {
			founded.Bio = userUpdate.Bio
		}
		if userUpdate.FullName != "" {
			founded.FullName = userUpdate.FullName
		}

		Id, err := primitive.ObjectIDFromHex(params["id"])
		if err != nil {
			log.Println(err)
		}

		result, err := models.UserColl.UpdateOne(
			context.Background(),
			bson.M{"_id": Id},
			bson.M{"$set": userUpdate},
		)

		if err != nil {
			log.Println("Failed to update user: ", err)
		}

		res, _ := json.Marshal(result)
		w.Write(res)
	} else {
		json.NewEncoder(w).Encode("You are not authorized to edit this profile!")
	}
}

// TODO delete issue

func DeleteUserById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	params := mux.Vars(r)

	id := utils.JwtUserIdUsername(w, r)
	founded, err := models.GetUserById(params["id"])
	if err != nil {
		log.Println(err)
	}

	if founded.ID.String() == id {
		user, books, err := models.DeleteUser(params["id"])
		if err != nil {
			log.Println(err)
		}
		deleteInfo := fmt.Sprintf("The user %v has been deleted along with his bookstore: %v", user.DeletedCount, books.DeletedCount)
		json.NewEncoder(w).Encode(deleteInfo)
	} else {
		json.NewEncoder(w).Encode("You are not authorized to delete this profile!")
	}
}
