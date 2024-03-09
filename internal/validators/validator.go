package validators

import (
	"regexp"
	"time"
)

const (
	emailRgx    = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	passwordRgx = `^[a-zA-Z0-9!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]*$`
)

func EmailIsValid(email string) (bool, error) {
	if len(email) < 3 || len(email) > 254 {
		return false, nil
	}

	return regexp.MatchString(emailRgx, email)
}

func PasswordIsValid(password string) (bool, error) {
	if len(password) < 8 || len(password) > 64 {
		return false, nil
	}

	return regexp.MatchString(passwordRgx, password)
}

func PasswordsMatch(password, passConf string) bool {
	return password == passConf
}

func BDayValidation(date time.Time) bool {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	return !date.After(today)

}

func PhoneE164Validation(phone string) (bool, error) {
	if len(phone) < 10 || len(phone) > 16 {
		return false, nil
	}

	match, _ := regexp.MatchString(`^\+\d{1,15}$`, phone)
	return match, nil
}
