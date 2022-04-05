package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"url-shortener/cache"
	"url-shortener/database"
	"url-shortener/log"
	"url-shortener/model"
	"url-shortener/route"
	"url-shortener/service"

	"github.com/gin-gonic/gin"
)

const (
	_srvAddr = ":80"
)

func main() {
	ctx := context.Background()
	logger := log.NewLogrus(ctx)
	cache := cache.NewCache(logger)
	// db := database.NewSqlite(_dbPath)
	db := database.NewPostgres()

	db.Migrate(&model.URLs{})
	srv := service.NewService(db, cache, logger)
	ser := gin.Default()
	route.SetupRoute(ctx, ser, srv)
	server := &http.Server{
		Addr:    _srvAddr,
		Handler: ser,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("listen: %s\n", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server Shutdown: ", err)
	}
	logger.Println("Server exiting")
}
