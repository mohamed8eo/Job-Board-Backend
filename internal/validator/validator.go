package validator

import (
	"errors"
	"unicode"
)

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool

	for _, ch := range password {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsDigit(ch):
			hasDigit = true
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !hasDigit {
		return errors.New("password must contain at least one number")
	}
	if !hasSpecial {
		return errors.New("password must contain at least one special character")
	}

	return nil
}

func ValidateName(name string) error {
	if len(name) < 2 {
		return errors.New("name must be at least 2 characters")
	}
	if len(name) > 50 {
		return errors.New("name must be at most 50 characters")
	}
	return nil
}

func ValidateEmail(email string) error {
	if len(email) < 5 {
		return errors.New("invalid email")
	}
	hasAt := false
	for _, ch := range email {
		if ch == '@' {
			hasAt = true
			break
		}
	}
	if !hasAt {
		return errors.New("invalid email")
	}
	return nil
}
