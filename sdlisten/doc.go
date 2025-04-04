// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: Copyright (c) 2025 Matthew Penner

// Package sdlisten provides a simple API for binding to systemd sockets, useful
// for [socket activation], binding to unix sockets configured by users, or
// binding to privileged ports without needing escalated privileges.
//
// NOTE: this package is only useful on `linux` operating systems. Calling any
// functions exposed by this package are a no-op on other operating systems.
//
// Services are usually configured as a single `<NAME>.service` and a matching
// `<NAME>.socket`, however multiple sockets may be used. Named sockets may be
// configured by setting [FileDescriptorName=] under the `[Socket]` section in
// a `.socket` file. See both [systemd.service(5)] and [systemd.socket(5)] for
// details.
//
// See the [sd_listen_fds] docs for more details.
//
// [sd_listen_fds]: https://www.freedesktop.org/software/systemd/man/latest/sd_listen_fds.html
// [socket activation]: https://0pointer.de/blog/projects/socket-activation.html
// [FileDescriptorName=]: https://www.freedesktop.org/software/systemd/man/latest/systemd.socket.html#FileDescriptorName=
// [systemd.service(5)]: https://www.freedesktop.org/software/systemd/man/latest/systemd.service.html
// [systemd.socket(5)]: https://www.freedesktop.org/software/systemd/man/latest/systemd.socket.html
package sdlisten
