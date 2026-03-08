package game

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const saveVersion = 1

// SaveFile wraps GameState with a version for forward-compatibility checks
type SaveFile struct {
	Version int        `json:"version"`
	State   *GameState `json:"state"`
}

// savePath returns the path to the save file, creating the directory if needed
func savePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not determine home directory: %w", err)
	}
	dir := filepath.Join(home, ".barbarianprince")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("could not create save directory: %w", err)
	}
	return filepath.Join(dir, "save.json"), nil
}

// SaveExists reports whether a save file is present and readable
func SaveExists() bool {
	path, err := savePath()
	if err != nil {
		return false
	}
	_, err = os.Stat(path)
	return err == nil
}

// Save writes the current game state to disk.
// Returns an error message suitable for display, or "" on success.
func Save(s *GameState) error {
	// Only save from a clean action-select state — not mid-combat or mid-event
	if s.Phase == PhaseCombat {
		return errors.New("cannot save during combat")
	}
	if len(s.PendingChoices) > 0 {
		return errors.New("cannot save while a choice is pending")
	}

	path, err := savePath()
	if err != nil {
		return err
	}

	sf := SaveFile{Version: saveVersion, State: s}
	data, err := json.MarshalIndent(sf, "", "  ")
	if err != nil {
		return fmt.Errorf("could not encode save data: %w", err)
	}

	// Write atomically via temp file
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return fmt.Errorf("could not write save file: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		return fmt.Errorf("could not finalize save file: %w", err)
	}
	return nil
}

// Load reads a saved game state from disk
func Load() (*GameState, error) {
	path, err := savePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read save file: %w", err)
	}

	var sf SaveFile
	if err := json.Unmarshal(data, &sf); err != nil {
		return nil, fmt.Errorf("save file is corrupt: %w", err)
	}
	if sf.Version != saveVersion {
		return nil, fmt.Errorf("save file version %d is incompatible (expected %d)", sf.Version, saveVersion)
	}
	if sf.State == nil {
		return nil, errors.New("save file contains no game state")
	}

	// Restore nil maps that JSON unmarshal leaves as nil
	if sf.State.HexFlags == nil {
		sf.State.HexFlags = make(map[HexID]*HexFlags)
	}
	if sf.State.VisitedHexes == nil {
		sf.State.VisitedHexes = make(map[HexID]bool)
	}
	if sf.State.AudienceBarred == nil {
		sf.State.AudienceBarred = make(map[HexID]int)
	}
	if sf.State.Flags.DragonSlain == nil {
		sf.State.Flags.DragonSlain = make(map[HexID]bool)
	}
	if sf.State.Flags.WizardsTowerVisited == nil {
		sf.State.Flags.WizardsTowerVisited = make(map[HexID]bool)
	}
	if sf.State.Flags.SecretFound == nil {
		sf.State.Flags.SecretFound = make(map[string]bool)
	}
	for id, hf := range sf.State.HexFlags {
		if hf != nil && hf.EventUsed == nil {
			sf.State.HexFlags[id].EventUsed = make(map[string]bool)
		}
	}

	// Reset any transient combat/event state that shouldn't persist across saves
	sf.State.ActiveEnemy = nil
	sf.State.CombatLog = nil
	sf.State.PendingChoices = nil
	sf.State.Phase = PhaseActionSelect

	return sf.State, nil
}

// DeleteSave removes the save file
func DeleteSave() error {
	path, err := savePath()
	if err != nil {
		return err
	}
	err = os.Remove(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}
