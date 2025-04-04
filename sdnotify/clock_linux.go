// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: Copyright (c) 2025 Matthew Penner

package sdnotify

import (
	"syscall"
	"time"
	"unsafe"
)

// nowMonotonic returns the current time from the CLOCK_MONOTONIC clock.
func nowMonotonic() (time.Time, error) {
	// This constant is from `golang.org/x/sys/unix` but placed in-line to
	// avoid a dependency on the entire package just for a single constant.
	const CLOCK_MONOTONIC = 0x1 //nolint:revive

	var ts syscall.Timespec
	if _, _, err := syscall.Syscall(syscall.SYS_CLOCK_GETTIME, CLOCK_MONOTONIC, uintptr(unsafe.Pointer(&ts)), 0); err != 0 {
		return time.Time{}, err
	}
	return time.Unix(int64(ts.Sec), int64(ts.Nsec)), nil //nolint:unconvert
}
