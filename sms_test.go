package silcomms_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/silcomms"
)

// ClientMock ...
type ClientMock struct {
	MockMakeRequestFn func(ctx context.Context, method, path string, queryParams map[string]string, body interface{}, authorised bool) (*http.Response, error)
}

// MakeRequest ...
func (c *ClientMock) MakeRequest(ctx context.Context, method, path string, queryParams map[string]string, body interface{}, authorised bool) (*http.Response, error) {
	return c.MockMakeRequestFn(ctx, method, path, queryParams, body, authorised)
}

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
			client := &ClientMock{}
			l := silcomms.NewSILCommsLib(client)

			if tt.name == "happy case: send bulk sms" {
				client.MockMakeRequestFn = func(ctx context.Context, method, path string, queryParams map[string]string, body interface{}, authorised bool) (*http.Response, error) {
					msg := silcomms.APIResponse{
						Status:  silcomms.StatusSuccess,
						Message: "success",
						Data: silcomms.BulkSMSResponse{
							GUID:       "",
							Sender:     "",
							Message:    "",
							Recipients: []string{},
							State:      "",
							SMS:        []string{},
							Created:    "",
							Updated:    "",
						},
					}

					payload, _ := json.Marshal(msg)

					return &http.Response{StatusCode: http.StatusAccepted, Body: io.NopCloser(bytes.NewBuffer(payload))}, nil
				}
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
