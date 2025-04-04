// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: Copyright (c) 2025 Matthew Penner

// Package sdnotify provides a simple API to notify systemd about start-up
// completion and other service status changes.
//
// NOTE: this package is only useful on `linux` operating systems. Calling any
// functions in this package are a no-op on other operating systems.
//
// See the [sd_notify] docs for more details.
//
// [sd_notify]: https://www.freedesktop.org/software/systemd/man/latest/sd_notify.html
package sdnotify
