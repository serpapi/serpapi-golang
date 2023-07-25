# serpapi-golang release package
#
version=`grep VERSION serpapi.go | cut -d'"' -f2`

.PHONY: test oobt

all: version lint test doc ready

# lint source code using go tools
lint:
	go vet .
	go fmt .

# run integration test suite
test:
	go test -v ./test

# create documentation
doc:
	go doc

# check that everything is pushed
ready:
	git status | grep "Nothing"

# out of box testing 
#  validate the pre-released library
oobt:
	mkdir -p /tmp/oobt
	cp oobt/demo.go /tmp/oobt
	cd /tmp/oobt ; \
		go mod init serpapi.com/golang/oobt ; \
		go get -u github.com/serpapi/serpapi-golang ; \
		go run demo.go

# show current version for golang and library
version:
  @echo "golang: " `go --version`
	@echo "current version: ${version}"

# display the current release information
release: oobt version
	git tag -a ${version}
	git push origin ${version}
	@echo "create release: ${version}"
