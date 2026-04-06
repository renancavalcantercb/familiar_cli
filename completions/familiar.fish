# familiar fish completions

set -l commands init status prompt xp stats roll inventory switch rename export version daemon help

# Main commands
complete -c familiar -f -n "not __fish_seen_subcommand_from $commands" -a init      -d "Summon your familiar"
complete -c familiar -f -n "not __fish_seen_subcommand_from $commands" -a status    -d "Show familiar with ASCII art"
complete -c familiar -f -n "not __fish_seen_subcommand_from $commands" -a prompt    -d "Short string for shell prompt"
complete -c familiar -f -n "not __fish_seen_subcommand_from $commands" -a stats     -d "XP breakdown and activity"
complete -c familiar -f -n "not __fish_seen_subcommand_from $commands" -a roll      -d "Spend a roll — drop species or cosmetic"
complete -c familiar -f -n "not __fish_seen_subcommand_from $commands" -a inventory -d "View your collection"
complete -c familiar -f -n "not __fish_seen_subcommand_from $commands" -a switch    -d "Switch active familiar"
complete -c familiar -f -n "not __fish_seen_subcommand_from $commands" -a rename    -d "Rename your familiar"
complete -c familiar -f -n "not __fish_seen_subcommand_from $commands" -a export    -d "Export familiar card"
complete -c familiar -f -n "not __fish_seen_subcommand_from $commands" -a version   -d "Show version"
complete -c familiar -f -n "not __fish_seen_subcommand_from $commands" -a daemon    -d "Run background daemon"

# familiar switch — complete with all species IDs
complete -c familiar -f -n "__fish_seen_subcommand_from switch" -a "capybara axolotl mushroom ghost dragon duck cat owl" -d "Species"
