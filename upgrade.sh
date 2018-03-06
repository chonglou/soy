#!/bin/sh
go get -u github.com/kardianos/govendor
rm -r vendor
govendor init
govendor fetch github.com/gorilla/feeds
govendor fetch github.com/ikeikeikeike/go-sitemap-generator/stm
govendor fetch gopkg.in/go-playground/validator.v9
govendor fetch github.com/go-playground/form
govendor fetch github.com/gorilla/mux
govendor fetch github.com/gorilla/csrf
govendor fetch golang.org/x/crypto/bcrypt
govendor fetch golang.org/x/text/language
govendor fetch github.com/google/uuid
govendor fetch github.com/urfave/cli
govendor fetch github.com/BurntSushi/toml
govendor fetch github.com/sirupsen/logrus
govendor fetch github.com/sirupsen/logrus/hooks/syslog
govendor fetch gopkg.in/gomail.v2
