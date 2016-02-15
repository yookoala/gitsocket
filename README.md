githook
======

**githook** is a lightweight server process that listen to a socket
request and triggers the command `git fetch [remote source]` and
`git rebase [remote source]/[branch]`.

With proper git hook setup, it should helps triggering auto deploy.

**githook** is written in [golang][golang].

[golang]: https://golang.org


## Build and Install

You need to install [golang][golang] first.

Go into the folder. Build with this command:

```bash
go build
```

## Usage

```
NAME:
   githook - on socket connection, update files from git

USAGE:
   githook [global options] command [command options] [arguments...]

VERSION:
   0.1.0

COMMANDS:
   help, h	Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --remote, -r "origin"		name of remote repository
   --branch, -b "master"		branch of remote repository
   --socket, -s "./githook.sock"	path to socket to listen for connection
   --help, -h				show help
   --version, -v			print the version
```

## License

This software is licensed under MIT license.

You can find a copy of [the license][license] in this repository.

[license]: /LICENSE
