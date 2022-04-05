package cache

import (
	"context"
	"testing"
	"time"
	"url-shortener/log"
	"url-shortener/model"

	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
)

var (
	_ctx    = context.Background()
	_logger = log.NewLogrus(_ctx)
	_cache  = NewCache(_logger)
)

func TestSetAndGet(t *testing.T) {
	testCases := []struct {
		desc string
		data *model.URLs
	}{
		{
			desc: "",
			data: &model.URLs{
				ShortenURL:  "s",
				OriginalURL: "o",
				ExpireAt:    time.Now().Add(1 * time.Second),
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			err := _cache.SetURL(_ctx, tC.data)
			if err != nil {
				t.Fatal(err)
			}

			assert.Nil(t, err, err)
			s, err := _cache.GetURL(_ctx, tC.data.ShortenURL)
			assert.Nil(t, err, err)
			assert.Equal(t, tC.data.OriginalURL, s, "The original url is not matched ")
			time.Sleep(time.Until(tC.data.ExpireAt))
			_, err = _cache.GetURL(_ctx, tC.data.ShortenURL)
			assert.Equal(t, redis.Nil, err, "The key is not expired!")
		})
	}
}
