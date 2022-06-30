# serpapi-golang release package
#
version=`grep VERSION serpapi.go | cut -d'"' -f2`

all: test

check:
	go vet .
	go fmt .

test:
	go test -v .

# check that everything is pushed
package:
	git status | grep "Nothing"

oobt:
	mkdir -p /tmp/oobt
	cp demo/demo.go /tmp/oobt
	cd /tmp/oobt ; \
		go get -u github.com/serpapi/serpapi-golang ; \
		go run demo.go

version:
	@echo "current version: ${version}"

release: oobt version
	git tag -a ${version}
	git push origin ${version}
	@echo "create release: ${version}"
