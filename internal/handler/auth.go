package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"blog-system/internal/auth"

	"github.com/gorilla/mux"
)

type AuthHandler struct {
	sessionManager *auth.SessionManager
	adminPassword  string
}

func NewAuthHandler(sessionManager *auth.SessionManager, adminPassword string) *AuthHandler {
	return &AuthHandler{
		sessionManager: sessionManager,
		adminPassword:  adminPassword,
	}
}

func (h *AuthHandler) RegisterRoutes(r *mux.Router) {
	(*r).HandleFunc("/auth/login", (*h).Login).Methods("POST")
	(*r).HandleFunc("/auth/logout", (*h).Logout).Methods("POST")
	(*r).HandleFunc("/auth/status", (*h).Status).Methods("GET")
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Password != (*h).adminPassword {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	session, err := (*h).sessionManager.CreateSession("admin")
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		Path:     "/",
		Expires:  session.ExpiresAt,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":    "Login successful",
		"expires_at": session.ExpiresAt,
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == nil {
		h.sessionManager.DeleteSession(cookie.Value)
	}

	clearCookie := &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, clearCookie)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Logout successful",
	})
}

func (h *AuthHandler) Status(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"authenticated": false,
		})
		return
	}

	session, exists := (*h).sessionManager.GetSession(cookie.Value)
	if !exists {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"authenticated": false,
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"authenticated": true,
		"user_id":       session.UserID,
		"expires_at":    session.ExpiresAt,
	})
}
