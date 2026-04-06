package cmd

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"os"
	"time"

	"github.com/renan-dev/familiar/internal/species"
	"github.com/renan-dev/familiar/internal/state"
	"github.com/renan-dev/familiar/internal/xp"
)

func runInit() {
	if state.Exists() {
		s, err := state.Load()
		if err == nil {
			fmt.Printf("You already have a familiar: %s %s (Lv.%d)\n", s.Emoji, s.Name, s.Level)
			fmt.Println("Use 'familiar status' to see them.")
			return
		}
	}

	username := os.Getenv("USER")
	if username == "" {
		username = os.Getenv("USERNAME") // Windows fallback
	}
	if username == "" {
		username = "adventurer"
	}

	hostname, _ := os.Hostname()
	if hostname == "" {
		hostname = "localhost"
	}

	seed := deterministicHash(username + "@" + hostname)
	sp := species.FromHash(seed)

	// Use a slightly different hash offset for the name so it's not always index 0.
	nameSeed := deterministicHash(hostname + "@" + username)
	name := species.NameFromHash(sp, nameSeed)

	s := &state.State{
		Species:       sp.ID,
		Emoji:         sp.Emoji,
		Name:          name,
		Level:         1,
		XP:            0,
		XPToNext:      xp.XPPerLevel,
		Attributes:    sp.Attributes,
		TotalCommands: 0,
		CreatedAt:     time.Now(),
		Inventory:     []string{sp.ID},
		XPByCategory:  make(map[string]int),
	}

	if err := state.Save(s); err != nil {
		fmt.Fprintf(os.Stderr, "error saving state: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nA wild familiar appeared!\n\n")
	fmt.Print(sp.AsciiArt)
	fmt.Printf("\nYour familiar: %s %s\n", s.Emoji, s.Name)
	fmt.Printf("Species:       %s (%s)\n", sp.Name, joinTraits(sp.Traits))
	fmt.Printf("Level:         %d\n", s.Level)
	fmt.Printf("\nRun 'familiar status' anytime to check on them.\n")
	fmt.Printf("Add the fish hooks from the README to start earning XP!\n\n")
}

// deterministicHash returns a uint64 derived from a string via SHA-256.
func deterministicHash(input string) uint64 {
	h := sha256.Sum256([]byte(input))
	return binary.BigEndian.Uint64(h[:8])
}

func joinTraits(traits []string) string {
	result := ""
	for i, t := range traits {
		if i > 0 {
			result += ", "
		}
		result += t
	}
	return result
}
