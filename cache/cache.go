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
var lc *LegisCache

func GetLegisCache() *LegisCache {
	once.Do(func() {
		c, err := bigcache.NewBigCache(bigcache.DefaultConfig(12 * time.Hour))
		if err != nil {
			log.WithError(err).Panic("couldn't create the cache: exiting")
		}
		lc = &LegisCache{c: c}
	})
	return lc
}

type LegisCache struct {
	c *bigcache.BigCache
}

func (l *LegisCache) Delete(key string) (err error) {
	if err = l.c.Delete(key); err != nil {
		err = fmt.Errorf("LegisCache delete error: %w", err)
	}
	return
}

func (l *LegisCache) AddToCache(key string, object interface{}) {
	objectJson, err := json.Marshal(object)
	if err != nil {
		log.WithError(err).WithField("legisCacheKey", key).Error("LegisCache marshal error")
		return
	}
	if err = l.c.Set(key, objectJson); err != nil {
		log.WithError(err).WithField("legisCacheKey", key).Error("LegisCache set error")
	}
	log.WithField("legisCacheKey", key).Info("added to LegisCache")
}

func (l *LegisCache) GetFromCache(key string, objectToUnmarshal interface{}) error {
	retrievedItem, err := l.c.Get(key)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(retrievedItem, objectToUnmarshal); err != nil {
		return fmt.Errorf("LegisCache unmmarshal error: %w")
	}
	//log.WithField("legisCacheKey", key).Info("hit LegisCache, returning unmarshaled object")
	return nil
}

func (l *LegisCache) NotFound(err error) bool {
	return errors.Is(err, bigcache.ErrEntryNotFound)
}
