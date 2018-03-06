package env

import (
	"fmt"
	"io"
	"time"

	"github.com/gorilla/feeds"
)

var (
	rss []RSSHandler
)

// RSSHandler rss handler
type RSSHandler func(l string) ([]*feeds.Item, error)

// RegisterRssHandler register rss handler
func RegisterRssHandler(args ...RSSHandler) {
	rss = append(rss, args...)
}

func rssRtom(wrt io.Writer, host, lang, title, dest string, author *feeds.Author) error {
	now := time.Now()
	feed := &feeds.Feed{
		Title:       title,
		Link:        &feeds.Link{Href: fmt.Sprintf("%s/?locale=%s", host, lang)},
		Description: dest,
		Author:      author,
		Created:     now,
		Items:       make([]*feeds.Item, 0),
	}

	for _, hnd := range rss {
		items, err := hnd(lang)
		if err != nil {
			return err
		}
		feed.Items = append(feed.Items, items...)
	}

	return feed.WriteAtom(wrt)
}
