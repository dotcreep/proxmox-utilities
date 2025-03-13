package apiidrac

import (
	"context"
	"net/http"

	"github.com/dotcreep/mondayreport/internal/services/idrac"
	"github.com/dotcreep/mondayreport/internal/utils"
)

func GetIDRacStatus(w http.ResponseWriter, r *http.Request) {
	response := utils.Response{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if list, err := idrac.GetStatus(ctx); err != nil {
		response.Json(false, w, list, err.Error(), http.StatusInternalServerError, err)
		return
	} else {
		response.Json(true, w, list, "success get idrac status", http.StatusOK, nil)
		return
	}
}
