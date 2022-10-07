package silcomms

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mitchellh/mapstructure"
	"github.com/savannahghi/serverutils"
)

var (
	// commsSenderID is the ID used to send the SMS
	commsSenderID = serverutils.MustGetEnvVar("SIL_COMMS_SENDER_ID")
)

// ICommsClient is the interface for the client to make API request to sil communications
type ICommsClient interface {
	MakeRequest(ctx context.Context, method, path string, queryParams map[string]string, body interface{}, authorised bool) (*http.Response, error)
}

// CommsLib is the SDK implementation for interacting with the sil communications API
type CommsLib struct {
	Client ICommsClient
}

// NewSILCommsLib initializes a new implementation of the SIL Comms SDK
func NewSILCommsLib(client ICommsClient) *CommsLib {
	l := &CommsLib{
		Client: client,
	}

	return l
}

// SendBulkSMS returns a 202 Accepted synchronous response while the API attempts to send the SMS in the background.
// An asynchronous call is made to the app's sms_callback URL with a notification that shows the Bulk SMS status.
// An asynchronous call is made to the app's sms_callback individually for each of the recipients with the SMS status.
func (l CommsLib) SendBulkSMS(ctx context.Context, message string, recipients []string) (*BulkSMSResponse, error) {
	path := "/v1/sms/bulk/"
	payload := struct {
		Sender     string   `json:"sender"`
		Message    string   `json:"message"`
		Recipients []string `json:"recipients"`
	}{
		Sender:     commsSenderID,
		Message:    message,
		Recipients: recipients,
	}

	response, err := l.Client.MakeRequest(ctx, http.MethodPost, path, nil, payload, true)
	if err != nil {
		return nil, fmt.Errorf("failed to make send bulk sms request: %w", err)
	}

	if response.StatusCode != http.StatusAccepted {
		err := fmt.Errorf("invalid send bulk sms response code, got: %d", response.StatusCode)
		return nil, err
	}

	var resp APIResponse
	err = json.NewDecoder(response.Body).Decode(&resp)
	if err != nil {
		return nil, fmt.Errorf("failed to decode send bulk sms api response: %w", err)
	}

	var bulkSMS BulkSMSResponse
	err = mapstructure.Decode(resp.Data, &bulkSMS)
	if err != nil {
		return nil, fmt.Errorf("failed to decode send bulk sms data in api response: %w", err)
	}

	return &bulkSMS, nil
}
