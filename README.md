githook
======

**githook** is a command line tool to help auto deploy with git.

It can start a lightweight server process that listen to a socket
request and triggers the command `git fetch [remote source]` and
`git checkout [remote source]/[branch]`.

It also help you to setup git hook so you can run command after
any triggered git checkout.

**githook** is written in [golang][golang]. It has only been tested
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
go get github.com/codegangsta/cli
go build
```

Just move it to any folder in your `$PATH`.

### Install with `go get`

If you have properly install `golang`, setup `$PATH` to include
`$GOPATH/bin`, you may just use `go get` to install:

```bash
go get github.com/yookoala/githook
```


## Usage

The tool supports 2 commands:

### A. Socket Server

Command:
```bash
githook server
```
This command creates a unix socket server. It updates the local git
repository to origin/master branch whenever there is socket input to
the socket file (default: `./githook.sock`).

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
   --listen, -l "./githook.sock"	path to socket to listen for connection
   --pidfile, -p 			path to pidfile. empty for no pidfile
   --output, -o 			log output of server. empty for displaying on stdout

Command:
```bash
githook client
```
This command connects to the githook server and trigger one git checkout

```manpage
NAME:
   client - connects to socket triggers the socket server then returns the output

USAGE:
   command client [command options] [arguments...]

OPTIONS:
   --conn, -c "./githook.sock"	socket or address to connect
```

### B. Setup Helper

Command:
```bash
githook setup
```

This command helps setup the post-checkout script file with vi.
Just a time-saver in case you don't want to read all about git hook.

The script created will be run whenever `git checkout` is run. It will
be triggered after each time `githook server` is triggered.

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

[issues]: https://github.com/yookoala/githook/issues


## License

This software is licensed under MIT license.

You can find a copy of [the license][license] in this repository.

[license]: /LICENSE
