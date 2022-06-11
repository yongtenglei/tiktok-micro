package scryptx

import (
	"fmt"
	"yunyandz.com/tiktok/logger"

	"golang.org/x/crypto/scrypt"
)

const (
	Salt = "66666"
)

func PasswordEncrypt(password string) string {
	dk, err := scrypt.Key([]byte(password),
		[]byte(Salt),
		32768, 8, 1, 32)
	if err != nil {
		logger.Suger().Errorw("PasswordEncrypt scrypt.Key failed",
			"err", err.Error())

	}
	return fmt.Sprintf("%x", string(dk))
}

func PasswordValidate(password, expectPassword string) bool {
	scryptedPassword := PasswordEncrypt(password)

	if scryptedPassword == expectPassword {
		return true
	}

	return false
}
