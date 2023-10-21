package models

import (
	"errors"
	"regexp"

	"github.com/shuaibu222/go-bookstore/config"
	"gorm.io/gorm"
)

var mydb *gorm.DB

type UserProfile struct {
	gorm.Model
	Username string `json:"username"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}

func init() {
	config.Connect()
	mydb = config.GetDb()
	config.GetDb().AutoMigrate(&UserProfile{})
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

func (u *UserProfile) CreateUser() *UserProfile {
	mydb.Create(&u)
	return u
}

func GetAllUsers() []UserProfile {
	var users []UserProfile
	mydb.Find(&users)
	return users
}

func GetUserById(id int64) (*UserProfile, *gorm.DB) {
	var user UserProfile
	db := mydb.Where("ID=?", id).Find(&user)
	return &user, db
}

// for authentication purposes only
func GetUserByUsername(username string) UserProfile {
	var user UserProfile
	mydb.Where("username=?", username).Find(&user)
	return user
}

func DeleteUser(id int64) UserProfile {
	var user UserProfile
	var books []Books
	mydb.Where("ID=?", id).Delete(&user)
	Mydb.Where("user_id=?", id).Delete(&books)
	return user
}
