package spotifind

import (
	"errors"
	"strings"
)

var (
	ErrNoResults    = errors.New("no results found")
	ErrTimeout      = errors.New("timeout")
	ErrTokenExpired = errors.New("token expired")
)

type ErrorHandling struct {
}

func (e *ErrorHandling) Handle(err error) error {
	if err == nil {
		return nil
	}

	if strings.Contains(err.Error(), "oauth2: token expired and refresh token is not set") {
		return ErrTokenExpired
	}

	switch err.Error() {
	case "spotify: couldn't decode error: (17) [Too many requests]":
		return ErrTimeout
	default:
		return err
	}
}
