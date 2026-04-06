package gacha

import (
	"math/rand"

	"github.com/renancavalcantercb/familiar_cli/internal/species"
)

var hats = []string{"🎩", "👒", "⛑️", "🪖", "👑", "🎓", "🧢", "🪄"}

// Result holds the outcome of a roll.
type Result struct {
	Type  string // "hat", "species", "shiny"
	Value string // hat emoji or species ID (shiny prefixed with "shiny_")
	IsNew bool
}

// Roll performs one gacha roll given the current inventory and owned hats.
// Inventory contains species IDs already unlocked (including "shiny_" prefixed ones).
func Roll(inventory, ownedHats []string) Result {
	n := rand.Float64()

	switch {
	case n < 0.60:
		// 60%: hat
		hat := hats[rand.Intn(len(hats))]
		return Result{Type: "hat", Value: hat, IsNew: !contains(ownedHats, hat)}

	case n < 0.90:
		// 30%: new base species
		available := availableSpecies(inventory, false)
		if len(available) == 0 {
			// Already have all — give a hat instead
			hat := hats[rand.Intn(len(hats))]
			return Result{Type: "hat", Value: hat, IsNew: !contains(ownedHats, hat)}
		}
		id := available[rand.Intn(len(available))]
		return Result{Type: "species", Value: id, IsNew: true}

	default:
		// 10%: shiny species
		available := availableSpecies(inventory, true)
		if len(available) == 0 {
			hat := hats[rand.Intn(len(hats))]
			return Result{Type: "hat", Value: hat, IsNew: !contains(ownedHats, hat)}
		}
		id := available[rand.Intn(len(available))]
		return Result{Type: "shiny", Value: "shiny_" + id, IsNew: true}
	}
}

// availableSpecies returns species IDs not yet in inventory.
// If shiny is true, checks for "shiny_<id>" entries.
func availableSpecies(inventory []string, shiny bool) []string {
	var result []string
	for _, sp := range species.All {
		id := sp.ID
		if shiny {
			id = "shiny_" + sp.ID
		}
		if !contains(inventory, id) {
			result = append(result, sp.ID)
		}
	}
	return result
}

func contains(slice []string, val string) bool {
	for _, s := range slice {
		if s == val {
			return true
		}
	}
	return false
}
