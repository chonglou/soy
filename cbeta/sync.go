package cbeta

import (
	"os"
	"path/filepath"

	"github.com/chonglou/soy/env"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func sync(c *cli.Context) error {
	return nil
}

func init() {
	env.Command(cli.Command{
		Name:    "epub",
		Aliases: []string{"e"},
		Usage:   "sync epub books",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "dir, d",
				Usage: "book's directory",
			},
		},
		Action: func(c *cli.Context) error {
			root := c.String("dir")
			if root == "" {
				cli.ShowSubcommandHelp(c)
				return nil
			}
			if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if info.IsDir() {
					return nil
				}
				log.Infof("find book %s", path[len(root)+1:])
				return nil
			}); err != nil {
				return err
			}
			return nil
		},
	})
}
