package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/shuaibu222/go-bookstore/models"
	"github.com/shuaibu222/go-bookstore/utils"
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
		if user.Email == userProfile.Email && user.Username == userProfile.Username {
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
		u := userProfile.CreateUser()
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
	ID, err := strconv.ParseInt(params["id"], 0, 0)
	if err != nil {
		log.Println("Failed to parse: ", err)
	}

	id, _ := utils.JwtUserIdUsername(w, r)
	founded, _ := models.GetUserById(ID)

	intId, err := strconv.ParseInt(id, 0, 0)
	if err != nil {
		log.Println("Failed to parse: ", err)
	}

	if founded.ID == uint(intId) {
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

	w.WriteHeader(http.StatusOK)

	userUpdate := &models.UserProfile{}
	utils.ParseBody(r, userUpdate)

	params := mux.Vars(r)
	ID, err := strconv.ParseInt(params["id"], 0, 0)
	if err != nil {
		log.Println("Error while parsing id: ", err)
	}

	id, _ := utils.JwtUserIdUsername(w, r)
	founded, db := models.GetUserById(ID)

	intId, err := strconv.ParseInt(id, 0, 0)
	if err != nil {
		log.Println("Failed to parse: ", err)
	}

	if founded.ID == uint(intId) {
		if userUpdate.Username != "" {
			founded.Username = userUpdate.Username
		}
		if userUpdate.Password != "" {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userUpdate.Password), bcrypt.DefaultCost)
			if err != nil {
				log.Println("Error generating hashed password: ", err)
			}

			founded.Password = string(hashedPassword)
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

		// TODO: update avatar image

		if err := db.Save(&founded).Error; err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatal("Error updating user: ", err)
			return
		}

		res, _ := json.Marshal(founded)
		w.Write(res)
	} else {
		json.NewEncoder(w).Encode("You are not authorized to edit this profile!")
	}
}

func DeleteUserById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	params := mux.Vars(r)
	ID, err := strconv.ParseInt(params["id"], 0, 0)
	if err != nil {
		log.Println("Error while parsing ", err)
	}

	id, _ := utils.JwtUserIdUsername(w, r)
	founded, _ := models.GetUserById(ID)

	intId, err := strconv.ParseInt(id, 0, 0)
	if err != nil {
		log.Println("Failed to parse: ", err)
	}

	if founded.ID == uint(intId) {
		deleted := models.DeleteUser(ID)
		deletedUser, err := json.Marshal(deleted.DeletedAt)
		if err != nil {
			log.Println("Error while marshaling the user: ", err)
		}
		deleteInfo := fmt.Sprintf("The user has been deleted along with his bookstore at: %s", deletedUser)
		w.Write([]byte(deleteInfo))
	} else {
		json.NewEncoder(w).Encode("You are not authorized to delete this profile!")
	}
}
