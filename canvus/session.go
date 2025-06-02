// Package canvus provides a Go SDK for the Canvus API.
package canvus

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"path"
)

// Authenticator applies authentication to an HTTP request.
type Authenticator interface {
	Authenticate(req *http.Request)
}

// APIKeyAuthenticator authenticates using a static API key and header.
type APIKeyAuthenticator struct {
	Header string
	APIKey string
}

// Authenticate sets the API key header on the request.
func (a *APIKeyAuthenticator) Authenticate(req *http.Request) {
	if a.Header != "" && a.APIKey != "" {
		req.Header.Set(a.Header, a.APIKey)
	}
}

// TokenAuthenticator authenticates using a bearer token.
type TokenAuthenticator struct {
	Token string
}

// Authenticate sets the Authorization header on the request.
func (a *TokenAuthenticator) Authenticate(req *http.Request) {
	if a.Token != "" {
		req.Header.Set("Private-Token", a.Token)
	}
}

// SessionOption configures a Session.
type SessionOption func(*Session)

// WithAPIKey configures the session to use a static API key.
func WithAPIKey(apiKey string) SessionOption {
	return func(s *Session) {
		s.authenticator = &APIKeyAuthenticator{Header: "Private-Token", APIKey: apiKey}
	}
}

// WithToken configures the session to use a bearer token.
func WithToken(token string) SessionOption {
	return func(s *Session) {
		s.authenticator = &TokenAuthenticator{Token: token}
	}
}

// Session is the main entry point for interacting with the Canvus API.
type Session struct {
	BaseURL       string
	HTTPClient    *http.Client
	authenticator Authenticator
	userID        int64 // ID of the authenticated user, if available
}

// NewSession creates a new Canvus API session.
// If httpClient is nil, http.DefaultClient is used.
func NewSession(baseURL string, opts ...SessionOption) *Session {
	s := &Session{
		BaseURL:    baseURL,
		HTTPClient: http.DefaultClient,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// doRequest is a helper for making HTTP requests to the Canvus API.
// queryParams is an optional map of query parameters to append to the URL.
// If rawResponse is true, the response body is returned as []byte in 'out' (must be *[]byte).
func (s *Session) doRequest(ctx context.Context, method, endpoint string, body interface{}, out interface{}, queryParams map[string]string, rawResponse bool) error {
	u, err := url.Parse(s.BaseURL)
	if err != nil {
		return err
	}
	u.Path = path.Join(u.Path, endpoint)

	// Add query parameters if provided
	if queryParams != nil && len(queryParams) > 0 {
		q := u.Query()
		for k, v := range queryParams {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
	}

	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), reqBody)
	if err != nil {
		return err
	}
	if s.authenticator != nil {
		s.authenticator.Authenticate(req)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return &APIError{StatusCode: resp.StatusCode, Message: string(b)}
	}

	if out != nil {
		if rawResponse {
			// out must be *[]byte
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			if ptr, ok := out.(*[]byte); ok {
				*ptr = b
			} else {
				return errors.New("out must be *[]byte when rawResponse is true")
			}
		} else {
			if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
				return err
			}
		}
	}
	return nil
}

// Login authenticates a user and stores the returned token and user ID for future requests.
func (s *Session) Login(ctx context.Context, email, password string) error {
	loginReq := map[string]string{
		"email":    email,
		"password": password,
	}
	var loginResp struct {
		Token string `json:"token"`
		User  struct {
			ID int64 `json:"id"`
		} `json:"user"`
	}
	err := s.doRequest(ctx, http.MethodPost, "users/login", loginReq, &loginResp, nil, false)
	if err != nil {
		return err
	}
	if loginResp.Token == "" {
		return errors.New("login: no token returned")
	}
	s.authenticator = &TokenAuthenticator{Token: loginResp.Token}
	s.userID = loginResp.User.ID
	return nil
}

// Logout invalidates the current token and clears authentication.
// It calls POST /users/logout and clears the authenticator on success.
func (s *Session) Logout(ctx context.Context) error {
	logoutReq := map[string]string{}
	var logoutResp map[string]interface{}
	err := s.doRequest(ctx, http.MethodPost, "users/logout", logoutReq, &logoutResp, nil, false)
	if err != nil {
		return err
	}
	s.authenticator = nil
	return nil
}

// Users provides access to user management methods.
func (s *Session) Users() *Session {
	return s
}

// UserID returns the authenticated user's ID, or 0 if not logged in.
func (s *Session) UserID() int64 {
	return s.userID
}
