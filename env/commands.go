package env

import (
	"context"
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/unrolled/render"
	"github.com/urfave/cli"
)

var commands []cli.Command

// Command register commands
func Command(args ...cli.Command) {
	commands = append(commands, args...)
}

func init() {
	Command(
		cli.Command{
			Name:    "routes",
			Aliases: []string{"rt"},
			Usage:   "print out all defined routes",
			Action: func(_ *cli.Context) error {
				tpl := "%-16s %s\n"
				fmt.Printf(tpl, "METHODS", "PATH")
				return router.Walk(func(route *mux.Route, _ *mux.Router, _ []*mux.Route) error {
					pat, err := route.GetPathTemplate()
					if err != nil {
						return err
					}
					mtd, err := route.GetMethods()
					if err != nil {
						return err
					}
					fmt.Printf(tpl, strings.Join(mtd, ","), pat)
					return nil
				})
			},
		},
		cli.Command{
			Name:    "server",
			Aliases: []string{"s"},
			Usage:   "start the app server",
			Action: func(_ *cli.Context) error {
				var env Env
				if _, err := toml.DecodeFile(Config(), &env); err != nil {
					return err
				}
				_env = &env

				// static files
				for k, v := range map[string]string{
					"/3rd/":    "node_modules",
					"/assets/": path.Join("themes", _env.Theme, "assets"),
					"/global/": path.Join("themes", "global"),
				} {
					router.PathPrefix(k).Handler(http.StripPrefix(k, http.FileServer(http.Dir(v))))
				}

				// open render
				_render = render.New(render.Options{
					Directory:  path.Join("themes", env.Theme, "views"),
					Layout:     "application/index",
					Extensions: []string{".html"},
					Funcs:      []template.FuncMap{},
				})

				addr := fmt.Sprintf(":%d", env.Port)
				secrets, err := base64.StdEncoding.DecodeString(env.Secrets)
				if err != nil {
					return err
				}

				log.Infof("listen %s", addr)
				CSRF := csrf.Protect(
					secrets,
					csrf.Secure(false),
					csrf.RequestHeader("Authenticity-Token"),
					csrf.FieldName("authenticity_token"),
				)
				// graceful shutdown
				srv := &http.Server{
					Addr:    addr,
					Handler: CSRF(router),
				}

				go func() {
					if err := srv.ListenAndServe(); err != nil {
						log.Error(err)
					}
				}()

				quit := make(chan os.Signal, 1)
				signal.Notify(quit, os.Interrupt)
				<-quit

				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				if err := srv.Shutdown(ctx); err != nil {
					return err
				}
				log.Info("shutting down")
				return nil
			},
		},
		cli.Command{
			Name:    "init",
			Aliases: []string{"i"},
			Usage:   "init config files",
			Action: func(_ *cli.Context) error {
				if err := os.MkdirAll(ROOT, 0755); err != nil {
					return err
				}
				// config.toml
				if err := generateConfigToml(); err != nil {
					return err
				}
				return nil
			},
		},
	)
}

func generateConfigToml() error {
	fn := Config()
	log.Infof("generate file %s\n", fn)
	fd, err := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	defer fd.Close()

	var header []Link
	for i := 1; i <= 6; i++ {
		var it Link
		it.Title = fmt.Sprintf("header %d", i)
		it.URL = fmt.Sprintf("/header-%d", i)
		for j := 1; j <= 3; j++ {
			it.Children = append(
				it.Children,
				Link{
					Title: fmt.Sprintf("header %d %d", i, j),
					URL:   fmt.Sprintf("/header-%d-%d", i, j),
				},
			)
		}
		header = append(header, it)
	}
	var footer []Link
	for i := 1; i <= 5; i++ {
		footer = append(
			footer,
			Link{
				Title: fmt.Sprintf("footer %d", i),
				URL:   fmt.Sprintf("/footer-%d", i),
			},
		)
	}

	secret, err := RandomBytes(32)
	if err != nil {
		return err
	}

	return toml.NewEncoder(fd).Encode(Env{
		Port:           8080,
		Theme:          "bootstrap",
		Secrets:        base64.StdEncoding.EncodeToString(secret),
		Administrators: []string{"change-me@localhost"},
		Google: Google{
			VerifyID: "google-site-verify-id",
			ReCaptcha: ReCaptcha{
				SiteKey:   "reCAPTCHA-site-key",
				SecretKey: "reCAPTCHA-secret-key",
			},
		},
		SMTP: SMTP{
			Host:     "smtp.gmail.com",
			Port:     465,
			User:     "who-am-i@gmail.com",
			Password: "change-me",
		},
		Site: map[string]string{
			"title":       "site title",
			"subhead":     "site subhead",
			"description": "site description",
		},
		Header: header,
		Footer: footer,
	})
}
