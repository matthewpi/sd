// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: Copyright (c) 2025 Matthew Penner

//go:build !linux

package sdnotify

func Notify([]byte) error            { return nil }
func Ready() error                   { return nil }
func Reloading() error               { return nil }
func Stopping() error                { return nil }
func Status(string) error            { return nil }
func StatusBytes([]byte) error       { return nil }
func Error(error, int) error         { return nil }
func ErrorMessage(string, int) error { return nil }
func ErrorBytes([]byte, int) error   { return nil }
