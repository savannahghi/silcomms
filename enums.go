package silcomms

// Status is the API response status
type Status string

const (
	// StatusSuccess returns a `success` response
	StatusSuccess Status = "success"
	// StatusFailure returns a `failure` response
	StatusFailure Status = "failure"
	// StatusError returns an `error` response
	StatusError Status = "error"
)

// IsValid returns true if a status is valid
func (s Status) IsValid() bool {
	switch s {
	case StatusSuccess, StatusError, StatusFailure:
		return true
	}
	return false
}

// String representation of status
func (s Status) String() string {
	return string(s)
}

// VariantsType is a list of all the variant types.
type VariantsType string

// Variants type constants
const (
	BeWellApp           VariantsType = "BeWellApp"
	BeWellProApp        VariantsType = "BeWellProApp"
	BeWellWeb           VariantsType = "BeWellWeb"
	MyCareHubApp        VariantsType = "MyCareHubApp"
	MyCareHubProApp     VariantsType = "MyCareHubProApp"
	UONAfyaApp          VariantsType = "UONAfyaApp"
	UONAfyaProApp       VariantsType = "UONAfyaProApp"
	AdvantageLiteProApp VariantsType = "AdvantageLiteProApp"
	UzaziSalamaProApp   VariantsType = "UzaziSalamaProApp"
	NCIApp              VariantsType = "NCIApp"
	AfyaMojaApp         VariantsType = "AfyaMojaApp"
)

// IsValid returns true if a variant is valid
func (m VariantsType) IsValid() bool {
	switch m {
	case BeWellApp,
		BeWellProApp,
		MyCareHubApp,
		MyCareHubProApp,
		UONAfyaApp,
		UONAfyaProApp,
		AdvantageLiteProApp,
		UzaziSalamaProApp,
		NCIApp,
		AfyaMojaApp:
		return true
	}

	return false
}
