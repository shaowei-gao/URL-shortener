package database

import (
	"context"
	"errors"
	"time"
	"url-shortener/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	ErrNotFound    = errors.New("query was not found")
	ErrTimeExpired = errors.New("the time is expired")
)

type DB struct {
	orm *gorm.DB
}
type Databaseer interface {
	GetURL(ctx context.Context, shortenURL string) (string, error)
	StoreURL(ctx context.Context, data *model.URLs) error
	Migrate(dst ...interface{}) error
}

func NewSqlite(path string) Databaseer {
	file := path + "/sqlite.db"
	orm, err := gorm.Open(sqlite.Open(file))
	if err != nil {
		panic(err)
	}
	return &DB{
		orm: orm,
	}
}

func (db *DB) GetURL(ctx context.Context, shortenURL string) (string, error) {
	var urls model.URLs
	result := db.orm.Where("shorten_url = ?", shortenURL).Find(&urls)

	if result.Error != nil {
		return "", result.Error
	}
	if result.RowsAffected == 0 {
		return "", ErrNotFound
	}
	if urls.ExpireAt.Before(time.Now()) {
		return "", ErrTimeExpired
	}
	return urls.OriginalURL, nil
}

func (db *DB) StoreURL(ctx context.Context, data *model.URLs) error {
	return db.orm.Create(data).Error
}
func (db *DB) Migrate(dst ...interface{}) error {
	return db.orm.AutoMigrate(dst...)
}
