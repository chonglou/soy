dist=dist
pkg=github.com/chonglou/soy/env
theme=moon

VERSION=`git rev-parse --short HEAD`
BUILD_TIME=`date -R`
AUTHOR_NAME=`git config --get user.name`
AUTHOR_EMAIL=`git config --get user.email`
COPYRIGHT=`head -n 1 LICENSE`
USAGE=`sed -n '3p' README.md`

build:
	go build -ldflags "-s -w -X ${pkg}.Version=${VERSION} -X '${pkg}.BuildTime=${BUILD_TIME}' -X '${pkg}.AuthorName=${AUTHOR_NAME}' -X ${pkg}.AuthorEmail=${AUTHOR_EMAIL} -X '${pkg}.Copyright=${COPYRIGHT}' -X '${pkg}.Usage=${USAGE}'" -o ${dist}/soy main.go
	-cp -r themes LICENSE README.md $(dist)/
	cd $(dist) && mkdir -p tmp public
	cd $(dist) && tar cfJ ../$(dist).tar.xz *

clean:
	-rm -r $(dist) $(dist).tar.xz

init:
	govendor sync
	npm install
