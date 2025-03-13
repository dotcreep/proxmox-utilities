package idrac

import (
	"context"
	"crypto/tls"
	"encoding/base64"
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

type Login struct {
	Server   string `json:"server"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func Get(ctx context.Context, login *Login) (*http.Response, error) {
	if login.Server == "" {
		return nil, errors.New("server is required")
	}
	if login.Username == "" {
		return nil, errors.New("username is required")
	}
	if login.Password == "" {
		return nil, errors.New("password is required")
	}
	url := fmt.Sprintf("https://%s/redfish/v1/Chassis/System.Embedded.1", login.Server)
	token := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", login.Username, login.Password)))
	if url == "" {
		return nil, errors.New("url is required")
	}
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	header := Header()
	for k, v := range header {
		req.Header.Set(k, v)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", token))
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
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
