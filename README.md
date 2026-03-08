# Barbarian Prince

A terminal TUI implementation of the 1981 solo board game **Barbarian Prince**, written in Go.

## About

You are Cal Arath, a deposed prince. You have **70 days** to accumulate **500 gold** — or achieve one of several alternate win conditions — to reclaim your throne. Travel a hex map, manage food and followers, fight monsters, and navigate hundreds of random events.

This project implements the full game in a playable terminal interface, including:

- Full 25×24 hex map with terrain, settlements, ruins, castles, and temples
- Travel system with lost checks, river crossings (requires raft), and road bonuses
- Combat with wound penalties, followers, special weapons, and undead modifiers
- ~170 random events (wilderness encounters, NPCs, treasures, audience/temple mechanics)
- 15 special possessions (Ring of Command, Magic Sword, Elven Boots, etc.)
- Win/lose conditions, food economy, follower wages, and starvation
- True Love mechanics — loyal companion with special desertion immunity and W&W bonus
- Save/load support (single slot, `~/.barbarianprince/save.json`)
- Interactive tutorial for new players

## Running

Requires Go 1.24+.

```
go run .
```

## Controls

Navigate menus with the listed key bindings shown on screen. The game runs in alternate screen mode with mouse support enabled.

## Project Structure

```
barbarianprince/
├── main.go            # entry point
├── game/              # all game logic (state, map, combat, events, economy)
└── ui/                # bubbletea TUI (map view, log pane, menus, styles)
```

## Built With

- [bubbletea](https://github.com/charmbracelet/bubbletea) — terminal UI framework
- [lipgloss](https://github.com/charmbracelet/lipgloss) — terminal styling

## Original Game

**Barbarian Prince** was designed by Arnold Hendrick and Robert J. Bartash, published by Dwarfstar Games in 1981. It is one of the most celebrated solitaire wargames ever made.

- Original rules and game files: https://dwarfstar.brainiac.com/ds_barbarianprince.html
- BoardGameGeek page: https://boardgamegeek.com/boardgame/1631/barbarian-prince

This project is a fan implementation for personal and educational use. All game content and mechanics are based on the original Dwarfstar Games publication.
