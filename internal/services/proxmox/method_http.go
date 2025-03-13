package proxmox

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
)

func Header() map[string]string {
	data := map[string]string{
		"Content-Type": "application/json",
		"User-Agent":   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36",
	}
	return data
}

func Get(ctx context.Context, url, token string) (*http.Response, error) {
	if url == "" {
		return nil, errors.New("url is required")
	}
	if token == "" {
		return nil, errors.New("token is required")
	}
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	header := Header()
	for k, v := range header {
		req.Header.Set(k, v)
	}
	req.Header.Set("Authorization", token)
	client := &http.Client{Transport: transport}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%v", resp.StatusCode)
	}
	return resp, nil
}
