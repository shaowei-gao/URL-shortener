package model

import "time"

type ReqGenerateShortenURL struct {
	URL      string    `json:"url" binding:"required"`
	ExpireAt time.Time `json:"expireAt" binding:"required"`
}

type ReqShortenURL struct {
	UrlID string `uri:"url_id" binding:"required"`
}
type URLs struct {
	ShortenURL  string `gorm:"primaryID"`
	OriginalURL string
	ExpireAt    time.Time
}
