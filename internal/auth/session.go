package auth

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

type Session struct {
	ID        string
	UserID    string
	CreatedAt time.Time
	ExpiresAt time.Time
}

type SessionManager struct {
	sessions map[string]*Session
	mutex    sync.RWMutex
	secret   string
}

func NewSessionManager(secret string) *SessionManager {
	sm := &SessionManager{
		sessions: make(map[string]*Session),
		secret:   secret,
	}

	// cleanup goroutine for expired sessions
	go sm.cleanupExpiredSessions()

	return sm
}

func (sm *SessionManager) CreateSession(userID string) (*Session, error) {
	sessionID, err := generateSessionID()
	if err != nil {
		return nil, err
	}

	session := &Session{
		ID:        sessionID,
		UserID:    userID,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour), //24 hour expiration, todo: export to env
	}

	(*sm).mutex.Lock()
	(*sm).sessions[sessionID] = session
	(*sm).mutex.Unlock()

	return session, nil
}

func (sm *SessionManager) GetSession(sessionID string) (*Session, bool) {
	(*sm).mutex.RLock()
	session, exists := (*sm).sessions[sessionID]
	(*sm).mutex.RUnlock()

	if !exists || session.ExpiresAt.Before(time.Now()) {
		if exists {
			(*sm).DeleteSession(sessionID)
		}
		return nil, false
	}

	return session, true
}

func (sm *SessionManager) DeleteSession(sessionID string) {
	(*sm).mutex.Lock()
	delete(sm.sessions, sessionID)
	(*sm).mutex.Unlock()
}

// -- helpers --
func (sm *SessionManager) cleanupExpiredSessions() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		(*sm).mutex.Lock()
		now := time.Now()
		for id, session := range (*sm).sessions {
			if session.ExpiresAt.Before(now) {
				delete((*sm).sessions, id)
			}
		}
		(*sm).mutex.Unlock()
	}
}

func generateSessionID() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
