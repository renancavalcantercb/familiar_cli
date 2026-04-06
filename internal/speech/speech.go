package speech

import (
	"math/rand"
	"strings"
)

var evolveLines = map[string]string{
	"capybara": "stillness achieved. i am everything.",
	"mushroom": "the mycelium expands beyond comprehension.",
	"ghost":    "i have transcended the void.",
	"dragon":   "I AM THE ELDER FLAME. BOW.",
	"axolotl":  "chaos... perfected.",
	"duck":     "QUAK. I KNOW ALL BUGS. QUAK.",
	"cat":      "...whatever. i was always this.",
	"owl":      "all is known. all is understood.",
}

var levelUpLines = map[string][]string{
	"capybara": {"still here. still growing.", "patience rewarded.", "level up. calm.", "another level. same vibes."},
	"mushroom": {"the mycelium expands.", "new spores released.", "deeper roots.", "growth is inevitable."},
	"ghost":    {"ascending... or descending?", "more undefined behavior unlocked.", "leveled up. still haunting.", "void expanded."},
	"dragon":   {"POWER INCREASES.", "THE FLAME GROWS.", "ANOTHER LEVEL. BOW.", "UNSTOPPABLE."},
	"axolotl":  {"regenerated to new level.", "chaos leveled up.", "bloop bloop bloop.", "~leveled~"},
	"duck":     {"QUAK QUAK QUAK.", "level up. quak.", "rubber duck approves.", "PROMOTED. QUAK."},
	"cat":      {"...fine.", "leveled up. still unimpressed.", "whatever.", "*yawns*"},
	"owl":      {"wisdom deepens.", "hoot. knowledge grows.", "another level of understanding.", "hoot hoot."},
}

var lines = map[string]map[string][]string{
	"capybara": {
		"devops":  {"it's running.", "cluster is chill.", "deploy done. calm.", "uptime is peace."},
		"error":   {"you erred. it happens.", "it's just a warning.", "try again, no rush.", "the error is part of the path."},
		"editor":  {"writing code with calm.", "each line in its time.", "save the file.", "good edit."},
		"default": {"here.", "present.", "all good.", "go with the flow."},
	},
	"mushroom": {
		"devops":  {"the cluster knows.", "spores released.", "mycelium connected.", "the network grows."},
		"error":   {"suffering is wisdom.", "errors feed the soil.", "failure is fertilizer.", "grow through it."},
		"editor":  {"words become spores.", "each line spreads.", "the code breathes.", "write and release."},
		"default": {"i am here.", "the mycelium listens.", "growing.", "patient."},
	},
	"ghost": {
		"devops":  {"deployed into the void.", "null pointer acquired.", "ghost in the cluster.", "undefined behavior expected."},
		"error":   {"404: soul not found.", "segfault. classic.", "error? i am the error.", "undefined. like me."},
		"editor":  {"writing in the dark.", "code haunts.", "undefined function. relatable.", "the cursor blinks. so do i."},
		"default": {"...", "still here.", "watching.", "boo."},
	},
	"dragon": {
		"devops":  {"FIRE DEPLOYED.", "the cluster bows.", "kubectl OBEYED.", "DOMINANCE ESTABLISHED."},
		"error":   {"UNACCEPTABLE.", "the flame was not enough.", "i will burn harder.", "retry with MORE POWER."},
		"editor":  {"WRITING CODE OF DESTRUCTION.", "each keystroke ignites.", "POWER FLOWS THROUGH VIM.", "the editor trembles."},
		"default": {"RAWR.", "i am here.", "POWER.", "bow."},
	},
	"axolotl": {
		"devops":  {"regenerating... deployed.", "chaos accepted.", "it works? chaos.", "undefined order achieved."},
		"error":   {"regenerating from error.", "chaos is home.", "errors heal.", "i regrow."},
		"editor":  {"regenerating code.", "rewriting reality.", "chaos typed.", "words regrow."},
		"default": {"bloop.", "~", "floating.", "chaos."},
	},
	"duck": {
		"devops":  {"quak. deployed.", "rubber duck approved.", "debug complete. quak.", "kubectl quak."},
		"error":   {"QUAK.", "did you explain to the duck?", "quak quak quak.", "rubber duck says: think harder."},
		"editor":  {"quak. writing.", "explain to the duck.", "quak quak.", "the duck watches."},
		"default": {"quak.", "..quak..", "quak?", "QUAK."},
	},
	"cat": {
		"devops":  {"deployed. whatever.", "it works. i don't care.", "cluster running. not impressed.", "fine."},
		"error":   {"told you.", "obviously.", "i knew this would happen.", "not my problem."},
		"editor":  {"writing. don't talk to me.", "code. silence.", "typing. go away.", "focused. leave."},
		"default": {"...", "hmm.", "whatever.", "*stares*"},
	},
	"owl": {
		"devops":  {"wisdom deployed.", "the cluster is balanced.", "knowledge propagates.", "hoot. running."},
		"error":   {"observe the error. learn.", "wisdom comes from failure.", "hoot. analyze.", "the pattern reveals itself."},
		"editor":  {"words carry knowledge.", "write with intention.", "each line teaches.", "hoot. good code."},
		"default": {"hoot.", "observing.", "wisdom accumulates.", "all is known."},
	},
}

// classify maps a command + exit code to a speech context.
func classify(cmd string, exitCode int) string {
	if exitCode != 0 {
		return "error"
	}
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return "default"
	}
	bin := strings.ToLower(parts[0])
	if idx := strings.LastIndexByte(bin, '/'); idx >= 0 {
		bin = bin[idx+1:]
	}
	switch bin {
	case "git", "kubectl", "terraform", "docker",
		"python", "python3", "go", "cargo", "npm", "yarn", "pnpm", "make", "mvn", "gradle":
		return "devops"
	case "vim", "nvim", "nano", "emacs", "hx":
		return "editor"
	default:
		return "default"
	}
}

// GetEvolve returns the evolution speech line for the given species.
func GetEvolve(species string) string {
	if line, ok := evolveLines[species]; ok {
		return line
	}
	return "i have evolved."
}

// GetLevelUp returns a random level-up speech line for the given species.
func GetLevelUp(species string) string {
	l, ok := levelUpLines[species]
	if !ok {
		l = levelUpLines["capybara"]
	}
	return l[rand.Intn(len(l))]
}

// Get returns a random speech line for the given species and command context.
func Get(species, cmd string, exitCode int) string {
	ctx := classify(cmd, exitCode)

	sp, ok := lines[species]
	if !ok {
		sp = lines["capybara"]
	}

	ctxLines, ok := sp[ctx]
	if !ok {
		ctxLines = sp["default"]
	}

	return ctxLines[rand.Intn(len(ctxLines))]
}
