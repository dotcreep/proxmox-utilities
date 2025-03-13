package idrac

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/dotcreep/mondayreport/internal/config"
)

type RedFishResponse struct {
	Status Status `json:"Status"`
}

type Status struct {
	Health       string `json:"Health"`
	HealthRollup string `json:"HealthRollup"`
	State        string `json:"State"`
}

type IDRacArrayResponse struct {
	IDRacResult []IDRacResult `json:"idrac"`
}

type IDRacResult struct {
	ServerIP string `json:"server_ip"`
	Status   Status `json:"status"`
}

func GetStatus(ctx context.Context) (*IDRacArrayResponse, error) {
	configYaml, err := config.OpenConfig()
	if err != nil {
		return nil, err
	}
	var data IDRacArrayResponse
	var mux sync.Mutex
	var wg sync.WaitGroup
	errChan := make(chan error, len(configYaml.IDRac))
	// var subdata RedFishResponse
	for _, server := range configYaml.IDRac {
		wg.Add(1)

		go func(server config.IDRac) {
			defer wg.Done()
			var subdata IDRacResult
			// Check server
			if server.Server == "" {
				errChan <- errors.New("server is required")
				return
			}
			// Check username
			if server.Username == "" {
				errChan <- errors.New("username is required")
				return
			}
			// Check password
			if server.Password == "" {
				errChan <- errors.New("password is required")
				return
			}
			// Initial login
			login := &Login{
				Server:   server.Server,
				Username: server.Username,
				Password: server.Password,
			}
			// Get status from server IDRac
			resp, err := Get(ctx, login)
			if err != nil {
				errChan <- fmt.Errorf("failed to get status from server %s: %w", server.Server, err)
				return
			}
			if resp != nil {
				defer resp.Body.Close()
				if resp.StatusCode == 404 {
					errChan <- fmt.Errorf("server %s not found: %w", server.Server, err)
					return
				}
			} else {
				errChan <- fmt.Errorf("failed to get status from server %s: %w", server.Server, err)
				return
			}
			// Decode response

			err = json.NewDecoder(resp.Body).Decode(&subdata)
			if err != nil {
				errChan <- fmt.Errorf("failed to decode from server %s: %w", server.Server, err)
				return
			}
			mux.Lock()
			data.IDRacResult = append(data.IDRacResult, IDRacResult{
				ServerIP: server.Server,
				Status:   subdata.Status,
			})
			mux.Unlock()
		}(server)
	}
	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}
	// Initial login
	// login := &Login{
	// 	Server:   configYaml.IDRac[i].Server,
	// 	Username: configYaml.IDRac[i].Username,
	// 	Password: configYaml.IDRac[i].Password,
	// }
	// if login.Server == "" {
	// 	return nil, errors.New("server is required")
	// }
	// if login.Username == "" {
	// 	return nil, errors.New("username is required")
	// }
	// if login.Password == "" {
	// 	return nil, errors.New("password is required")
	// }
	// // Get status from server IDRac
	// resp, err := Get(ctx, login)
	// if err != nil {
	// 	return nil, err
	// }
	// defer resp.Body.Close()
	// // Decode response
	// err = json.NewDecoder(resp.Body).Decode(&subdata)
	// if err != nil {
	// 	return nil, err
	// }

	// data.IDRacResult = append(data.IDRacResult, IDRacResult{
	// 	ServerIP: server.Server,
	// 	Status:   subdata.Status,
	// })
	if len(data.IDRacResult) == 0 {
		return nil, errors.New("no data")
	}
	return &data, nil
}
