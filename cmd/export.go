package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/renancavalcantercb/familiar_cli/internal/species"
	"github.com/renancavalcantercb/familiar_cli/internal/state"
	"github.com/renancavalcantercb/familiar_cli/internal/xp"
)

const cardWidth = 44

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
	var displaySpeciesName string
	if sp != nil {
		rawArt := sp.AsciiArt
		labelName := strings.ToLower(sp.Name)
		displaySpeciesName = sp.Name
		if s.Evolved {
			rawArt = sp.EvolvedArt
			labelName = strings.ToLower(sp.EvolvedName)
			displaySpeciesName = sp.EvolvedName + " ✨"
			traits = strings.Join(sp.EvolvedTraits, ", ")
		} else {
			traits = strings.Join(sp.Traits, ", ")
		}
		raw := strings.TrimSpace(rawArt)
		allLines := strings.Split(raw, "\n")
		// Drop trailing label line if it matches the species name
		for len(allLines) > 0 && strings.TrimSpace(allLines[len(allLines)-1]) == labelName {
			allLines = allLines[:len(allLines)-1]
		}
		asciiLines = allLines
	}

	bar := xp.ProgressBar(s.XP, s.XPToNext)
	daysActive := len(s.DaysActive)

	border := strings.Repeat("═", cardWidth-2)
	fmt.Printf("╔%s╗\n", border)
	printCardLine("      my familiar", cardWidth)
	fmt.Printf("╠%s╣\n", border)
	printCardLine("", cardWidth)

	// Find minimum indentation across non-empty lines to normalize art alignment
	minIndent := 999
	for _, line := range asciiLines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		indent := len(line) - len(strings.TrimLeft(line, " \t"))
		if indent < minIndent {
			minIndent = indent
		}
	}
	if minIndent == 999 {
		minIndent = 0
	}
	for _, line := range asciiLines {
		normalized := line
		if len(line) >= minIndent {
			normalized = line[minIndent:]
		}
		printCardLine("  "+normalized, cardWidth)
	}

	printCardLine("", cardWidth)
	printCardLine(fmt.Sprintf("  %s %s  •  %s", s.Emoji, s.Name, displaySpeciesName), cardWidth)
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
		switch {
		case r >= 0x1F000:
			// Emoji blocks (includes ✨ 0x2728 is below, handled separately)
			w += 2
		case r >= 0x2600 && r <= 0x27FF:
			// Misc symbols, dingbats — includes ✨ (U+2728)
			w += 2
		case r >= 0x2E80 && r <= 0x9FFF:
			// CJK
			w += 2
		case r >= 0xAC00 && r <= 0xD7AF:
			// Hangul
			w += 2
		default:
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
