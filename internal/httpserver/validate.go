package httpserver

import (
	"net/mail"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var slugPattern = regexp.MustCompile(`^[a-z0-9-]+$`)

func validationErr(msg string) error { return ValidationError{Message: msg} }

func requireEmail(v string) (string, error) {
	v = strings.TrimSpace(v)
	if _, err := mail.ParseAddress(v); err != nil {
		return "", validationErr("Ungültige E-Mail-Adresse")
	}
	return v, nil
}

func requireMinLen(field, v string, min int) (string, error) {
	if len(v) < min {
		return "", validationErr(field + " muss mindestens " + strconv.Itoa(min) + " Zeichen haben")
	}
	return v, nil
}

func optionalMinMaxLen(field, v string, min, max int) (*string, error) {
	if v == "" {
		return nil, nil
	}
	if len(v) < min || len(v) > max {
		return nil, validationErr(field + " muss zwischen " + strconv.Itoa(min) + " und " + strconv.Itoa(max) + " Zeichen haben")
	}
	return &v, nil
}

func requireMinMaxLen(field, v string, min, max int) (string, error) {
	if len(v) < min || len(v) > max {
		return "", validationErr(field + " muss zwischen " + strconv.Itoa(min) + " und " + strconv.Itoa(max) + " Zeichen haben")
	}
	return v, nil
}

func requireSlug(v string) (string, error) {
	if len(v) < 2 || !slugPattern.MatchString(v) {
		return "", validationErr("Slug muss aus Kleinbuchstaben, Ziffern und Bindestrichen bestehen")
	}
	return v, nil
}

func requireURL(field, v string) (string, error) {
	v = strings.TrimSpace(v)
	u, err := url.ParseRequestURI(v)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return "", validationErr(field + " muss eine gültige URL sein")
	}
	return v, nil
}

func optionalURL(field, v string) (*string, error) {
	v = strings.TrimSpace(v)
	if v == "" {
		return nil, nil
	}
	s, err := requireURL(field, v)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func requireEnum(field, v string, allowed ...string) (string, error) {
	for _, a := range allowed {
		if v == a {
			return v, nil
		}
	}
	return "", validationErr(field + " muss einer von " + strings.Join(allowed, ", ") + " sein")
}

func coerceInt(field, v string) (int, error) {
	n, err := strconv.Atoi(strings.TrimSpace(v))
	if err != nil {
		return 0, validationErr(field + " muss eine Zahl sein")
	}
	return n, nil
}

func coerceOptionalInt(v string) (*int, error) {
	v = strings.TrimSpace(v)
	if v == "" {
		return nil, nil
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return nil, validationErr("Wert muss eine Zahl sein")
	}
	return &n, nil
}

func requireIntMin(field string, n, min int) error {
	if n < min {
		return validationErr(field + " muss mindestens " + strconv.Itoa(min) + " sein")
	}
	return nil
}
