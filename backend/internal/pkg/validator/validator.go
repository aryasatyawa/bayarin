package validator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	validate   *validator.Validate
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	phoneRegex = regexp.MustCompile(`^(\+62|62|0)[0-9]{9,12}$`)
)

func init() {
	validate = validator.New()

	// Register custom validators
	_ = validate.RegisterValidation("indonesian_phone", validateIndonesianPhone)
}

func New() *validator.Validate {
	return validate
}

// ValidateStruct validates a struct
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

// ValidateEmail validates email format
func ValidateEmail(email string) error {
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

// ValidatePhone validates Indonesian phone number
func ValidatePhone(phone string) error {
	if !phoneRegex.MatchString(phone) {
		return fmt.Errorf("invalid phone number format")
	}
	return nil
}

// ValidatePIN validates 6-digit PIN
func ValidatePIN(pin string) error {
	if len(pin) != 6 {
		return fmt.Errorf("PIN must be 6 digits")
	}

	for _, char := range pin {
		if char < '0' || char > '9' {
			return fmt.Errorf("PIN must contain only digits")
		}
	}

	return nil
}

// ValidatePassword validates password strength
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}
	return nil
}

// ValidateAmount validates transaction amount (must be positive)
func ValidateAmount(amount int64) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}
	return nil
}

// Custom validator for Indonesian phone
func validateIndonesianPhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	return phoneRegex.MatchString(phone)
}

// NormalizePhone normalizes Indonesian phone number to E.164 format
func NormalizePhone(phone string) string {
	phone = strings.TrimSpace(phone)

	// Remove all non-digit characters except +
	cleaned := ""
	for _, char := range phone {
		if (char >= '0' && char <= '9') || char == '+' {
			cleaned += string(char)
		}
	}

	// Convert to +62 format
	if strings.HasPrefix(cleaned, "0") {
		return "+62" + cleaned[1:]
	} else if strings.HasPrefix(cleaned, "62") {
		return "+" + cleaned
	} else if strings.HasPrefix(cleaned, "+62") {
		return cleaned
	}

	return "+62" + cleaned
}
