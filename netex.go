// Copyright 2020 Vedran Vuk. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package netex provides net related utilities.
package netex

import (
	"crypto/tls"

	"github.com/vedranvuk/errorex"
)

var (
	// ErrNetEx is the base error of netex package.
	ErrNetex = errorex.New("netex")

	// ErrAlreadyRunning is returned on a listen/serve request when server is
	// already running.
	ErrAlreadyRunning = ErrNetex.Wrap("server already running")
	// ErrNotRunning is returned when a close is issued on a server that is
	// not running.
	ErrNotRunning = ErrNetex.Wrap("server not running")
)

// TlsConfigWithCertificate returns a new tls.Config with a certificate loaded
// from specified key and cert files or returns an error.
func TlsConfigWithCertificate(cert, key string) (*tls.Config, error) {
	certificate, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}
	tlscfg := &tls.Config{
		Certificates: []tls.Certificate{certificate},
	}
	return tlscfg, nil
}
