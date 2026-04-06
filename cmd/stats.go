package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/renan-dev/familiar/internal/state"
	"github.com/renan-dev/familiar/internal/xp"
)

func runStats() {
	s, err := state.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	totalXP := 0
	for _, v := range s.XPByCategory {
		totalXP += v
	}

	fmt.Println()
	fmt.Printf("📊 %s — Stats\n", s.Name)
	fmt.Println("─────────────────────────────")
	fmt.Printf("\n  Total commands:  %d\n", s.TotalCommands)
	fmt.Printf("  Total XP earned: %d\n", totalXP)
	fmt.Printf("  Days active:     %d\n", len(s.DaysActive))

	if len(s.XPByCategory) > 0 {
		fmt.Println("\n  XP by category:")

		// Find max for bar scaling
		maxXP := 0
		for _, v := range s.XPByCategory {
			if v > maxXP {
				maxXP = v
			}
		}

		// Sort categories by XP descending
		type catEntry struct {
			name string
			xp   int
		}
		var entries []catEntry
		for k, v := range s.XPByCategory {
			entries = append(entries, catEntry{k, v})
		}
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].xp > entries[j].xp
		})

		for _, e := range entries {
			bar := categoryBar(e.xp, maxXP, 10)
			fmt.Printf("    %-8s %s  %dxp\n", e.name, bar, e.xp)
		}
	}

	bar := xp.ProgressBar(s.XP, s.XPToNext)
	collectionSize := countBaseSpecies(s.Inventory)

	fmt.Printf("\n  Familiar:   %s %s Lv.%d  %s\n", s.Emoji, s.Name, s.Level, bar)
	fmt.Printf("  Rolls left: %d\n", s.Rolls)
	fmt.Printf("  Collection: %d/8 species\n", collectionSize)
	fmt.Println()
}

// categoryBar renders a filled bar with given width, scaled to max.
func categoryBar(val, max, width int) string {
	if max <= 0 {
		return strings.Repeat("░", width)
	}
	filled := (val * width) / max
	if filled > width {
		filled = width
	}
	return strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
}
