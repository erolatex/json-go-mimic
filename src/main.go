package main

import (
	"encoding/json"
	"fmt"
	"io"
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
		ContentType  string `json:"contentType"`
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
				var token string
				if e.AuthType == "Bearer" {
					token = strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
				} else {
					token = r.Header.Get(e.AuthType)
				}

				if token != e.AuthKey {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusUnauthorized)
					jsonResponse := `{"error": "Unauthorized"}`
					if _, err := io.WriteString(w, jsonResponse); err != nil {
						log.Printf("Error sending unauthorized response for path %s: %v", e.Path, err)
					}
					return
				}
			}

			if e.ContentType == "" {
				w.Header().Set("Content-Type", "application/json")
			} else {
				w.Header().Set("Content-Type", e.ContentType)
			}
			if _, err := w.Write(jsonFile); err != nil {
				log.Printf("Error writing JSON response for path %s: %v", e.Path, err)
			}
		})
	}

	addr := fmt.Sprintf(":%d", config.Port)
	log.Printf("Starting server on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
