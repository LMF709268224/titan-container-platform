package kubesphere

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

// generatePassword generates a random password of the specified length.
// It returns an error if the length is not between 8 and 64.
func generatePassword(length int) (string, error) {
	if length < 8 || length > 64 {
		return "", fmt.Errorf("密码长度必须在8到64之间")
	}

	// 定义字符集
	numbers := "0123456789"
	lowerCase := "abcdefghijklmnopqrstuvwxyz"
	upperCase := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	special := "~!@#$%^&*()-_=+\\|[{}];:'\",<.>/? "

	password := make([]string, 4)
	var err error

	password[0], err = randomChar(numbers)
	if err != nil {
		return "", err
	}

	password[1], err = randomChar(lowerCase)
	if err != nil {
		return "", err
	}

	password[2], err = randomChar(upperCase)
	if err != nil {
		return "", err
	}

	password[3], err = randomChar(special)
	if err != nil {
		return "", err
	}

	allChars := numbers + lowerCase + upperCase + special

	for i := 4; i < length; i++ {
		char, err := randomChar(allChars)
		if err != nil {
			return "", err
		}
		password = append(password, char)
	}

	shuffled, err := shuffleSlice(password)
	if err != nil {
		return "", err
	}

	return strings.Join(shuffled, ""), nil
}

func randomChar(chars string) (string, error) {
	max := big.NewInt(int64(len(chars)))
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return string(chars[n.Int64()]), nil
}

func shuffleSlice(slice []string) ([]string, error) {
	result := make([]string, len(slice))
	copy(result, slice)

	for i := len(result) - 1; i > 0; i-- {
		max := big.NewInt(int64(i + 1))
		j, err := rand.Int(rand.Reader, max)
		if err != nil {
			return nil, err
		}
		result[i], result[j.Int64()] = result[j.Int64()], result[i]
	}

	return result, nil
}

// // ValidatePassword 验证密码是否符合要求
// func ValidatePassword(password string) bool {
// 	if len(password) < 8 || len(password) > 64 {
// 		return false
// 	}

// 	var (
// 		hasNumber   bool
// 		hasLower    bool
// 		hasUpper    bool
// 		hasSpecial  bool
// 		specialChar = "~!@#$%^&*()-_=+\\|[{}];:'\",<.>/? "
// 	)

// 	for _, char := range password {
// 		switch {
// 		case strings.ContainsRune("0123456789", char):
// 			hasNumber = true
// 		case strings.ContainsRune("abcdefghijklmnopqrstuvwxyz", char):
// 			hasLower = true
// 		case strings.ContainsRune("ABCDEFGHIJKLMNOPQRSTUVWXYZ", char):
// 			hasUpper = true
// 		case strings.ContainsRune(specialChar, char):
// 			hasSpecial = true
// 		}
// 	}

// 	return hasNumber && hasLower && hasUpper && hasSpecial
// }
