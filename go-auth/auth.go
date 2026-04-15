package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// Jangan pakai global config langsung (env belum tentu kebaca)
func getGoogleConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

// Redirect ke Google
func GoogleLogin(w http.ResponseWriter, r *http.Request) {
	config := getGoogleConfig()
	url := config.AuthCodeURL("random-state")

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// Callback dari Google
func GoogleCallback(w http.ResponseWriter, r *http.Request) {
	config := getGoogleConfig()

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	// tukar code ke token
	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Token exchange failed", http.StatusInternalServerError)
		return
	}

	client := config.Client(context.Background(), token)

	// ambil data user
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var user map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		http.Error(w, "Decode error", http.StatusInternalServerError)
		return
	}

	// amanin type assertion
	email, ok := user["email"].(string)
	if !ok {
		http.Error(w, "Invalid email", http.StatusInternalServerError)
		return
	}

	name, _ := user["name"].(string)
	googleID, _ := user["id"].(string)

	// UPSERT (insert atau update)
	query := `
	INSERT INTO users (email, name, google_id)
	VALUES (?, ?, ?)
	ON DUPLICATE KEY UPDATE
		name = VALUES(name),
		google_id = VALUES(google_id)
	`

	_, err = DB.Exec(query, email, name, googleID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Login sukses 🚀"))
}
