package validators

import (
	"assetra/pb"
	"errors"
	"regexp"
	"strings"
	"unicode"
)

var (
	// ErrInvalidEmail indicates that the provided email is not valid.
	ErrInvalidUserId      = errors.New("invalid user ID format")
	ErrInvalidEmail       = errors.New("invalid email format")
	ErrInvalidName        = errors.New("invalid name format")
	ErrInvalidPassword    = errors.New("invalid password format")
	ErrInvalidCredentials = errors.New("cannot update user with invalid credentials")
	ErrPasswordTooWeak    = errors.New("password must be at least 8 characters and contain uppercase, lowercase, number and special character")
	ErrNameTooShort       = errors.New("name must be at least 2 characters")
	ErrNameTooLong        = errors.New("name must be at most 50 characters")
	ErrEmailTooLong       = errors.New("email must be at most 100 characters")
)

// Email validation regex - more comprehensive
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// ValidateSignUp validates the sign-up data.
func ValidateSignUp(user *pb.User) error {
	if err := ValidateEmail(user.Email); err != nil {
		return err
	}
	if err := ValidateName(user.Name); err != nil {
		return err
	}
	if err := ValidatePassword(user.Password); err != nil {
		return err
	}
	return nil
}

func ValidateUpdateUser(user *pb.User) error {
	// At least one field must be provided
	if user.Email == "" && user.Name == "" && user.Password == "" {
		return ErrInvalidCredentials
	}
	
	// Validate each field if provided
	if user.Email != "" {
		if err := ValidateEmail(user.Email); err != nil {
			return err
		}
	}
	if user.Name != "" {
		if err := ValidateName(user.Name); err != nil {
			return err
		}
	}
	if user.Password != "" {
		if err := ValidatePassword(user.Password); err != nil {
			return err
		}
	}
	return nil
}

// ValidateEmail validates email format and length
func ValidateEmail(email string) error {
	if email == "" {
		return ErrInvalidEmail
	}
	if len(email) > 100 {
		return ErrEmailTooLong
	}
	if !emailRegex.MatchString(email) {
		return ErrInvalidEmail
	}
	return nil
}

// ValidateName validates name format and length
func ValidateName(name string) error {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return ErrInvalidName
	}
	if len(trimmed) < 2 {
		return ErrNameTooShort
	}
	if len(trimmed) > 50 {
		return ErrNameTooLong
	}
	// Check for suspicious characters that might indicate XSS attempts
	if strings.ContainsAny(trimmed, "<>&\"'") {
		return errors.New("name contains invalid characters")
	}
	return nil
}

// ValidatePassword enforces strong password policy
func ValidatePassword(password string) error {
	if password == "" {
		return ErrInvalidPassword
	}
	if len(password) < 8 {
		return ErrPasswordTooWeak
	}
	if len(password) > 128 {
		return errors.New("password is too long")
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		return ErrPasswordTooWeak
	}

	return nil
}

// NormalizeEmail normalizes the email address.
func NormalizeEmail(email string) string {
	return strings.TrimSpace(strings.ToLower(email))
}

// SanitizeName removes potentially dangerous characters from name
func SanitizeName(name string) string {
	// Remove any HTML/script tags or dangerous characters
	name = strings.TrimSpace(name)
	name = strings.ReplaceAll(name, "<", "")
	name = strings.ReplaceAll(name, ">", "")
	name = strings.ReplaceAll(name, "&", "")
	name = strings.ReplaceAll(name, "\"", "")
	name = strings.ReplaceAll(name, "'", "")
	return name
}
