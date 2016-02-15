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

The tool supports 2 commands:

### Socket Server

Command:
```
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
   --socket, -s "./githook.sock"	path to socket to listen for connection
```


### Setup Help

Command:
```
githook setup
```

This command helps setup the post-checkout script file with vi.
Just a time-saver in case you don't want to read all about git hook.

The script created will be run whenever `git checkout` is run. It will
be triggered after each time `githook server` is triggered.

Note that this command only works if you run it inside a git repository.


## License

This software is licensed under MIT license.

You can find a copy of [the license][license] in this repository.

[license]: /LICENSE
