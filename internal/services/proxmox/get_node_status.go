package proxmox

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/dotcreep/mondayreport/internal/config"
)

// Proxmox Status
type ProxmoxStatus struct {
	Data Data `json:"data"`
}

type Data struct {
	LoadAvg    []string `json:"loadavg"`
	Idle       int      `json:"idle"`
	CPU        float64  `json:"cpu"`
	KSM        KSM      `json:"ksm"`
	RootFS     RootFS   `json:"rootfs"`
	PveVersion string   `json:"pveversion"`
	KVersion   string   `json:"kversion"`
	Uptime     int      `json:"uptime"`
	Wait       float64  `json:"wait"`
	Memory     Memory   `json:"memory"`
	CPUInfo    CPUInfo  `json:"cpuinfo"`
	Swap       Swap     `json:"swap"`
}

type KSM struct {
	Shared int64 `json:"shared"`
}

type RootFS struct {
	Avail int64 `json:"avail"`
	Used  int64 `json:"used"`
	Free  int64 `json:"free"`
	Total int64 `json:"total"`
}

type Memory struct {
	Free  int64 `json:"free"`
	Used  int64 `json:"used"`
	Total int64 `json:"total"`
}

type CPUInfo struct {
	Sockets int    `json:"sockets"`
	CPUs    int    `json:"cpus"`
	Cores   int    `json:"cores"`
	UserHz  int    `json:"user_hz"`
	MHz     string `json:"mhz"`
	HVM     string `json:"hvm"`
	Model   string `json:"model"`
	Flags   string `json:"flags"`
}

type Swap struct {
	Free  int64 `json:"free"`
	Used  int64 `json:"used"`
	Total int64 `json:"total"`
}

// Proxmox Tasks
type ProxmoxTasks struct {
	Total int `json:"total"`
	Data  []struct {
		Node      string `json:"node"`
		Starttime int64  `json:"starttime"`
		Status    string `json:"status"`
		User      string `json:"user"`
		Type      string `json:"type"`
		PID       int64  `json:"pid"`
		EndTime   int64  `json:"endtime"`
		UpID      string `json:"upid"`
		ID        string `json:"id"`
		PStart    int64  `json:"pstart"`
	} `json:"data"`
}

type ProxmoxResultArray struct {
	ProxmoxResult []ProxmoxResult `json:"proxmox_result"`
}

type ProxmoxResult struct {
	ProxmoxIP string `json:"proxmox_ip"`
	CPUUsage  string `json:"cpu_usage"`
	MemUsage  string `json:"mem_usage"`
	DiskUsage string `json:"disk_usage"`
	Backup    string `json:"backup"`
}

func GetNodeStatus(ctx context.Context) (*ProxmoxResultArray, error) {
	configYaml, err := config.OpenConfig()
	if err != nil {
		return nil, err
	}
	var data ProxmoxResultArray
	var mux sync.Mutex
	var wg sync.WaitGroup
	errChan := make(chan error, len(configYaml.Proxmox))
	// Access per Node
	for _, server := range configYaml.Proxmox {
		wg.Add(1)

		go func(server config.Proxmox) {
			defer wg.Done()
			var subdata ProxmoxResult
			// Check server
			if server.Server == "" {
				errChan <- errors.New("server is required")
				return
			}
			// Check node
			if server.Node == "" {
				errChan <- errors.New("node is required")
				return
			}
			// Check token
			if server.Token == "" {
				errChan <- errors.New("token is required")
				return
			}
			// Get status of each node
			baseURL := config.GetBaseURLProxmox(server.Server)
			url := fmt.Sprintf("%s/nodes/%s/status", baseURL, server.Node)
			respTask, err := Get(ctx, url, server.Token)
			if err != nil {
				errChan <- err
				return
			}
			defer respTask.Body.Close()
			// Decode response
			var dataStatus ProxmoxStatus
			err = json.NewDecoder(respTask.Body).Decode(&dataStatus)
			if err != nil {
				errChan <- err
				return
			}
			// Response for status each node
			subdata.ProxmoxIP = server.Server
			subdata.CPUUsage = fmt.Sprintf("%.2f", dataStatus.Data.CPU*100)
			subdata.MemUsage = fmt.Sprintf("%.2f", float64(dataStatus.Data.Memory.Used)/float64(dataStatus.Data.Memory.Total)*100)
			subdata.DiskUsage = fmt.Sprintf("%.2f", float64(dataStatus.Data.RootFS.Used)/float64(dataStatus.Data.RootFS.Total)*100)

			// Get tasks for backup
			baseURL = config.GetBaseURLProxmox(server.Server)
			url = fmt.Sprintf("%s/nodes/%s/tasks", baseURL, server.Node)
			respTask, err = Get(ctx, url, server.Token)
			if err != nil {
				errChan <- err
				return
			}
			defer respTask.Body.Close()
			// Decode response
			var dataTasks ProxmoxTasks
			err = json.NewDecoder(respTask.Body).Decode(&dataTasks)
			if err != nil {
				errChan <- err
				return
			}
			// Response for tasks
			// Filter by time and type backup
			var backup []string
			for _, task := range dataTasks.Data {
				if task.Type == "vzdump" {
					oneDay := task.Starttime - 86400
					if task.Starttime >= oneDay {
						backup = append(backup, task.Status)
					}
				}
			}

			// Set backup status
			subdata.Backup = "OK"
			for _, status := range backup {
				if status != "OK" {
					subdata.Backup = status
					break
				}
			}
			mux.Lock()
			data.ProxmoxResult = append(data.ProxmoxResult, subdata)
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
	// 	}
	// 	// Set default value
	// 	subdata.ProxmoxIP = server.Server
	// 	// Get status of each node
	// 	baseURL := config.GetBaseURLProxmox(server.Server)
	// 	url := fmt.Sprintf("%s/nodes/%s/status", baseURL, server.Node)
	// 	respTask, err := Get(ctx, url, server.Token)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	defer respTask.Body.Close()
	// 	// Decode response
	// 	var dataStatus ProxmoxStatus
	// 	err = json.NewDecoder(respTask.Body).Decode(&dataStatus)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	// Response for status each node
	// 	subdata.CPUUsage = fmt.Sprintf("%.2f", dataStatus.Data.CPU*100)
	// 	subdata.MemUsage = fmt.Sprintf("%.2f", float64(dataStatus.Data.Memory.Used)/float64(dataStatus.Data.Memory.Total)*100)
	// 	subdata.DiskUsage = fmt.Sprintf("%.2f", float64(dataStatus.Data.RootFS.Used)/float64(dataStatus.Data.RootFS.Total)*100)

	// 	// Get tasks for backup
	// 	taskURL := fmt.Sprintf("%s/nodes/%s/status", baseURL, server.Node)
	// 	taskResp, err := Get(ctx, taskURL, server.Token)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	defer taskResp.Body.Close()
	// 	// Decode response
	// 	var dataTasks ProxmoxTasks
	// 	err = json.NewDecoder(taskResp.Body).Decode(&data)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	// Filter by time and type backup
	// 	// Get tasks by type backup
	// 	var backup []string
	// 	for _, task := range dataTasks.Data {
	// 		if task.Type == "vzdump" {
	// 			oneDay := task.Starttime - 86400
	// 			if task.Starttime >= oneDay {
	// 				backup = append(backup, task.Status)
	// 			}
	// 		}
	// 	}
	// 	// Set default value
	// 	subdata.Backup = "OK"
	// 	for _, status := range backup {
	// 		if status != "OK" {
	// 			subdata.Backup = status
	// 			break
	// 		}
	// 	}
	// 	// Append data
	// 	data.ProxmoxResult = append(data.ProxmoxResult, subdata)
	// }
	if len(data.ProxmoxResult) == 0 {
		return nil, errors.New("node not found")
	}
	return &data, nil
}
