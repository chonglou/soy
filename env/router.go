package env

import (
	"net/http"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	_render "github.com/unrolled/render"
)

var (
	router = mux.NewRouter()
	render *_render.Render
)

// H hash
type H map[string]interface{}

func abort(w http.ResponseWriter, e error) {
	render.HTML(w, http.StatusInternalServerError, "error", H{"reason": e.Error()})
}

// POST http post
func POST(pat string, hnd func(*http.Request) (H, error)) {
	router.HandleFunc(pat, func(wrt http.ResponseWriter, req *http.Request) {
		val, err := hnd(req)
		if err != nil {
			abort(wrt, err)
			return
		}
		render.JSON(wrt, http.StatusOK, val)
	}).Methods(http.MethodPost)
}

// GET http get
func GET(pat string, tpl string, hnd func(*http.Request) (H, error)) {
	router.HandleFunc(pat, func(wrt http.ResponseWriter, req *http.Request) {
		var env Env
		if _, err := toml.DecodeFile(Config(), &env); err != nil {
			abort(wrt, err)
			return
		}
		data, err := hnd(req)
		if err != nil {
			abort(wrt, err)
			return
		}
		data["env"] = env
		render.HTML(wrt, http.StatusOK, tpl, data)
	}).Methods(http.MethodGet)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		log.Infof("%s %s %s", r.Proto, r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
		log.Info(time.Now().Sub(now))
	})
}

func init() {
	router.Use(loggingMiddleware)
}
