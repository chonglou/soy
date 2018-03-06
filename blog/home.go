package blog

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/chonglou/soy/env"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

const (
	ext = ".md"
)

func root() string {
	return filepath.Join(env.ROOT, "blog")
}

func home(*http.Request) (env.H, error) {
	return readMD("主页.md")
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
		if filepath.Ext(info.Name()) != ext {
			return nil
		}
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

func show(wrt http.ResponseWriter, req *http.Request, arg *env.Env, rdr *render.Render) {
	vars := mux.Vars(req)
	name := vars["name"]
	if filepath.Ext(name) == ext {
		data, err := readMD(name)
		if err != nil {
			env.Abort(wrt, err)
			return
		}
		data["env"] = arg
		rdr.HTML(wrt, http.StatusOK, "blog/show", data)
		return
	}

	buf, err := ioutil.ReadFile(filepath.Join(root(), vars["name"]))
	if err != nil {
		rdr.Text(wrt, http.StatusInternalServerError, err.Error())
		return
	}
	typ := http.DetectContentType(buf)
	wrt.Header().Set("Content-type", typ)
	rdr.Data(wrt, http.StatusOK, buf)
}

func readMD(fn string) (env.H, error) {
	buf, err := ioutil.ReadFile(filepath.Join(root(), fn))
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
	// env.GET(`/blog/{name:[\w\W]*}.md`, "blog/show", showMD)
	env.HANDLE(http.MethodGet, `/blog/{name:[\w\W]*}`, show)
}
