package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Configuration struct {
	Config     Config     `json:"config"`
	Proxmox    []Proxmox  `json:"proxmox"`
	IDRac      []IDRac    `json:"idrac"`
	Cloudflare Cloudflare `json:"cloudflare"`
}

type Config struct {
	Token string `json:"token"`
}

type Proxmox struct {
	Server string `json:"server"`
	Node   string `json:"node"`
	Token  string `json:"token"`
}

type IDRac struct {
	Server   string `json:"server"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Cloudflare struct {
	BaseURL   string `json:"base_url"`
	Key       string `json:"key"`
	Email     string `json:"email"`
	AccountID string `json:"account_id"`
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func OpenConfig() (*Configuration, error) {
	var config Configuration
	files := []string{"config.yml", "config.yaml"}
	for _, file := range files {
		if PathExists(file) {
			yamlFile, err := os.ReadFile(file)
			if err != nil {
				return nil, err
			}
			err = yaml.Unmarshal(yamlFile, &config)
			if err != nil {
				return nil, err
			}
		}
	}
	return &config, nil
}
