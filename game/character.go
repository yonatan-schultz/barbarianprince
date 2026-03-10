package game

// CharacterType distinguishes the prince from followers
type CharacterType int

const (
	TypePrince CharacterType = iota
	TypeMercenary
	TypeAmazon
	TypeDwarf
	TypeElf
	TypeHalfling
	TypePriest
	TypeWizard
	TypeSwordsman
	TypeBandit
	TypeGuide
	TypeEscapee
	TypeGeneric
)

// MountType represents the type of mount a character has
type MountType int

const (
	MountNone    MountType = iota
	MountHorse
	MountWarhorse
	MountPegasus
)

// PossessionType represents special items
type PossessionType int

const (
	PossNone           PossessionType = iota
	PossRoyalHelm                     // win condition item
	PossStaffOfCommand                // win condition item
	PossRingOfCommand                 // +2 CS in all combat
	PossAmuletOfPower                 // +1 CS
	PossElvenBoots                    // ignore lost checks in forest
	PossRopeAndGrapnel                // cross 1 mountain hex/day extra
	PossLantern                       // search ruins night bonus
	PossPoisonAntidote                // cure poison wounds
	PossHealingPotion                 // restore 1d6 endurance
	PossMagicSword                    // +2 CS, counts as wound on 9+
	PossHolySymbol                    // repel undead in ruins
	PossNobleParchment                // proof of noble ally
	PossMap                           // reveals adjacent hexes
	PossRaft                          // river travel
	PossGoldenCrown                   // royal helm substitute
)

// Character represents the prince or a follower
type Character struct {
	Name           string
	Type           CharacterType
	CombatSkill    int
	MaxEndurance   int
	Wounds         int
	PoisonWounds   int
	WitWiles       int    // for avoiding traps / social rolls
	WealthCode     int    // for treasure rolls (followers that share)
	DailyWage      int    // in gold pieces
	HasMount       bool
	MountType      MountType
	StarvationDays int    // consecutive days without food
	Possessions    []PossessionType
	IsGuide        bool
	IsTrueLove     bool
	IsEscapee      bool  // freed prisoner, will desert if in town
	MustDesert     bool  // flagged for desertion
	DaysHired      int   // how many days hired for
	Morale         int   // 1-6, lower = deserts sooner
	MaxMorale      int   // starting morale cap; recovered toward this value each week
	LoadCapacity   int   // units of food/gold this follower can carry (porters: 5, others: 0)
	IsUndead         bool // true for undead enemies (affected by Holy Symbol)
	PlagueDustActive bool // r227: deals wounds daily until recovery roll (1d6 >= 4)
}

// CurrentEndurance returns remaining endurance
func (c *Character) CurrentEndurance() int {
	return c.MaxEndurance - c.Wounds - c.PoisonWounds
}

// IsAlive returns true if the character still has endurance
func (c *Character) IsAlive() bool {
	return c.CurrentEndurance() > 0
}

// IsDead returns true if wounds exceed endurance
func (c *Character) IsDead() bool {
	return !c.IsAlive()
}

// IsUnconscious returns true when wounds = endurance-1 (one step from death, r221b)
func (c *Character) IsUnconscious() bool {
	return c.Wounds == c.MaxEndurance-1 && c.Wounds > 0
}

// EffectiveCombatSkill applies wound, starvation, and possession bonuses (prince only)
func (c *Character) EffectiveCombatSkill() int {
	cs := c.CombatSkill
	// Wound penalty: -1 CS per wound (r220c: -1 if any wound, -2 if seriously wounded)
	if c.Wounds > 0 {
		cs--
	}
	if c.Wounds*2 >= c.MaxEndurance {
		cs-- // seriously wounded: additional -1
	}
	// Starvation penalty (r216b): -1 CS per consecutive day without food
	cs -= c.StarvationDays
	// Poison wounds reduce CS too
	cs -= c.PoisonWounds / 3
	// Possession bonuses — only the prince carries these items
	if c.HasPossession(PossRingOfCommand) {
		cs += 2
	}
	if c.HasPossession(PossAmuletOfPower) {
		cs += 1
	}
	if c.HasPossession(PossMagicSword) {
		cs += 2
	}
	if cs < 1 {
		cs = 1
	}
	return cs
}

// HasPossession checks if a character has a specific item
func (c *Character) HasPossession(p PossessionType) bool {
	for _, pos := range c.Possessions {
		if pos == p {
			return true
		}
	}
	return false
}

// AddPossession adds an item to the character's inventory
func (c *Character) AddPossession(p PossessionType) {
	c.Possessions = append(c.Possessions, p)
}

// RemovePossession removes the first instance of an item
func (c *Character) RemovePossession(p PossessionType) bool {
	for i, pos := range c.Possessions {
		if pos == p {
			c.Possessions = append(c.Possessions[:i], c.Possessions[i+1:]...)
			return true
		}
	}
	return false
}

// NewPrince creates the starting character for Cal Arath
func NewPrince() Character {
	return Character{
		Name:         "Cal Arath",
		Type:         TypePrince,
		CombatSkill:  5,
		MaxEndurance: 9,
		WitWiles:     4,
		WealthCode:   5,
		Morale:       6,
	}
}

// Cache represents a hidden stash of treasure
type Cache struct {
	Location HexID
	Gold     int
	Found    bool
}

// PossessionName returns the display name of a possession
func PossessionName(p PossessionType) string {
	switch p {
	case PossRoyalHelm:
		return "Royal Helm"
	case PossStaffOfCommand:
		return "Staff of Command"
	case PossRingOfCommand:
		return "Ring of Command"
	case PossAmuletOfPower:
		return "Amulet of Power"
	case PossElvenBoots:
		return "Elven Boots"
	case PossRopeAndGrapnel:
		return "Rope & Grapnel"
	case PossLantern:
		return "Lantern"
	case PossPoisonAntidote:
		return "Poison Antidote"
	case PossHealingPotion:
		return "Healing Potion"
	case PossMagicSword:
		return "Magic Sword"
	case PossHolySymbol:
		return "Holy Symbol"
	case PossNobleParchment:
		return "Noble Parchment"
	case PossMap:
		return "Ancient Map"
	case PossRaft:
		return "Raft"
	case PossGoldenCrown:
		return "Golden Crown"
	}
	return "Unknown Item"
}
