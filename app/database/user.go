package database

import (
	"crypto/sha1"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const HashKye = "chuck norris"

type User struct {
	ID           int    `gorm:"primary_key" json:"id"`
	Email        string `json:"email" valid:"email~В поле 'емайл' нужно указать настоящий емайл,required~'Емайл' обязательное поле"`
	PasswordHash string
	Password     string `sql:"-" json:"password" valid:"required~'Пароль' обязательное поле"`
	Token        string `json:"token"`

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

func (User) TableName() string {
	return "user"
}

func (u *User) BeforeSave() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	if err != nil {
		return err
	}

	u.PasswordHash = string(bytes)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

func (u *User) GenerateToken() (string, error) {
	h := sha1.New()
	h.Write([]byte(u.Email))
	h.Write([]byte(u.Password))
	h.Write([]byte(HashKye))

	bs := h.Sum(nil)

	return fmt.Sprintf("%x", bs), nil
}
