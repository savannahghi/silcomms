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

			l := silcomms.NewSILCommsLib()

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
