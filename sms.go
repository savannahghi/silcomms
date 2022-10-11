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
	// SenderID is the ID used to send the SMS
	SenderID = serverutils.MustGetEnvVar("SIL_COMMS_SENDER_ID")
)

// CommsLib is the SDK implementation for interacting with the sil communications API
type CommsLib struct {
	client *client
}

// NewSILCommsLib initializes a new implementation of the SIL Comms SDK
func NewSILCommsLib() *CommsLib {
	client := newClient()

	l := &CommsLib{
		client: client,
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
		Sender:     SenderID,
		Message:    message,
		Recipients: recipients,
	}

	response, err := l.client.MakeRequest(ctx, http.MethodPost, path, nil, payload, true)
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
