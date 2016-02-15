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

## License

This software is licensed under MIT license.

You can find a copy of [the license][license] in this repository.

[license]: /LICENSE
