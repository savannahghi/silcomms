package silcomms

// APIResponse is the base response from sil communications API
type APIResponse struct {
	Status  Status      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ResultsResponse is the base response from a paginated list of results
type ResultsResponse struct {
	Count    int           `json:"count"`
	Next     *string       `json:"next"`
	Previous *string       `json:"previous"`
	Results  []interface{} `json:"results"`
}

// TokenResponse is the data in the API response when logging in
// The access token is used as the bearer token when making API requests
// The refresh token is used to obtain a new access token when it expires
type TokenResponse struct {
	Refresh string `json:"refresh"`
	Access  string `json:"access"`
}

// ErrorMessage is the message in the ErrorResponse
type ErrorMessage struct {
	Message string `json:"message,omitempty"`
}

// ErrorResponse is the data in the API response when an error is encountered
type ErrorResponse struct {
	Detail  string         `json:"detail,omitempty"`
	Code    string         `json:"code,omitempty"`
	Message []ErrorMessage `json:"message,omitempty"`
}

// BulkSMSResponse is the data in the API response that is returned after making a request to send bulk sms
type BulkSMSResponse struct {
	GUID       string   `json:"guid"`
	Sender     string   `json:"sender"`
	Message    string   `json:"message"`
	Recipients []string `json:"recipients"`
	State      string   `json:"state"`
	SMS        []string `json:"sms"`
	Created    string   `json:"created"`
	Updated    string   `json:"updated"`
}
