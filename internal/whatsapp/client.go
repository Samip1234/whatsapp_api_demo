package whatsapp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"whatsapp_api/pkg/errorx"
)

type Config struct {
	APIVersion     string
	PhoneNumberID  string
	Token          string
	BaseURL        string
	HTTPClient     *http.Client
	RequestTimeout time.Duration
}

type Client struct {
	apiVersion    string
	phoneNumberID string
	token         string
	baseURL       string
	httpClient    *http.Client
}

type MessageRequest struct {
	MessagingProduct string         `json:"messaging_product"`
	To               string         `json:"to"`
	Type             string         `json:"type"`
	Text             *TextComponent `json:"text,omitempty"`
}

type TextComponent struct {
	PreviewURL bool   `json:"preview_url"`
	Body       string `json:"body"`
}

type MessageResponse struct {
	MessagingProduct string           `json:"messaging_product"`
	Contacts         []ContactPayload `json:"contacts"`
	Messages         []MessagePayload `json:"messages"`
}

type ContactPayload struct {
	Input string `json:"input"`
	WaID  string `json:"wa_id"`
}

type MessagePayload struct {
	ID string `json:"id"`
}

type errorEnvelope struct {
	Error *apiError `json:"error"`
}

type apiError struct {
	Message   string       `json:"message"`
	Type      string       `json:"type"`
	Code      int          `json:"code"`
	ErrorData *errorDetail `json:"error_data"`
}

type errorDetail struct {
	Details string `json:"details"`
}

const defaultBaseURL = "https://graph.facebook.com"

func NewClient(cfg Config) (*Client, error) {
	if cfg.APIVersion == "" {
		return nil, errorx.WithDetail(errorx.ErrInvalidConfig, "api version is required")
	}
	if cfg.PhoneNumberID == "" {
		return nil, errorx.WithDetail(errorx.ErrInvalidConfig, "phone number id is required")
	}
	if cfg.Token == "" {
		return nil, errorx.WithDetail(errorx.ErrInvalidConfig, "token is required")
	}

	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	httpClient := cfg.HTTPClient
	if httpClient == nil {
		timeout := cfg.RequestTimeout
		if timeout <= 0 {
			timeout = 10 * time.Second
		}
		httpClient = &http.Client{Timeout: timeout}
	}

	return &Client{
		apiVersion:    cfg.APIVersion,
		phoneNumberID: cfg.PhoneNumberID,
		token:         cfg.Token,
		baseURL:       baseURL,
		httpClient:    httpClient,
	}, nil
}

func (c *Client) SendTextMessage(ctx context.Context, to, body string) (*MessageResponse, error) {
	if c == nil {
		return nil, errorx.WithDetail(errorx.ErrInvalidConfig, "client is not initialized")
	}

	requestPayload := &MessageRequest{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "text",
		Text: &TextComponent{
			PreviewURL: false,
			Body:       body,
		},
	}

	jsonPayload, err := json.Marshal(requestPayload)
	if err != nil {
		return nil, errorx.Wrap(errorx.ErrMarshalRequest, err)
	}

	url := fmt.Sprintf("%s/%s/%s/messages", c.baseURL, c.apiVersion, c.phoneNumberID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonPayload))
	if err != nil {
		return nil, errorx.Wrap(errorx.ErrBuildRequest, err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errorx.Wrap(errorx.ErrSendRequest, err)
	}
	defer resp.Body.Close()

	responseBody, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, errorx.Wrap(errorx.ErrSendRequest, readErr)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, c.parseAPIError(responseBody)
	}

	var messageResponse MessageResponse
	if err = json.Unmarshal(responseBody, &messageResponse); err != nil {
		return nil, errorx.Wrap(errorx.ErrMarshalRequest, err)
	}

	return &messageResponse, nil
}

func (c *Client) parseAPIError(body []byte) error {
	var envelope errorEnvelope
	if err := json.Unmarshal(body, &envelope); err != nil {
		return errorx.Wrap(errorx.ErrAPIResponse, err)
	}

	if envelope.Error == nil {
		return errorx.ErrAPIResponse
	}

	detail := envelope.Error.Message
	if envelope.Error.ErrorData != nil && envelope.Error.ErrorData.Details != "" {
		detail = detail + " (" + envelope.Error.ErrorData.Details + ")"
	}

	return errorx.WithDetail(errorx.ErrAPIResponse, detail)
}
