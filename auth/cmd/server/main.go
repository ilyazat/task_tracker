package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
	"sync"

	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
)

type UserStorage struct {
	mu sync.RWMutex
	m  map[string]User
}

func NewUserStorage() *UserStorage {
	m := make(map[string]User, 100)

	return &UserStorage{m: m}
}

type User struct {
	UserID string `json:"name"`
	Role   string `json:"role"`
}

func (us *UserStorage) GetByName(name string) (User, bool) {
	us.mu.RLock()
	defer us.mu.RUnlock()

	user, ok := us.m[name]

	return user, ok
}

func (us *UserStorage) GetOrCreate(user User) (User, error) {
	us.mu.Lock()
	defer us.mu.Unlock()

	var u User

	if _, ok := us.m[user.UserID]; !ok {
		us.m[user.UserID] = user
		u = user
	} else {
		u = us.m[user.UserID]
	}

	return u, nil

}

var (
	Port         = ":5555"
	clientID     = "228"
	clientSecret = "fuckimperialism"
)

func main() {
	manager := manage.NewDefaultManager()
	manager.MustTokenStorage(store.NewMemoryTokenStore())

	clientStore := store.NewClientStore()
	if err := clientStore.Set(clientID, &models.Client{ID: clientID, Secret: clientSecret, Domain: "http://localhost:5554"}); err != nil {
		panic(err)
	}
	manager.MapClientStorage(clientStore)
	manager.MapAccessGenerate(generates.NewJWTAccessGenerate("", []byte("00000000"), jwt.SigningMethodHS512))

	userStorage := NewUserStorage()

	srv := server.NewDefaultServer(manager)
	srv.SetAllowGetAccessRequest(true)
	srv.SetClientInfoHandler(server.ClientFormHandler)

	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

	srv.SetUserAuthorizationHandler(userAuthorizeHandler)

	http.HandleFunc("/authorize", func(w http.ResponseWriter, r *http.Request) {
		if err := srv.HandleAuthorizeRequest(w, r); err != nil {
			log.Println(err)
		}
	})

	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		if err := srv.HandleTokenRequest(w, r); err != nil {
			log.Println(err)
		}
	})

	http.HandleFunc("/verify", func(w http.ResponseWriter, r *http.Request) {
		_, err := srv.ValidationBearerToken(r)
		accessToken, ok := srv.BearerAuth(r)

		if err != nil || !ok {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		token, err := jwt.ParseWithClaims(accessToken, &generates.JWTAccessClaims{}, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("parse error")
			}
			return []byte("00000000"), nil
		})
		if err != nil {
			http.Error(w, "bad token", http.StatusBadRequest)
		}

		var claims interface{}

		if token.Valid {
			claims = token.Claims.(*generates.JWTAccessClaims)
		} else {
			http.Error(w, "bad token", http.StatusBadRequest)
			return
		}

		user, ok := userStorage.GetByName(claims.(*generates.JWTAccessClaims).Subject)
		if !ok {
			http.Error(w, "failed to find user with this token", http.StatusInternalServerError)
			return
		}

		jsonBytes, err := json.Marshal(user)
		if err != nil {
			http.Error(w, "failed to marshal json", http.StatusInternalServerError)
			return
		}

		_, err = w.Write(jsonBytes)
		if err != nil {
			log.Println(err)
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

func userAuthorizeHandler(_ http.ResponseWriter, r *http.Request) (string, error) {
	userID := r.URL.Query().Get("user_id")

	if userID == "" {
		return userID, errors.New("userID is empty")
	}

	return userID, nil
}
