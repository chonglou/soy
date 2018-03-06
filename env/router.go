package env

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/unrolled/render"
)

var (
	router  = mux.NewRouter()
	_render *render.Render
	_env    *Env
)

// H hash
type H map[string]interface{}

func abort(w http.ResponseWriter, e error) {
	log.Error(e)
	_render.Text(w, http.StatusInternalServerError, e.Error())
}

// POST http post
func POST(pat string, hnd func(*http.Request) (interface{}, error)) {
	router.HandleFunc(pat, func(wrt http.ResponseWriter, req *http.Request) {
		val, err := hnd(req)
		if err != nil {
			abort(wrt, err)
			return
		}
		_render.JSON(wrt, http.StatusOK, val)
	}).Methods(http.MethodPost)
}

// GET http get
func GET(pat string, tpl string, hnd func(*http.Request) (H, error)) {
	router.HandleFunc(pat, func(wrt http.ResponseWriter, req *http.Request) {
		data, err := hnd(req)
		if err != nil {
			abort(wrt, err)
			return
		}
		data["env"] = _env
		data[csrf.TemplateTag] = csrf.TemplateField(req)
		_render.HTML(wrt, http.StatusOK, tpl, data)
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

// Home home url
func Home(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s", scheme, r.Host)
}

func init() {
	router.Use(loggingMiddleware)
	router.HandleFunc("/robots.txt", func(wrt http.ResponseWriter, req *http.Request) {
		tpl := `User-agent: *
Disallow:
Crawl-delay: 10
Sitemap: %s/sitemap.xml.gz`
		_render.Text(wrt, http.StatusOK, fmt.Sprintf(tpl, Home(req)))
	}).Methods(http.MethodGet)

	router.HandleFunc("/sitemap.xml.gz", func(wrt http.ResponseWriter, req *http.Request) {
		if err := sitemapXMLGz(wrt, Home(req)); err != nil {
			abort(wrt, err)
		}
	}).Methods(http.MethodGet)

	router.HandleFunc("/rss.atom", func(wrt http.ResponseWriter, req *http.Request) {
		site := _env.Site
		if err := rssRtom(
			wrt,
			Home(req),
			site["title"],
			site["description"],
			site["author"],
		); err != nil {
			abort(wrt, err)
		}
	}).Methods(http.MethodGet)

	router.HandleFunc(`/google{id:[\w]+}.html`, func(wrt http.ResponseWriter, req *http.Request) {
		id := _env.Google.VerifyID
		vars := mux.Vars(req)
		code := http.StatusNotFound
		if vars["id"] == id {
			code = http.StatusOK
		}

		tpl := "google-site-verification: google%s.html"
		_render.Text(wrt, code, fmt.Sprintf(tpl, id))
	}).Methods(http.MethodGet)

}
