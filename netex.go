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
	// ErrNetex is the base error of netex package.
	ErrNetex = errorex.New("netex")
)

// TLSConfigFromCertificateFile returns a new tls.Config with a single
// certificate loaded from specified key and cert files or returns an error.
func TLSConfigFromCertificateFile(cert, key string) (*tls.Config, error) {
	certificate, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}
	tlscfg := &tls.Config{
		Certificates: []tls.Certificate{certificate},
	}
	return tlscfg, nil
}
