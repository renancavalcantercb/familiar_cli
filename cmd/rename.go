package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/renan-dev/familiar/internal/state"
)

func runRename() {
	args := os.Args[2:]
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: familiar rename <name>")
		os.Exit(1)
	}

	name := strings.TrimSpace(strings.Join(args, " "))
	if name == "" {
		fmt.Fprintln(os.Stderr, "error: name cannot be empty")
		os.Exit(1)
	}
	if len(name) > 20 {
		fmt.Fprintln(os.Stderr, "error: name must be 20 characters or fewer")
		os.Exit(1)
	}

	s, err := state.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	state.MigrateState(s) //nolint:errcheck

	old := s.Name
	s.Name = name

	if err := state.Save(s); err != nil {
		fmt.Fprintln(os.Stderr, "error saving state:", err)
		os.Exit(1)
	}

	fmt.Printf("%s %s renamed to %s!\n", s.Emoji, old, s.Name)
}
