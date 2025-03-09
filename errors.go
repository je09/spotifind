package spotifind

import (
	"errors"
	"strings"
)

var (
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

	// HACK: This is a workaround for the error message "com.spotify.hermes.service.RequestTimeoutException: Request timed out"
	// it feels like it needs to be retry instead, but I'm not sure rn, and it kinda breaks the scan at the moment.
	// TODO: Investigate further.
	if strings.Contains(strings.ToLower(err.Error()), "timeout") {
		return ErrTimeout
	}

	switch err.Error() {
	case "spotify: couldn't decode error: (17) [Too many requests]":
		return ErrTimeout
	default:
		return err
	}
}
