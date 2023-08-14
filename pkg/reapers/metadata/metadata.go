// Package metadata contains the reaper for http://muratbuffalo.blogspot.com/
package metadata

import (
	"context"
	"fmt"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/gocolly/colly/v2"
	"github.com/k3a/html2text"

	"github.com/lsytj0413/ena/xslog"
	"github.com/lsytj0413/viss/pkg/reapers"
)

type metadataReaper struct {
	listCollector    *colly.Collector
	articleCollector *colly.Collector
}

// NewReaper will return the metadata reaper implement
func NewReaper() reapers.Reaper {
	return &metadataReaper{}
}

func (r *metadataReaper) Name() string {
	return "metadata"
}

func (r *metadataReaper) Initialize(ctx context.Context, c reapers.InitializeConfig) (err error) {
	r.listCollector = c.NewCollector()
	r.articleCollector = c.NewCollector()

	xslog.FromContext(ctx).InfoCtx(ctx, "Initialize success")
	return nil
}

func (r *metadataReaper) processArticleListItem(h *colly.HTMLElement, s reapers.Storage) {
	title := h.ChildText("h3.post-title")
	if title == "" {
		panic("empty title")
	}

	path := h.ChildAttr("h3.post-title a", "href")
	if path == "" {
		panic("empty path")
	}
	publishedAt := h.ChildAttr("div.post-header time.published", "datetime")
	t, err := time.Parse("2006-01-02T15:04:05-07:00", publishedAt)
	if err != nil {
		panic(err)
	}

	err = s.SaveArticle(reapers.Article{
		Title:       title,
		Path:        path,
		PublishedAt: t.UTC(),
	})
	if err != nil {
		panic(err)
	}

	err = r.articleCollector.Visit(path)
	if err != nil {
		panic(err)
	}
}

func (r *metadataReaper) processArticleListPager(h *colly.HTMLElement) {
	path := h.ChildAttr("a.blog-pager-older-link", "href")
	if path != "" {
		err := r.listCollector.Visit(path)
		if err != nil {
			panic(err)
		}
	}
}

func (r *metadataReaper) processArticle(h *colly.HTMLElement) {
	html, err := h.DOM.Html()
	if err != nil {
		panic(err)
	}

	text := html2text.HTML2Text(html)
	fmt.Printf("%v\n", text)

	converter := md.NewConverter("", true, nil)
	text, err = converter.ConvertString(html)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", text)
}

func (r *metadataReaper) Process(_ context.Context, s reapers.Storage) error {
	r.listCollector.OnHTML("div.post", func(h *colly.HTMLElement) {
		r.processArticleListItem(h, s)
	})

	r.listCollector.OnHTML("div.blog-pager", func(h *colly.HTMLElement) {
		r.processArticleListPager(h)
	})

	r.articleCollector.OnHTML("div.post div.post-body", func(h *colly.HTMLElement) {
		r.processArticle(h)
	})

	err := r.listCollector.Visit("http://muratbuffalo.blogspot.com/")
	if err != nil {
		return err
	}

	r.articleCollector.Wait()
	return nil
}
