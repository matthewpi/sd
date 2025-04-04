// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: Copyright (c) 2025 Matthew Penner

//go:build !linux

package sdlisten

import "os"

// Files is a NO-OP on platforms other than `linux`.
func Files(unsetEnvironment ...bool) []*os.File {
	return nil
}
