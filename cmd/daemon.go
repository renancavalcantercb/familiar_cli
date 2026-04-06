package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/renancavalcantercb/familiar_cli/internal/state"
)

func runDaemon() {
	s, err := state.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("[familiar daemon] watching over %s %s (Lv.%d)\n", s.Emoji, s.Name, s.Level)
	fmt.Println("[familiar daemon] press Ctrl+C to stop")

	// In MVP, the daemon is a simple keeper process.
	// XP is written directly by `familiar xp`, so the daemon just stays alive
	// and periodically prints a heartbeat (useful for future event processing).
	tick := time.NewTicker(60 * time.Second)
	defer tick.Stop()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-tick.C:
			s, err := state.Load()
			if err == nil {
				fmt.Printf("[familiar daemon] %s %s — Lv.%d %dxp — %d commands\n",
					s.Emoji, s.Name, s.Level, s.XP, s.TotalCommands)
			}
		case sig := <-sigs:
			fmt.Printf("\n[familiar daemon] received %s, shutting down\n", sig)
			return
		}
	}
}
