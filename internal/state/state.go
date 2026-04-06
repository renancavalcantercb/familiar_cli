package state

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ErrNotInitialized is returned when no state file exists.
var ErrNotInitialized = errors.New("familiar not initialized — run 'familiar init' first")

// State holds all persistent data for the user's familiar.
type State struct {
	Species       string         `json:"species"`
	Emoji         string         `json:"emoji"`
	Name          string         `json:"name"`
	Level         int            `json:"level"`
	XP            int            `json:"xp"`
	XPToNext      int            `json:"xp_to_next"`
	Attributes    map[string]int `json:"attributes"`
	TotalCommands int            `json:"total_commands"`
	CreatedAt     time.Time      `json:"created_at"`

	// Gacha system
	Rolls     int      `json:"rolls"`
	Inventory []string `json:"inventory"`  // species IDs unlocked (includes "shiny_" prefix for shiny)
	Hats      []string `json:"hats"`       // cosmetic hat emojis unlocked

	// Stats tracking
	XPByCategory map[string]int `json:"xp_by_category"`
	DaysActive   []string       `json:"days_active"` // YYYY-MM-DD
}

// Dir returns the path to ~/.familiar.
func Dir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot find home directory: %w", err)
	}
	return filepath.Join(home, ".familiar"), nil
}

// Path returns the path to ~/.familiar/state.json.
func Path() (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "state.json"), nil
}

// Exists reports whether a state file already exists.
func Exists() bool {
	p, err := Path()
	if err != nil {
		return false
	}
	_, err = os.Stat(p)
	return err == nil
}

// Load reads and parses the state file.
func Load() (*State, error) {
	p, err := Path()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNotInitialized
		}
		return nil, fmt.Errorf("reading state: %w", err)
	}

	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("corrupted state file: %w", err)
	}
	return &s, nil
}

// legacyIDs maps old Portuguese species IDs to their English equivalents.
var legacyIDs = map[string]string{
	"capivara": "capybara",
	"axolote":  "axolotl",
	"cogumelo": "mushroom",
	"fantasma": "ghost",
	"dragao":   "dragon",
	"pato":     "duck",
	"gato":     "cat",
	"coruja":   "owl",
}

// MigrateState converts any legacy Portuguese species IDs to English in-place.
// It saves the state if any migration occurred. Safe to call even if state is current.
func MigrateState(s *State) error {
	migrated := false

	if newID, ok := legacyIDs[s.Species]; ok {
		s.Species = newID
		migrated = true
	}

	for i, id := range s.Inventory {
		prefix := ""
		bare := id
		if len(id) > 6 && id[:6] == "shiny_" {
			prefix = "shiny_"
			bare = id[6:]
		}
		if newID, ok := legacyIDs[bare]; ok {
			s.Inventory[i] = prefix + newID
			migrated = true
		}
	}

	if migrated {
		return Save(s)
	}
	return nil
}

// Save writes the state atomically via a temp file + rename.
func Save(s *State) error {
	dir, err := Dir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating state directory: %w", err)
	}

	p, err := Path()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("serializing state: %w", err)
	}

	tmp := p + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return fmt.Errorf("writing temp state: %w", err)
	}

	if err := os.Rename(tmp, p); err != nil {
		os.Remove(tmp)
		return fmt.Errorf("committing state: %w", err)
	}

	return nil
}
