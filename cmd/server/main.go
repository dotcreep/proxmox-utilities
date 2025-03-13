package main

import (
	"fmt"
	"log"
	"net/http"

	apiidrac "github.com/dotcreep/mondayreport/internal/apis/idrac"
	apiproxmox "github.com/dotcreep/mondayreport/internal/apis/proxmox"
	"github.com/dotcreep/mondayreport/internal/config"
	"github.com/dotcreep/mondayreport/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := utils.Response{}
		ymlconf, err := config.OpenConfig()
		if err != nil {
			log.Println(err)
			response.Json(false, w, nil, "config not found", http.StatusInternalServerError, nil)
			return
		}
		apiKey := ymlconf.Config.Token
		if apiKey == "" {
			log.Println("token not found")
			response.Json(false, w, nil, "token not found", http.StatusUnauthorized, nil)
			return
		}
		if apiKey != r.Header.Get("X-API-Key") {
			log.Println("unauthorized")
			response.Json(false, w, nil, "unauthorized", http.StatusUnauthorized, nil)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	_, err := config.OpenConfig()
	if err != nil {
		panic(err)
	}
	response := utils.Response{}

	r := chi.NewRouter()
	cors := cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-API-Key"},
		AllowCredentials: true,
	})
	r.Use(cors)
	r.Route("/status", func(r chi.Router) {
		r.Use(Middleware)
		r.Get("/proxmox", apiproxmox.GetNodeStatus)
		r.Get("/idrac", apiidrac.GetIDRacStatus)
	})
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		response.Json(false, w, nil, "not found", http.StatusNotFound, nil)
	})
	fmt.Println("running on http://localhost:8888")
	log.Println(http.ListenAndServe(":8888", r))

}
