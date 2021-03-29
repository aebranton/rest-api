package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// Email regex for validation
var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Phone number regex for validation (not acutally cleaning them, probabably would in production)
var phoneRegex = regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`)

// Service - the user service. Holds a DB connection pointer and has
// methods attached for working with user objects
type Service struct {
	DB *gorm.DB
}

// User - defines the user model/structure
type User struct {
	gorm.Model
	Username  string `gorm:"unique"`
	Password  string
	FirstName string
	LastName  string
	Email     string `gorm:"unique"`
	Telephone string
}

// IsValid - this is called within a BeforeCreate hook on the gorm User model.
// This will check to make sure fields are entered, sized correctly, are valid data, etc.
func (u *User) IsValid() (bool, string) {
	if len(u.FirstName) < 2 || len(u.FirstName) > 255 {
		return false, "Length of FirstName is not between 2-255 characters"
	}
	if len(u.LastName) < 2 || len(u.LastName) > 255 {
		return false, "Length of LastName is not between 2-255 characters"
	}
	if len(u.Password) < 8 || len(u.Password) > 255 {
		return false, "Length of Password is not between 8-255 characters"
	}
	if len(u.Email) < 5 || len(u.Email) > 255 {
		return false, "Length of the email is not between 5-255 characters"
	}
	if !emailRegex.MatchString(u.Email) {
		return false, "Email is not a valid address"
	}
	if len(u.Telephone) < 5 || len(u.Telephone) > 50 {
		return false, "Length of the telephone number is not between 5-50 characters"
	}
	if !phoneRegex.MatchString(u.Telephone) {
		return false, "Telephone number is not a valid number"
	}
	return true, ""
}

// BeforeCreate - User hook before it is created to check if it is valid. This runs the User struct's IsValid method
// If any of the tests fail, an error is returned that will be shown on the page with a 400 Bad Request status
func (u *User) BeforeCreate(tx *gorm.DB) error {
	valid, errorMsg := u.IsValid()
	if !valid {
		err := errors.New(errorMsg)
		return err
	}
	return nil
}

// Users - Custom slice of User objects with some attached methods
// that are necessary to match a User object, ie ToJSON
type Users []User

// ToJSON - converts a Users object (slice of User) to JSON and writes to the given
// responseWriter with a header status of Ok
func (u *Users) ToJSON(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	return encoder.Encode(u)
}

// ToJSON - converts a User to JSON and writes to the given
// responseWriter with a header status of Ok
func (u *User) ToJSON(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	return encoder.Encode(u)
}

// HashPassword - given a string password, ex: "testpassword", converts it to a hash
// for storage in the database.
func HashPassword(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return string(hash), err
}

// ComparePassword - Given a password and a pre-hashed string, see if the given
// password matches the hash supplied. Return nil if they match, error if not.
func ComparePassword(pwd, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))
	return err == nil
}

// UserService - interface for our user service
type UserService interface {
	GetUser(ID uint) (User, error)
	GetUserByUsername(username string) (User, error)
	CreateUser(user User) (User, error)
	UpdateUser(ID uint, updatedUser User) (User, error)
	DeleteUser(ID uint) error
	GetAllUsers() ([]User, error)
}

// NewService - returns a new user service
func NewService(db *gorm.DB) *Service {
	return &Service{
		DB: db,
	}
}

// GetUser - retreives a user by ID from the database
func (s *Service) GetUser(ID uint) (User, error) {
	var user User
	if result := s.DB.First(&user, ID); result.Error != nil {
		return User{}, result.Error
	}
	return user, nil
}

// GetUserByUsername - retreives users by username from the database
func (s *Service) GetUserByUsername(username string) (User, error) {
	var user User
	if result := s.DB.Find(user).Where("username = ?", username); result.Error != nil {
		return User{}, result.Error
	}
	return user, nil
}

// CreateUser - creates a user in the database. Users do have a BeforeCreate hook to validate
// and make sure the data coming in is sufficient, ie the email is valid, phone is valid, username is unique, etc.
// Errors will be returned if anything is invalid
func (s *Service) CreateUser(user User) (User, error) {
	if result := s.DB.Save(&user); result.Error != nil {
		return User{}, result.Error
	}
	return user, nil
}

// UpdateUser - updates a user in the database by ID.
func (s *Service) UpdateUser(ID uint, updatedUser User) (User, error) {
	user, err := s.GetUser(ID)
	if err != nil {
		return User{}, err
	}

	if result := s.DB.Model(&user).Updates(updatedUser); result.Error != nil {
		return User{}, result.Error
	}
	return user, nil
}

// DeleteUser - Deletes a user object from the database
func (s *Service) DeleteUser(ID uint) error {
	if result := s.DB.Delete(&User{}, ID); result.Error != nil {
		return result.Error
	}
	return nil
}

// GetAllUsers - returns all users from the database as a Users object
func (s *Service) GetAllUsers() (Users, error) {
	var users Users
	if result := s.DB.Find(&users); result.Error != nil {
		return []User{}, result.Error
	}
	return users, nil
}
