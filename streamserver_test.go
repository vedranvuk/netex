// Copyright 2020 Vedran Vuk. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package netex

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"testing"
	"time"
)

type Handler struct{}

func (h *Handler) HandleConn(conn net.Conn) {
	for {
		buf := make([]byte, 1024)
		nr, err := conn.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			log.Println(err)
			return
		}
		fmt.Println(string(buf[:nr]))
		conn.Write([]byte("<- pong"))
	}
}

func TestStreamServer(t *testing.T) {

	handler := &Handler{}

	server := NewStreamServer("tcp", "0.0.0.0:9001", handler)
	cfg, err := TLSConfigFromCertificateFile("cert.pem", "key.unencrypted.pem")
	if err != nil {
		t.Fatal(err)
		return
	}
	cfg.InsecureSkipVerify = true

	go func() {
		if err := server.ListenAndServeTLS(cfg); err != nil {
			t.Fatal(err)
			return
		}
	}()

	time.Sleep(1 * time.Millisecond)

	go func() {
		tlscfg := &tls.Config{}
		tlscfg.InsecureSkipVerify = true
		conn, err := tls.Dial("tcp", "0.0.0.0:9001", tlscfg)
		if err != nil {
			t.Fatal(err)
			return
		}
		for i := 0; i < 10000; i++ {
			conn.Write([]byte("-> ping"))
			buf := make([]byte, 1024)
			nr, err := conn.Read(buf)
			if err != nil {
				t.Fatal(err)
			}
			fmt.Println(string(buf[:nr]))
		}
		fmt.Println("done")
		conn.Close()
	}()

	time.Sleep(1 * time.Second)
	if err := server.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestStreamServerStates(t *testing.T) {

	server := NewStreamServer("no,u", ":9001", nil)
	if server.State() != StateIdle {
		t.Fatal("State failed.")
	}

	if err := server.ListenAndServe(); err != nil {
		if server.State() != StateIdle {
			t.Fatal("STate failed.")
		}
	}

	server = NewStreamServer("tcp", ":9182", ConnHandlerFunc(func(c net.Conn) {

	}))

	go func() {
		if err := server.ListenAndServe(); err != nil {
			t.Fatal(err)
		}
	}()
	time.Sleep(1 * time.Millisecond)

	if server.State() != StateRunning {
		t.Fatal("State failed.")
	}
	if err := server.Close(); err != nil {
		t.Fatal(err)
	}
	if server.State() != StateIdle {
		t.Fatal("State failed.")
	}
}
