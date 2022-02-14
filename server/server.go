package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func Start(port, host string, keyLength int) error {
	cache, err := newCache()
	if err != nil {
		log.Fatal(err)
	}

	persistence, err := newPersistence("url-shortener")
	if err != nil {
		log.Fatal(err)
	}

	serve := &service{
		port:      port,
		host:      host,
		keyLength: keyLength,
		cache:     *cache,
		persist:   *persistence,
	}

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Post("/create", serve.addEntryHandler)
	r.Get("/{key}", serve.redirectHandler)
	r.Get("/{key}/details", serve.decodeHandler)

	srv := &http.Server{Addr: ":" + serve.port, Handler: r}

	return srv.ListenAndServe()
}
