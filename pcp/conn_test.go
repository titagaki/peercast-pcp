package pcp

import (
	"context"
	"errors"
	"io"
	"net"
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// NewConn — magic atom is sent correctly
// ---------------------------------------------------------------------------

func TestNewConn_SendsMagicAtom(t *testing.T) {
	server, client := net.Pipe()
	defer server.Close()

	// NewConn writes the magic atom on the client side.
	done := make(chan error, 1)
	var conn *Conn
	go func() {
		var err error
		conn, err = NewConn(client)
		done <- err
	}()

	// Read the magic atom from the server side.
	atom, err := ReadAtom(server)
	if err != nil {
		t.Fatalf("ReadAtom: %v", err)
	}
	if atom.Tag != PCPConnect {
		t.Errorf("Tag: got %v, want %v", atom.Tag, PCPConnect)
	}
	if atom.IsParent() {
		t.Error("expected data atom, got parent")
	}
	if len(atom.Data()) != 0 {
		t.Errorf("expected empty payload, got %d bytes", len(atom.Data()))
	}

	if err := <-done; err != nil {
		t.Fatalf("NewConn: %v", err)
	}
	if conn == nil {
		t.Fatal("NewConn returned nil Conn")
	}
}

// ---------------------------------------------------------------------------
// Dial — failure path (unreachable address)
// ---------------------------------------------------------------------------

func TestDial_Failure(t *testing.T) {
	// Use an address that will definitely fail to connect.
	_, err := Dial("127.0.0.1:0")
	if err == nil {
		t.Fatal("Dial: expected error for unreachable address")
	}
	if !strings.Contains(err.Error(), "pcp: dial") {
		t.Errorf("expected contextual error, got %v", err)
	}
}

func TestDialContext_Canceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := DialContext(ctx, "127.0.0.1:65535")
	if err == nil {
		t.Fatal("DialContext: expected canceled error")
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", err)
	}
	if !strings.Contains(err.Error(), "pcp: dial") {
		t.Errorf("expected contextual error, got %v", err)
	}
}

// ---------------------------------------------------------------------------
// NewConn — write failure closes the connection
// ---------------------------------------------------------------------------

// brokenWriter is a net.Conn that always fails on Write.
type brokenWriter struct {
	net.Conn
	closed bool
}

func (bw *brokenWriter) Write([]byte) (int, error) {
	return 0, io.ErrClosedPipe
}

func (bw *brokenWriter) Close() error {
	bw.closed = true
	return nil
}

func TestNewConn_WriteFailure_ClosesConn(t *testing.T) {
	bw := &brokenWriter{}
	_, err := NewConn(bw)
	if err == nil {
		t.Fatal("NewConn: expected error on write failure")
	}
	if !bw.closed {
		t.Error("connection should have been closed after write failure")
	}
	if !strings.Contains(err.Error(), "sending connect atom") {
		t.Errorf("expected contextual error, got %v", err)
	}
}

func TestNewConn_NilConn(t *testing.T) {
	_, err := NewConn(nil)
	if err == nil {
		t.Fatal("NewConn: expected error for nil net.Conn")
	}
	if !strings.Contains(err.Error(), "nil net.Conn") {
		t.Errorf("expected contextual error, got %v", err)
	}
}

// ---------------------------------------------------------------------------
// ReadAtom / WriteAtom convenience methods
// ---------------------------------------------------------------------------

func TestConn_ReadWriteAtom(t *testing.T) {
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	conn := &Conn{Conn: client}

	// Write an atom from conn, read it from the server side.
	orig := NewIntAtom(PCPHeloVersion, 1218)
	done := make(chan error, 1)
	go func() {
		done <- conn.WriteAtom(orig)
	}()

	got, err := ReadAtom(server)
	if err != nil {
		t.Fatalf("ReadAtom: %v", err)
	}
	v, _ := got.GetInt()
	if v != 1218 {
		t.Errorf("value: got %d, want 1218", v)
	}

	if err := <-done; err != nil {
		t.Fatalf("WriteAtom: %v", err)
	}

	// Write from server side, read via conn.ReadAtom.
	reply := NewStringAtom(PCPHeloAgent, "TestAgent")
	go func() {
		done <- reply.Write(server)
	}()

	got2, err := conn.ReadAtom()
	if err != nil {
		t.Fatalf("conn.ReadAtom: %v", err)
	}
	if got2.GetString() != "TestAgent" {
		t.Errorf("got %q, want %q", got2.GetString(), "TestAgent")
	}

	if err := <-done; err != nil {
		t.Fatalf("server Write: %v", err)
	}
}

func TestConn_WriteAtom_NilAtom(t *testing.T) {
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	conn := &Conn{Conn: client}
	err := conn.WriteAtom(nil)
	if err == nil {
		t.Fatal("WriteAtom: expected error for nil atom")
	}
	if !strings.Contains(err.Error(), "nil atom") {
		t.Errorf("expected contextual error, got %v", err)
	}
}

func TestConn_ReadAtom_NilConn(t *testing.T) {
	var conn *Conn
	_, err := conn.ReadAtom()
	if err == nil {
		t.Fatal("ReadAtom: expected error for nil receiver")
	}
	if !strings.Contains(err.Error(), "nil connection") {
		t.Errorf("expected contextual error, got %v", err)
	}
}
