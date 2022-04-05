package route

import (
	"context"
	"errors"
	"net/http"
	"time"
	"url-shortener/model"
	"url-shortener/service"

	"github.com/gin-gonic/gin"
)

const domain = "http://localhost/"

var errEarierTime = errors.New("expireAt time is earlier than current time")

func SetupRoute(ctx context.Context, r *gin.Engine, srv service.Servicer) {
	r.GET("/:url_id", redirectOriginalURL(ctx, srv))
	api := r.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.POST("/urls", generateShortenURL(ctx, srv))
		}
	}
}

func generateShortenURL(ctx context.Context, srv service.Servicer) gin.HandlerFunc {
	var bind *model.ReqGenerateShortenURL
	return func(c *gin.Context) {
		if err := c.ShouldBindJSON(&bind); err != nil {
			srv.Error(err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		if bind.ExpireAt.Before(time.Now()) {
			srv.Error(errEarierTime)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		urlID, err := srv.GenerateShortenURL(ctx, bind)
		if err != nil {
			srv.Error(err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"id":       urlID,
			"shortUrl": domain + urlID,
		})
	}
}

func redirectOriginalURL(ctx context.Context, srv service.Servicer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var bind *model.ReqShortenURL
		if err := c.ShouldBindUri(&bind); err != nil {
			srv.Error(err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		originURL, err := srv.GetOriginalURL(ctx, bind)
		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.Redirect(http.StatusFound, originURL)
	}
}
