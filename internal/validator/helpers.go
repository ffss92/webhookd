package validator

import (
	"net/mail"
	"net/url"
	"strings"
	"unicode/utf8"
)

func NotBlank(value string) bool {
	value = strings.TrimSpace(value)
	return utf8.RuneCountInString(value) > 0
}

func MinLength(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

func MaxLength(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

func ValidEmail(value string) bool {
	addr, err := mail.ParseAddress(value)
	return err == nil && addr.Address == value
}

func HTTPUrl(value string) bool {
	u, err := url.ParseRequestURI(value)
	if err != nil {
		return false
	}
	if u.Scheme == "http" || u.Scheme == "https" {
		return u.Host != ""
	}
	return false
}
