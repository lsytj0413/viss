// Package main is the entrance of project
package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gocolly/colly/v2"
	"github.com/lsytj0413/viss/pkg/reapers"
	"github.com/lsytj0413/viss/pkg/reapers/metadata"
)

type viss struct {
}

func (v *viss) NewCollector() *colly.Collector {
	return colly.NewCollector()
}

func (v *viss) SaveArticle(a reapers.Article) error {
	data, err := json.MarshalIndent(a, "  ", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", string(data))

	return nil
}

func main() {
	v := &viss{}
	metadataReaper := metadata.NewReaper()
	err := metadataReaper.Initialize(context.Background(), v)
	if err != nil {
		panic(err)
	}

	err = metadataReaper.Process(context.Background(), v)
	if err != nil {
		panic(err)
	}
}
