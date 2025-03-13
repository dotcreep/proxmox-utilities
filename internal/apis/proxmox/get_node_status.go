package apiproxmox

import (
	"context"
	"net/http"

	"github.com/dotcreep/mondayreport/internal/services/proxmox"
	"github.com/dotcreep/mondayreport/internal/utils"
)

func GetNodeStatus(w http.ResponseWriter, r *http.Request) {
	response := utils.Response{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var data []string
	if list, err := proxmox.GetNodeStatus(ctx); err != nil {
		response.Json(false, w, data, err.Error(), http.StatusInternalServerError, err)
		return
	} else {
		response.Json(true, w, list, "success get node status", http.StatusOK, nil)
		return
	}
}
