package json

import (
	"errors"
	"fmt"
)

var SoftwareIdentifiers = newSoftwareIdentifierRegistry()

type SoftwareIdentifier struct {
	// SKU is the value that refers to a specific license type in Supermicro documentation.
	SKU string
	// ID is the numeric identifier that refers to a license type.
	ID byte
}

func (sid *SoftwareIdentifier) Byte() byte {
	return sid.ID
}

func newSoftwareIdentifierRegistry() *softwareIdentifierRegistry {
	oob := &SoftwareIdentifier{
		SKU: "SFT-OOB-LIC",
		ID:  1,
	}
	dcms := &SoftwareIdentifier{
		SKU: "SFT-DCMS-SINGLE",
		ID:  2,
	}
	spm := &SoftwareIdentifier{
		SKU: "SFT-SPM-LIC",
		ID:  3,
	}
	svc := &SoftwareIdentifier{
		SKU: "SFT-DCMS-SVC-KEY",
		ID:  4,
	}
	sddc := &SoftwareIdentifier{
		SKU: "SFT-SDDC-SINGLE",
		ID:  5,
	}

	return &softwareIdentifierRegistry{
		OOB:  oob,
		DCMS: dcms,
		SPM:  spm,
		SVC:  svc,
		SDDC: sddc,
		softwareIdentifiers: []*SoftwareIdentifier{
			oob,
			dcms,
			spm,
			svc,
			sddc,
		},
	}
}

type softwareIdentifierRegistry struct {
	OOB  *SoftwareIdentifier
	DCMS *SoftwareIdentifier
	SPM  *SoftwareIdentifier
	SVC  *SoftwareIdentifier
	SDDC *SoftwareIdentifier

	softwareIdentifiers []*SoftwareIdentifier
}

// List returns all SoftwareIdentifier entries.
func (sidr *softwareIdentifierRegistry) List() []*SoftwareIdentifier {
	return sidr.softwareIdentifiers
}

// BySKU returns a SoftwareIdentifier that contains the given SKU, or an error if one
// was not found.
func (sidr *softwareIdentifierRegistry) BySKU(sku string) (*SoftwareIdentifier, error) {
	for _, softwareIdentifier := range sidr.List() {
		if softwareIdentifier.SKU == sku {
			return softwareIdentifier, nil
		}
	}
	msg := fmt.Sprintf("software identifier SKU '%s' not found in the registry", sku)
	return nil, errors.New(msg)
}

// ByID returns a SoftwareIdentifier that contains the given ID, or an error if one
// was not found.
func (sidr *softwareIdentifierRegistry) ByID(id byte) (*SoftwareIdentifier, error) {
	for _, softwareIdentifier := range sidr.List() {
		if softwareIdentifier.ID == id {
			return softwareIdentifier, nil
		}
	}
	msg := fmt.Sprintf("software identifier ID 0x%02X not found in the registry", id)
	return nil, errors.New(msg)
}
