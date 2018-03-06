package dict

import (
	"net/http"
	"path/filepath"

	"github.com/chonglou/soy/env"
	"github.com/kapmahc/stardict"
)

func dict(fn func(*stardict.Dictionary) error) error {
	dict, err := stardict.Open(filepath.Join(env.ROOT, "dict"))
	if err != nil {
		return err
	}
	for _, it := range dict {
		if err = fn(it); err != nil {
			return err
		}
	}
	return nil
}

func get(*http.Request) (env.H, error) {
	items := make(map[string]uint64)
	if err := dict(func(d *stardict.Dictionary) error {
		items[d.GetBookName()] = d.GetWordCount()
		return nil
	}); err != nil {
		return nil, err
	}
	return env.H{
		env.TITLE: "字典查询",
		"dict":    items,
	}, nil
}

func post(req *http.Request) (interface{}, error) {
	req.ParseForm()
	var items []env.H
	if err := dict(func(dt *stardict.Dictionary) error {
		senses := dt.Translate(req.Form.Get("keywords"))
		for _, seq := range senses {
			for _, p := range seq.Parts {
				items = append(items, env.H{"type": p.Type, "data": string(p.Data)})
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return items, nil
}

func init() {
	env.GET("/dict", "dict/search", get)
	env.POST("/dict", post)
}
