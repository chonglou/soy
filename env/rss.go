package env

import (
	"io"
	"time"

	"github.com/gorilla/feeds"
)

var (
	rss []RSSHandler
)

// RSSHandler rss handler
type RSSHandler func() ([]*feeds.Item, error)

// RegisterRssHandler register rss handler
func RegisterRssHandler(args ...RSSHandler) {
	rss = append(rss, args...)
}

func rssRtom(wrt io.Writer, host, title, dest, author string) error {
	now := time.Now()
	feed := &feeds.Feed{
		Title:       title,
		Link:        &feeds.Link{Href: host},
		Description: dest,
		Author:      &feeds.Author{Name: author},
		Created:     now,
		Items:       make([]*feeds.Item, 0),
	}

	for _, hnd := range rss {
		items, err := hnd()
		if err != nil {
			return err
		}
		feed.Items = append(feed.Items, items...)
	}

	return feed.WriteAtom(wrt)
}
