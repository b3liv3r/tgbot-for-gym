package telegram

import (
	"errors"
	"github.com/b3liv3r/tgbot-for-gym/modules/gym/models"
	"regexp"
	"strconv"
	"strings"
)

func isValidUsername(username string) error {
	words := strings.Fields(username)

	// Проверка каждого слова на длину и наличие только букв
	for _, word := range words {
		if len(word) < 2 {
			return errors.New("each word in username should contain at least two letters")
		}
		if !isAlpha(word) {
			return errors.New("username can only contain letters")
		}
	}

	return nil
}

func isValidPhone(phone string) error {
	// Проверка длины номера телефона и наличие только цифр
	phoneRegex := regexp.MustCompile(`^\+\d{11}$`)
	if !phoneRegex.MatchString(phone) {
		return errors.New("phone number should start with '+' and consist of 11 digits. Example: +78005553535")
	}
	if phone[1] != '7' {
		if phone[1] == '8' {
			return errors.New("use number entry without 8. example: +78005553535")
		}
		return errors.New("currently only Russian phone numbers are supported")
	}
	return nil
}

func isValidEmail(email string) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("incorrect email address. example: edgarusa22@gmail.com")
	}
	return nil
}

func isAlpha(s string) bool {
	// Проверка, что строка состоит только из букв (латинских или кириллических)
	alphaRegex := regexp.MustCompile(`^[a-zA-Zа-яА-Я]+$`)
	return alphaRegex.MatchString(s)
}

func isDephis(s string) bool {
	if len(s) == 1 && s[0] == '-' {
		return true
	}
	return false
}

func isSubLvl(s string) (int, error) {
	switch s {
	case "0":
		return 0, nil
	case "1":
		return 1, nil
	case "2":
		return 2, nil
	case "3":
		return 3, nil
	default:
		return 0, errors.New("invalid sub level use /help to see info for all commands")
	}
}

func isPositiveNumber(s string) (int, error) {
	num, err := strconv.Atoi(s)
	if err != nil {
		return 0, errors.New("invalid input: not a number")
	}
	if num <= 0 {
		return 0, errors.New("not a positive number")
	}
	return num, nil
}

func isPositiveFloat(s string) (float64, error) {
	num, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, errors.New("invalid input: not a number")
	}
	if num <= 0 {
		return 0, errors.New("not a positive number")
	}
	return num, nil
}

func isFoundGym(arr []models.Gym, gymID int) (int, bool) {
	found := true
	sublvl := 0
	for _, g := range arr {
		if g.Id == gymID {
			found = false
			sublvl = g.SubLvl
			break
		}
	}
	return sublvl, found
}
