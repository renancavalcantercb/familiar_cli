package cmd

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/renan-dev/familiar/internal/speech"
	"github.com/renan-dev/familiar/internal/state"
	"github.com/renan-dev/familiar/internal/xp"
)

func runXP() {
	// Usage: familiar xp <command> [exit_code]
	// Called from fish_postexec: familiar xp $argv[1] $status
	args := os.Args[2:]

	var cmd string
	exitCode := 0

	if len(args) >= 1 {
		cmd = args[0]
	}
	if len(args) >= 2 {
		if code, err := strconv.Atoi(args[1]); err == nil {
			exitCode = code
		}
	}

	s, err := state.Load()
	if err != nil {
		// Silent fail — must never break the shell.
		os.Exit(0)
	}

	gained := xp.Earned(cmd, exitCode)
	newXP, newLevel, leveledUp := xp.Apply(s.XP, s.Level, gained)

	s.XP = newXP
	s.Level = newLevel
	s.XPToNext = xp.XPPerLevel
	s.TotalCommands++

	// Update XP by category
	category := classifyCommand(cmd)
	if exitCode != 0 {
		category = "error"
	}
	if s.XPByCategory == nil {
		s.XPByCategory = make(map[string]int)
	}
	s.XPByCategory[category] += gained

	// Track active day
	today := time.Now().Format("2006-01-02")
	if !containsStr(s.DaysActive, today) {
		s.DaysActive = append(s.DaysActive, today)
	}

	// Roll on level-up
	if leveledUp {
		s.Rolls++
	}

	if err := state.Save(s); err != nil {
		// Silent fail.
		os.Exit(0)
	}

	if leveledUp {
		bar := xp.ProgressBar(s.XP, s.XPToNext)
		fmt.Printf("\n✨ %s %s leveled up to Lv.%d! %s\n", s.Emoji, s.Name, s.Level, bar)
		fmt.Printf("   🎰 +1 roll available! (%d total)\n", s.Rolls)
	}

	// Occasional speech (1 in 8 chance)
	if rand.Intn(8) == 0 {
		line := speech.Get(s.Species, cmd, exitCode)
		fmt.Printf("\n  %s %s: %s\n", s.Emoji, s.Name, line)
	}

	// Emit feedback if command was interesting (optional debug mode via env).
	if os.Getenv("FAMILIAR_DEBUG") != "" {
		mul := xp.CommandMultiplier(cmd)
		cat := classifyCommand(cmd)
		fmt.Fprintf(os.Stderr, "[familiar] +%dxp (%.1fx %s) → %d/%d Lv.%d\n",
			gained, mul, cat, s.XP, s.XPToNext, s.Level)
	}
}

func classifyCommand(cmd string) string {
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return "unknown"
	}
	bin := strings.ToLower(parts[0])
	if idx := strings.LastIndexByte(bin, '/'); idx >= 0 {
		bin = bin[idx+1:]
	}
	switch bin {
	case "git", "kubectl", "terraform", "docker":
		return "devops"
	case "python", "python3", "go", "cargo", "npm", "yarn", "pnpm", "make", "mvn", "gradle":
		return "build"
	case "vim", "nvim", "nano", "emacs", "hx":
		return "editor"
	case "cd", "ls", "echo", "cat", "pwd", "clear", "exit", "history":
		return "shell"
	default:
		return "misc"
	}
}

func containsStr(slice []string, val string) bool {
	for _, s := range slice {
		if s == val {
			return true
		}
	}
	return false
}
