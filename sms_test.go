package silcomms_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/jarcoal/httpmock"
	"github.com/savannahghi/silcomms"
)

func TestSILCommsLib_SendBulkSMS(t *testing.T) {
	type args struct {
		ctx        context.Context
		message    string
		recipients []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: send bulk sms",
			args: args{
				ctx:     context.Background(),
				message: "This is a test",
				recipients: []string{
					gofakeit.Phone(),
				},
			},
			wantErr: false,
		},
		{
			name: "sad case: invalid status code",
			args: args{
				ctx:     context.Background(),
				message: "This is a test",
				recipients: []string{
					gofakeit.Phone(),
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: invalid API response",
			args: args{
				ctx:     context.Background(),
				message: "This is a test",
				recipients: []string{
					gofakeit.Phone(),
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: invalid bulk SMS data response",
			args: args{
				ctx:     context.Background(),
				message: "This is a test",
				recipients: []string{
					gofakeit.Phone(),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			silcomms.MockLogin()

			l := silcomms.MustNewSILCommsLib()

			if tt.name == "happy case: send bulk sms" {
				httpmock.RegisterResponder(http.MethodPost, fmt.Sprintf("%s/v1/sms/bulk/", silcomms.BaseURL), func(r *http.Request) (*http.Response, error) {
					resp := silcomms.APIResponse{
						Status:  silcomms.StatusSuccess,
						Message: "success",
						Data: silcomms.BulkSMSResponse{
							GUID: gofakeit.UUID(),
						},
					}
					return httpmock.NewJsonResponse(http.StatusAccepted, resp)
				})
			}
			if tt.name == "sad case: invalid status code" {
				httpmock.RegisterResponder(http.MethodPost, fmt.Sprintf("%s/v1/sms/bulk/", silcomms.BaseURL), func(r *http.Request) (*http.Response, error) {
					return httpmock.NewJsonResponse(http.StatusUnauthorized, nil)
				})
			}

			if tt.name == "sad case: invalid API response" {
				httpmock.RegisterResponder(http.MethodPost, fmt.Sprintf("%s/v1/sms/bulk/", silcomms.BaseURL), func(r *http.Request) (*http.Response, error) {
					resp := map[string]interface{}{
						"status":  1234,
						"message": 1234,
					}

					return httpmock.NewJsonResponse(http.StatusAccepted, resp)
				})
			}

			if tt.name == "sad case: invalid bulk SMS data response" {
				httpmock.RegisterResponder(http.MethodPost, fmt.Sprintf("%s/v1/sms/bulk/", silcomms.BaseURL), func(r *http.Request) (*http.Response, error) {
					resp := silcomms.APIResponse{
						Status:  silcomms.StatusSuccess,
						Message: "success",
						Data: map[string]interface{}{
							"guid":    123456,
							"sender":  123456,
							"message": 123456,
						},
					}
					return httpmock.NewJsonResponse(http.StatusAccepted, resp)
				})
			}

			got, err := l.SendBulkSMS(tt.args.ctx, tt.args.message, tt.args.recipients)
			if (err != nil) != tt.wantErr {
				t.Errorf("SILCommsLib.SendBulkSMS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("SILCommsLib.SendBulkSMS() expected response not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestSILCommsLib_SendPremiumSMS(t *testing.T) {
	type args struct {
		ctx          context.Context
		message      string
		msisdn       string
		subscription string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: send premium sms",
			args: args{
				ctx:          context.Background(),
				message:      "test premium sms",
				msisdn:       gofakeit.Phone(),
				subscription: "01262626626",
			},
			wantErr: false,
		},
		{
			name: "Sad case: invalid status code",
			args: args{
				ctx:          context.Background(),
				message:      "test premium sms",
				msisdn:       gofakeit.Phone(),
				subscription: "01262626626",
			},
			wantErr: true,
		},
		{
			name: "Sad case: invalid API response",
			args: args{
				ctx:          context.Background(),
				message:      "test premium sms",
				msisdn:       gofakeit.Phone(),
				subscription: "01262626626",
			},
			wantErr: true,
		},
		{
			name: "Sad case: invalid premium SMS data response",
			args: args{
				ctx:          context.Background(),
				message:      "test premium sms",
				msisdn:       gofakeit.Phone(),
				subscription: "01262626626",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			silcomms.MockLogin()

			l := silcomms.MustNewSILCommsLib()

			if tt.name == "Happy case: send premium sms" {
				httpmock.RegisterResponder(http.MethodPost, fmt.Sprintf("%s/v1/sms/sms/", silcomms.BaseURL), func(r *http.Request) (*http.Response, error) {
					resp := silcomms.APIResponse{
						Status:  silcomms.StatusSuccess,
						Message: "success",
						Data: silcomms.PremiumSMSResponse{
							GUID: gofakeit.UUID(),
						},
					}
					return httpmock.NewJsonResponse(http.StatusOK, resp)
				})
			}

			if tt.name == "Sad case: invalid status code" {
				httpmock.RegisterResponder(http.MethodPost, fmt.Sprintf("%s/v1/sms/sms/", silcomms.BaseURL), func(r *http.Request) (*http.Response, error) {
					return httpmock.NewJsonResponse(http.StatusUnauthorized, nil)
				})
			}

			if tt.name == "Sad case: invalid API response" {
				httpmock.RegisterResponder(http.MethodPost, fmt.Sprintf("%s/v1/sms/sms/", silcomms.BaseURL), func(r *http.Request) (*http.Response, error) {
					resp := map[string]interface{}{
						"status":  1234,
						"message": 1234,
					}

					return httpmock.NewJsonResponse(http.StatusOK, resp)
				})
			}

			if tt.name == "Sad case: invalid premium SMS data response" {
				httpmock.RegisterResponder(http.MethodPost, fmt.Sprintf("%s/v1/sms/sms/", silcomms.BaseURL), func(r *http.Request) (*http.Response, error) {
					resp := silcomms.APIResponse{
						Status:  silcomms.StatusSuccess,
						Message: "success",
						Data: map[string]interface{}{
							"guid":    123456,
							"carrier": 123456,
							"message": 123456,
						},
					}
					return httpmock.NewJsonResponse(http.StatusOK, resp)
				})
			}

			got, err := l.SendPremiumSMS(tt.args.ctx, tt.args.message, tt.args.msisdn, tt.args.subscription)
			if (err != nil) != tt.wantErr {
				t.Errorf("SILCommsLib.SendPremiumSMS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("SILCommsLib.SendPremiumSMS() expected response not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestSILCommsLib_ActivateSubscription(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx      context.Context
		offer    string
		msisdn   string
		activate bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: activate subscription",
			args: args{
				ctx:    ctx,
				offer:  "01262626626",
				msisdn: gofakeit.Phone(),
			},
			wantErr: false,
		},
		{
			name: "Happy case: activate subscription bypass sdp",
			args: args{
				ctx:      ctx,
				offer:    "01262626626",
				msisdn:   gofakeit.Phone(),
				activate: false,
			},
			wantErr: false,
		},
		{
			name: "Sad case: invalid status code",
			args: args{
				ctx:    ctx,
				offer:  "01262626626",
				msisdn: gofakeit.Phone(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			silcomms.MockLogin()

			l := silcomms.MustNewSILCommsLib()

			if tt.name == "Happy case: activate subscription" {
				httpmock.RegisterResponder(http.MethodPost, fmt.Sprintf("%s/v1/sms/subscriptions/", silcomms.BaseURL), func(r *http.Request) (*http.Response, error) {
					resp := silcomms.APIResponse{
						Status:  silcomms.StatusSuccess,
						Message: "success",
						Data: map[string]interface{}{
							"guid":   "123456",
							"offer":  "123456",
							"msisdn": "123456",
						},
					}
					return httpmock.NewJsonResponse(http.StatusOK, resp)
				})
			}

			if tt.name == "Happy case: activate subscription bypass sdp" {
				httpmock.RegisterResponder(http.MethodPost, fmt.Sprintf("%s/v1/sms/subscriptions/", silcomms.BaseURL), func(r *http.Request) (*http.Response, error) {
					resp := silcomms.APIResponse{
						Status:  silcomms.StatusSuccess,
						Message: "success",
						Data: map[string]interface{}{
							"guid":   "123456",
							"offer":  "123456",
							"msisdn": "123456",
						},
					}
					return httpmock.NewJsonResponse(http.StatusOK, resp)
				})
			}

			if tt.name == "Sad case: invalid status code" {
				httpmock.RegisterResponder(http.MethodPost, fmt.Sprintf("%s/v1/sms/subscriptions/", silcomms.BaseURL), func(r *http.Request) (*http.Response, error) {
					return httpmock.NewJsonResponse(http.StatusUnauthorized, nil)
				})
			}

			got, err := l.ActivateSubscription(tt.args.ctx, tt.args.offer, tt.args.msisdn, tt.args.activate)
			if (err != nil) != tt.wantErr {
				t.Errorf("SILCommsLib.ActivateSubscription() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == false {
				t.Errorf("SILCommsLib.ActivateSubscription() expected response not to be false for %v", tt.name)
				return
			}
		})
	}
}

func TestSILCommsLib_GetSubscriptions(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx         context.Context
		queryParams map[string]string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: get subscription",
			args: args{
				ctx: ctx,
				queryParams: map[string]string{
					"msisdn": gofakeit.Phone(),
					"offer":  "01262626626",
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: invalid status code",
			args: args{
				ctx: ctx,
				queryParams: map[string]string{
					"msisdn": gofakeit.Phone(),
					"offer":  "01262626626",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			silcomms.MockLogin()

			l := silcomms.MustNewSILCommsLib()

			if tt.name == "Happy case: get subscription" {
				httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf("%s/v1/sms/subscriptions/", silcomms.BaseURL), func(r *http.Request) (*http.Response, error) {
					resp := silcomms.APIResponse{
						Status:  silcomms.StatusSuccess,
						Message: "success",
						Data: map[string]interface{}{
							"count":    1,
							"next":     nil,
							"previous": nil,
							"results": []map[string]interface{}{
								{
									"guid":              "e602b8b8-9591-4526-915d-57ef2579d8c4",
									"gateway":           "SAFARICOM",
									"offer":             "0022345234",
									"msisdn":            "+254722345678",
									"link_id":           "123123123123123",
									"activation_date":   "2022-08-04 14:11:17.206377+03:00",
									"deactivation_date": nil,
									"deactivation_type": "USER_INITIATED",
									"sms":               []string{},
									"created":           "2022-08-04 14:11:17.206377+03:00",
									"updated":           "2022-08-04 14:11:17.206377+03:00",
								},
							},
						},
					}
					return httpmock.NewJsonResponse(http.StatusOK, resp)
				})
			}

			if tt.name == "Sad case: invalid status code" {
				httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf("%s/v1/sms/subscriptions/", silcomms.BaseURL), func(r *http.Request) (*http.Response, error) {
					return httpmock.NewJsonResponse(http.StatusUnauthorized, nil)
				})
			}

			got, err := l.GetSubscriptions(tt.args.ctx, tt.args.queryParams)
			if (err != nil) != tt.wantErr {
				t.Errorf("SILCommsLib.GetSubscriptions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("SILCommsLib.GetSubscriptions() expected response not to be nil for %v", tt.name)
				return
			}
		})
	}
}
