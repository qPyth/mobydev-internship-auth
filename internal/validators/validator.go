package validators

import "regexp"

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
