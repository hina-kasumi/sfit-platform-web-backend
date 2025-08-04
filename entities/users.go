package entities

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Users struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Username  string    `gorm:"unique;not null"`
	Email     string    `gorm:"unique;not null"`
	Password  string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func NewUser(username, email, password string) *Users {
	passwordHash, _ := hashPassword(password)

	return &Users{
		ID:        uuid.New(),
		Username:  username,
		Email:     email,
		Password:  passwordHash,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (u *Users) IsValidPasswrod(password string) error {
	bytePassword := []byte(password)
	byteHashedPassword := []byte(u.Password)
	return bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)
}

func (user *Users) SetPassword(password string) error {
	passwordHash, err := hashPassword(password)
	if err != nil {
		return err
	}
	user.Password = passwordHash
	return nil
}

func hashPassword(password string) (string, error) {
	bytePassword := []byte(password)
	passwordHash, _ := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	return string(passwordHash), nil
}
