// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: Copyright (c) 2025 Matthew Penner

// Package monotime provides a fast monotonic clock source.
package monotime

import (
	"time"
	_ "unsafe"
)

// NOTE: I tried using the example (see docs of [runtime.nanotime]) of calling
// [time.Now] during init and using [time.Since] with the previously fetched
// [time.Now] value as an alternative to [nanotime]. It "worked", but only if
// you want to measure elapsed time monotonically. If you need an actual
// monotonic clock value from the system, you need to use either [nanotime]
// or the SYS_CLOCK_GETTIME syscall with the CLOCK_MONOTONIC parameter.
//
// I preferred using [nanotime] in this case as it is internally optimized by
// Go and can take advantage of vDSOs on supported systems instead of always
// being a syscall.

//go:noescape
//go:linkname nanotime runtime.nanotime
func nanotime() int64

// Now returns the current time in nanoseconds from a monotonic clock.
func Now() int64 {
	return nanotime()
}

// Since returns the amount of time that has elapsed since t. t should be
// the result of a call to [Now] on the same machine.
func Since(t int64) time.Duration {
	return time.Duration(Now() - t)
}
