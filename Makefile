export PATH:=$(PWD):$(PATH)

build: gitsocket

all: clean build test

clean:
	rm -f gitsocket

gitsocket:
	@echo
	@echo "Building gitsocket"
	@echo "------------------"
	@echo
	go build -o gitsocket

test-repo:
	@echo
	@echo "Creating test reopsitory"
	@echo "------------------------"
	@echo
	rm -Rf _test
	mkdir -p _test/remote
	cd _test/remote && git init --bare
	mkdir -p _test/local
	cd _test/local && git init
	cd _test/local && git config user.name  "Test User"
	cd _test/local && git config user.email "user@localhost"
	echo "# README" > _test/local/README.md
	cd _test/local && git add README.md && git commit -m "Initial commit"
	cd _test/local && git checkout -b other_branch
	echo "# OTHER" > _test/local/OTHER.md
	cd _test/local && git add OTHER.md && git commit -m "Unwanted commit"
	cd _test/local && git remote add origin ../remote
	cd _test/local && git push -u origin master

test-stop-gitsocket:
	kill `cat "test.pid"`

test-test-result:
	sleep 1
	cd _test/local && ls
	cd _test/local && if [ -f "OTHER.md" ]; then make test-stop-gitsocket; exit 1; fi
	cd _test/local && git status | head -1 > status.txt
	cd _test/local && if ! grep "HEAD detached at origin/master" status.txt; then make test-stop-gitsocket; exit 1 ; fi

test: test-default test-socket test-port test-ip-port

test-default:
	@echo
	@echo "Functional Tests (default)"
	@echo "--------------------------"
	@echo
	@make test-repo
	gitsocket server --daemon --gitrepo "_test/local" --pidfile "test.pid"
	gitsocket client
	@make test-test-result
	@make test-stop-gitsocket
	@echo
	@echo "Test Passed"
	@echo

test-socket:
	@echo
	@echo "Functional Tests (socket)"
	@echo "-------------------------"
	@echo
	@make test-repo
	gitsocket server --daemon --gitrepo "_test/local" --listen "_test/test.sock" --pidfile "test.pid"
	gitsocket client --conn "_test/test.sock"
	@make test-test-result
	@make test-stop-gitsocket
	@echo
	@echo "Test Passed"
	@echo

test-port:
	@echo
	@echo "Functional Tests (port)"
	@echo "-----------------------"
	@echo
	@make test-repo
	gitsocket server --daemon --gitrepo "_test/local" --listen 9301 --pidfile "test.pid"
	gitsocket client --conn 9301
	@make test-test-result
	@make test-stop-gitsocket
	@echo
	@echo "Test Passed"
	@echo

test-ip-port:
	@echo
	@echo "Functional Tests (ip:port)"
	@echo "--------------------------"
	@echo
	@make test-repo
	gitsocket server --daemon --gitrepo "_test/local" --listen 9301 --pidfile "test.pid"
	gitsocket client --conn 127.0.0.1:9301
	@make test-test-result
	@make test-stop-gitsocket
	@echo
	@echo "Test Passed"
	@echo

.PHONY: all build clean test
.PHONY: test-repo test-start-gitsocket test-stop-gitsocket test-test-result
.PHONY: test-default test-socket test-port test-ip-port
