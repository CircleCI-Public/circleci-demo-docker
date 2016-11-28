package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Client defines the interface exposed by our API.
type Client interface {
	AddContact(contact AddContactRequest) (*Contact, error)
	GetContactByEmail(email string) (*Contact, error)
}

// ErrorResponse is returned by our service when an error occurs.
type ErrorResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func (e ErrorResponse) Error() string {
	return fmt.Sprintf("%v: %v", e.StatusCode, e.Message)
}

// NewClient creates a Client that accesses a service at the given base URL.
func NewClient(baseURL string) Client {
	httpClient := http.DefaultClient
	return &DefaultClient{
		http:    httpClient,
		BaseURL: baseURL,
	}
}

// ===== DefaultClient =================================================================================================

// DefaultClient provides an implementation of the Client interface.
type DefaultClient struct {
	http    *http.Client
	BaseURL string
}

// performRequestMethod constructs a request and uses `performRequest` to execute it.
func (c *DefaultClient) performRequestMethod(method string, path string, headers map[string]string, data interface{}, response interface{}) error {

	req, err := c.newRequest(method, path, headers, data)
	if err != nil {
		return err
	}

	return c.performRequest(req, response)
}

// performRequest executes the given request, and uses `response` to parse the JSON response.
func (c *DefaultClient) performRequest(req *http.Request, response interface{}) error {
	// perform the request
	httpResponse, err := c.http.Do(req)
	if err != nil {
		return err
	}

	defer httpResponse.Body.Close()

	// read the response
	var responseBody []byte
	if responseBody, err = ioutil.ReadAll(httpResponse.Body); err != nil {
		return err
	}

	if httpResponse.StatusCode >= 400 {
		contentTypeHeader := httpResponse.Header["Content-Type"]
		if len(contentTypeHeader) != 0 && contentTypeHeader[0] == "application/json" {
			var errResponse ErrorResponse
			err := json.Unmarshal(responseBody, &errResponse)
			if err == nil {
				return errResponse
			}
		}

		return &ErrorResponse{
			StatusCode: httpResponse.StatusCode,
			Message:    httpResponse.Status,
		}
	}

	// map the response to an object value
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return err
	}

	return nil
}

// newRequest builds a new request using the given parameters.
func (c *DefaultClient) newRequest(method string, path string, headers map[string]string, data interface{}) (*http.Request, error) {

	// Construct request body
	var body io.Reader
	if data != nil {
		requestJSON, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}

		body = bytes.NewReader(requestJSON)
	}

	// construct the request
	req, err := http.NewRequest(method, c.BaseURL+path, body)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return req, nil
}



// ----- Add Contact ---------------------------------------------------------------------------------------------------

type AddContactRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type ContactResponse struct {
	Contact *Contact `json:"contact"`
}

func (c *DefaultClient) AddContact(contact AddContactRequest) (*Contact, error) {
	var response ContactResponse
	err := c.performRequestMethod(http.MethodPost, "/contacts", nil, contact, &response)
	if err != nil {
		return nil, err
	}

	return response.Contact, nil
}

func (c *DefaultClient) GetContactByEmail(email string) (*Contact, error) {
	var response ContactResponse
	var path = fmt.Sprintf("/contacts/%v", url.QueryEscape(email))
	err := c.performRequestMethod(http.MethodGet, path, nil, nil, &response)
	if err != nil {
		return nil, err
	}

	return response.Contact, nil
}
