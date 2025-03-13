package config

import "fmt"

func GetBaseURLProxmox(server string) string {
	if server == "" {
		return ""
	}
	url := fmt.Sprintf("https://%s:8006/api2/json", server)
	return url
}
