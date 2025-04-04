// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: Copyright (c) 2025 Matthew Penner

package sdlisten

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"slices"
)

// Listener is a wrapper around a [net.Listener] used to attach additional data
// to the listener.
type Listener struct {
	// Listener is the underlying [net.Listener].
	net.Listener

	// Name of the listener, provided by systemd.
	//
	// You can use [FileDescriptorName=] property in [systemd.socket(5)] units
	// associated with this application to set this value. Keep in mind that the
	// name will apply to all listeners defined within the same [systemd.socket(5)]
	// unit. In order to have separate names for listeners, you need to use
	// multiple separate [systemd.socket(5)] units with the [systemd.service(5)]
	// the application is being ran by.
	//
	// NOTE: Name is not guaranteed to be unique. With newer versions of systemd
	// it will default to the name of the `.socket` unit the listener came from.
	// If systemd does not provide us a name, Name will be set to `LISTEN_FD_${FD}`,
	// where `FD` is the listeners file descriptor number.
	//
	// [systemd.service(5)]: https://www.freedesktop.org/software/systemd/man/latest/systemd.service.html
	// [systemd.socket(5)]: https://www.freedesktop.org/software/systemd/man/latest/systemd.socket.html
	// [FileDescriptorName=]: https://www.freedesktop.org/software/systemd/man/latest/systemd.socket.html#FileDescriptorName=
	Name string
}

// Listeners opens [Listener]s on the file descriptors provided by [Files].
func Listeners() ([]Listener, error) {
	files := Files(true)
	listeners := make([]Listener, 0, len(files))
	var errs error
	for _, f := range files {
		name := f.Name()
		l, err := net.FileListener(f)
		if err != nil {
			errs = errors.Join(errs, fmt.Errorf("sdlisten: unable to open listener (%s): %w", name, err))
			continue
		}
		_ = f.Close()
		listeners = append(listeners, Listener{
			Listener: l,
			Name:     name,
		})
	}
	return slices.Clip(listeners), errs
}

// TLSListeners is the same as [Listeners] except that it will wrap all TCP
// [net.Listener]s using [tls.NewListener] and the provided [*tls.Config].
//
// If the provided [*tls.Config] is nil, the result of [Listeners] will be
// returned as-is without being modified.
func TLSListeners(tlsConfig *tls.Config) ([]Listener, error) {
	listeners, err := Listeners()
	if err != nil {
		return nil, err
	}

	if listeners == nil || tlsConfig == nil {
		return listeners, nil
	}

	for i, l := range listeners {
		// Activate TLS only for TCP sockets
		if l.Addr().Network() == "tcp" {
			listeners[i].Listener = tls.NewListener(l, tlsConfig)
		}
	}

	return listeners, nil
}

// PacketConn is a wrapper around a [net.PacketConn] used to attach additional
// data to the connection.
type PacketConn struct {
	// PacketConn is the underlying [net.PacketConn].
	net.PacketConn

	// Name of the listener, provided by systemd.
	//
	// You can use [FileDescriptorName=] property in [systemd.socket(5)] units
	// associated with this application to set this value. Keep in mind that the
	// name will apply to all listeners defined within the same [systemd.socket(5)]
	// unit. In order to have separate names for listeners, you need to use
	// multiple separate [systemd.socket(5)] units with the [systemd.service(5)]
	// the application is being ran by.
	//
	// NOTE: Name is not guaranteed to be unique. With newer versions of systemd
	// it will default to the name of the `.socket` unit the listener came from.
	// If systemd does not provide us a name, Name will be set to `LISTEN_FD_${FD}`,
	// where `FD` is the listeners file descriptor number.
	//
	// [systemd.service(5)]: https://www.freedesktop.org/software/systemd/man/latest/systemd.service.html
	// [systemd.socket(5)]: https://www.freedesktop.org/software/systemd/man/latest/systemd.socket.html
	// [FileDescriptorName=]: https://www.freedesktop.org/software/systemd/man/latest/systemd.socket.html#FileDescriptorName=
	Name string
}

// PacketConns opens [PacketConn]s on the file descriptors provided by [Files].
func PacketConns() ([]PacketConn, error) {
	files := Files(true)
	conns := make([]PacketConn, 0, len(files))
	var errs error
	for _, f := range files {
		name := f.Name()
		pc, err := net.FilePacketConn(f)
		if err != nil {
			errs = errors.Join(errs, fmt.Errorf("sdlisten: unable to open packet conn (%s): %w", name, err))
			continue
		}
		_ = f.Close()
		conns = append(conns, PacketConn{
			PacketConn: pc,
			Name:       name,
		})
	}
	return slices.Clip(conns), errs
}
