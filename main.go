package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

//
type Store interface {
	GetUserByUsername(username string) (*User, error)
	CreateSession() (token string, err error)
}

//
type handler struct {
	store Store
}

//
func NewHandler(store Store) *handler {
	return &handler{
		store: store,
	}
}

//
func (h *handler) login(w http.ResponseWriter, r *http.Request) {
	//
	var authData struct{
		Username string `json:"username"`
		Password string `json:"password"`
	}
	//
	if err := json.NewDecoder(r.Body).Decode(&authData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//
	user, err := h.store.GetUserByUsername(authData.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	//
	hasher := sha256.New()
	//
	hasher.Write([]byte(authData.Password))
	passwordHash := hasher.Sum(nil)

	//
	if subtle.ConstantTimeCompare(user.PasswordHash, passwordHash) == 1 {
		//
		token, err := h.store.CreateSession()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		//
		resp := struct {
			Token string `json:"token"`
		}{
			Token: token,
		}
		//
		b, _ := json.Marshal(resp)

		//
		w.WriteHeader(http.StatusOK)
		//
		if _, err := w.Write(b); err != nil {
			log.Panicln("Error while trying to write! ", err)
			return
		}
	} else {
		//
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
}

//
func main() {
	//
	db, err := sql.Open("postgres", "user=postgres dbname=simple_app")
	if err != nil {
		panic(err)
	}
	//
	store := NewSimpleStore(db)
	//
	h := NewHandler(store)

	//
	http.HandleFunc("/login", h.login)
	//
	http.ListenAndServe(":8090", nil)
}