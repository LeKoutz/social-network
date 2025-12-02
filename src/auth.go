package forum

import (
	"crypto/rand"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func SaltGenerator(length int) ([]byte, error) {
	bytes := make([]byte, length)
	rand.Read(bytes)
	fmt.Printf("%x\n", bytes)
	return bytes, nil
}

func CompareRegistrationPasswords(pass1, pass2 string) bool {
	return pass1 == pass2
}

func IsUniqueUsername(username string) bool {
	// Check against all usernames?!
	return false
}

func IsUniqueEmail(email string) bool {
	// Check against all emails?!
	return false
}

func IsEmailRegistered(email string) bool {
	// Check against all emails?!
	return false
}

func Auth(email, password string) error {
	// I guess in order to authenticate against an email and a password, we will
	// need to first check if email is registered
	if !IsEmailRegistered(email) {
		// Return an error... this should be sent back to umh... places...?!
		return ErrorNotRegistered
	}
	// Since we passed the initial test, we can fetch the user's salt
	salt := GetUserSalt(email)
	hashStored := GetUserHash(email)
	// Salt password
	saltedPassword := salt+password
	hash, err := bcrypt.GenerateFromPassword([]byte(saltedPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	if string(hash) != hashStored {
		return ErrorWrongPassword
	}
	return nil
}
