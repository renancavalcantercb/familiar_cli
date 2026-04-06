package species

// Species defines a familiar species with all its characteristics.
type Species struct {
	ID         string
	Name       string
	Emoji      string
	Attributes map[string]int
	AsciiArt   string
	Names      []string
	Traits     []string
}

// All contains every available species for the MVP.
var All = []Species{
	{
		ID:    "capybara",
		Name:  "Capybara",
		Emoji: "🦦",
		Attributes: map[string]int{
			"patience": 9,
			"calm":     8,
		},
		Traits: []string{"PATIENCE", "CALM"},
		AsciiArt: `
  ╭─────╮
  │ ◕ ◕ │
  ╰──▀──╯─~
 capybara
`,
		Names: []string{"Nim", "Kapi", "Roux", "Mochi", "Taro", "Zara", "Paca", "Boto"},
	},
	{
		ID:    "axolotl",
		Name:  "Axolotl",
		Emoji: "🦎",
		Attributes: map[string]int{
			"regeneration": 9,
			"chaos":        7,
		},
		Traits: []string{"REGENERATION", "CHAOS"},
		AsciiArt: `
  ><((°>  ~
  |||||
 axolotl
`,
		Names: []string{"Axl", "Wooper", "Gel", "Nix", "Bloop", "Aqua", "Gill", "Regen"},
	},
	{
		ID:    "mushroom",
		Name:  "Mushroom",
		Emoji: "🍄",
		Attributes: map[string]int{
			"wisdom": 9,
			"poison": 6,
		},
		Traits: []string{"WISDOM", "POISON"},
		AsciiArt: `
   _( )_
  (  *  )
   |   |
  mushroom
`,
		Names: []string{"Spore", "Myco", "Fungi", "Shiitake", "Trufa", "Velvet", "Cap", "Hifa"},
	},
	{
		ID:    "ghost",
		Name:  "Ghost",
		Emoji: "👻",
		Attributes: map[string]int{
			"sarcasm": 8,
			"stealth": 9,
		},
		Traits: []string{"SARCASM", "STEALTH"},
		AsciiArt: `
   .-.
  ( o o)
  |     |
  '~~~~~'
   ghost
`,
		Names: []string{"Shade", "Boo", "Wisp", "Null", "Vex", "Phantom", "Mist", "Void"},
	},
	{
		ID:    "dragon",
		Name:  "Dragon",
		Emoji: "🐉",
		Attributes: map[string]int{
			"power": 10,
			"chaos": 9,
		},
		Traits: []string{"POWER", "CHAOS"},
		AsciiArt: `
  /\___/\
 ( >   < )
  \  ^  /
   '---'
   dragon
`,
		Names: []string{"Ignis", "Vyrn", "Ember", "Drako", "Scorch", "Rex", "Ryu", "Blaze"},
	},
	{
		ID:    "duck",
		Name:  "Duck",
		Emoji: "🦆",
		Attributes: map[string]int{
			"debugging":    8,
			"stubbornness": 9,
		},
		Traits: []string{"DEBUGGING", "STUBBORNNESS"},
		AsciiArt: `
   _
 <(o )___
  ( ._> /
   '---'
    duck
`,
		Names: []string{"Duck", "Quak", "Rubber", "Debug", "Waddler", "Quill", "Pip", "Ducky"},
	},
	{
		ID:    "cat",
		Name:  "Cat",
		Emoji: "🐱",
		Attributes: map[string]int{
			"independence": 10,
			"irony":        8,
		},
		Traits: []string{"INDEPENDENCE", "IRONY"},
		AsciiArt: `
  /\_/\
 ( o.o )
  > ^ <
   cat
`,
		Names: []string{"Nox", "Pixel", "Byte", "Misty", "Sable", "Zephyr", "Rune", "Kira"},
	},
	{
		ID:    "owl",
		Name:  "Owl",
		Emoji: "🦉",
		Attributes: map[string]int{
			"wisdom":   10,
			"patience": 7,
		},
		Traits: []string{"WISDOM", "PATIENCE"},
		AsciiArt: `
  ,___,
 (o) (o)
  -beak-
  m   m
   owl
`,
		Names: []string{"Hoot", "Sage", "Luna", "Athena", "Orion", "Wren", "Noctua", "Dusk"},
	},
}

// ByID returns a species by its ID, or nil if not found.
func ByID(id string) *Species {
	for i := range All {
		if All[i].ID == id {
			return &All[i]
		}
	}
	return nil
}

// FromHash picks a species deterministically from a hash value.
func FromHash(h uint64) *Species {
	return &All[h%uint64(len(All))]
}

// NameFromHash picks a name from the species name list deterministically.
func NameFromHash(sp *Species, h uint64) string {
	return sp.Names[h%uint64(len(sp.Names))]
}
