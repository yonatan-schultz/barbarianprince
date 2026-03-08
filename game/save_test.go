package game

import (
	"os"
	"path/filepath"
	"testing"
)

// overrideSavePath points saves at a temp dir for testing
func withTempSave(t *testing.T, fn func()) {
	t.Helper()
	dir := t.TempDir()
	orig := os.Getenv("HOME")
	os.Setenv("HOME", dir)
	defer os.Setenv("HOME", orig)
	fn()
	// Ensure save dir was created if Save was called
	_ = filepath.Join(dir, ".barbarianprince")
}

func TestSaveAndLoad(t *testing.T) {
	withTempSave(t, func() {
		s := NewGameState()
		s.Gold = 123
		s.Day = 15
		s.FoodUnits = 42
		s.Prince.Wounds = 3
		s.AddFollower(Character{Name: "Guard", CombatSkill: 4, MaxEndurance: 8, DailyWage: 3})
		s.VisitedHexes[NewHexID(2, 2)] = true
		s.Flags.HasRoyalHelm = true

		if err := Save(s); err != nil {
			t.Fatalf("Save() error: %v", err)
		}

		if !SaveExists() {
			t.Fatal("SaveExists() should return true after saving")
		}

		loaded, err := Load()
		if err != nil {
			t.Fatalf("Load() error: %v", err)
		}

		if loaded.Gold != 123 {
			t.Errorf("Gold = %d, want 123", loaded.Gold)
		}
		if loaded.Day != 15 {
			t.Errorf("Day = %d, want 15", loaded.Day)
		}
		if loaded.FoodUnits != 42 {
			t.Errorf("FoodUnits = %d, want 42", loaded.FoodUnits)
		}
		if loaded.Prince.Wounds != 3 {
			t.Errorf("Prince.Wounds = %d, want 3", loaded.Prince.Wounds)
		}
		if len(loaded.Party) != 1 || loaded.Party[0].Name != "Guard" {
			t.Errorf("Party = %v, want 1 follower named Guard", loaded.Party)
		}
		if !loaded.VisitedHexes[NewHexID(2, 2)] {
			t.Error("VisitedHexes not preserved")
		}
		if !loaded.Flags.HasRoyalHelm {
			t.Error("Flags.HasRoyalHelm not preserved")
		}
	})
}

func TestSaveBlockedDuringCombat(t *testing.T) {
	withTempSave(t, func() {
		s := NewGameState()
		s.Phase = PhaseCombat
		if err := Save(s); err == nil {
			t.Error("Save() should fail during combat")
		}
	})
}

func TestSaveBlockedWithPendingChoice(t *testing.T) {
	withTempSave(t, func() {
		s := NewGameState()
		s.PendingChoices = []string{"Fight", "Flee"}
		if err := Save(s); err == nil {
			t.Error("Save() should fail with pending choices")
		}
	})
}

func TestLoadMissingFile(t *testing.T) {
	withTempSave(t, func() {
		if SaveExists() {
			t.Error("SaveExists() should be false with no save file")
		}
		_, err := Load()
		if err == nil {
			t.Error("Load() should fail when no save file exists")
		}
	})
}

func TestLoadRestoresNilMaps(t *testing.T) {
	withTempSave(t, func() {
		s := NewGameState()
		if err := Save(s); err != nil {
			t.Fatalf("Save() error: %v", err)
		}

		loaded, err := Load()
		if err != nil {
			t.Fatalf("Load() error: %v", err)
		}

		// All maps must be non-nil after load
		if loaded.HexFlags == nil {
			t.Error("HexFlags is nil after load")
		}
		if loaded.VisitedHexes == nil {
			t.Error("VisitedHexes is nil after load")
		}
		if loaded.AudienceBarred == nil {
			t.Error("AudienceBarred is nil after load")
		}
		if loaded.Flags.DragonSlain == nil {
			t.Error("Flags.DragonSlain is nil after load")
		}
	})
}

func TestLoadResetsTransientState(t *testing.T) {
	withTempSave(t, func() {
		s := NewGameState()
		s.Phase = PhaseCombat
		enemy := MakeEnemy("Orc", 4, 8, 2)
		s.ActiveEnemy = &enemy
		s.CombatLog = []string{"Strike!", "Hit!"}
		// Force save by temporarily clearing phase
		s.Phase = PhaseActionSelect
		if err := Save(s); err != nil {
			t.Fatalf("Save() error: %v", err)
		}

		loaded, err := Load()
		if err != nil {
			t.Fatalf("Load() error: %v", err)
		}

		if loaded.ActiveEnemy != nil {
			t.Error("ActiveEnemy should be nil after load")
		}
		if len(loaded.CombatLog) != 0 {
			t.Error("CombatLog should be empty after load")
		}
		if loaded.Phase != PhaseActionSelect {
			t.Errorf("Phase = %v, want PhaseActionSelect", loaded.Phase)
		}
	})
}
