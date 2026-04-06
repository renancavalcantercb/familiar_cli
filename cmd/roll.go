package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/renancavalcantercb/familiar_cli/internal/gacha"
	"github.com/renancavalcantercb/familiar_cli/internal/species"
	"github.com/renancavalcantercb/familiar_cli/internal/state"
)

func runRoll() {
	s, err := state.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	state.MigrateState(s) //nolint:errcheck

	if s.Rolls == 0 {
		fmt.Println("no rolls available. keep using the terminal.")
		return
	}

	s.Rolls--

	result := gacha.Roll(s.Inventory, s.Hats)

	// Animation
	fmt.Println()
	fmt.Println("🎰 rolling...")
	time.Sleep(300 * time.Millisecond)
	fmt.Println()
	fmt.Println("  ╔═══════════╗")
	fmt.Println("  ║  dropped! ║")
	fmt.Println("  ╚═══════════╝")
	fmt.Println()

	switch result.Type {
	case "hat":
		if result.IsNew {
			s.Hats = append(s.Hats, result.Value)
			fmt.Printf("%s Top Hat — new cosmetic unlocked!\n", result.Value)
		} else {
			fmt.Printf("%s Hat — already owned, but it's still stylish.\n", result.Value)
		}

	case "species":
		sp := species.ByID(result.Value)
		if sp == nil {
			fmt.Println("unknown species dropped.")
			break
		}
		if result.IsNew {
			s.Inventory = append(s.Inventory, result.Value)
		}
		fmt.Printf("%s %s — new familiar unlocked! (use: familiar switch %s)\n", sp.Emoji, sp.Name, sp.ID)

	case "shiny":
		baseID := strings.TrimPrefix(result.Value, "shiny_")
		sp := species.ByID(baseID)
		if sp == nil {
			fmt.Println("unknown species dropped.")
			break
		}
		if result.IsNew {
			s.Inventory = append(s.Inventory, result.Value)
		}
		fmt.Printf("✨ %s %s — ✨ SHINY familiar unlocked! (use: familiar switch %s)\n", sp.Emoji, sp.Name, baseID)
	}

	fmt.Println()

	if err := state.Save(s); err != nil {
		fmt.Fprintf(os.Stderr, "error saving state: %v\n", err)
		os.Exit(1)
	}
}

func runSwitch() {
	args := os.Args[2:]
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: familiar switch <species_id>")
		os.Exit(1)
	}

	targetID := strings.ToLower(args[0])

	s, err := state.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	state.MigrateState(s) //nolint:errcheck

	// Check inventory
	if !containsStr(s.Inventory, targetID) {
		fmt.Printf("'%s' is not in your collection. use 'familiar inventory' to see what you have.\n", targetID)
		os.Exit(1)
	}

	sp := species.ByID(targetID)
	if sp == nil {
		fmt.Fprintf(os.Stderr, "unknown species: %s\n", targetID)
		os.Exit(1)
	}

	if s.Species == targetID {
		fmt.Printf("you're already using %s %s!\n", sp.Emoji, sp.Name)
		return
	}

	prevName := s.Name
	s.Species = sp.ID
	s.Emoji = sp.Emoji
	s.Attributes = sp.Attributes
	// Keep level, XP, name (the name belongs to the trainer, not species)
	// But assign a new name from the new species
	s.Name = sp.Names[0] // default first name; deterministic

	if err := state.Save(s); err != nil {
		fmt.Fprintf(os.Stderr, "error saving state: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n%s switched from %s to %s %s!\n", prevName, prevName, sp.Emoji, sp.Name)
	fmt.Print(sp.AsciiArt)
	fmt.Printf("\nNow using: %s %s (Lv.%d)\n\n", s.Emoji, s.Name, s.Level)
}

func runInventory() {
	s, err := state.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	state.MigrateState(s) //nolint:errcheck

	fmt.Println()
	fmt.Printf("📦 %s's Collection\n", s.Name)
	fmt.Println("─────────────────────────────")

	fmt.Printf("\n  Familiars (%d/8):\n", countBaseSpecies(s.Inventory))
	if len(s.Inventory) == 0 {
		fmt.Println("    none yet — roll to unlock!")
	} else {
		for _, id := range s.Inventory {
			shiny := strings.HasPrefix(id, "shiny_")
			baseID := strings.TrimPrefix(id, "shiny_")
			sp := species.ByID(baseID)
			if sp == nil {
				continue
			}
			active := ""
			if s.Species == baseID {
				active = " ← active"
			}
			if shiny {
				fmt.Printf("    ✨ %s %s (shiny)%s  → familiar switch %s\n", sp.Emoji, sp.Name, active, baseID)
			} else {
				fmt.Printf("    %s %s%s  → familiar switch %s\n", sp.Emoji, sp.Name, active, baseID)
			}
		}
	}

	fmt.Printf("\n  Hats (%d):\n", len(s.Hats))
	if len(s.Hats) == 0 {
		fmt.Println("    none yet")
	} else {
		fmt.Printf("    ")
		fmt.Println(strings.Join(s.Hats, " "))
	}

	fmt.Printf("\n  Rolls available: %d\n", s.Rolls)
	fmt.Println()
}

func countBaseSpecies(inventory []string) int {
	count := 0
	for _, id := range inventory {
		if !strings.HasPrefix(id, "shiny_") {
			count++
		}
	}
	return count
}
