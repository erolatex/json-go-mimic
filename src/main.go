package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

type Config struct {
	Port      int `json:"port"`
	Endpoints []struct {
		Path         string `json:"path"`
		JsonFilePath string `json:"jsonFilePath"`
		AuthType     string `json:"authType"`
		AuthKey      string `json:"authKey"`
	} `json:"endpoints"`
}

func main() {
	configFile, err := os.ReadFile("configs/config.json")
	if err != nil {
		log.Fatal("Error reading config file: ", err)
	}

	var config Config
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatal("Error parsing config file: ", err)
	}

	for _, endpoint := range config.Endpoints {
		e := endpoint
		jsonFile, err := os.ReadFile(e.JsonFilePath)
		if err != nil {
			log.Fatalf("Error reading JSON file for path %s: %v", e.Path, err)
		}

		http.HandleFunc(e.Path, func(w http.ResponseWriter, r *http.Request) {
			if e.AuthType != "None" {
				if e.AuthType == "Bearer" {
					token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
					if token != e.AuthKey {
						http.Error(w, "Unauthorized", http.StatusUnauthorized)
						return
					}
				} else if r.Header.Get(e.AuthType) != e.AuthKey {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}
			}

			w.Header().Set("Content-Type", "application/json")
			if _, writeErr := w.Write(jsonFile); writeErr != nil {
				log.Printf("Error writing JSON response for path %s: %v", e.Path, writeErr)
			}
		})
	}

	addr := fmt.Sprintf(":%d", config.Port)
	log.Printf("Starting server on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
