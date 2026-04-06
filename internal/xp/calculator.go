package xp

import (
	"strings"
)

const (
	BaseXP      = 10
	XPPerLevel  = 100
	MaxLevel    = 10
	SufferingXP = 5
)

// CommandMultiplier returns the XP multiplier for a given command string.
func CommandMultiplier(cmd string) float64 {
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return 1.0
	}

	// Strip any path prefix (e.g. /usr/bin/git → git)
	bin := strings.ToLower(parts[0])
	if idx := strings.LastIndexByte(bin, '/'); idx >= 0 {
		bin = bin[idx+1:]
	}

	switch bin {
	case "git", "kubectl", "terraform", "docker":
		return 3.0
	case "python", "python3", "go", "cargo", "npm", "yarn", "pnpm", "make", "mvn", "gradle":
		return 2.0
	case "vim", "nvim", "nano", "emacs", "hx":
		return 1.5
	case "cd", "ls", "echo", "cat", "pwd", "clear", "exit", "history":
		return 0.5
	default:
		return 1.0
	}
}

// Earned returns the raw XP earned for a command execution.
// exitCode != 0 adds SufferingXP.
func Earned(cmd string, exitCode int) int {
	gained := int(float64(BaseXP) * CommandMultiplier(cmd))
	if exitCode != 0 {
		gained += SufferingXP
	}
	return gained
}

// Apply adds earned XP to the current state values and handles level-ups.
// Returns (newXP, newLevel, leveledUp).
func Apply(currentXP, currentLevel, gained int) (newXP, newLevel int, leveledUp bool) {
	if currentLevel >= MaxLevel {
		return currentXP, currentLevel, false
	}

	newXP = currentXP + gained
	newLevel = currentLevel

	for newXP >= XPPerLevel && newLevel < MaxLevel {
		newXP -= XPPerLevel
		newLevel++
		leveledUp = true
	}

	return newXP, newLevel, leveledUp
}

// ProgressBar returns a 5-cell bar string like "▓▓▓░░" for the given xp/xpToNext.
func ProgressBar(xp, xpToNext int) string {
	const cells = 5
	if xpToNext <= 0 {
		return "▓▓▓▓▓"
	}
	filled := (xp * cells) / xpToNext
	if filled > cells {
		filled = cells
	}
	bar := ""
	for i := 0; i < cells; i++ {
		if i < filled {
			bar += "▓"
		} else {
			bar += "░"
		}
	}
	return bar
}
