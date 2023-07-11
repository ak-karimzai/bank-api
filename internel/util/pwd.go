package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPasswrod(pwd string) (string, error) {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password")
	}
	return string(hashedPwd), nil
}

func CheckPwd(pwd, hashedPwd string) error {
	return bcrypt.
		CompareHashAndPassword([]byte(hashedPwd), []byte(pwd))
}
