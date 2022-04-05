package cache

import (
	"context"
	"time"
	"url-shortener/model"

	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
)

var (
	_addr     = "localhost:6379"
	_password = ""
	_database = 0
)

type Cache struct {
	db  *redis.Client
	log *log.Logger
}
type Cacheer interface {
	SetURL(ctx context.Context, data *model.URLs) error
	GetURL(ctx context.Context, shortenURL string) (string, error)
}

func NewCache(logger *log.Logger) Cacheer {
	redis := redis.NewClient(&redis.Options{
		Addr:     _addr,
		Password: _password,
		DB:       _database,
	})
	return &Cache{
		db:  redis,
		log: logger,
	}
}

func (c *Cache) SetURL(ctx context.Context, data *model.URLs) error {
	duration := time.Until(data.ExpireAt)
	return c.db.SetNX(data.ShortenURL, data.OriginalURL, duration).Err()
}
func (c *Cache) GetURL(ctx context.Context, shortenURL string) (string, error) {
	return c.db.Get(shortenURL).Result()
}
