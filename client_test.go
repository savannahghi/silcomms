package silcomms

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/jarcoal/httpmock"
	"github.com/savannahghi/authutils"
)

// AuthServerServiceMock mocks onboarding implementations
type AuthServerServiceMock struct {
	MockLoginUserFn    func(ctx context.Context, input *authutils.LoginUserPayload) (*authutils.OAUTHResponse, error)
	MockRefreshTokenFn func(ctx context.Context, refreshToken string) (*authutils.OAUTHResponse, error)
}

// NewAuthServerServiceMock initializes our client mocks
func NewAuthServerServiceMock() *AuthServerServiceMock {
	return &AuthServerServiceMock{
		MockLoginUserFn: func(ctx context.Context, input *authutils.LoginUserPayload) (*authutils.OAUTHResponse, error) { //nolint: revive
			return &authutils.OAUTHResponse{
				AccessToken:  "access",
				RefreshToken: "refresh",
			}, nil
		},
		MockRefreshTokenFn: func(ctx context.Context, refreshToken string) (*authutils.OAUTHResponse, error) { //nolint:revive
			return &authutils.OAUTHResponse{
				AccessToken:  "access",
				RefreshToken: "refresh",
			}, nil
		},
	}
}

// LoginUser mocks the implementation of proxying login requests for users to authserver
func (oc AuthServerServiceMock) LoginUser(ctx context.Context, input *authutils.LoginUserPayload) (*authutils.OAUTHResponse, error) {
	return oc.MockLoginUserFn(ctx, input)
}

// MockRefreshToken mocks the implementation of getting refresh tokens
func (oc AuthServerServiceMock) RefreshToken(ctx context.Context, refreshToken string) (*authutils.OAUTHResponse, error) {
	return oc.MockRefreshTokenFn(ctx, refreshToken)
}

var authServer = NewAuthServerServiceMock()

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

			s := mustNewClient(authServer)

			if tt.name == "happy case: make authenticated POST request" {
				httpmock.RegisterResponder(http.MethodPost, "/v1/sms/bulk/", func(_ *http.Request) (*http.Response, error) {
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
				httpmock.RegisterResponder(http.MethodGet, "/v1/sms/bulk/", func(_ *http.Request) (*http.Response, error) {
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

			_, err := s.MakeRequest(tt.args.ctx, tt.args.method, tt.args.path, tt.args.queryParams, tt.args.body, tt.args.authorised) //nolint: bodyclose
			if (err != nil) != tt.wantErr {
				t.Errorf("SILclient.MakeRequest() error = %v, wantErr %v", err, tt.wantErr)
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
			name:    "sad case: error occurs",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			s := mustNewClient(authServer)

			if tt.name == "sad case: error occurs" {
				authServer.MockRefreshTokenFn = func(ctx context.Context, refreshToken string) (*authutils.OAUTHResponse, error) { //nolint:all
					return nil, fmt.Errorf("error")
				}
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
			name:    "sad case: error case",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			s, err := newClient(authServer)
			if err != nil {
				t.Errorf("failed to initialize: %v", err)
				return
			}

			if tt.name == "sad case: error case" {
				authServer.MockLoginUserFn = func(ctx context.Context, input *authutils.LoginUserPayload) (*authutils.OAUTHResponse, error) { //nolint:all
					return nil, fmt.Errorf("error")
				}
			}

			if err := s.login(); (err != nil) != tt.wantErr {
				t.Errorf("client.login() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
