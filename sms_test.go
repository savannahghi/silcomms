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
