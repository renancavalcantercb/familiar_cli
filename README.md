# familiar

A terminal companion/pet that lives in your shell, gains XP as you work, and evolves over time — like a familiar from an RPG.

```
🦦 Nim Lv.3 ▓▓▓░░ 42xp
$ _
```

## What is it?

Every time you run a command, your familiar gains experience. Different tools give different multipliers (because `git commit` takes more brain than `ls`). Level up by using your terminal. Max level: 10.

## Installation

```sh
go install github.com/renan-dev/familiar@latest
```

Or build from source:

```sh
git clone https://github.com/renan-dev/familiar
cd familiar
go build -o familiar .
sudo mv familiar /usr/local/bin/
```

## Setup

### Initialize your familiar

```sh
familiar init
```

Your familiar is generated deterministically from your `username@hostname` — you always get the same species and name on the same machine.

### Fish shell hooks

Add to `~/.config/fish/config.fish`:

```fish
# Register XP after each command (runs in background so it never delays your prompt)
function fish_postexec --on-event fish_postexec
    familiar xp $argv[1] $status &
end

# Show your familiar in the prompt
function fish_prompt
    set -l familiar_line (familiar prompt 2>/dev/null)
    if test -n "$familiar_line"
        echo $familiar_line
    end
    echo "❯ "
end
```

### Minimal prompt (just the familiar line)

```fish
function fish_prompt
    echo (familiar prompt 2>/dev/null)
    echo "$ "
end
```

## Commands

| Command | Description |
|---------|-------------|
| `familiar init` | Create your familiar (deterministic by username+hostname) |
| `familiar status` | Show your familiar with ASCII art and full stats |
| `familiar prompt` | Return a short string for the shell prompt |
| `familiar xp <cmd> [exit_code]` | Register a command and award XP |
| `familiar daemon` | Run the background keeper process |

## XP System

| Category | Commands | Multiplier |
|----------|----------|------------|
| DevOps | `git`, `kubectl`, `terraform`, `docker` | 3× |
| Build | `go`, `python`, `cargo`, `npm`, `make` | 2× |
| Editor | `vim`, `nvim`, `nano`, `emacs` | 1.5× |
| Shell | `cd`, `ls`, `echo`, `cat` | 0.5× |
| Other | everything else | 1× |

**Error bonus:** any command with exit code ≠ 0 gives +5 XP ("learning by suffering").

Base XP per command: **10**. Level up every **100 XP**. Max level: **10**.

## Species

| Emoji | Species | Traits |
|-------|---------|--------|
| 🦦 | Capivara | PACIÊNCIA, CALMA |
| 🦎 | Axolote | REGENERAÇÃO, CAOS |
| 🍄 | Cogumelo | SABEDORIA, VENENO |
| 👻 | Fantasma | SARCASMO, FURTIVIDADE |
| 🐉 | Dragão | PODER, CAOS |
| 🦆 | Pato | DEBUGGING, TEIMOSIA |
| 🐱 | Gato | INDEPENDÊNCIA, IRONIA |
| 🦉 | Coruja | SABEDORIA, PACIÊNCIA |

## Example output

```
familiar status

  /\_/\
 ( o.o )
  > ^ <
   gato

─────────────────────────────
  🐱  Pixel
  Species:   Gato
  Traits:    INDEPENDÊNCIA, IRONIA
  Level:     3 / 10
  XP:        ▓▓▓░░  42 / 100
  Commands:  156
  Born:      2026-04-05
─────────────────────────────

  Attributes:
    independence    ██████████ 10
    irony           ████████░░ 8
```

## Architecture

The prompt integration uses a read-only pattern to stay fast:

```
fish_prompt   → reads ~/.familiar/state.json  (<1ms, pure fish)
fish_postexec → calls familiar xp <cmd> &     (background, non-blocking)
familiar xp   → updates ~/.familiar/state.json atomically
```

`familiar daemon` is optional in the MVP — XP is written directly by `familiar xp`. Start it for future event processing:

```sh
familiar daemon &
```

## Debug mode

```sh
FAMILIAR_DEBUG=1 familiar xp "git commit -m fix" 0
# [familiar] +30xp (3.0x devops) → 72/100 Lv.3
```

## State file

Located at `~/.familiar/state.json`:

```json
{
  "species": "gato",
  "emoji": "🐱",
  "name": "Pixel",
  "level": 3,
  "xp": 42,
  "xp_to_next": 100,
  "attributes": {"independence": 10, "irony": 8},
  "total_commands": 156,
  "created_at": "2026-04-05T00:00:00Z"
}
```
