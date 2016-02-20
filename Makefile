# Makefile
# --------
# This file is mostly for behaviour test.
# Check README.md for installation options.

# add pwd to path so the built binary can be used
export PATH:=$(PWD):$(PATH)

# colors and special characters for terminal output
export RED:=\033[0;31m
export GREEN:=\033[0;32m
export NC:=\033[0m
export LF:="\033[999D"

# signs for results
export OK:="${GREEN}✔${NC} "
export FAIL:="${RED}❌${NC} "

# quick definition for ok or fail actions
export DOOK:=echo ${LF}${OK}
export DOFAIL:=echo ${LF}${FAIL}; exit 1

# ------------

build: gitsocket

all: clean build test

clean:
	rm -f gitsocket

gitsocket:
	@echo
	@echo "Building gitsocket"
	@echo "=================="
	@echo
	go build -o gitsocket

test-repo:
	@echo
	@echo "Creating test reopsitory"
	@echo "------------------------"
	@echo
	@rm -Rf _test
	@echo -n "-  create remote repo"
	@mkdir -p _test/remote
	@cd _test/remote && git init --bare 1>/dev/null && echo ${LF}${OK}
	@echo -n "-  create local repo with user config"
	@mkdir -p _test/local
	@cd _test/local && git init 1>/dev/null
	@cd _test/local && git config user.name  "Test User" 1>/dev/null
	@cd _test/local && git config user.email "user@localhost" 1>/dev/null
	@echo ${LF}${OK}
	@echo "-  add README.md and commit as \"Initial commit\""
	@echo "# README" > _test/local/README.md
	@cd _test/local && git add README.md && git commit -m "Initial commit" 1>/dev/null
	@cd _test/local && git checkout -b other_branch 1>/dev/null
	@echo "-  add README.md and commit as \"Initial commit\"" ${LF}${OK}
	@echo "-  add OTHER.md and commit as \"Unwanted commit\""
	@echo "# OTHER" > _test/local/OTHER.md
	@cd _test/local && git add OTHER.md && git commit -m "Unwanted commit" 1>/dev/null
	@cd _test/local && git remote add origin ../remote 1>/dev/null
	@cd _test/local && git push -u origin master 1>/dev/null
	@echo "-  add OTHER.md and commit as \"Unwanted commit\"" ${LF}${OK}
	@echo "-  test repository created" ${LF}${OK}

test-stop-gitsocket:
	@echo -n "stop gitsocket: "
	@kill `cat "test.pid"`
	@echo "[ok]"

test-server-test-result:
	@echo "-  resulting local directory" ${LF}${OK}
	@echo -n "   - "
	@cd _test/local && ls
	@echo -n "-  test if the file \"OTHER.md\" is gone"
	@cd _test/local && if [ -f "OTHER.md" ]; then ${DOFAIL}; make test-stop-gitsocket; exit 1; fi && ${DOOK}
	@echo -n "-  test if the current commit is \"Initial commit\""
	@cd _test/local && git log --format=oneline | cut -d " " -f2- > current.txt
	@cd _test/local && if ! grep -q "Initial commit" current.txt; then ${DOFAIL}; make test-stop-gitsocket; exit 1; fi && ${DOOK}

test: test-server test-once test-setup

test-server: test-server-default test-server-socket test-server-port test-server-ip-port

test-server-default:
	@echo
	@echo "Test: gitsocket server (default)"
	@echo "================================"
	@echo
	@make test-repo
	@echo
	@echo "Test command 1"
	@echo "--------------"
	gitsocket server --daemon --gitrepo "_test/local" --pidfile "test.pid"
	@echo
	@echo "Test command 2"
	@echo "--------------"
	gitsocket client
	@echo
	@echo "Examine the result"
	@echo "------------------"
	@make test-server-test-result
	@make test-stop-gitsocket
	@echo
	@echo "--------------"
	@echo "Test Passed "${OK}
	@echo "--------------"
	@echo

test-server-socket:
	@echo
	@echo "Test: gitsocket server (socket)"
	@echo "==============================="
	@echo
	@make test-repo
	@echo
	@echo "Test command 1"
	@echo "--------------"
	gitsocket server --daemon --gitrepo "_test/local" --listen "_test/test.sock" --pidfile "test.pid"
	@echo
	@echo "Test command 2"
	@echo "--------------"
	gitsocket client --conn "_test/test.sock"
	@echo
	@echo "Examine the result"
	@echo "------------------"
	@make test-server-test-result
	@make test-stop-gitsocket
	@echo
	@echo "--------------"
	@echo "Test Passed "${OK}
	@echo "--------------"
	@echo

test-server-port:
	@echo
	@echo "Test: gitsocket server (port)"
	@echo "============================="
	@echo
	@make test-repo
	@echo
	@echo "Test command 1"
	@echo "--------------"
	gitsocket server --daemon --gitrepo "_test/local" --listen 9301 --pidfile "test.pid"
	@echo
	@echo "Test command 2"
	@echo "--------------"
	gitsocket client --conn 9301
	@echo
	@echo "Examine the result"
	@echo "------------------"
	@make test-server-test-result
	@make test-stop-gitsocket
	@echo
	@echo "--------------"
	@echo "Test Passed "${OK}
	@echo "--------------"
	@echo

test-server-ip-port:
	@echo
	@echo "Test: gitsocket server (ip:port)"
	@echo "================================"
	@echo
	@make test-repo
	@echo
	@echo "Test command 1"
	@echo "--------------"
	gitsocket server --daemon --gitrepo "_test/local" --listen 9301 --pidfile "test.pid"
	@echo
	@echo "Test command 2"
	@echo "--------------"
	gitsocket client --conn 127.0.0.1:9301
	@echo
	@echo "Examine the result"
	@echo "------------------"
	@make test-server-test-result
	@make test-stop-gitsocket
	@echo
	@echo "--------------"
	@echo "Test Passed "${OK}
	@echo "--------------"
	@echo

test-once:
	@echo
	@echo "Test: gitsocket once"
	@echo "===================="
	@echo
	@make test-repo
	@echo
	@echo "Test command"
	@echo "------------"
	gitsocket once --gitrepo "_test/local"
	@echo
	@echo "Examine the result"
	@echo "------------------"
	@make test-server-test-result
	@echo
	@echo "--------------"
	@echo "Test Passed "${OK}
	@echo "--------------"
	@echo

test-setup:
	@echo
	@echo "Test: gitsocket setup"
	@echo "====================="
	@echo
	@make test-repo
	@echo
	@echo "Test command"
	@echo "------------"
	gitsocket setup --gitrepo "_test/local" --command "echo hello world"
	@echo
	@echo "Examine the result"
	@echo "------------------"
	@echo -n "-  test if the post-checkout script created according to --command"
	@if [ "`tail -1 _test/local/.git/hooks/post-checkout`" != "echo hello world" ]; then ${DOFAIL}; exit 1; fi && ${DOOK}
	@echo
	@echo "--------------"
	@echo "Test Passed "${OK}
	@echo "--------------"
	@echo

.PHONY: all build clean test
.PHONY: test-repo test-start-gitsocket test-stop-gitsocket test-server-test-result
.PHONY: test-server test-once test-setup
.PHONY: test-server-default test-server-socket test-server-port test-server-ip-port
