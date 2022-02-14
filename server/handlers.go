package server

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
	"time"
)

type service struct {
	port      string
	host      string
	keyLength int
	cache     cacher
	persist   persister
}

func (s *service) redirectHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	key := chi.URLParam(r, "key")
	now := time.Now()

	entity, err := s.getEntity(ctx, key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	addTTLHeader(w, now, entity.ExpiresAt)

	http.Redirect(w, r, entity.OriginalURL, http.StatusSeeOther)

}

func (s *service) decodeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	key := chi.URLParam(r, "key")
	now := time.Now()

	entity, err := s.getEntity(ctx, key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	addTTLHeader(w, now, entity.ExpiresAt)

	render.JSON(w, r, entity)

}

func (s *service) addEntryHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var obj create
	err := json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	now := time.Now()

	key, err := s.validateKey(ctx, now)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	entity := &entity{
		ID:          key,
		ShortURL:    s.host + "/" + key,
		OriginalURL: obj.OriginalURL,
		CreatedAt:   now,
		ExpiresAt:   now.Local().Add(time.Hour * time.Duration(24)),
	}

	err = s.saveEntity(ctx, key, entity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, entity)

}
