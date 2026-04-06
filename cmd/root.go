package cmd

import (
	"fmt"
	"os"
)

func Execute() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(0)
	}

	switch os.Args[1] {
	case "init":
		runInit()
	case "status":
		runStatus()
	case "prompt":
		runPrompt()
	case "xp":
		runXP()
	case "daemon":
		runDaemon()
	case "roll":
		runRoll()
	case "switch":
		runSwitch()
	case "inventory":
		runInventory()
	case "stats":
		runStats()
	case "help", "--help", "-h":
		printHelp()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		printHelp()
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Print(`familiar — your terminal companion

Usage:
  familiar <command> [args]

Commands:
  init                    Create your familiar (deterministic by username+hostname)
  status                  Show your familiar with ASCII art and stats
  prompt                  Return short string for fish/zsh prompt
  xp <cmd> [code]         Register a command and calculate XP
  daemon                  Run the background daemon
  roll                    Spend a roll to unlock species or cosmetics
  switch <species_id>     Switch your active familiar (must be in inventory)
  inventory               List unlocked species and hats
  stats                   Show usage statistics
  help                    Show this help

Examples:
  familiar init
  familiar status
  familiar prompt
  familiar xp "git commit -m fix" 0
  familiar roll
  familiar switch mushroom
  familiar inventory
  familiar stats
`)
}
