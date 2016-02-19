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

test-server-test-result:
	sleep 1
	cd _test/local && ls
	cd _test/local && if [ -f "OTHER.md" ]; then make test-stop-gitsocket; exit 1; fi
	cd _test/local && git status | head -1 > status.txt
	cd _test/local && if ! grep "HEAD detached at origin/master" status.txt; then make test-stop-gitsocket; exit 1 ; fi

test: test-server test-once test-setup

test-server: test-server-default test-server-socket test-server-port test-server-ip-port

test-server-default:
	@echo
	@echo "gitsocket server (default)"
	@echo "--------------------------"
	@echo
	@make test-repo
	gitsocket server --daemon --gitrepo "_test/local" --pidfile "test.pid"
	gitsocket client
	@make test-server-test-result
	@make test-stop-gitsocket
	@echo
	@echo "Test Passed"
	@echo

test-server-socket:
	@echo
	@echo "gitsocket server (socket)"
	@echo "-------------------------"
	@echo
	@make test-repo
	gitsocket server --daemon --gitrepo "_test/local" --listen "_test/test.sock" --pidfile "test.pid"
	gitsocket client --conn "_test/test.sock"
	@make test-server-test-result
	@make test-stop-gitsocket
	@echo
	@echo "Test Passed"
	@echo

test-server-port:
	@echo
	@echo "gitsocket server (port)"
	@echo "-----------------------"
	@echo
	@make test-repo
	gitsocket server --daemon --gitrepo "_test/local" --listen 9301 --pidfile "test.pid"
	gitsocket client --conn 9301
	@make test-server-test-result
	@make test-stop-gitsocket
	@echo
	@echo "Test Passed"
	@echo

test-server-ip-port:
	@echo
	@echo "gitsocket server (ip:port)"
	@echo "--------------------------"
	@echo
	@make test-repo
	gitsocket server --daemon --gitrepo "_test/local" --listen 9301 --pidfile "test.pid"
	gitsocket client --conn 127.0.0.1:9301
	@make test-server-test-result
	@make test-stop-gitsocket
	@echo
	@echo "Test Passed"
	@echo

test-once:
	@echo
	@echo "gitsocket once"
	@echo "--------------"
	@echo
	@make test-repo
	gitsocket once --gitrepo "_test/local"
	@make test-server-test-result
	@echo
	@echo "Test Passed"
	@echo

test-setup:
	@echo
	@echo "gitsocket setup"
	@echo "---------------"
	@echo
	@make test-repo
	gitsocket setup --gitrepo "_test/local" --command "echo hello world"
	cat _test/local/.git/hooks/post-checkout
	@echo
	if [ "`tail -1 _test/local/.git/hooks/post-checkout`" != "echo hello world" ]; then exit 1; fi
	@echo
	@echo "Test Passed"
	@echo

.PHONY: all build clean test
.PHONY: test-repo test-start-gitsocket test-stop-gitsocket test-server-test-result
.PHONY: test-server test-once test-setup
.PHONY: test-server-default test-server-socket test-server-port test-server-ip-port
