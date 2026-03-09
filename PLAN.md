# Barbarian Prince ASCII Game — Go Implementation Plan

## Context
Implement a playable terminal TUI version of the 1981 solo board game "Barbarian Prince" in Go. The player is Cal Arath, a deposed prince who has 70 days to accumulate 500 gold (or achieve an alternate win condition) to reclaim his throne. The game is driven by a hex-map travel system, random event tables, combat, and resource management (food, gold, followers).

**Location**: `~/go/src/barbarianprince/` — standalone Go module
**Scope**: Full playable MVP — all core systems + ~50 highest-impact events
**UI**: Terminal TUI using [bubbletea](https://github.com/charmbracelet/bubbletea) + [lipgloss](https://github.com/charmbracelet/lipgloss)

---

## Status

### ✅ MVP Complete
All phases below marked ✅ are implemented and passing tests.

---

## Project Structure

```
barbarianprince/
├── go.mod                    # module barbarianprince
├── main.go                   # entry point, bubbletea setup
├── game/
│   ├── state.go             # GameState, turn loop, win/lose checks
│   ├── hex.go               # Hex, HexID, TerrainType, adjacency
│   ├── map.go               # full 25×24 hex map definition
│   ├── character.go         # Character, Prince, followers
│   ├── combat.go            # strike resolution, combat loop
│   ├── travel.go            # movement, lost checks, event lookup
│   ├── tables.go            # all lookup tables (r207, r220c, r226, r231-r281)
│   ├── events.go            # event dispatcher
│   ├── events_e001.go       # events e001–e050
│   ├── events_e051.go       # events e051–e100
│   ├── events_e100.go       # events e101–e170 (airborne, ruins, audience, secrets)
│   ├── events_e180.go       # special possessions e180–e194
│   ├── dice.go              # Roll1d6(), Roll2d6(), Roll1d3()
│   ├── economy.go           # food, lodging, starvation, treasure
│   ├── actions.go           # daily action handlers (travel, rest, search, hire, audience, offering)
│   └── flags.go             # global/per-hex persistent flags
└── ui/
    ├── model.go             # bubbletea Model, Update, View
    ├── map_view.go          # ASCII hex grid renderer
    ├── log_view.go          # narrative text scroll pane
    ├── menu_view.go         # action menu + prompts
    └── styles.go            # lipgloss color/style definitions
```

---

## Implementation Phases

### ✅ Phase 1 — Project Scaffold
- Go module, bubbletea/lipgloss deps, main.go

### ✅ Phase 2 — Core Data
- All type enums (Terrain, Structure, Direction, CharacterType, etc.)
- Full 25×24 hex map (~600 hexes) with terrain, structures, rivers, roads
- Travel Table r207, Combat Table r220c, Treasure Table r226, Event refs r231–r281

### ✅ Phase 3 — Dice & Economy
- Roll1d6/2d6/1d3, food consumption, hunting, starvation, lodging, follower desertion, treasure

### ✅ Phase 4 — Character & Combat
- Wound penalties, strike resolution, combat loop, rout/escape, surprise, unconscious/death

### ✅ Phase 5 — Travel & Actions
- Adjacency, lost checks, event triggers, road/river travel, all daily action handlers

### ✅ Phase 6 — Event System (~50 events)
- e001–e050: start, NPCs, settlements, common encounters
- e051–e100: wilderness, monsters, environmental hazards, big monsters
- e100–e161: airborne, ruins, secrets, audience, temple offerings
- e180–e194: all 15 special possessions

### ✅ Phase 7 — Win/Lose & Turn Loop
- All 4 win conditions, 3 lose conditions, AdvanceDay with food/wages/desertion

### ✅ Phase 8 — TUI
- Map grid, log pane, status bar, action/combat/travel/event menus, game over screen

### ✅ Phase 9 — Gameplay Polish
- **Possession CS bonuses**: Ring of Command +2, Amulet of Power +1, Magic Sword +2 in `EffectiveCombatSkill()`
- **Magic Sword extra wound**: +1 wound when net roll ≥ 9
- **Elven Boots**: suppresses lost rolls in Forest terrain
- **Holy Symbol**: +2 CS vs undead enemies (`IsUndead` flag on Character)
- **ActionUseItem**: Healing Potion and Poison Antidote usable from action menu (`[U]`)
- **ActionHunt**: Hunt for food in non-farmland hexes (`[G]`)
- **Follower combat**: `TotalCombatSkill()` (prince + all followers) used as party attack CS
- **Poison display**: `StatusLine()` shows poison wounds when > 0
- **Poison recovery**: `HealRest()` reduces poison wounds over time
- **e151 gold bug fixed**: duel prize only awarded on combat victory via `PendingDuelGold`
- **Enemy fled message**: "The enemy flees! You find nothing of value." instead of "Victory! Gained 0 gold."
- **River crossing enforcement**: `DoTravel` checks `RiverSides[dir]`; blocks travel without `PossRaft`; raft has 1-in-6 chance of being wrecked on crossing

---

## Future Work (Post-MVP)

### ✅ Tutorial Mode
Interactive guided introduction for new players. Runs before the main game (or can be selected from a start menu).

**Implemented:**
- 18-slide reference guide (`ui/tutorial.go`) shown in the menu panel before gameplay
- In-game hint bar (4 guided steps: buy food → travel → see event → rest) stored in `GameState.Tutorial *TutorialState`
- Hint displayed at top of action menu while tutorial is active
- Skippable at any time with `[X]` Skip
- Completion persisted in `~/.barbarianprince/tutorial_done` — returning players skip hints by default
- `--tutorial` CLI flag forces in-game hints even after completion
- `game/tutorial_state.go` contains all TutorialState logic; `ui/tutorial.go` contains slide text

**Not yet implemented:**
- Pre-seeded RNG for deterministic tutorial encounters

### ✅ Save / Load Game
- JSON serialization of `GameState` to `~/.barbarianprince/save.json`
- `[S]ave` and `[L]oad` options on the start menu and in-game pause menu
- Single save slot for simplicity

### ✅ River Rendering on Map
- `initRiversAndRoads()` populates `RiverSides` on all world-map hexes
- **Nesser River** (north-south, between cols 12–13): blue `~` separator in every row of RenderMapGrid
- **Tragoth River** (east-west, between rows 11–12): thin `~~~~` blue band inserted between rows 11 and 12
- `linesUsed` counter prevents the separator from exceeding `maxHeight`

### 🔲 Raft Travel (downstream fast-travel)
- **Blocking done** — river crossings now require `PossRaft`; raft wrecked on roll of 1
- Remaining: downstream fast-travel action (`ActionTakeRaft`) — moves 1–3 river hexes in one day
- Need: define downstream direction per river (Tragoth flows east→west, Nesser flows north→south)
- New action visible only when on a river-boundary hex carrying `PossRaft`

### 🔲 Flat-Top Hex Map Rendering
- The real BP map uses **flat-top hexes** (⬡): flat N and S sides, 6 neighbors N/NE/SE/S/SW/NW, no E/W
- Our direction system is already correct; only the visual renderer uses a square grid
- `RenderMapGrid` needs a rewrite to stagger columns: odd columns displayed half a row lower than even columns (2 display lines per hex row)
- Requires verifying the actual map data in `map.go` matches the original game before investing in a pixel-accurate renderer
- **Dependency**: Map Data Verification (below) should come first

### 🔲 Map Data Verification
- Our `map.go` has 25×24 hexes but the original BP map is 20 columns × 23 rows (CCRR hex numbering 0101–2023)
- Terrain, structure, and name data in `map.go` was hand-entered and likely has errors
- Options: (a) manually cross-reference against the original map image/BGG hex data, (b) import a fan-created CSV if one exists
- Rivers (Tragoth, Nesser) and settlement locations should be verified before adding the hex renderer
- A fan-made Google Sheet of all BP hexes exists on BGG — importing from that would be the cleanest path

### ✅ Full Event Coverage
- events_e110.go: 40 new events (e110–e119, e122–e127, e139–e142, e162–e179)
- Swamp terrain: dedicated sw1–sw6 event tables (no longer reuses Hills tables)
- Desert: r281/r282 alternating tables; ruins: expanded from 1d6 to 2d6 (11 outcomes)
- News, audience, farmland/countryside/forest/hills/mountains tables all updated with new NPC outcomes

### ✅ True Love Mechanics (r228)
- NPC follower flagged `IsTrueLove`
- Never deserts (starvation or unpaid wages), +1 W&W while present
- Prince unconscious: True Love guarantees rescue with no roll
- Met via seek-news (roll 11) or seek-followers (roll 12), one-shot per game (e195)

### 🔲 Full Follower Economy
- Porter load tracking (each porter carries N food/gold units)
- Lancer daily tracking (separate mount food cost)
- Morale system: followers check morale weekly based on food, pay, wounds

### 🔲 Magic Items (Complex State)
- Alcove of Sending, Arch of Travel, Gateway to Darkness, Mirror of Reversal
- Require multi-turn state and new UI flows

### 🔲 Test Coverage Expansion
Current: **58.7%**. Target: **70%+**

Priority additions:
- `TestResolveCombatRound` — wounds applied, enemy death, loot
- `TestAttemptFlee` — succeeds on roll ≥4, damages on failure
- `TestAdvanceDay` — food consumed, day advances, starvation accrues, wages deducted
- `TestBuyFood` / `TestHuntForFood` — gold deducted, food added
- `TestAvailableActions` — correct actions per hex structure
- `TestDoTravel` — valid/invalid moves, visited hex tracking

### 🔲 Sound / Accessibility
- Terminal bell on combat, death, win
- Colorblind-friendly palette option (`--no-color`)

### 🔲 Difficulty Modes
- **Easy**: start with 20 gold + 20 food; 80-day limit
- **Hard**: start with 5 gold + 7 food; harsher event tables
- Selectable from start menu; stored in `GameState`
