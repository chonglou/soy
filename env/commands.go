package env

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
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
				return nil
			},
		},
		cli.Command{
			Name:    "server",
			Aliases: []string{"s"},
			Usage:   "start the app server",
			Action: func(_ *cli.Context) error {
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
		it.URL = fmt.Sprintf("footer %d", i)
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
	return toml.NewEncoder(fd).Encode(Env{
		Port:           8080,
		Theme:          "bootstrap",
		Administrators: []string{"change-me@localhost"},
		ReCaptcha: ReCaptcha{
			SiteKey:   "reCAPTCHA-site-key",
			SecretKey: "reCAPTCHA-secret-key",
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
