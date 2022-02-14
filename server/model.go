package server

import "time"

type create struct {
	OriginalURL string    `json:"original_url"`
	ExpiresAt   time.Time `json:"expires_at"`
}

type entity struct {
	ID          string    `json:"id" redis:"id" firestore:"id"`
	ShortURL    string    `json:"short_url" redis:"short_url" firestore:"short_url"`
	OriginalURL string    `json:"original_url" redis:"original_url" firestore:"original_url"`
	CreatedAt   time.Time `json:"created_at" redis:"created_at" firestore:"created_at"`
	ExpiresAt   time.Time `json:"expires_at" redis:"expires_at" firestore:"expires_at"`
}
