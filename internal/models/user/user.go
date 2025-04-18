package user

import (
	"encoding/hex"
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// CreateUserRequest represents the request body for user registration
type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// User represents a user in the system
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Username  string         `gorm:"uniqueIndex;not null" json:"username"`
	Email     string         `gorm:"uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"not null" json:"-"`
}

// HashPassword hashes the password using bcrypt
func (u *User) HashPassword() error {
	log.Printf("Hashing password for user (length: %d)", len(u.Password))
	if u.Password == "" {
		return errors.New("password cannot be empty")
	}
	
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return err
	}
	
	u.Password = string(hashedPassword)
	log.Printf("Password hashed successfully (hash length: %d, hex: %s)", 
		len(u.Password), 
		hex.EncodeToString([]byte(u.Password)))
	return nil
}

// CheckPassword checks if the provided password matches the hashed password
func (u *User) CheckPassword(password string) bool {
	log.Printf("Checking password for user ID: %d", u.ID)
	
	if u.Password == "" {
		log.Printf("Error: Stored password hash is empty for user ID: %d", u.ID)
		return false
	}
	
	if password == "" {
		log.Printf("Error: Input password is empty for user ID: %d", u.ID)
		return false
	}
	
	log.Printf("Stored hash (hex): %s", hex.EncodeToString([]byte(u.Password)))
	log.Printf("Input password (hex): %s", hex.EncodeToString([]byte(password)))
	
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		log.Printf("Password comparison failed for user ID %d: %v", u.ID, err)
		return false
	}
	
	log.Printf("Password comparison successful for user ID: %d", u.ID)
	return true
} 