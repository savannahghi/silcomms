package silcomms

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
)

// MockLogin mocks a mock login request to obtain a token
func MockLogin() {
	httpmock.RegisterResponder(http.MethodPost, fmt.Sprintf("%s/auth/token/", BaseURL), func(r *http.Request) (*http.Response, error) {
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

func TestSILComms_Login(t *testing.T) {
	type fields struct {
		client *http.Client
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "happy case: successful login",
			fields: fields{
				client: &http.Client{},
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

			s := &client{
				client: tt.fields.client,
			}
			s.login()
		})
	}
}

func TestSILclient_refreshAccessToken(t *testing.T) {
	type fields struct {
		client       *http.Client
		refreshToken string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "happy case: refresh access token",
			fields: fields{
				client: &http.Client{},
			},
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

			s := &client{
				client:       tt.fields.client,
				refreshToken: tt.fields.refreshToken,
			}
			s.refreshAccessToken()
		})
	}
}

func TestSILclient_MakeRequest(t *testing.T) {
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
		args    args
		wantErr bool
	}{
		{
			name: "happy case: make authenticated request",
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
			name: "happy case: make unauthenticated request",
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
			MockLogin()

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

			s := newClient()
			got, err := s.MakeRequest(tt.args.ctx, tt.args.method, tt.args.path, tt.args.queryParams, tt.args.body, tt.args.authorised)
			if (err != nil) != tt.wantErr {
				t.Errorf("SILclient.MakeRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("SILclient.MakeRequest() expected response not to be nil for %v", tt.name)
				return
			}
		})
	}
}
