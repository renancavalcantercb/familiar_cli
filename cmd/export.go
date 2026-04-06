package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/renan-dev/familiar/internal/species"
	"github.com/renan-dev/familiar/internal/state"
	"github.com/renan-dev/familiar/internal/xp"
)

const cardWidth = 36

func runExport() {
	s, err := state.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	state.MigrateState(s) //nolint:errcheck

	sp := species.ByID(s.Species)
	var asciiLines []string
	var traits string
	if sp != nil {
		raw := strings.TrimSpace(sp.AsciiArt)
		// Remove last line (species label baked into art)
		allLines := strings.Split(raw, "\n")
		// Drop trailing label line if it's just the species name
		for len(allLines) > 0 && strings.TrimSpace(allLines[len(allLines)-1]) == strings.ToLower(sp.Name) {
			allLines = allLines[:len(allLines)-1]
		}
		asciiLines = allLines
		traits = strings.Join(sp.Traits, ", ")
	}

	bar := xp.ProgressBar(s.XP, s.XPToNext)
	daysActive := len(s.DaysActive)

	border := strings.Repeat("═", cardWidth-2)
	fmt.Printf("╔%s╗\n", border)
	printCardLine("      my familiar", cardWidth)
	fmt.Printf("╠%s╣\n", border)
	printCardLine("", cardWidth)

	for _, line := range asciiLines {
		printCardLine("  "+line, cardWidth)
	}

	printCardLine("", cardWidth)
	printCardLine(fmt.Sprintf("  %s %s  •  %s", s.Emoji, s.Name, sp.Name), cardWidth)
	printCardLine(fmt.Sprintf("  Lv.%d  •  %s", s.Level, traits), cardWidth)
	printCardLine(fmt.Sprintf("  %s  %d/%d XP", bar, s.XP, s.XPToNext), cardWidth)
	printCardLine(fmt.Sprintf("  %d commands  •  %d days active", s.TotalCommands, daysActive), cardWidth)
	printCardLine("", cardWidth)

	fmt.Printf("╚%s╝\n", border)
	fmt.Println()
	fmt.Println("copied to clipboard? use: familiar export | pbcopy")
}

// visualWidth returns the display width of a string, counting wide chars (emojis, CJK) as 2.
func visualWidth(s string) int {
	w := 0
	for _, r := range s {
		if r >= 0x1F000 || (r >= 0x2E80 && r <= 0x9FFF) || (r >= 0xAC00 && r <= 0xD7AF) {
			w += 2
		} else {
			w++
		}
	}
	return w
}

func printCardLine(content string, width int) {
	inner := width - 2
	vw := visualWidth(content)
	padding := inner - vw
	if padding < 0 {
		padding = 0
	}
	fmt.Printf("║%s%s║\n", content, strings.Repeat(" ", padding))
}
