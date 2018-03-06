package cbeta

import (
	"net/http"
	"path/filepath"

	"github.com/chonglou/soy/env"
	"github.com/unrolled/render"
)

func root() string {
	return filepath.Join(env.ROOT, "books")
}

func index(*http.Request) (env.H, error) {
	return env.H{}, nil
}

func show(wrt http.ResponseWriter, req *http.Request, arg *env.Env, rdr *render.Render) {

}
