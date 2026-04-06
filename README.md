# familiar 🐾

A terminal companion that lives in your shell, gains XP as you work, and evolves over time — like a familiar from an RPG.

```
🍄 Trufa Lv.5 ▓▓░░░ 35xp ❯
```

Every command you run gives your familiar experience. Deploy to Kubernetes? Big XP. Hit `ls`? Tiny XP. Make an error? +5 suffering XP. Level up, unlock new species through gacha rolls, and collect them all.

---

## Install

```sh
go install github.com/renancavalcantercb/familiar_cli@latest
```

Or build from source:

```sh
git clone https://github.com/renancavalcantercb/familiar_cli
cd familiar_cli
go build -o familiar .
mv familiar ~/bin/  # anywhere in your $PATH
```

**Requirements:** Go 1.21+, no external dependencies.

---

## Quick start

```sh
familiar init      # summon your familiar (deterministic — same machine = same familiar)
familiar status    # see your companion
```

---

## Shell integration

### Fish

Add to `~/.config/fish/config.fish`:

```fish
# Award XP after every command (background, never delays prompt)
function fish_postexec --on-event fish_postexec
    $HOME/bin/familiar xp $argv[1] $status &
    disown
end
```

### Starship prompt

Add to `~/.config/starship.toml`:

```toml
format = """
... your other modules ...
$line_break\
${custom.familiar}\
❯ """

[custom.familiar]
command = "$HOME/bin/familiar prompt"
when = "true"
format = "[$output ](fg:#c4a7e7)"
shell = ["bash", "--noprofile", "--norc"]
```

### Zsh / Bash

```sh
# ~/.zshrc or ~/.bashrc
preexec() { familiar xp "$1" & }
PS1='$(familiar prompt) $ '
```

---

## Autocomplete

### Fish
```sh
cp completions/familiar.fish ~/.config/fish/completions/
```

---

## Commands

| Command | Description |
|---------|-------------|
| `familiar init` | Summon your familiar |
| `familiar status` | Full stats with ASCII art |
| `familiar prompt` | Short string for shell prompt |
| `familiar xp <cmd> [exit_code]` | Register a command and award XP |
| `familiar stats` | XP breakdown by category + activity |
| `familiar roll` | Spend a roll — drop a new species or cosmetic |
| `familiar inventory` | View your collection |
| `familiar switch <species>` | Switch your active familiar |
| `familiar rename <name>` | Rename your familiar |
| `familiar export` | Print a shareable ASCII card |
| `familiar version` | Show version |

---

## XP system

| Category | Commands | Multiplier |
|----------|----------|------------|
| DevOps | `git`, `kubectl`, `terraform`, `docker` | **3×** |
| Build | `go`, `python`, `cargo`, `npm`, `make` | **2×** |
| Editor | `vim`, `nvim`, `nano`, `emacs`, `hx` | **1.5×** |
| Shell | `cd`, `ls`, `echo`, `cat` | **0.5×** |
| Other | everything else | **1×** |

- Base XP per command: **10**
- Error bonus (exit ≠ 0): **+5 XP** *(learning by suffering)*
- Level up every **100 XP** — max level **10**
- Each level up grants **1 roll** 🎰

---

## Species

Each species is generated **deterministically** from your `username@hostname` — you always hatch the same familiar on the same machine.

| Emoji | Species | Traits |
|-------|---------|--------|
| 🦦 | Capybara | PATIENCE, CALM |
| 🦎 | Axolotl | REGENERATION, CHAOS |
| 🍄 | Mushroom | WISDOM, POISON |
| 👻 | Ghost | SARCASM, STEALTH |
| 🐉 | Dragon | POWER, CHAOS |
| 🦆 | Duck | DEBUGGING, STUBBORNNESS |
| 🐱 | Cat | INDEPENDENCE, IRONY |
| 🦉 | Owl | WISDOM, PATIENCE |

Unlock more species through **gacha rolls**. Rare drop: ✨ **shiny** variants.

---

## Gacha

Spend rolls earned from leveling up:

```
❯ familiar roll

🎰 rolling...

  ╔═══════════╗
  ║  dropped! ║
  ╚═══════════╝

🐉 Dragon — new familiar unlocked! (use: familiar switch dragon)
```

Drop rates:
- **60%** — cosmetic hat (🎩 👒 ⛑️ 🪖 👑 🎓 🧢 🪄)
- **30%** — new species
- **10%** — ✨ shiny species

---

## Speech

Your familiar occasionally comments on what you're doing. Each species has its own personality:

```
❯ kubectl get pods
NAME                    READY   STATUS
...

  🍄 Trufa: the cluster knows.
```

```
❯ git push origin main --force  # exit 1

  🐉 Ignis: UNACCEPTABLE.
```

---

## Status output

```
familiar status

   _( )_
  (  *  )
   |   |
  mushroom

─────────────────────────────
  🍄  Trufa
  Species:   Mushroom
  Traits:    WISDOM, POISON
  Level:     5 / 10
  XP:        ▓▓░░░  35 / 100
  Commands:  312
  Born:      2026-04-05
─────────────────────────────

  Attributes:
    wisdom          █████████░ 9
    poison          ██████░░░░ 6
```

---

## Stats

```
familiar stats

📊 Trufa — Stats

  Total commands:  312
  Total XP earned: 2840
  Days active:     4

  XP by category:
    devops  ████████░░  1420xp
    build   █████░░░░░   890xp
    misc    ███░░░░░░░   380xp
    editor  ██░░░░░░░░   150xp

  Familiar:   🍄 Trufa Lv.5
  Rolls left: 2
  Collection: 3/8 species
```

---

## Architecture

```
fish_prompt    →  reads ~/.familiar/state.json   (<1ms, no subprocess)
fish_postexec  →  familiar xp <cmd> &            (background, non-blocking)
familiar xp    →  updates state.json atomically  (tmp + rename, no corruption)
```

State is stored at `~/.familiar/state.json`. No daemon required in the MVP.

---

## Debug mode

```sh
FAMILIAR_DEBUG=1 familiar xp "kubectl apply -f deploy.yaml" 0
# [familiar] +30xp (3.0x devops) → 65/100 Lv.3
```

---

## License

MIT
