package authservice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is an HTTP client for the Auth Service
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new Auth Service client
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// ValidateTokenRequest represents a token validation request
type ValidateTokenRequest struct {
	Token string `json:"token"`
}

// ValidateTokenResponse represents a token validation response
type ValidateTokenResponse struct {
	Valid   bool   `json:"valid"`
	UserID  string `json:"userId,omitempty"`
	Email   string `json:"email,omitempty"`
	Role    string `json:"role,omitempty"`
	Message string `json:"message,omitempty"`
}

// ValidateToken validates a JWT token with the Auth Service
func (c *Client) ValidateToken(token string) (*ValidateTokenResponse, error) {
	url := fmt.Sprintf("%s/api/v1/auth/validate", c.baseURL)
	
	reqBody := ValidateTokenRequest{Token: token}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return &ValidateTokenResponse{
			Valid:   false,
			Message: string(body),
		}, nil
	}

	var response ValidateTokenResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// GetUserRequest represents a get user request
type GetUserRequest struct {
	UserID string `json:"userId"`
}

// User represents a user from Auth Service
type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// GetUserResponse represents a get user response
type GetUserResponse struct {
	User    *User  `json:"user,omitempty"`
	Message string `json:"message,omitempty"`
}

// GetUser gets user information from the Auth Service
func (c *Client) GetUser(userID string) (*GetUserResponse, error) {
	url := fmt.Sprintf("%s/api/v1/auth/users/%s", c.baseURL, userID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return &GetUserResponse{
			Message: string(body),
		}, nil
	}

	var response GetUserResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

