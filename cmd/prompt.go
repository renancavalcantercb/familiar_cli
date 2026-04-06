package cmd

import (
	"fmt"
	"os"

	"github.com/renan-dev/familiar/internal/state"
	"github.com/renan-dev/familiar/internal/xp"
)

func runPrompt() {
	s, err := state.Load()
	if err != nil {
		// Silent fail — prompt must never break the shell.
		os.Exit(0)
	}
	state.MigrateState(s) //nolint:errcheck

	bar := xp.ProgressBar(s.XP, s.XPToNext)
	fmt.Printf("%s %s Lv.%d %s %dxp", s.Emoji, s.Name, s.Level, bar, s.XP)
}
