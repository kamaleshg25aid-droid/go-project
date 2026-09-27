package bloodbank

import (
	"path/filepath"
	"testing"
)

func TestBloodBankPersistsUnits(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bloodbank.json")
	bank := &BloodBank{}
	want := BloodBag{
		UnitID:         "U001",
		BloodGroup:     "O+",
		VolumeML:       450,
		DonorAge:       25,
		CollectionDate: "2026-09-27",
	}

	if got := bank.CollectUnit(want); got != "Blood bag collected" {
		t.Fatalf("CollectUnit() = %q, want success", got)
	}
	if err := bank.SaveToFile(path); err != nil {
		t.Fatalf("SaveToFile() error = %v", err)
	}

	loaded, err := LoadFromFile(path)
	if err != nil {
		t.Fatalf("LoadFromFile() error = %v", err)
	}
	if got := loaded.SearchUnit(want.UnitID); got == nil || *got != want {
		t.Fatalf("loaded unit = %v, want %+v", got, want)
	}
}
