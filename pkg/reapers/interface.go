// Package reapers contains all in-tree reaper to collect article.
package reapers

import (
	"context"
	"time"

	"github.com/gocolly/colly/v2"
)

// Article is the represent of acticle metadata
type Article struct {
	// Title is the brief introduction to the article
	Title string

	// Path is the absolute url path for article
	Path string

	// PublishedAt is the timestamp when this article is created.
	PublishedAt time.Time
}

// InitializeConfig provide the config for initialize reaper
type InitializeConfig interface {
	NewCollector() *colly.Collector
}

// Storage will persistent the article
type Storage interface {
	// SaveArticle will persistent the article content
	SaveArticle(a Article) error
}

// Reaper is used to collect article from each domain.
type Reaper interface {
	// Name return the reaper unique name
	Name() string

	// Initialize is called when reaper should initialize before process.
	Initialize(ctx context.Context, c InitializeConfig) error

	// Process will start to collect article
	Process(ctx context.Context, s Storage) error
}
