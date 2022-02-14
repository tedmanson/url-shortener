package server

import (
	"context"
	"fmt"
	"github.com/dchest/uniuri"
	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func isExpired(now, check time.Time) bool {
	return now.After(check)
}

func (s *service) saveEntity(ctx context.Context, key string, e *entity) error {
	err := s.persist.set(ctx, key, e)
	if err != nil {
		log.WithError(err).Error("could not save to cache")
	}

	err = s.cache.set(key, e)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) getEntity(ctx context.Context, key string) (*entity, error) {
	obj, err := s.cache.get(key)
	if err != nil && err != redis.ErrNil {
		return nil, err
	}

	if obj == nil {
		obj, err = s.persist.get(ctx, key)
		if err != nil {
			return nil, err
		}

		now := time.Now()
		if !isExpired(now, obj.ExpiresAt) {
			err = s.cache.set(key, obj)
			if err != nil {
				return nil, err
			}
		}
	}
	return obj, err
}

func (s *service) validateKey(ctx context.Context, now time.Time) (string, error) {
	var key string
	for {
		key = uniuri.NewLen(s.keyLength)

		e, err := s.persist.exists(ctx, key)
		if err != nil {
			log.WithError(err).Error("could not validate from cache layer")
		}

		if err == nil && !e {
			break
		}

		obj, err := s.persist.get(ctx, key)
		if err != nil {
			return "", err
		}

		if isExpired(now, obj.ExpiresAt) {
			break
		}
	}

	return key, nil
}

func addTTLHeader(w http.ResponseWriter, now, expiresAt time.Time) {
	var ttl int
	if expiresAt.After(now.Add(1 * time.Hour)) {
		ttl = int((1 * time.Hour).Seconds())
	} else {
		ttl = int(expiresAt.Sub(now).Seconds())
	}
	w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", ttl))
}
