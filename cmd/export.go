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

func printCardLine(content string, width int) {
	// width includes the two border chars ║ ... ║
	inner := width - 2
	// Pad or truncate content to inner width
	runes := []rune(content)
	if len(runes) > inner {
		runes = runes[:inner]
	}
	padding := inner - len(runes)
	fmt.Printf("║%s%s║\n", string(runes), strings.Repeat(" ", padding))
}
