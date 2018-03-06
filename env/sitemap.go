package env

import (
	"compress/gzip"
	"io"

	"github.com/ikeikeikeike/go-sitemap-generator/stm"
)

var (
	sitemap []SitemapHandler
)

// SitemapHandler sitemap handler
type SitemapHandler func() ([]stm.URL, error)

// RegisterSitemapHandler register sitemap handler
func RegisterSitemapHandler(args ...SitemapHandler) {
	sitemap = append(sitemap, args...)
}

func sitemapXMLGz(w io.Writer, h string) error {
	sm := stm.NewSitemap()
	sm.Create()
	sm.SetDefaultHost(h)
	for _, hnd := range sitemap {
		items, err := hnd()
		if err != nil {
			return err
		}
		for _, it := range items {
			sm.Add(it)
		}
	}
	buf := sm.XMLContent()

	wrt := gzip.NewWriter(w)
	defer wrt.Close()
	wrt.Write(buf)
	return nil
}
