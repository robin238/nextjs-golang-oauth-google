package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleConfig = &oauth2.Config{
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
	Scopes:       []string{"email", "profile"},
	Endpoint:     google.Endpoint,
}

func GoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := googleConfig.AuthCodeURL("state")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func GoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	token, _ := googleConfig.Exchange(context.Background(), code)

	client := googleConfig.Client(context.Background(), token)

	resp, _ := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")

	defer resp.Body.Close()

	var user map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&user)

	email := user["email"].(string)
	name := user["name"].(string)
	googleID := user["id"].(string)

	// insert ke DB
	DB.Exec("INSERT IGNORE INTO users(email, name, google_id) VALUES(?,?,?)",
		email, name, googleID)

	w.Write([]byte("Login sukses"))
}
