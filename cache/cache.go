package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/allegro/bigcache"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

var once sync.Once
var lc *legisCache

func GetLegisCache() *legisCache {
	once.Do(func() {
		c, err := bigcache.NewBigCache(bigcache.DefaultConfig(12 * time.Hour))
		if err != nil {
			panic("WE NEED DA CACHE" + err.Error())
		}
		lc = &legisCache{c: c}
	})
	return lc
}

type legisCache struct {
	c *bigcache.BigCache
}

func (l *legisCache) Delete(key string) (err error) {
	if err = l.c.Delete(key); err != nil {
		err = fmt.Errorf("legisCache delete error: %w", err)
	}
	return
}

func (l *legisCache) AddToCache(key string, object interface{}) {
	objectJson, err := json.Marshal(object)
	if err != nil {
		log.WithError(err).WithField("legisCacheKey", key).Error("legisCache marshal error")
		return
	}
	if err = l.c.Set(key, objectJson); err != nil {
		log.WithError(err).WithField("legisCacheKey", key).Error("legisCache set error")
	}
}

func (l *legisCache) GetFromCache(key string, objectToUnmarshal interface{}) error {
	retrievedItem, err := l.c.Get(key)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(retrievedItem, objectToUnmarshal); err != nil {
		return fmt.Errorf("legisCache unmmarshal error: %w")
	}
	log.WithField("legisCacheKey", key).Info("hit legisCache, returning unmarshaled object")
	return nil
}

func (l *legisCache) NotFound(err error) bool {
	return errors.Is(err, bigcache.ErrEntryNotFound)
}
