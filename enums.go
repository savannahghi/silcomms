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
	case "":
		return true
	}
	return false
}

// String representation of status
func (s Status) String() string {
	return string(s)
}
