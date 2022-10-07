package silcomms

import (
	"context"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestSILComms_Login(t *testing.T) {
	type fields struct {
		client http.Client
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "happy case: successful login",
			fields: fields{
				client: http.Client{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			if tt.name == "happy case: successful login" {
				httpmock.RegisterResponder(http.MethodPost, "/auth/token/", func(r *http.Request) (*http.Response, error) {
					resp := APIResponse{
						Status:  StatusSuccess,
						Message: "success",
						Data: TokenResponse{
							Refresh: "refresh",
							Access:  "access",
						},
					}
					return httpmock.NewJsonResponse(http.StatusOK, resp)
				})
			}

			s := &CommsClient{
				client: tt.fields.client,
			}
			s.login()
		})
	}
}

func TestSILCommsClient_refreshAccessToken(t *testing.T) {
	type fields struct {
		client       http.Client
		refreshToken string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name:   "happy case: refresh access token",
			fields: fields{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			if tt.name == "happy case: refresh access token" {
				httpmock.RegisterResponder(http.MethodPost, "/auth/token/refresh/", func(r *http.Request) (*http.Response, error) {
					resp := APIResponse{
						Status:  StatusSuccess,
						Message: "success",
						Data: TokenResponse{
							Refresh: "refresh",
							Access:  "access",
						},
					}
					return httpmock.NewJsonResponse(http.StatusOK, resp)
				})
			}

			s := &CommsClient{
				client:       tt.fields.client,
				refreshToken: tt.fields.refreshToken,
			}
			s.refreshAccessToken()
		})
	}
}

func TestSILCommsClient_MakeRequest(t *testing.T) {
	type fields struct {
		client      http.Client
		accessToken string
	}
	type args struct {
		ctx         context.Context
		method      string
		path        string
		queryParams map[string]string
		body        interface{}
		authorised  bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "happy case: make authenticated request",
			fields: fields{},
			args: args{
				ctx:         context.Background(),
				method:      http.MethodPost,
				path:        "/auth/token/",
				queryParams: nil,
				body:        nil,
				authorised:  false,
			},
			wantErr: false,
		},
		{
			name:   "happy case: make unauthenticated request",
			fields: fields{},
			args: args{
				ctx:         context.Background(),
				method:      http.MethodPost,
				path:        "/v1/sms/bulk/",
				queryParams: nil,
				body:        nil,
				authorised:  true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			httpmock.RegisterResponder(http.MethodPost, "/v1/sms/bulk/", func(r *http.Request) (*http.Response, error) {
				resp := APIResponse{
					Status:  StatusSuccess,
					Message: "success",
					Data: BulkSMSResponse{
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
				return httpmock.NewJsonResponse(http.StatusOK, resp)
			})

			if tt.name == "happy case: make authenticated request" {
				httpmock.RegisterResponder(http.MethodPost, "/auth/token/", func(r *http.Request) (*http.Response, error) {
					resp := APIResponse{
						Status:  StatusSuccess,
						Message: "success",
						Data: TokenResponse{
							Refresh: "refresh",
							Access:  "access",
						},
					}
					return httpmock.NewJsonResponse(http.StatusOK, resp)
				})
			}

			s := &CommsClient{
				client:      tt.fields.client,
				accessToken: tt.fields.accessToken,
			}
			got, err := s.MakeRequest(tt.args.ctx, tt.args.method, tt.args.path, tt.args.queryParams, tt.args.body, tt.args.authorised)
			if (err != nil) != tt.wantErr {
				t.Errorf("SILCommsClient.MakeRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("SILCommsClient.MakeRequest() expected response not to be nil for %v", tt.name)
				return
			}
		})
	}
}
