package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type User struct {
	UserID string `json:"name"`
	Role   string `json:"role"`
}

var (
	Port         = ":5554"
	clientID     = "228"
	redirectURI  = "http://localhost:5554/code"
	grantType    = "authorization_code"
	clientSecret = "fuckimperialism"
)

func main() {
	http.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		var inputUser User

		if err := json.NewDecoder(r.Body).Decode(&inputUser); err != nil {
			http.Error(w, "failed to decode JSON body", http.StatusBadRequest)
			return
		}

		fullURL := getURL(inputUser.UserID)

		resp, err := http.Post(fullURL, "application/json", r.Body)
		if err != nil {
			fmt.Println("failed to get auth_code", err)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("failed to read response body", err)
			return
		}

		w.WriteHeader(http.StatusOK)

		w.Write(body)
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		var inputUser User

		if err := json.NewDecoder(r.Body).Decode(&inputUser); err != nil {
			http.Error(w, "Failed to decode JSON body", http.StatusBadRequest)
			return
		}

	})

	http.HandleFunc("/code", func(w http.ResponseWriter, r *http.Request) {
		authCode := r.URL.Query().Get("code")

		if authCode == "" {
			http.Error(w, "failed to get authorization code", http.StatusBadRequest)
			return
		}

		redirectURL, err := url.Parse("http://localhost:5555/token")
		if err != nil {
			http.Error(w, "failed to parse redirect URI", http.StatusInternalServerError)
			return
		}

		q := redirectURL.Query()
		q.Add("grant_type", grantType)
		q.Add("client_id", clientID)
		q.Add("client_secret", clientSecret)
		q.Add("redirect_uri", redirectURI)
		q.Add("scope", "read")
		q.Add("code", authCode)

		redirectURL.RawQuery = q.Encode()

		newReq, err := http.NewRequest(http.MethodPost, redirectURL.String(), r.Body)
		if err != nil {
			http.Error(w, "Failed to create a new request", http.StatusInternalServerError)
			return
		}
		defer newReq.Body.Close()

		newReq.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(newReq)
		if err != nil {
			http.Error(w, "failed to make POST request", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
		w.WriteHeader(resp.StatusCode)

		if _, err = io.Copy(w, resp.Body); err != nil {
			http.Error(w, "failed to transfer body", http.StatusInternalServerError)
			return
		}

	})

	handler := corsMiddleware(http.DefaultServeMux)

	log.Printf("starting to serve at %s", Port)
	log.Fatal(http.ListenAndServe(Port, handler))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getURL(userID string) string {
	baseURL := "http://localhost:5555/authorize"
	params := url.Values{}
	params.Set("client_id", clientID)
	params.Set("response_type", "code")
	params.Set("user_id", userID)
	params.Set("redirect_uri", redirectURI)
	params.Set("client_secret", clientSecret)
	params.Set("grant_type", grantType)

	return fmt.Sprintf("%s?%s", baseURL, params.Encode())
}
