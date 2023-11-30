package validator

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString
	isValidName = regexp.MustCompile(`^[a-zA-Z\\s]+$`).MatchString
)

func ValidateStringLenght(value string, min int, max int) error {
	n := len(value)
	if n < min || n > max {
		return fmt.Errorf("lenght must be between %d and %d", min, max)
	}
	return nil
}

func ValidateUsername(username string) error {
	if err := ValidateStringLenght(username, 3, 100); err != nil {
		return fmt.Errorf("username %v", err)
	}
	if !isValidUsername(username) {
		return fmt.Errorf("username must only contain letters, digits and underscores")
	}
	return nil
}

func ValidateName(name string) error {
	if err := ValidateStringLenght(name, 2, 20); err != nil {
		return fmt.Errorf("name %v", err)
	}
	if !isValidName(name) {
		return fmt.Errorf("name must only contain letters or spaces")
	}
	return nil
}

func ValidatePassword(password string) error {
	if err := ValidateStringLenght(password, 8, 255); err != nil {
		return fmt.Errorf("password %v", err)
	}
	return nil
}

func ValidateEmail(email string) error {
	if err := ValidateStringLenght(email, 5, 255); err != nil {
		return fmt.Errorf("email %v", err)
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return fmt.Errorf("email provided is not valid %v", err)
	}
	return nil
}