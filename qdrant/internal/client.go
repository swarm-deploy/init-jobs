package internal

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
)

type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
}

type collectionCreateRequest struct {
	Vectors vectorParams `json:"vectors"`
}

type vectorParams struct {
	Size     int    `json:"size"`
	Distance string `json:"distance"`
}

type Config struct {
	// Addr is the Qdrant HTTP endpoint.
	Addr string `env:"ADDR,notEmpty,required"`
	// Insecure if true request sends without TLS verify.
	Insecure bool `env:"INSECURE"`
}

func NewClient(cfg Config) (*Client, error) {
	baseURL, err := readAddress(cfg.Addr)
	if err != nil {
		return nil, err
	}

	transport := http.DefaultTransport
	if cfg.Insecure {
		transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Transport: transport,
		},
	}, nil
}

func (c *Client) CollectionExists(ctx context.Context, collectionName string) (bool, error) {
	collectionURL := c.collectionURL(collectionName)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, collectionURL, nil)
	if err != nil {
		return false, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	return false, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
}

func (c *Client) CreateCollection(ctx context.Context, collectionName string, vectorSize int, vectorDistance string) error {
	collectionURL := c.collectionURL(collectionName)

	payload, err := json.Marshal(collectionCreateRequest{
		Vectors: vectorParams{
			Size:     vectorSize,
			Distance: vectorDistance,
		},
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, collectionURL, bytes.NewReader(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	return nil
}

func (c *Client) collectionURL(collectionName string) string {
	collectionURL := *c.baseURL
	collectionURL.Path = path.Join(collectionURL.Path, "collections", collectionName)

	return collectionURL.String()
}

func readAddress(addr string) (*url.URL, error) {
	if !strings.HasPrefix(addr, "http://") && !strings.HasPrefix(addr, "https://") {
		addr = "http://" + addr
	}

	parsed, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	if parsed.Host == "" {
		return nil, fmt.Errorf("addr is empty")
	}

	return parsed, nil
}
