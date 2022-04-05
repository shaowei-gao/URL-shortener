package service

import (
	"context"
	"errors"
	"testing"
	"time"
	"url-shortener/cache"
	"url-shortener/database"
	"url-shortener/log"
	"url-shortener/model"

	"github.com/stretchr/testify/assert"
)

type Result struct {
	URL string
	Err error
}

var (
	ctx    = context.Background()
	db     = database.NewSqlite("../database")
	logger = log.NewLogrus(ctx)
	c      = cache.NewCache(logger)
	srv    = NewService(db, c, logger)
)
var (
	errNotMatch = errors.New("the value is not expected")
)

func TestGenerateShortenURLAndGetOriginalURL(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		desc        string
		originalURL string
		shortenURL  string
		expireAt    time.Time
		wait        time.Duration
		result      Result
	}{
		{
			desc:        "General",
			originalURL: "url1",
			expireAt:    now.Add(1 * time.Second),
			wait:        time.Duration(0),
			result: Result{
				URL: "url1",
				Err: nil,
			},
		},
		{
			desc:        "Time expired",
			originalURL: "url2",
			expireAt:    now.Add(1 * time.Second),
			wait:        time.Duration(1 * time.Second),
			result: Result{
				URL: "",
				Err: database.ErrTimeExpired,
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			request := &model.ReqGenerateShortenURL{
				URL:      tC.originalURL,
				ExpireAt: tC.expireAt,
			}
			shortenURL, err := srv.GenerateShortenURL(ctx, request)
			assert.Nil(t, err, err)
			reqShortenURL := &model.ReqShortenURL{UrlID: shortenURL}
			time.Sleep(tC.wait)
			originalURL, err := srv.GetOriginalURL(ctx, reqShortenURL)
			assert.Equal(t, tC.result.URL, originalURL, errNotMatch)
			assert.Equal(t, tC.result.Err, err, errNotMatch)
		})
	}
}
