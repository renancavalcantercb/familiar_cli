package cmd

import (
	"fmt"
	"os"

	"github.com/renancavalcantercb/familiar_cli/internal/state"
	"github.com/renancavalcantercb/familiar_cli/internal/xp"
)

func runPrompt() {
	s, err := state.Load()
	if err != nil {
		// Silent fail — prompt must never break the shell.
		os.Exit(0)
	}
	state.MigrateState(s) //nolint:errcheck

	bar := xp.ProgressBar(s.XP, s.XPToNext)
	prefix := ""
	if s.Evolved {
		prefix = "✨"
	}
	fmt.Printf("%s%s %s Lv.%d %s %dxp", prefix, s.Emoji, s.Name, s.Level, bar, s.XP)
}
