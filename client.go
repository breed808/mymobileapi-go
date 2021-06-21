package mymobileapi

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

// Client represents a client to the mymobile REST API.
type Client struct {
	authorization   string
	authToken       string
	AuthTokenExpiry time.Time
	debug           bool
	endpoint        string
	httpClient      *http.Client
}

// NewClient returns a client.
func NewClient(clientID, clientSecret string, debug bool) (*Client, error) {
	// Create Authorization header from client ID & secret
	authorization := "BASIC " + base64.StdEncoding.EncodeToString([]byte(clientID+":"+clientSecret))

	client := Client{authorization: authorization, endpoint: "https://rest.mymobileapi.com/v1/", debug: debug, httpClient: &http.Client{}}

	err := client.Authenticate()
	if err != nil {
		return nil, err
	}

	return &client, err
}

// Authenticate to API, generating an Authentication token for use in future requests.
func (c *Client) Authenticate() error {
	authResponse := struct {
		Token            string `json:"token"`
		Schema           string `json:"schema"`
		ExpiresInMinutes int    `json:"expiresInMinutes"`
	}{}

	_, err := c.get("Authentication", nil, &authResponse)
	if err != nil {
		return err
	}

	duration, err := time.ParseDuration(fmt.Sprint(authResponse.ExpiresInMinutes) + "m")
	if err != nil {
		return err
	}

	c.AuthTokenExpiry = time.Now().Add(duration)
	c.authToken = "Bearer " + authResponse.Token

	return nil
}

func (c *Client) ask(method, path string, params, recipient interface{}) (http.Header, error) {
	resp, err := c.do(method, path, params, nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(recipient)

	return resp.Header, nil
}

type StandardResponse struct {
	Code    int             `json:"code,omitempty"`
	Message string          `json:"message,omitempty"`
	UUID    string          `json:"uuid,omitempty"`
	Object  string          `json:"object,omitempty"`
	Cause   string          `json:"cause,omitempty"`
	Status  string          `json:"status,omitempty"`
	Errors  []StandardError `json:"errors,omitempty"`
}

type StandardError struct {
	Location    string `json:"location"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (c *Client) do(method, path string, p interface{}, extraHeaders [][2]string) (*http.Response, error) {
	var (
		err error
		req *http.Request
	)

	params, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	suffix := ""

	if params != nil && string(params) != "null" {
		req, err = http.NewRequest(method, c.endpoint+path+suffix, bytes.NewReader(params))
	} else {
		req, err = http.NewRequest(method, c.endpoint+path+suffix, nil)
	}

	if err != nil {
		return nil, err
	}

	// If authentication token is empty, use client ID & secret for authentication.
	// This is only valid for the Authentication endpoint.
	if c.authToken == "" {
		req.Header.Add("Authorization", c.authorization)
	} else {
		req.Header.Add("Authorization", c.authToken)
	}
	req.Header.Add("Content-Type", "application/json")

	for _, header := range extraHeaders {
		req.Header.Add(header[0], header[1])
	}

	if c.debug {
		dump, _ := httputil.DumpRequestOut(req, true)
		fmt.Println("=======================================\nREQUEST:")
		fmt.Println(string(dump))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return resp, err
	}

	if c.debug {
		dump, _ := httputil.DumpResponse(resp, true)

		fmt.Println("=======================================\nRESPONSE:")
		fmt.Println(string(dump))
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		defer resp.Body.Close()

		var message StandardResponse

		decoder := json.NewDecoder(resp.Body)
		decoder.Decode(&message)
		if message.Message != "" {
			err = fmt.Errorf("%d: %s", resp.StatusCode, message.Message)
		} else if len(message.Errors) > 0 {
			var errors []string

			for _, oneError := range message.Errors {
				errors = append(errors, fmt.Sprintf("%s: %s", oneError.Name, oneError.Description))
			}

			err = fmt.Errorf(strings.Join(errors, ", "))
		} else {
			err = fmt.Errorf("%d", resp.StatusCode)
		}
	}

	return resp, err
}

func (c *Client) get(path string, params, recipient interface{}) (http.Header, error) {
	return c.ask(http.MethodGet, path, params, recipient)
}

func (c *Client) post(path string, params, recipient interface{}) (http.Header, error) {
	return c.ask(http.MethodPost, path, params, recipient)
}

func (c *Client) put(path string, params, recipient interface{}) (http.Header, error) {
	return c.ask(http.MethodPut, path, params, recipient)
}

func (c *Client) patch(path string, params, recipient interface{}) (http.Header, error) {
	return c.ask(http.MethodPatch, path, params, recipient)
}

func (c *Client) delete(path string, params, recipient interface{}) (http.Header, error) {
	return c.ask(http.MethodDelete, path, params, recipient)
}
