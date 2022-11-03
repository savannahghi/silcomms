package silcomms

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit"
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
			name: "happy case: make unauthenticated request",
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
			name: "happy case: make authenticated POST request",
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
		{
			name: "happy case: make authenticated GET request",
			args: args{
				ctx:    context.Background(),
				method: http.MethodGet,
				path:   "/v1/sms/bulk/",
				queryParams: map[string]string{
					"app": gofakeit.UUID(),
				},
				body:       nil,
				authorised: true,
			},
			wantErr: false,
		},
		{
			name: "sad case: make unsupported protocol request",
			args: args{
				ctx:    context.Background(),
				method: http.MethodOptions,
				path:   "/v1/sms/bulk/",
				queryParams: map[string]string{
					"app": gofakeit.UUID(),
				},
				body:       nil,
				authorised: true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			MockLogin()
			s := mustNewClient()

			if tt.name == "happy case: make authenticated POST request" {
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
			}

			if tt.name == "happy case: make authenticated GET request" {
				httpmock.RegisterResponder(http.MethodGet, "/v1/sms/bulk/", func(r *http.Request) (*http.Response, error) {
					resp := APIResponse{
						Status:  StatusSuccess,
						Message: "success",
						Data: []BulkSMSResponse{
							{
								GUID:       "",
								Sender:     "",
								Message:    "",
								Recipients: []string{},
								State:      "",
								SMS:        []string{},
								Created:    "",
								Updated:    "",
							},
						},
					}
					return httpmock.NewJsonResponse(http.StatusOK, resp)
				})
			}

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

func Test_client_refreshAccessToken(t *testing.T) {

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "happy case: refresh access token",
			wantErr: false,
		},
		{
			name:    "sad case: invalid status code",
			wantErr: true,
		},
		{
			name:    "sad case: invalid api response",
			wantErr: true,
		},
		{
			name:    "sad case: invalid token response",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			MockLogin()
			s := mustNewClient()

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

			if tt.name == "sad case: invalid status code" {
				httpmock.RegisterResponder(http.MethodPost, "/auth/token/refresh/", func(r *http.Request) (*http.Response, error) {
					resp := APIResponse{
						Status:  StatusSuccess,
						Message: "success",
						Data:    nil,
					}
					return httpmock.NewJsonResponse(http.StatusBadRequest, resp)
				})
			}

			if tt.name == "sad case: invalid api response" {
				httpmock.RegisterResponder(http.MethodPost, fmt.Sprintf("%s/auth/token/refresh/", BaseURL), func(r *http.Request) (*http.Response, error) {
					resp := map[string]interface{}{
						"status":  1234,
						"message": 1234,
					}
					return httpmock.NewJsonResponse(http.StatusOK, resp)
				})
			}

			if tt.name == "sad case: invalid token response" {
				httpmock.RegisterResponder(http.MethodPost, fmt.Sprintf("%s/auth/token/refresh/", BaseURL), func(r *http.Request) (*http.Response, error) {
					resp := APIResponse{
						Status:  StatusSuccess,
						Message: "success",
						Data: map[string]interface{}{
							"refresh": 1234,
							"access":  1234,
						},
					}
					return httpmock.NewJsonResponse(http.StatusOK, resp)
				})
			}

			if err := s.refreshAccessToken(); (err != nil) != tt.wantErr {
				t.Errorf("client.refreshAccessToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_client_login(t *testing.T) {

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "happy case: successful login",
			wantErr: false,
		},
		{
			name:    "sad case: invalid status code",
			wantErr: true,
		},
		{
			name:    "sad case: invalid api response",
			wantErr: true,
		},
		{
			name:    "sad case: invalid token response",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

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

			s, err := newClient()
			if err != nil {
				t.Errorf("failed to initialize: %v", err)
				return
			}

			if tt.name == "happy case: successful login" {
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

			if tt.name == "sad case: invalid status code" {
				httpmock.RegisterResponder(http.MethodPost, fmt.Sprintf("%s/auth/token/", BaseURL), func(r *http.Request) (*http.Response, error) {
					resp := APIResponse{
						Status:  StatusSuccess,
						Message: "success",
						Data:    nil,
					}
					return httpmock.NewJsonResponse(http.StatusBadRequest, resp)
				})
			}

			if tt.name == "sad case: invalid api response" {
				httpmock.RegisterResponder(http.MethodPost, fmt.Sprintf("%s/auth/token/", BaseURL), func(r *http.Request) (*http.Response, error) {
					resp := map[string]interface{}{
						"status":  1234,
						"message": 1234,
					}
					return httpmock.NewJsonResponse(http.StatusOK, resp)
				})
			}

			if tt.name == "sad case: invalid token response" {
				httpmock.RegisterResponder(http.MethodPost, fmt.Sprintf("%s/auth/token/", BaseURL), func(r *http.Request) (*http.Response, error) {
					resp := APIResponse{
						Status:  StatusSuccess,
						Message: "success",
						Data: map[string]interface{}{
							"refresh": 1234,
							"access":  1234,
						},
					}
					return httpmock.NewJsonResponse(http.StatusOK, resp)
				})
			}

			if err := s.login(); (err != nil) != tt.wantErr {
				t.Errorf("client.login() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
