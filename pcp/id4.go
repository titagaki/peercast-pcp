// Package pcp implements the PeerCast Protocol (PCP) binary protocol
// for P2P streaming.
//
// PCP is a Tag-Length-Value (TLV) protocol operating over TCP.
// All multi-byte integers are encoded in Little Endian byte order.
// The fundamental unit is an "Atom", which is either a container
// (parent) holding child atoms, or a data atom holding a byte payload.
package pcp

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

// ID4 is a 4-byte tag identifier used in PCP atom headers.
// Tags shorter than 4 characters are null-padded on the right,
// matching the C++ ID4 union behavior where the int member is
// initialized to 0 before copying characters from a string literal.
type ID4 [4]byte

// NewID4 creates an ID4 from a string literal.
// Characters are copied into the 4-byte array; if the string is
// shorter than 4 bytes, the remaining positions are zero (null) padded.
// If longer than 4 bytes, only the first 4 are used.
//
// This mirrors the C++ constructor:
//
//	ID4(const char *id) : iv(0) {
//	    for (int i=0; i<4; i++)
//	        if ((cv[i]=id[i])==0) break;
//	}
func NewID4(s string) ID4 {
	var id ID4
	for i := 0; i < 4 && i < len(s); i++ {
		if s[i] == 0 {
			break
		}
		id[i] = s[i]
	}
	return id
}

// Uint32 returns the ID4 as a little-endian uint32 value.
// This corresponds to the C++ ID4's implicit operator int() which
// returns the union's int member on a little-endian platform.
func (id ID4) Uint32() uint32 {
	return binary.LittleEndian.Uint32(id[:])
}

// ID4FromUint32 converts a uint32 value back to an ID4 in little-endian order.
func ID4FromUint32(v uint32) ID4 {
	var id ID4
	binary.LittleEndian.PutUint32(id[:], v)
	return id
}

// String returns a human-readable representation of the ID4.
// Printable ASCII characters are shown as-is; null bytes are omitted;
// other bytes are shown in \xNN escape notation.
func (id ID4) String() string {
	buf := make([]byte, 0, 8)
	for _, b := range id {
		if b >= 0x20 && b <= 0x7e {
			buf = append(buf, b)
		} else if b == '\n' {
			buf = append(buf, '\\', 'n')
		} else if b != 0 {
			buf = append(buf, fmt.Sprintf("\\x%02x", b)...)
		}
	}
	return string(buf)
}

// GnuID is a 16-byte unique identifier used as session IDs,
// channel IDs, and broadcast IDs throughout the PCP protocol.
type GnuID [16]byte

// IsEmpty returns true if all 16 bytes are zero.
func (id GnuID) IsEmpty() bool {
	return id == GnuID{}
}

// Clear sets all bytes to zero.
func (id *GnuID) Clear() {
	*id = GnuID{}
}

// String returns the lowercase hex representation of the GnuID.
func (id GnuID) String() string {
	return hex.EncodeToString(id[:])
}
