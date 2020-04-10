// Copyright 2019 Vedran Vuk. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package netex

import (
	"crypto/tls"
	"net"
	"sync/atomic"
)

// ConnHandler is a connection handler.
type ConnHandler interface {
	// Handle must handle the specified net.Conn.
	// It is responsible for closing the specified net.Conn.
	HandleConn(net.Conn)
}

// ConnHandlerFunc is a prototype of a connection handler function.
// Appropriate funcs can be cast to this type to implement ConnHandler.
type ConnHandlerFunc func(net.Conn)

// HandleConn implements a ConnHandler on ConnHandlerFunc.
func (chf ConnHandlerFunc) HandleConn(conn net.Conn) {
	chf(conn)
}

const (
	// StateInvalid is an invalid/undefined server state.
	StateInvalid int32 = iota
	// StateIdle is idle server state.
	StateIdle
	// StateRunning is the idle server state.
	StateRunning
	// StateShuttingDown is the shutting-down server state.
	StateShuttingDown
)

// StreamServer is a blocking stream connection server. It dispatches new
// connections to a ConnHandler in goroutines.
type StreamServer struct {
	// network
	network string
	// addr is the address on which Server will listen for connections.
	addr string
	// listener
	listener net.Listener
	// handler
	handler ConnHandler
	// state indicates server state.
	state int32 // atomic access.
}

// NewStreamServer creates a new stream connection server that listens on
// specified network and addr and dispatches events to specified handler.
func NewStreamServer(network, addr string, handler ConnHandler) *StreamServer {
	p := &StreamServer{
		network: network,
		addr:    addr,
		handler: handler,
	}
	atomic.StoreInt32(&p.state, StateIdle)
	return p
}

// isReady reports if server is idle and ready to run.
func (s *StreamServer) isReady() bool {
	return atomic.LoadInt32(&s.state) == StateIdle
}

// serve is the implementation of Serve().
func (s *StreamServer) serve(l net.Listener) (err error) {
	s.listener = l
	for {
		var conn net.Conn
		conn, err = l.Accept()
		if err != nil {
			// Ignore errors durring shutdown as there is no way to
			// gracefully unblock an Accept call.
			if atomic.LoadInt32(&s.state) == StateShuttingDown {
				err = nil
			}
			break
		}
		go s.handler.HandleConn(conn)
	}
	atomic.StoreInt32(&s.state, StateIdle)
	return
}

// Serve serves on specified listener l.
// It blocks until listener l is closed or an error occurs.
func (s *StreamServer) Serve(l net.Listener) (err error) {
	if !s.isReady() {
		return ErrAlreadyRunning
	}
	atomic.StoreInt32(&s.state, StateRunning)
	return s.serve(l)
}

// ListenAndServe listens on defined Server ListenAddr and blocks until
// underlying listener returns by Close() or an error occurs.
func (s *StreamServer) ListenAndServe() error {
	if !s.isReady() {
		return ErrAlreadyRunning
	}
	atomic.StoreInt32(&s.state, StateRunning)
	l, err := net.Listen(s.network, s.addr)
	if err != nil {
		return err
	}
	return s.serve(l)
}

// ListenAndServeTLS listens on defined Server ListenAddr and blocks until
// underlying listener returns by Close() or an error occurs.
func (s *StreamServer) ListenAndServeTLS(tlsconfig *tls.Config) error {
	if !s.isReady() {
		return ErrAlreadyRunning
	}
	atomic.StoreInt32(&s.state, StateRunning)
	l, err := tls.Listen(s.network, s.addr, tlsconfig)
	if err != nil {
		return err
	}
	return s.serve(l)
}

// Close closes the listener.
func (s *StreamServer) Close() error {
	if atomic.LoadInt32(&s.state) != StateRunning {
		return ErrNotRunning
	}
	atomic.StoreInt32(&s.state, StateShuttingDown)
	return s.listener.Close()
}
