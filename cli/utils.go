package cli

import (
	"errors"
	"github.com/je09/spotifind"
	"github.com/je09/spotifind/cli/durationFmt"
	"os"
	"time"
)

func errHandler(err error) {
	if errors.Is(err, spotifind.ErrTimeout) {
		rootCmd.Printf(Red + "\ntimeout while searching playlist\n" + Reset)
		os.Exit(1)
	}

	rootCmd.Printf(Red+"\nerror while searching playlist: %v\n"+Reset, err)
	os.Exit(1)
}

func firstThree(items []string) []string {
	if len(items) > 3 {
		return items[:3]
	}
	return items
}

func shortDur(d time.Duration) string {
	r, _ := durationFmt.Format(d, "%0h:%0m:%0s")
	return r
}
