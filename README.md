# Systemd

[![Godoc Reference][pkg.go.dev_img]][pkg.go.dev]
[![Pipeline Status][pipeline_img ]][pipeline ]

Go package that provides functions to integrate with systemd features.

[pkg.go.dev]: https://pkg.go.dev/github.com/matthewpi/sd
[pkg.go.dev_img]: https://img.shields.io/badge/%E2%80%8B-reference-007d9c?logo=go&logoColor=white&style=flat-square
[pipeline]: https://github.com/matthewpi/sd/actions/workflows/ci.yaml
[pipeline_img]: https://img.shields.io/github/actions/workflow/status/matthewpi/sd/ci.yaml?style=flat-square&label=tests

## Features

- systemd notify - `sd_notify` (`Type=notify` and `Type=notify-reload`)
  - Allows applications to notify systemd about its status, useful for ensuring systemd knows when a service is actually started or indicating status details.
  - Support for watchdogs to ensure applications are still alive, similar to a Kubernetes readiness probe.
- systemd sockets
  - Allows applications to bind to privileged ports without privileges.
  - Support for socket-activation to allow applications to be started automatically when an incoming connection comes in.
  - Simplifies binding to unix sockets as an application doesn't need special logic to handle it, instead just binds to a listener, the same as if a port was being used.

## Installation

```bash
go get github.com/matthewpi/sd
```

## Usage

### sdlisten

See [`sdlisten/example_test.go`](./sdlisten/example_test.go) or the [Godoc reference](https://pkg.go.dev/github.com/matthewpi/sd/sdlisten) for examples and usage.

### sdnotify

See [`sdnotify/example_test.go`](./sdnotify/example_test.go) or the [Godoc reference](https://pkg.go.dev/github.com/matthewpi/sd/sdnotify) for examples and usage.

## Licensing

All code in this repository is licensed under the [MIT license](./LICENSE).

This package includes **ZERO** external dependencies, including any `golang.org/x` packages.
