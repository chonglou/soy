package blog

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/chonglou/soy/env"
	"github.com/gorilla/mux"
)

const (
	ext = ".md"
)

func root() string {
	return filepath.Join(env.ROOT, "blog")
}

func home(*http.Request) (env.H, error) {
	return readMD("主页")
}

func index(*http.Request) (env.H, error) {
	rt := root()
	items := make(map[string]string)
	if err := filepath.Walk(rt, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		// name := info.Name()
		items["/blog"+path[len(rt):]] = path[len(rt)+1 : len(path)-len(ext)]
		return nil
	}); err != nil {
		return nil, err
	}
	return env.H{
		env.TITLE: time.Now().Format(time.RFC822),
		"links":   items,
	}, nil
}

func show(req *http.Request) (env.H, error) {
	vars := mux.Vars(req)
	return readMD(vars["name"])
}

func readMD(fn string) (env.H, error) {
	buf, err := ioutil.ReadFile(filepath.Join(root(), fn+ext))
	if err != nil {
		return nil, err
	}
	return env.H{
		env.TITLE: fn[:len(fn)],
		"body":    string(buf),
	}, nil
}

func init() {
	env.GET("/", "blog/show", home)
	env.GET("/blog", "blog/index", index)
	env.GET(`/blog/{name:[\w\W]*}.md`, "blog/show", show)
}
