package blog

import (
	"io/ioutil"
	"net/http"
	"path"

	"github.com/chonglou/soy/env"
)

const (
	ext = ".md"
)

func root() string {
	return path.Join(env.ROOT, "blog")
}

func home(*http.Request) (env.H, error) {
	buf, err := ioutil.ReadFile(path.Join(root(), "readme.md"))
	if err != nil {
		return nil, err
	}
	return env.H{
		"title": "主页",
		"body":  string(buf),
	}, nil
}

func index(*http.Request) (env.H, error) {
	return env.H{}, nil
}

func show(*http.Request) (env.H, error) {
	return env.H{}, nil
}

func init() {
	env.GET("/", "home", home)
	env.GET("/blog", "blog/index", index)
	env.GET(`/blog/{name:[\w]*}`, "blog/show", show)
}
