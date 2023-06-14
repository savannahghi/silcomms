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

// PremiumSMSResponse is the response returned after making a request to SILCOMMS to send a premium SMS
type PremiumSMSResponse struct {
	GUID         string `json:"guid"`
	Body         string `json:"body"`
	Msisdn       string `json:"msisdn"`
	SMSType      string `json:"sms_type"`
	Gateway      string `json:"gateway"`
	Carrier      string `json:"carrier"`
	Subscription string `json:"subscription"`
	Direction    string `json:"direction"`
	State        string `json:"state"`
}

// Subscription represents the response that is returned when activating a subscription to an offer
type Subscription struct {
	GUID             string `json:"guid"`
	Gateway          string `json:"gateway"`
	Offer            string `json:"offer"`
	Msisdn           string `json:"msisdn"`
	LinkID           string `json:"link_id"`
	ActivationDate   string `json:"activation_date"`
	DeactivationDate any    `json:"deactivation_date"`
	DeactivationType string `json:"deactivation_type"`
	Sms              []any  `json:"sms"`
	Created          string `json:"created"`
	Updated          string `json:"updated"`
}
