package pcp

import (
	"context"
	"fmt"
	"net"
)

// Conn is a thin wrapper around net.Conn for PCP communication.
// It embeds net.Conn so it satisfies the net.Conn interface and
// can be used anywhere a plain connection is expected.
type Conn struct {
	net.Conn
}

// Dial establishes a TCP connection to the given address and sends
// the PCP magic atom ("pcp\n" with zero-length payload) to initiate
// the protocol handshake.
//
// On failure the underlying connection (if any) is closed before
// the error is returned.
func Dial(address string) (*Conn, error) {
	return DialContext(context.Background(), address)
}

// DialContext establishes a TCP connection using the provided context and sends
// the PCP magic atom ("pcp\n" with zero-length payload) to initiate the
// protocol handshake.
//
// On failure the underlying connection (if any) is closed before
// the error is returned.
func DialContext(ctx context.Context, address string) (*Conn, error) {
	var d net.Dialer
	c, err := d.DialContext(ctx, "tcp", address)
	if err != nil {
		return nil, fmt.Errorf("pcp: dial %s: %w", address, err)
	}

	conn, err := NewConn(c)
	if err != nil {
		return nil, fmt.Errorf("pcp: initialize connection %s: %w", address, err)
	}
	return conn, nil
}

// NewConn wraps an existing net.Conn and sends the PCP magic atom.
// This is useful when you already have a connection (e.g. from a
// test, a TLS wrapper, or an accepted socket).
//
// If the magic atom cannot be written the connection is closed and
// an error is returned.
func NewConn(c net.Conn) (*Conn, error) {
	if c == nil {
		return nil, fmt.Errorf("pcp: NewConn: nil net.Conn")
	}
	conn := &Conn{Conn: c}
	if err := NewEmptyAtom(PCPConnect).Write(c); err != nil {
		c.Close()
		return nil, fmt.Errorf("pcp: sending connect atom: %w", err)
	}
	return conn, nil
}

// ReadAtom reads a single PCP atom from the connection.
func (c *Conn) ReadAtom() (*Atom, error) {
	if c == nil || c.Conn == nil {
		return nil, fmt.Errorf("pcp: ReadAtom: nil connection")
	}
	return ReadAtom(c.Conn)
}

// WriteAtom writes a single PCP atom to the connection.
func (c *Conn) WriteAtom(a *Atom) error {
	if c == nil || c.Conn == nil {
		return fmt.Errorf("pcp: WriteAtom: nil connection")
	}
	if a == nil {
		return fmt.Errorf("pcp: WriteAtom: nil atom")
	}
	return a.Write(c.Conn)
}
