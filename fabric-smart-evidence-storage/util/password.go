package util

import "golang.org/x/crypto/bcrypt"

const PassWordCost = 12

func PasswordEncrypt(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), PassWordCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func PasswordMatch(password string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
