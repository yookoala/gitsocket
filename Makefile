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
	cd _test/local && kill `cat "test.pid"`

test: test-socket test-port

test-socket:
	@make test-repo
	@echo
	@echo "Functional Tests (socket)"
	@echo "-------------------------"
	@echo
	##
	## start gitsocket on local
	cd _test/local && gitsocket server --daemon --pidfile "test.pid"
	@echo
	##
	## trigger the gitsocket
	cd _test/local && ls
	cd _test/local && gitsocket client
	@echo
	##
	## verify the checkout
	sleep 1
	cd _test/local && ls
	@echo
	##
	## test if the OTHER.md is still here
	cd _test/local && if [ -f "OTHER.md" ]; then make test-stop-gitsocket; exit 1; fi
	@echo
	##
	## test if the git status is HEAD on origin/master
	cd _test/local && git status | head -1 > status.txt
	cd _test/local && if ! grep "HEAD detached at origin/master" status.txt; then make test-stop-gitsocket; exit 1 ; fi
	@echo
	##
	## end gitsocket
	make test-stop-gitsocket
	@echo
	@echo "Test Passed"
	@echo

test-port:
	@make test-repo
	@echo
	@echo "Functional Tests (port)"
	@echo "-----------------------"
	@echo
	##
	## start gitsocket on local
	cd _test/local && gitsocket server --daemon --listen 9301 --pidfile "test.pid"
	@echo
	##
	## trigger the gitsocket
	cd _test/local && ls
	cd _test/local && gitsocket client --conn 9301
	@echo
	##
	## verify the checkout
	sleep 1
	cd _test/local && ls
	@echo
	##
	## test if the OTHER.md is still here
	cd _test/local && if [ -f "OTHER.md" ]; then make test-stop-gitsocket; exit 1; fi
	@echo
	##
	## test if the git status is HEAD on origin/master
	cd _test/local && git status | head -1 > status.txt
	cd _test/local && if ! grep "HEAD detached at origin/master" status.txt; then make test-stop-gitsocket; exit 1 ; fi
	@echo
	##
	## end gitsocket
	make test-stop-gitsocket
	@echo
	@echo "Test Passed"
	@echo

.PHONY: all build clean test
.PHONY: test-repo test-start-gitsocket test-stop-gitsocket
.PHONY: test-socket test-port
