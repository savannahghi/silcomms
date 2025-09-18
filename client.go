package silcomms

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/savannahghi/authutils"
	"github.com/savannahghi/serverutils"
	"github.com/sirupsen/logrus"
)

var (
	// BaseURL represents the SIL-Comms base URL
	BaseURL = serverutils.MustGetEnvVar("SIL_COMMS_BASE_URL")

	// email is used for authentication against the SIL comms API
	email = serverutils.MustGetEnvVar("SIL_COMMS_EMAIL")

	// password is used for authentication against the SIL comms API
	password = serverutils.MustGetEnvVar("SIL_COMMS_PASSWORD")

	// accessTokenTimeout shows the access token expiry time.
	// After the access token expires, one is required to obtain a new one
	accessTokenTimeout = 59 * time.Minute
)

// AuthServerImpl defines the methods provided by
// the auth server library
type AuthServerImpl interface {
	LoginUser(ctx context.Context, input *authutils.LoginUserPayload) (*authutils.OAUTHResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*authutils.OAUTHResponse, error)
}

// It is the client used to make API request to sil communications API
type client struct {
	authServer AuthServerImpl
	client     *http.Client

	refreshToken string

	accessToken       string
	accessTokenTicker *time.Ticker

	authFailed bool
}

// newClient initializes a new SIL comms client instance
func newClient(authServer AuthServerImpl) (*client, error) {
	s := &client{
		client: &http.Client{
			Timeout: time.Second * 10,
		},
		authServer:   authServer,
		accessToken:  "",
		refreshToken: "",
		authFailed:   false,
	}

	err := s.login()
	if err != nil {
		return nil, err
	}

	// set up background routine to update tokens
	go s.background()

	return s, nil
}

// mustNewClient initializes a new SIL comms client instance
func mustNewClient(authServer AuthServerImpl) *client {
	client, err := newClient(authServer)
	if err != nil {
		panic(err)
	}

	return client
}

// executed as a go routine to update access and refresh token
func (s *client) background() {
	for t := range s.accessTokenTicker.C {
		logrus.Println("SIL Comms Access Token updated at: ", t)

		err := s.refreshAccessToken()
		if err != nil {
			s.authFailed = true
		} else {
			s.authFailed = false
		}
	}
}

// setAccessToken sets the access token and updates the ticker timer
func (s *client) setRefreshAndAccessToken(token *TokenResponse) {
	s.accessToken = token.Access

	s.refreshToken = token.Refresh
	if s.accessTokenTicker != nil {
		s.accessTokenTicker.Reset(accessTokenTimeout)
	} else {
		s.accessTokenTicker = time.NewTicker(accessTokenTimeout)
	}
}

// login uses the provided credentials to login to the authserver backend
// It obtains the necessary tokens required to make authenticated requests
func (s *client) login() error {
	ctx := context.Background()

	loginInput := authutils.LoginUserPayload{
		Email:    email,
		Password: password,
	}

	resp, err := s.authServer.LoginUser(ctx, &loginInput)
	if err != nil {
		return err
	}

	tokens := TokenResponse{
		Access:  resp.AccessToken,
		Refresh: resp.RefreshToken,
	}

	s.setRefreshAndAccessToken(&tokens)

	return nil
}

// refreshAccessToken makes a request to get
// new access and refresh tokens
func (s *client) refreshAccessToken() error {
	ctx := context.Background()

	resp, err := s.authServer.RefreshToken(ctx, s.refreshToken)
	if err != nil {
		return err
	}

	tokens := TokenResponse{
		Access:  resp.AccessToken,
		Refresh: resp.RefreshToken,
	}

	s.setRefreshAndAccessToken(&tokens)

	return nil
}

// MakeRequest performs a HTTP request to the provided path and parameters
func (s *client) MakeRequest(ctx context.Context, method, path string, queryParams map[string]string, body interface{}, authorised bool) (*http.Response, error) {
	// background refresh failed and the tokens are not valid
	if s.authFailed {
		return nil, fmt.Errorf("invalid credentials, cannot make request please update")
	}

	urlPath := fmt.Sprintf("%s%s", BaseURL, path)

	var request *http.Request

	switch method {
	case http.MethodGet:
		req, err := http.NewRequestWithContext(ctx, method, urlPath, nil)
		if err != nil {
			return nil, err
		}

		request = req

	case http.MethodPost:
		encoded, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}

		payload := bytes.NewBuffer(encoded)

		req, err := http.NewRequestWithContext(ctx, method, urlPath, payload)
		if err != nil {
			return nil, err
		}

		request = req

	default:
		return nil, fmt.Errorf("s.MakeRequest() unsupported http method: %s", method)
	}

	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")

	if authorised {
		request.Header.Set("Authorization", fmt.Sprintf("X-Bearer %s", s.accessToken))
	}

	if queryParams != nil {
		q := url.Values{}

		for key, value := range queryParams {
			q.Add(key, value)
		}

		request.URL.RawQuery = q.Encode()
	}

	return s.client.Do(request)
}
