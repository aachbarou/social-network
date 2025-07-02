package utils

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"social-network/pkg/models"
)

// validate all fields when user registers
func ValidateNewUser(user models.User) error {
	if err := validateFirstName(user.FirstName); err != nil {
		return err
	}
	if err := validateLastName(user.LastName); err != nil {
		return err
	}
	if err := validateBirth(user.DateOfBirth); err != nil {
		return err
	}
	if err := validatePassword(user.Password); err != nil {
		return err
	}
	if err := validateEmail(user.Email); err != nil {
		fmt.Println("Email validation failed: %w", err)
		return err
	}
	return nil
}

func validateFirstName(name string) error {
	if fieldEmpty(name) {
		return errors.New("Validation error")
	}
	return nil
}

func validateLastName(name string) error {
	if fieldEmpty(name) {
		return errors.New("Validation error")
	}
	return nil
}

func validateBirth(birthDate string) error {
	if fieldEmpty(birthDate) {
		return errors.New("Validation error")
	}
	if err := validateDate(birthDate); err != nil {
		return fmt.Errorf("Invalid date format: %w", err)
	}
	return nil
}

func validatePassword(password string) error {
	if fieldEmpty(password) {
		return errors.New("Validation error")
	}
	return nil
}

func validateEmail(email string) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if fieldEmpty(email) {
		return errors.New("Email cannot be empty")
	}
	if !emailRegex.MatchString(email) {
		return errors.New("Invalid email format")
	}
	return nil
}

func validateDate(dateStr string) error {
	layout := "2006-01-02"
	_, err := time.Parse(layout, dateStr)
	if err != nil {
		return fmt.Errorf("invalid date format: %s", err)
	}
	return nil
}

func fieldEmpty(value string) bool {
	return len(value) == 0
}
