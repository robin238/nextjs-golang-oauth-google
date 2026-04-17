package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var revokedTokens = make(map[string]time.Time)
var revokedTokensMu sync.Mutex

func RevokeToken(token string) {
	revokedTokensMu.Lock()
	defer revokedTokensMu.Unlock()
	revokedTokens[token] = time.Now()
}

func IsTokenRevoked(token string) bool {
	revokedTokensMu.Lock()
	defer revokedTokensMu.Unlock()
	_, revoked := revokedTokens[token]
	return revoked
}

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

	// ambil user info
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

	email, ok := user["email"].(string)
	if !ok {
		http.Error(w, "Invalid email", http.StatusInternalServerError)
		return
	}

	name, _ := user["name"].(string)
	googleID, _ := user["id"].(string)

	// UPSERT + default role user
	query := `
	INSERT INTO users (email, name, google_id, role)
	VALUES (?, ?, ?, 'user')
	ON DUPLICATE KEY UPDATE
		name = VALUES(name),
		google_id = VALUES(google_id)
	`

	_, err = DB.Exec(query, email, name, googleID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// ambil role dari DB
	var role string
	err = DB.QueryRow("SELECT role FROM users WHERE email = ?", email).Scan(&role)
	if err != nil {
		http.Error(w, "Failed to get role", http.StatusInternalServerError)
		return
	}

	// generate JWT
	jwtToken, err := GenerateJWT(email, role, name)
	if err != nil {
		http.Error(w, "JWT error", http.StatusInternalServerError)
		return
	}

	// redirect ke frontend
	redirectURL := "http://localhost:3000/dashboard?token=" + jwtToken
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

func GenerateJWT(email string, role string, name string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"name":  name,
		"role":  role,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
