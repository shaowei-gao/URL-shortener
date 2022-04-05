package route

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"url-shortener/cache"
	"url-shortener/database"
	"url-shortener/log"
	"url-shortener/service"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	db     = database.NewPostgres()
	ctx    = context.Background()
	logger = log.NewLogrus(ctx)
	c      = cache.NewCache(logger)
	srv    = service.NewService(db, c, logger)
	r      = gin.Default()
)

var (
	errNotMatch = errors.New("not match")
)

func TestGenerateShortUrlApi(t *testing.T) {
	const (
		apiURL      = `/api/v1/urls`
		dataFormat  = `{"url": "%s","expireAt": "%s"}`
		method      = `POST`
		originalURL = `url`
	)
	var (
		w          *httptest.ResponseRecorder
		req        *http.Request
		statusCode int
		data       string
		err        error
	)
	SetupRoute(ctx, r, srv)

	testCases := []struct {
		desc       string
		url        string
		expireAt   string
		statusCode int
	}{
		{
			desc:       "Normal request",
			url:        originalURL,
			expireAt:   time.Now().Add(2 * time.Second).Format(time.RFC3339),
			statusCode: http.StatusOK,
		},
		{
			desc:       "Empty url",
			url:        "",
			expireAt:   time.Now().Add(2 * time.Second).Format(time.RFC3339),
			statusCode: http.StatusBadRequest,
		},
		{
			desc:       "Empty expireAt",
			url:        originalURL,
			expireAt:   "",
			statusCode: http.StatusBadRequest,
		},
		{
			desc:       "The expireAt time is earlier than current time",
			url:        originalURL,
			expireAt:   time.Now().Add(-2 * time.Second).Format(time.RFC3339),
			statusCode: http.StatusBadRequest,
		},
		{
			desc:       "Fatal format of expireAt",
			url:        originalURL,
			expireAt:   "2022-04-0T09:20:41Z",
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			data = fmt.Sprintf(dataFormat, tC.url, tC.expireAt)
			req, err = http.NewRequest(method, apiURL, strings.NewReader(data))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Add("Content-Type", "application/json")
			w = httptest.NewRecorder()
			r.ServeHTTP(w, req)
			statusCode = w.Result().StatusCode
			assert.Equal(t, tC.statusCode, statusCode, errNotMatch)
		})
	}
}
