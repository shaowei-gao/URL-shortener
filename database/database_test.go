package database

import (
	"context"
	b64 "encoding/base64"
	"testing"
	"time"
	"url-shortener/model"
)

var (
	_ctx = context.Background()
	_db  = NewSqlite(".")
)

func TestInsertURLs(t *testing.T) {
	expireAt, _ := time.Parse(time.RFC3339, "2021-02-08T09:20:41Z")
	testCases := []struct {
		desc string
		data *model.URLs
	}{
		{
			desc: "",
			data: &model.URLs{
				OriginalURL: "url",
				ExpireAt:    expireAt,
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			tC.data.ShortenURL = b64.URLEncoding.EncodeToString([]byte(tC.data.OriginalURL))
			_db.StoreURL(_ctx, tC.data)
		})
	}
}
