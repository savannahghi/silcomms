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
	client   *client
	senderID string
}

// NewSILCommsLib initializes a new implementation of the SIL Comms SDK
func NewSILCommsLib() (*CommsLib, error) {
	client, err := newClient()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize SIL Comms SMS SDK: %w", err)
	}

	l := &CommsLib{
		client:   client,
		senderID: SenderID,
	}

	return l, nil
}

// MustNewSILCommsLib initializes a new implementation of the SIL Comms SDK
func MustNewSILCommsLib() *CommsLib {
	sdk, err := NewSILCommsLib()
	if err != nil {
		panic(err)
	}

	return sdk
}

// SendBulkSMS returns a 202 Accepted synchronous response while the API attempts to send the SMS in the background.
// An asynchronous call is made to the app's sms_callback URL with a notification that shows the Bulk SMS status.
// An asynchronous call is made to the app's sms_callback individually for each of the recipients with the SMS status.
// message - message to be sent via the Bulk SMS
// recipients - phone number(s) to receive the Bulk SMS
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

// SendPremiumSMS is used to send a premium SMS using SILCOMMS gateway.
// message - message to be sent via the premium SMS.
// msisdn - phone number to receive the premium SMS.
// subscription - subscription/offer associated with the premium SMS.
func (l CommsLib) SendPremiumSMS(ctx context.Context, message, msisdn, subscription string) (*PremiumSMSResponse, error) {
	path := "/v1/sms/sms/"
	payload := struct {
		Body         string `json:"body"`
		Msisdn       string `json:"msisdn"`
		Subscription string `json:"subscription"`
	}{
		Body:         message,
		Msisdn:       msisdn,
		Subscription: subscription,
	}

	response, err := l.client.MakeRequest(ctx, http.MethodPost, path, nil, payload, true)
	if err != nil {
		return nil, fmt.Errorf("failed to make send premium sms request: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid send premium sms response code, got: %d", response.StatusCode)
	}

	var resp APIResponse
	err = json.NewDecoder(response.Body).Decode(&resp)
	if err != nil {
		return nil, fmt.Errorf("failed to decode send premium sms api response: %w", err)
	}

	var premiumSMS PremiumSMSResponse
	err = mapstructure.Decode(resp.Data, &premiumSMS)
	if err != nil {
		return nil, fmt.Errorf("failed to decode send premium sms data in api response: %w", err)
	}

	return &premiumSMS, nil
}

// ActivateSubscription is used activate a subscription to an offer on SILCOMMS.
// msisdn - phone number to be to activate a subscription to an offer.
// offer - offercode used to create a subscription.
func (l CommsLib) ActivateSubscription(ctx context.Context, offer string, msisdn string) (bool, error) {
	path := "/v1/sms/subscriptions/"
	payload := struct {
		Offer  string `json:"offer"`
		Msisdn string `json:"msisdn"`
	}{
		Offer:  offer,
		Msisdn: msisdn,
	}

	response, err := l.client.MakeRequest(ctx, http.MethodPost, path, nil, payload, true)
	if err != nil {
		return false, fmt.Errorf("failed to make activate subscription request: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return false, fmt.Errorf("invalid activate subscription response code, got: %d", response.StatusCode)
	}

	return true, nil
}

// GetSubscriptions fetches subscriptions from SILCOMMs based on provided query params
// params - query params used to get a subscription to an offer.
func (l CommsLib) GetSubscriptions(ctx context.Context, queryParams map[string]string) ([]*Subscription, error) {
	path := "/v1/sms/subscriptions/"

	response, err := l.client.MakeRequest(ctx, http.MethodGet, path, queryParams, nil, true)
	if err != nil {
		return nil, fmt.Errorf("failed to make get subscriptions request: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid get subscriptions response code, got: %d", response.StatusCode)
	}

	var resp APIResponse
	err = json.NewDecoder(response.Body).Decode(&resp)
	if err != nil {
		return nil, fmt.Errorf("failed to decode get subscriptions api response: %w", err)
	}

	var resultResponse ResultsResponse
	err = mapstructure.Decode(resp.Data, &resultResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode result response data in api response: %w", err)
	}

	var subscriptions []*Subscription
	err = mapstructure.Decode(resultResponse.Results, &subscriptions)
	if err != nil {
		return nil, fmt.Errorf("failed to decode subscriptions data in api response: %w", err)
	}

	return subscriptions, nil
}
