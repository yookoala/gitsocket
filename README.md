gitsocket
=========

**gitsocket** is a command line tool to help auto deploy with git.

It can start a lightweight server process that listen to a socket
request and triggers the command `git fetch [remote source]` and
`git checkout [remote source]/[branch]`.

It also help you to setup git hook so you can run command after
any triggered git checkout.

**gitsocket** is written in [golang][golang]. It has only been tested
and used on Linux. However, any POSIX environment (e.g. Mac OSX)
with [git][git] and [vi][vi] installed should be fine.

[golang]: https://golang.org
[git]: https://git-scm.com/
[vi]: http://www.vim.org


## Install

### Manually Compile and Install

You need to install [golang][golang] first.

Go into the folder. Build with this command:

```bash
go get ./...
go build
```

Just move it to any folder in your `$PATH`.

### Install with `go get`

If you have properly install `golang`, setup `$PATH` to include
`$GOPATH/bin`, you may just use `go get` to install:

```bash
go get github.com/yookoala/gitsocket
```


## Usage

The tool supports 4 commands:

### A. Socket Server

Command:
```bash
gitsocket server
```
This command creates a unix socket server. It updates the local git
repository to origin/master branch whenever there is socket input to
the socket file (default: `./gitsocket.sock`).

You may change the remote repository name, branch or socket file path
by the command options:

```manpage
NAME:
   server - socket server. listen to unix socket and update local git repository accordingly

USAGE:
   command server [command options] [arguments...]

OPTIONS:
   --remote, -r "origin"		name of remote repository
   --branch, -b "master"		branch of remote repository
   --listen, -l "./gitsocket.sock"	path to socket to listen for connection
   --pidfile, -p 			path to pidfile. empty for no pidfile
   --output, -o 			log output of server. empty for displaying on stdout
   --daemon, -d				run server as daemon. will discard all output unless you have output flag set.```
```

### B. Run Once

Command:
```bash
gitsocket once
```
This command run as the server is triggered once

```manpage
NAME:
   once - run as the server is triggered once

USAGE:
   command once [command options] [arguments...]

OPTIONS:
   --remote, -r "origin"	name of remote repository
   --branch, -b "master"	branch of remote repository
   --output, -o 		log output of server. empty for displaying on stdout
```

### C. Client

Command:
```bash
gitsocket client
```
This command connects to the gitsocket server and trigger one git checkout

```manpage
NAME:
   client - connects to socket triggers the socket server then returns the output

USAGE:
   command client [command options] [arguments...]

OPTIONS:
   --conn, -c "./gitsocket.sock"	socket or address to connect
```

### D. Setup Helper

Command:
```bash
gitsocket setup
```

This command helps setup the post-checkout script file with vi.
Just a time-saver in case you don't want to read all about git hook.

The script created will be run whenever `git checkout` is run. It will
be triggered after each time `gitsocket server` is triggered.

```manpage
NAME:
   setup - help setting up the post-checkout hook in the current repository folder. depends on vi

USAGE:
   command setup [command options] [arguments...]

OPTIONS:
   --command, -c 	Use shell command to be run in post-checkout hook. By default it starts vi to edit it. If the file exists, it fails (unless you have -f flag).
   --force, -f		Overwrites the current file with -c flag set. Default not set
```


## Report Bug

You are welcomed to report issue of this software.

Please use our [issue tracker][issues] to report problem.

[issues]: https://github.com/yookoala/gitsocket/issues


## License

This software is licensed under MIT license.

You can find a copy of [the license][license] in this repository.

[license]: /LICENSE
