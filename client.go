package silcomms

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/mitchellh/mapstructure"
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
	accessTokenTimeout = 30 * time.Minute

	// refreshTokenTimeout shows the refresh token expiry time
	refreshTokenTimeout = 24 * time.Hour
)

// It is the client used to make API request to sil communications API
type client struct {
	client *http.Client

	refreshToken       string
	refreshTokenTicker *time.Ticker

	accessToken       string
	accessTokenTicker *time.Ticker

	authFailed bool
}

// newClient initializes a new SIL comms client instance
func newClient() (*client, error) {
	s := &client{
		client: &http.Client{
			Timeout: time.Second * 10,
		},
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
func mustNewClient() *client {
	client, err := newClient()
	if err != nil {
		panic(err)
	}

	return client
}

// executed as a go routine to update the api tokens when they timeout
func (s *client) background() {
	for {
		select {
		case t := <-s.refreshTokenTicker.C:
			logrus.Println("SIL Comms Refresh Token updated at: ", t)
			err := s.login()
			if err != nil {
				s.authFailed = true
			}

		case t := <-s.accessTokenTicker.C:
			logrus.Println("SIL Comms Access Token updated at: ", t)
			err := s.refreshAccessToken()
			if err != nil {
				s.authFailed = true
			}

		}
	}
}

// setAccessToken sets the access token and updates the ticker timer
func (s *client) setAccessToken(token string) {
	s.accessToken = token
	if s.accessTokenTicker != nil {
		s.accessTokenTicker.Reset(accessTokenTimeout)
	} else {
		s.accessTokenTicker = time.NewTicker(accessTokenTimeout)
	}
}

// setRefreshToken sets the access token and updates the ticker timer
func (s *client) setRefreshToken(token string) {
	s.refreshToken = token
	if s.refreshTokenTicker != nil {
		s.refreshTokenTicker.Reset(refreshTokenTimeout)
	} else {
		s.refreshTokenTicker = time.NewTicker(refreshTokenTimeout)
	}
}

// login uses the provided credentials to login to the SIL communications backend
// It obtains the necessary tokens required to make authenticated requests
func (s *client) login() error {
	path := "/auth/token/"
	payload := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		Email:    email,
		Password: password,
	}

	response, err := s.MakeRequest(context.Background(), http.MethodPost, path, nil, payload, false)
	if err != nil {
		err = fmt.Errorf("failed to make login request: %w", err)
		return err
	}

	if response.StatusCode != http.StatusOK {
		err := fmt.Errorf("invalid login response code, got: %d", response.StatusCode)
		return err
	}

	var resp APIResponse
	err = json.NewDecoder(response.Body).Decode(&resp)
	if err != nil {
		err = fmt.Errorf("failed to decode login api response: %w", err)
		return err
	}

	var tokens TokenResponse
	err = mapstructure.Decode(resp.Data, &tokens)
	if err != nil {
		err = fmt.Errorf("failed to decode login data in api response: %w", err)
		return err
	}

	s.setRefreshToken(tokens.Refresh)
	s.setAccessToken(tokens.Access)

	return nil

}

func (s *client) refreshAccessToken() error {
	path := "/auth/token/refresh/"
	payload := struct {
		Refresh string `json:"refresh"`
	}{
		Refresh: s.refreshToken,
	}

	response, err := s.MakeRequest(context.Background(), http.MethodPost, path, nil, payload, false)
	if err != nil {
		err = fmt.Errorf("failed to make refresh request: %w", err)
		return err
	}

	if response.StatusCode != http.StatusOK {
		err := fmt.Errorf("invalid refresh token response code, got: %d", response.StatusCode)
		return err
	}

	var resp APIResponse
	err = json.NewDecoder(response.Body).Decode(&resp)
	if err != nil {
		err = fmt.Errorf("failed to decode refresh token api response: %w", err)
		return err
	}

	var tokens TokenResponse
	err = mapstructure.Decode(resp.Data, &tokens)
	if err != nil {
		err = fmt.Errorf("failed to decode refresh token data in api response: %w", err)
		return err
	}

	s.setAccessToken(tokens.Access)

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
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
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
