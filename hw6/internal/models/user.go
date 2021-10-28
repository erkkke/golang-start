package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
)

var userNextID = 0

type User struct {
	ID                int    `json:"id"`
	Email             string `json:"email"`
	PhoneNumber       string `json:"phone_number"`
	Password          string `json:"password"`
	EncryptedPassword string `json:"-"`
	Name              string `json:"name"`
	Surname           string `json:"surname"`
	BirthDate         string `json:"birth_date"`
	City              string `json:"city"`
}

func (u *User) NextId() {
	u.ID = userNextID
	userNextID++
}

func (u *User) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password, validation.By(RequiredIf(u.EncryptedPassword == "")), validation.Length(6, 50)),
	)
}

func (u *User) BeforeCreating() error {
	if len(u.Password) > 0 {
		enc, err := encryptString(u.Password)
		if err != nil {
			return err
		}

		u.EncryptedPassword = enc
	}

	return nil
}

func (u *User) Sanitize() {
	u.Password = ""
}

func (u *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(password)) == nil
}

func encryptString(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
