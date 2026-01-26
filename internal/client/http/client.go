package http

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type HttpClient struct {
	Client    *http.Client
	BaseURL   string
	AuthToken string
}

func CreateHttpClient() *HttpClient {
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		fmt.Fprint(os.Stderr, "No base url env variable found\n")
		return &HttpClient{}
	}
	client := &HttpClient{
		Client: &http.Client{
			Timeout: 20 * time.Second,
		},
		BaseURL: baseURL,
	}
	token := loadToken()
	client.SetSession(token)
	return client
}

func (c *HttpClient) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *HttpClient) Delete(url string) (*http.Response, error) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *HttpClient) Post(url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return c.Do(req)
}

func (c *HttpClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", "Bearer "+c.AuthToken)
	return c.Client.Do(req)
}

func (c *HttpClient) SetAuthToken(token string) {
	c.AuthToken = token
}

func (c *HttpClient) SetSession(token string) {
	c.AuthToken = token
	saveToken(token)
}
