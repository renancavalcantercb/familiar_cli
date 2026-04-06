package cmd

import (
	"fmt"
	"os"

	"github.com/renan-dev/familiar/internal/species"
	"github.com/renan-dev/familiar/internal/state"
	"github.com/renan-dev/familiar/internal/xp"
)

func runStatus() {
	s, err := state.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	state.MigrateState(s) //nolint:errcheck

	sp := species.ByID(s.Species)
	if sp == nil {
		fmt.Fprintf(os.Stderr, "unknown species %q in state file\n", s.Species)
		os.Exit(1)
	}

	bar := xp.ProgressBar(s.XP, s.XPToNext)

	fmt.Print(sp.AsciiArt)
	fmt.Println("─────────────────────────────")
	fmt.Printf("  %s  %s\n", s.Emoji, s.Name)
	fmt.Printf("  Species:   %s\n", sp.Name)
	fmt.Printf("  Traits:    %s\n", joinTraits(sp.Traits))
	fmt.Printf("  Level:     %d / %d\n", s.Level, xp.MaxLevel)
	fmt.Printf("  XP:        %s  %d / %d\n", bar, s.XP, s.XPToNext)
	fmt.Printf("  Commands:  %d\n", s.TotalCommands)
	fmt.Printf("  Born:      %s\n", s.CreatedAt.Format("2006-01-02"))
	fmt.Println("─────────────────────────────")
	fmt.Println()

	// Print attributes
	fmt.Println("  Attributes:")
	for attr, val := range s.Attributes {
		attrBar := attrBar(val)
		fmt.Printf("    %-15s %s %d\n", attr, attrBar, val)
	}
	fmt.Println()
}

// attrBar renders a small 10-cell bar for an attribute value (0-10).
func attrBar(val int) string {
	const max = 10
	bar := ""
	for i := 0; i < max; i++ {
		if i < val {
			bar += "█"
		} else {
			bar += "░"
		}
	}
	return bar
}
