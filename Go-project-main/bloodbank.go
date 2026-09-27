package bloodbank

import (
	"encoding/json"
	"errors"
	"os"
)

type BloodBag struct {
	UnitID         string `json:"unit_id"`
	BloodGroup     string `json:"blood_group"`
	VolumeML       int    `json:"volume_ml"`
	DonorAge       int    `json:"donor_age"`
	CollectionDate string `json:"collection_date"`
}

type BloodBank struct {
	bloodBags []BloodBag
}

type bloodBankFile struct {
	BloodBags []BloodBag `json:"blood_bags"`
}

func LoadFromFile(path string) (*BloodBank, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return &BloodBank{}, nil
	}
	if err != nil {
		return nil, err
	}

	var stored bloodBankFile
	if err := json.Unmarshal(data, &stored); err != nil {
		return nil, err
	}

	return &BloodBank{bloodBags: stored.BloodBags}, nil
}

func (bank *BloodBank) SaveToFile(path string) error {
	data, err := json.MarshalIndent(bloodBankFile{BloodBags: bank.bloodBags}, "", "  ")
	if err != nil {
		return err
	}

	temporaryPath := path + ".tmp"
	if err := os.WriteFile(temporaryPath, data, 0600); err != nil {
		return err
	}
	if err := os.Rename(temporaryPath, path); err != nil {
		_ = os.Remove(temporaryPath)
		return err
	}

	return nil
}

func (bank *BloodBank) CollectUnit(b BloodBag) string {

	if b.DonorAge < 18 {
		return "Invalid age"
	}

	if b.BloodGroup != "A+" &&
		b.BloodGroup != "A-" &&
		b.BloodGroup != "B+" &&
		b.BloodGroup != "B-" &&
		b.BloodGroup != "AB+" &&
		b.BloodGroup != "AB-" &&
		b.BloodGroup != "O+" &&
		b.BloodGroup != "O-" {
		return "Invalid blood group"
	}

	bank.bloodBags = append(bank.bloodBags, b)

	return "Blood bag collected"
}

func (bank *BloodBank) SearchUnit(unitID string) *BloodBag {

	for i := 0; i < len(bank.bloodBags); i++ {

		if bank.bloodBags[i].UnitID == unitID {
			return &bank.bloodBags[i]
		}
	}

	return nil
}

func (bank *BloodBank) UpdateVolume(unitID string, volume int) string {

	for i := 0; i < len(bank.bloodBags); i++ {

		if bank.bloodBags[i].UnitID == unitID {
			bank.bloodBags[i].VolumeML = volume
			return "Volume updated"
		}
	}

	return "Blood unit not found"
}

func (bank *BloodBank) DiscardUnit(unitID string) string {
	for i := 0; i < len(bank.bloodBags); i++ {
		if bank.bloodBags[i].UnitID == unitID {
			bank.bloodBags = append(bank.bloodBags[:i], bank.bloodBags[i+1:]...)
			return "Blood bag discarded"
		}
	}

	return "Blood unit not found"
}

func (bank *BloodBank) ListUnits() []BloodBag {
	return append([]BloodBag(nil), bank.bloodBags...)
}
