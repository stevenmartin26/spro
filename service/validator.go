package service

import (
	"strings"

	"github.com/SawitProRecruitment/UserService/common"
)

func validateFullName(fullName string) *common.CustomError {
	errDetails := []string{}

	if len(fullName) < 3 || len(fullName) > 60 {
		errDetails = append(errDetails, "full name must be between 3 and 60 characters")
	}

	if len(errDetails) != 0 {
		return common.NewCustomError(common.ErrInvalidInput, "invalid request params", errDetails...)
	}
	return nil
}

func validatePhoneNumber(phoneNumber string) *common.CustomError {
	errDetails := []string{}

	if len(phoneNumber) < 10 || len(phoneNumber) > 13 {
		errDetails = append(errDetails, "phone number must be between 10 and 13 characters")
	}

	if !strings.HasPrefix(phoneNumber, "+62") {
		errDetails = append(errDetails, "phone number must start with +62")
	}

	if !containsOnlyNumber(phoneNumber[1:]) { // skip the +
		errDetails = append(errDetails, "phone number must only contain number")
	}

	if len(errDetails) != 0 {
		return common.NewCustomError(common.ErrInvalidInput, "invalid request params", errDetails...)
	}
	return nil
}

func validatePassword(password string) *common.CustomError {
	errDetails := []string{}

	if len(password) < 6 || len(password) > 64 {
		errDetails = append(errDetails, "password must be between 6 and 64 characters")
	}

	if !hasCapitalLetter(password) || !hasNumber(password) || !hasSpecialCharacter(password) {
		errDetails = append(errDetails, "password must contain at least 1 capital letter, 1 number and 1 special character")
	}

	if len(errDetails) != 0 {
		return common.NewCustomError(common.ErrInvalidInput, "invalid request params", errDetails...)
	}
	return nil
}

func hasSpecialCharacter(s string) bool {
	specialChars := "!@#$%^&*()_-+{}[];:,.<>"

	for _, c := range s {
		if strings.ContainsRune(specialChars, c) {
			return true
		}
	}
	return false
}

func hasNumber(s string) bool {
	for _, c := range s {
		if c >= '0' && c <= '9' {
			return true
		}
	}
	return false
}

func hasCapitalLetter(s string) bool {
	for _, c := range s {
		if c >= 'A' && c <= 'Z' {
			return true
		}
	}
	return false
}

func containsOnlyNumber(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
