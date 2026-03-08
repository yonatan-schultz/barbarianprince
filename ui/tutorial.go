package ui

import (
	"fmt"
	"strings"
)

type tutorialEntry struct {
	Title string
	Body  string
}

var tutorialSteps = []tutorialEntry{
	{
		Title: "Welcome, Cal Arath",
		Body: `You are a deposed prince. You have 70 days to
reclaim your throne. There are four ways to win:

  • Amass 500 gold and return north of the
    Tragoth River (shown as ~~~~ on the map)
  • Find the Royal Helm and bring it to
    Ogon or Weshor in the north
  • Secure a noble ally and march north
  • Obtain the Staff of Command and go north

You begin in Ogon with 10 gold and 14 food.`,
	},
	{
		Title: "Reading the Map",
		Body: `Named places show their first 3 letters:
  "Ogo" = Ogon (town)   "Hul" = Hulora (castle)

Unnamed structures show a symbol:
  [T] Town    [C] Castle  [K] Keep
  [R] Ruins   [+] Temple  [v] Village

Other symbols:
  [*] You   >>> Travel target   ~~~~ River
  =   Road (brown separator between hexes)

Terrain (affects travel safety & food):
  . Farmland   ~ Countryside   f Forest
  ^ Hills      M Mountains     s Swamp
  o Desert

Dangerous terrain increases the chance of
getting lost and triggering encounters.`,
	},
	{
		Title: "The Status Bar",
		Body: `The bar at the bottom tracks your resources:

  Day (1/70)  — time is always running out
  Gold        — pay wages, buy food, bribe
  Food        — 1 unit feeds 1 person/day
  HP bar (♥)  — endurance; 0 = death

  POISON:N    — active poison wounds

Starvation: each day without food, your
Combat Skill drops by 1 and followers may
desert. The penalty resets when you eat.

Poison wounds accumulate separately from
normal wounds. Use a Poison Antidote [U]
or rest to slowly purge it.`,
	},
	{
		Title: "Daily Actions — Anywhere",
		Body: `Each day you choose exactly one action.
Actions available in any hex:

  [T] Travel     Move to an adjacent hex
  [R] Rest       Recover wounds (slow outdoors,
                 faster in settlements)
  [G] Hunt       Forage for food
                 (blocked in mountains, desert,
                 and settlements; roll 12 = wounds)
  [U] Use Item   Use a consumable you carry
                 (Healing Potion or Poison Antidote)`,
	},
	{
		Title: "Daily Actions — Settlements",
		Body: `Additional actions in towns, castles, temples:

  [N] Seek News    Rumours, rewards, encounters
  [H] Hire         Recruit followers and mounts
  [B] Buy Food     2 gold per food unit
  [P] Buy Raft     15 gold — needed to cross rivers
                   (see: Rivers & Rafts)
  [A] Audience     Seek a lord's favour
                   (castles and keeps only)
  [O] Offering     Submit 10 gold at a temple for
                   a blessing (or a curse...)`,
	},
	{
		Title: "Daily Actions — Special Locations",
		Body: `Actions that appear only in certain hexes:

  [S] Search Ruins   Explore an unsearched ruin
                     for treasure, events, or danger.
                     Each ruin can only be searched
                     once. A lantern improves odds.

  [C] Search Cache   Recover gold you hid earlier.
                     A cache can be lost (roll of 5)
                     or already looted (roll of 6).

Caches are created by certain events during
play — they appear as an action when you
return to the hiding location.`,
	},
	{
		Title: "Travelling",
		Body: `Press [T] to enter travel mode.
Use [↑↓] or [1-6] to select a direction.
Press [Enter] to travel there.

Every journey can trigger an event:
  bandits, merchants, magic, weather, traps.

Getting lost: bad terrain (mountains, swamp,
desert, forest) rolls against your skill.
A guide follower reduces this chance.
Elven Boots suppress getting lost in forest.

Road travel: brown = separators on the map
mark roads. Roads suppress getting-lost rolls
and use the road event table (usually safer).

Mounted parties move two hexes per day.
Press [Esc] after the first hop to make camp.`,
	},
	{
		Title: "Rivers & Rafts",
		Body: `Rivers appear as ~ borders between hexes.
The Tragoth River (~~~~~) divides the map
east-west — you must cross it to win.

Crossing a river on a road uses the bridge —
no raft needed, no lost roll. The main road
crosses the Tragoth at column 8 (look for =
on the map leading to the ~~~~ separator).

Crossing a river off-road requires a raft.
Without one, you cannot cross.
Buy a raft with [P] at any settlement (15 gold).

Plan ahead: if you go south without a raft,
head for the col-8 bridge or buy one before
crossing. Running out of options by day 60
with no raft and no road access is fatal.`,
	},
	{
		Title: "Combat",
		Body: `Encounters may start a fight:

  [F] Fight   — one round of combat
  [R] Retreat — attempt to flee (risky;
                enemy gets a free attack)

Combat Skill (CS) determines the outcome.
Each round, the weaker side takes 1 wound.
Special rolls can deal 2 wounds at once
or cause a rout (enemy flees, no loot).

Followers add their CS to your total.

At 0 endurance, Cal Arath dies — game over.
If unconscious, followers may carry you or
desert, taking all your gold with them.`,
	},
	{
		Title: "Wounds, Poison & Healing",
		Body: `Wounds reduce your effective Combat Skill:
  1+ wound    — CS −1
  Seriously wounded (wounds ≥ half max) — CS −2

Healing:
  [R] Rest in a settlement — 1 wound/day
  [R] Rest in the wilderness — 1 wound on 4+
  Healing Potion [U] — recover 1d6 wounds

Poison wounds come from traps and enemies.
They accumulate on top of normal wounds.
  Poison Antidote [U] — clears all poison
  Rest in settlement   — −1 poison/day
  Rest in wilderness   — −1 poison on 4+

Plague dust (rare trap) deals wounds each
day until your body fights it off (1d6, 4+).`,
	},
	{
		Title: "Followers",
		Body: `Hire followers in towns with [H] (2d6 table).
Each follower costs a daily wage in gold.
Unpaid followers may desert overnight.

Follower types:
  Guide      — reduces getting-lost chance
  Mercenary  — strong fighters (CS 4, wage 2)
  Lancer     — mounted, high CS (5), wage 3
  Henchman   — cheap fighters (CS 3, wage 1)

Mounts eat 2 food units per day in rough
terrain (mountains, desert, swamp, sea).
In open terrain (farmland, countryside,
forest, hills) mounts forage for free.

Lodging costs 1 gold per person per night
in settlements, plus 1 per mount stabled.
If you can't pay, followers may desert.`,
	},
	{
		Title: "Settlements: Towns & Castles",
		Body: `Towns (marker T) offer the most services:
  Buy food, hire followers, seek news,
  buy a raft, and rest with faster healing.

Castles and Keeps (C / K):
  Seek Audience [A] — roll 2d6 on a table
  unique to each lord. Outcomes range from
  being thrown out or barred, to meeting the
  lord's family, bribing a seneschal, or
  winning a duel and earning gold or a
  noble ally. Each castle has different odds.

Villages (v) are small settlements — they
offer the same services as towns but may
have fewer options available.`,
	},
	{
		Title: "Settlements: Temples",
		Body: `Temples (marker +) offer the Submit Offering
action [O]. Cost: 10 gold.

Outcomes are rolled on a 1d6 table and
vary widely:
  • A divine blessing (bonus to next roll)
  • Gold reward or food
  • A cursed item or wound
  • A unique magical event

Hidden temples can be discovered through
events and exploration. They work the same
way once found.

Offering at the wrong time or place can
attract bad fortune — but the rewards can
be significant.`,
	},
	{
		Title: "Possessions & Items",
		Body: `Items are gained through events, ruins,
audience outcomes, and purchases.

Usable consumables ([U] Use Item):
  Healing Potion    — recover 1d6 wounds
  Poison Antidote   — clear all poison wounds

Passive combat bonuses (always active):
  Magic Sword       — CS +2; extra wound on 9+
  Ring of Command   — CS +2
  Amulet of Power   — CS +1
  Holy Symbol       — CS +2 vs undead enemies

Passive exploration bonuses:
  Elven Boots  — no getting-lost in forest
  Lantern      — improves ruin search rolls

Key quest items: Royal Helm, Staff of
Command, Golden Crown — win conditions
when returned to the right place.`,
	},
	{
		Title: "Economy & Time Management",
		Body: `Gold flows out fast. Track your spending:

Daily costs:
  Food       — 2 gold per 10 units bought
  Wages      — 1–3 gold per follower
  Lodging    — 1 gold per person + mount

Income sources:
  Combat loot, ruins, seek news, offerings,
  audience rewards, caches, events.

Time: you have 70 days. Travel takes 1 day
per hex (or 2 hexes if all mounted). Rest,
hire, buy, and audience each take 1 day.

Getting lost adds a day. Bad events can
cost days. Plan a route north early — the
journey south and back takes many days.

Save often with [Ctrl+S].`,
	},
	{
		Title: "You Are Ready",
		Body: `You start in Ogon (north) with 10 gold,
14 food, and Cal Arath alone.

Key references during play:
  [?]        Key reference in the log
  [F]        Field notes (discoveries recorded)
  [Ctrl+S]   Save your progress

The Tragoth River (~~~~~) runs east-west
across the middle of the map. You start
north of it. If you go south to find wealth,
plan your return crossing: the road bridge
at column 8 crosses free (no raft needed),
or buy a raft (15 gold) at any settlement.

Discoveries are recorded in your field
notes automatically when found.

Good luck, Cal Arath. The throne awaits.`,
	},
}

// RenderTutorial renders one tutorial slide for the menu panel.
// innerWidth is the available content width (chars) inside the panel borders.
func RenderTutorial(step, innerWidth int) string {
	if step < 0 || step >= len(tutorialSteps) {
		return ""
	}
	t := tutorialSteps[step]
	var lines []string
	lines = append(lines, StyleTitle.Render("── Tutorial ──"))
	lines = append(lines, "")
	lines = append(lines, StyleValue.Render(t.Title))
	lines = append(lines, "")
	for _, line := range strings.Split(t.Body, "\n") {
		for _, wrapped := range wrapLine(line, innerWidth) {
			lines = append(lines, StyleMenuText.Render(wrapped))
		}
	}
	lines = append(lines, "")
	lines = append(lines, StyleLabel.Render(
		fmt.Sprintf("[%d/%d]  Space/Enter next  Esc skip", step+1, len(tutorialSteps))))
	return strings.Join(lines, "\n")
}

