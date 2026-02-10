// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: Copyright (c) 2025 Matthew Penner

//go:build !linux

package sdlisten

import "os"

func Files(unsetEnvironment ...bool) []*os.File { return nil }
