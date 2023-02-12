package nonjson

import (
	"errors"
	"fmt"
)

var SoftwareIdentifiers = newSoftwareIdentifierRegistry()

type SoftwareIdentifier struct {
	// SKU is the value that refers to a specific license in Supermicro documentation.
	SKU string
	// DisplayName is the string identifier refers to the license type, and is stored
	// in the product key.
	DisplayName string
	// ID is the numeric identifier that refers to the license type, and is stored
	// in the product key.
	ID byte
}

func (sid *SoftwareIdentifier) Byte() byte {
	return sid.ID
}

func newSoftwareIdentifierRegistry() *softwareIdentifierRegistry {
	reserved := &SoftwareIdentifier{
		SKU:         "",
		DisplayName: "Reserved",
		ID:          0,
	}
	ssmServer := &SoftwareIdentifier{
		SKU:         "",
		DisplayName: "SSM",
		ID:          1,
	}
	sd5 := &SoftwareIdentifier{
		SKU:         "",
		DisplayName: "SD5",
		ID:          2,
	}
	sum := &SoftwareIdentifier{
		SKU:         "SFT-SUM-LIC",
		DisplayName: "SUM",
		ID:          3,
	}
	spm := &SoftwareIdentifier{
		SKU:         "SFT-SPM-LIC",
		DisplayName: "SPM",
		ID:          4,
	}
	scm := &SoftwareIdentifier{
		SKU:         "SFT-SCM-LIC",
		DisplayName: "SCM",
		ID:          5,
	}
	all := &SoftwareIdentifier{
		SKU:         "SFT-DCMS-SINGLE",
		DisplayName: "ALL",
		ID:          6,
	}
	site := &SoftwareIdentifier{
		SKU:         "SFT-DCMS-SITE",
		DisplayName: "SITE",
		ID:          7,
	}
	callHome := &SoftwareIdentifier{
		SKU:         "SFT-DCMS-CALL-HOME",
		DisplayName: "DCMS-CALL-HOME",
		ID:          8,
	}
	svc := &SoftwareIdentifier{
		SKU:         "SFT-DCMS-SVC-KEY",
		DisplayName: "SFT-DCMS-SVC-KEY",
		ID:          9,
	}
	sddc := &SoftwareIdentifier{
		SKU:         "SFT-SDDC-SINGLE",
		DisplayName: "SFT-SDDC-SINGLE",
		ID:          210,
	}

	return &softwareIdentifierRegistry{
		Reserved:  reserved,
		SSMServer: ssmServer,
		SD5:       sd5,
		SUM:       sum,
		SPM:       spm,
		SCM:       scm,
		ALL:       all,
		SITE:      site,
		CallHome:  callHome,
		SVC:       svc,
		SDDC:      sddc,
		softwareIdentifiers: []*SoftwareIdentifier{
			reserved,
			ssmServer,
			sd5,
			sum,
			spm,
			scm,
			all,
			site,
			callHome,
			svc,
			sddc,
		},
	}
}

type softwareIdentifierRegistry struct {
	Reserved  *SoftwareIdentifier
	SSMServer *SoftwareIdentifier
	SD5       *SoftwareIdentifier
	SUM       *SoftwareIdentifier
	SPM       *SoftwareIdentifier
	SCM       *SoftwareIdentifier
	ALL       *SoftwareIdentifier
	SITE      *SoftwareIdentifier
	CallHome  *SoftwareIdentifier
	SVC       *SoftwareIdentifier
	SDDC      *SoftwareIdentifier

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

// ByDisplayName returns a SoftwareIdentifier that contains the given display name, or an error
// if one was not found.
func (sidr *softwareIdentifierRegistry) ByDisplayName(displayName string) (*SoftwareIdentifier, error) {
	for _, softwareIdentifier := range sidr.List() {
		if softwareIdentifier.DisplayName == displayName {
			return softwareIdentifier, nil
		}
	}
	msg := fmt.Sprintf("software identifier display name '%s' not found in the registry", displayName)
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
