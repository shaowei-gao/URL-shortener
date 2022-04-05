package service

import (
	"context"
	b64 "encoding/base64"
	"url-shortener/cache"
	"url-shortener/database"
	"url-shortener/model"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type Service struct {
	db    database.Databaseer
	cache cache.Cacheer
	log   *log.Logger
}
type Servicer interface {
	GenerateShortenURL(ctx context.Context, data *model.ReqGenerateShortenURL) (string, error)
	GetOriginalURL(ctx context.Context, data *model.ReqShortenURL) (string, error)
	Error(value ...interface{})
}

func NewService(db database.Databaseer, cache cache.Cacheer, logger *log.Logger) Servicer {
	srv := &Service{
		db:    db,
		cache: cache,
		log:   logger,
	}
	return srv
}

const sizeLimit = 10

func (srv *Service) GenerateShortenURL(ctx context.Context, req *model.ReqGenerateShortenURL) (string, error) {
	token, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	tokenLimit := token.String()[:sizeLimit]
	shortenURL := b64.StdEncoding.EncodeToString([]byte(tokenLimit))
	urlsObj := &model.URLs{
		ShortenURL:  shortenURL,
		OriginalURL: req.URL,
		ExpireAt:    req.ExpireAt,
	}
	srv.db.StoreURL(ctx, urlsObj)
	err = srv.cache.SetURL(ctx, urlsObj)
	if err != nil {
		return "", err
	}
	return shortenURL, nil
}

func (srv *Service) GetOriginalURL(ctx context.Context, shortenUrlObj *model.ReqShortenURL) (string, error) {
	shortenURL := shortenUrlObj.UrlID
	originalURL, err := srv.cache.GetURL(ctx, shortenURL)
	if err != nil {
		originalURL, err = srv.db.GetURL(ctx, shortenURL)
	}
	return originalURL, err

}

func (srv *Service) Error(value ...interface{}) {
	srv.log.Error(value...)
}
