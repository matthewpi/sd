// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: Copyright (c) 2025 Matthew Penner

//go:build linux

package sdlisten

import (
	"os"
	"strconv"
	"strings"
	"syscall"
)

// listenFdsStart corresponds to [SD_LISTEN_FDS_START].
//
// [SD_LISTEN_FDS_START]: https://github.com/systemd/systemd/blob/v257.5/src/systemd/sd-daemon.h#L56
const listenFdsStart = 3

// Files returns the file descriptors passed to the application by systemd.
//
// Optionally, a single boolean argument with a value of `true` will cause us
// to unconditionally unset the environment variables used to get the file
// descriptors passed to us by systemd. These environment variables are:
//
// - LISTEN_PID
// - LISTEN_FDS
// - LISTEN_FDNAMES
func Files(unsetEnvironment ...bool) []*os.File {
	if len(unsetEnvironment) == 1 && unsetEnvironment[0] {
		defer func() {
			os.Unsetenv("LISTEN_PID")
			os.Unsetenv("LISTEN_FDS")
			os.Unsetenv("LISTEN_FDNAMES")
		}()
	}

	// Ensure `LISTEN_PID` matches our PID.
	pid, err := strconv.Atoi(os.Getenv("LISTEN_PID"))
	if err != nil || pid != os.Getpid() {
		return nil
	}

	// Get the number of file descriptors we need to open.
	nfds, err := strconv.Atoi(os.Getenv("LISTEN_FDS"))
	if err != nil || nfds == 0 {
		return nil
	}

	// Get the name of the file descriptors.
	names := strings.Split(os.Getenv("LISTEN_FDNAMES"), ":")

	// Open all the file descriptors.
	files := make([]*os.File, nfds)
	for i := range nfds {
		// Get the file descriptor ID, we need to account for listenFdsStart here.
		fd := i + listenFdsStart

		// Ensure the file descriptors are not passed to any child processes the
		// application spawns.
		syscall.CloseOnExec(fd)

		// Get the name of the file descriptor.
		var name string
		if i < len(names) && len(names[i]) > 0 {
			name = names[i]
		} else {
			name = "LISTEN_FD_" + strconv.Itoa(fd)
		}

		// Open the file descriptor and add it to the file slice.
		files[i] = os.NewFile(uintptr(fd), name)
	}

	return files
}
