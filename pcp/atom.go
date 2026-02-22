package pcp

import (
	"encoding/binary"
	"fmt"
	"io"
)

// Atom represents a PCP atom — the fundamental unit of the protocol.
//
// An atom is one of two kinds:
//   - Container (parent): holds zero or more child atoms. On the wire,
//     the length field has the MSB set and the lower 31 bits encode the
//     number of children.
//   - Data (leaf): holds a byte payload. On the wire, the length field
//     is the payload byte count.
//
// Use the constructor functions (NewParentAtom, NewIntAtom, etc.) to
// create atoms, and ReadAtom to decode from a stream.
type Atom struct {
	// Tag is the 4-byte identifier for this atom.
	Tag ID4

	// isParent is true for container atoms. When true, children is
	// meaningful; when false, data is meaningful.
	isParent bool

	// children holds sub-atoms for container atoms.
	children []*Atom

	// data holds the byte payload for data atoms.
	data []byte
}

// ---------------------------------------------------------------------------
// Predicates and accessors
// ---------------------------------------------------------------------------

// IsParent returns true if this atom is a container (parent) node.
func (a *Atom) IsParent() bool {
	return a.isParent
}

// Children returns the child atoms. Returns nil for data atoms.
func (a *Atom) Children() []*Atom {
	return a.children
}

// Data returns the raw byte payload. Returns nil for parent atoms.
func (a *Atom) Data() []byte {
	return a.data
}

// NumChildren returns the number of children for a container atom, or 0.
func (a *Atom) NumChildren() int {
	return len(a.children)
}

// ---------------------------------------------------------------------------
// Type-safe data getters
// ---------------------------------------------------------------------------

// GetInt returns the payload as a little-endian uint32.
func (a *Atom) GetInt() (uint32, error) {
	if len(a.data) != 4 {
		return 0, fmt.Errorf("pcp: GetInt: expected 4 bytes, got %d", len(a.data))
	}
	return binary.LittleEndian.Uint32(a.data), nil
}

// GetShort returns the payload as a little-endian uint16.
func (a *Atom) GetShort() (uint16, error) {
	if len(a.data) != 2 {
		return 0, fmt.Errorf("pcp: GetShort: expected 2 bytes, got %d", len(a.data))
	}
	return binary.LittleEndian.Uint16(a.data), nil
}

// GetByte returns the payload as a single byte.
func (a *Atom) GetByte() (byte, error) {
	if len(a.data) != 1 {
		return 0, fmt.Errorf("pcp: GetByte: expected 1 byte, got %d", len(a.data))
	}
	return a.data[0], nil
}

// GetString returns the payload as a string, stripping any trailing null byte.
func (a *Atom) GetString() string {
	d := a.data
	for len(d) > 0 && d[len(d)-1] == 0 {
		d = d[:len(d)-1]
	}
	return string(d)
}

// GetID returns the payload as a 16-byte GnuID.
func (a *Atom) GetID() (GnuID, error) {
	var id GnuID
	if len(a.data) != 16 {
		return id, fmt.Errorf("pcp: GetID: expected 16 bytes, got %d", len(a.data))
	}
	copy(id[:], a.data)
	return id, nil
}

// GetID4 returns the payload as a 4-byte ID4.
func (a *Atom) GetID4() (ID4, error) {
	var id ID4
	if len(a.data) != 4 {
		return id, fmt.Errorf("pcp: GetID4: expected 4 bytes, got %d", len(a.data))
	}
	copy(id[:], a.data)
	return id, nil
}

// ---------------------------------------------------------------------------
// Child lookup helpers
// ---------------------------------------------------------------------------

// FindChild returns the first child atom whose tag matches, or nil.
func (a *Atom) FindChild(tag ID4) *Atom {
	for _, c := range a.children {
		if c.Tag == tag {
			return c
		}
	}
	return nil
}

// FindChildren returns all child atoms whose tag matches.
func (a *Atom) FindChildren(tag ID4) []*Atom {
	var result []*Atom
	for _, c := range a.children {
		if c.Tag == tag {
			result = append(result, c)
		}
	}
	return result
}

// ---------------------------------------------------------------------------
// Atom constructors
// ---------------------------------------------------------------------------

// NewParentAtom creates a container atom with the given children.
func NewParentAtom(tag ID4, children ...*Atom) *Atom {
	return &Atom{
		Tag:      tag,
		isParent: true,
		children: children,
	}
}

// NewIntAtom creates a data atom with a 4-byte little-endian uint32 payload.
func NewIntAtom(tag ID4, v uint32) *Atom {
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, v)
	return &Atom{Tag: tag, data: data}
}

// NewShortAtom creates a data atom with a 2-byte little-endian uint16 payload.
func NewShortAtom(tag ID4, v uint16) *Atom {
	data := make([]byte, 2)
	binary.LittleEndian.PutUint16(data, v)
	return &Atom{Tag: tag, data: data}
}

// NewByteAtom creates a data atom with a single byte payload.
func NewByteAtom(tag ID4, v byte) *Atom {
	return &Atom{Tag: tag, data: []byte{v}}
}

// NewStringAtom creates a data atom with a null-terminated string payload.
// A trailing 0x00 byte is appended automatically.
func NewStringAtom(tag ID4, s string) *Atom {
	data := make([]byte, len(s)+1)
	copy(data, s)
	// data[len(s)] is already 0 from make
	return &Atom{Tag: tag, data: data}
}

// NewBytesAtom creates a data atom with an arbitrary byte payload.
// The data is copied.
func NewBytesAtom(tag ID4, src []byte) *Atom {
	data := make([]byte, len(src))
	copy(data, src)
	return &Atom{Tag: tag, data: data}
}

// NewIDAtom creates a data atom with a 16-byte GnuID payload.
func NewIDAtom(tag ID4, id GnuID) *Atom {
	data := make([]byte, 16)
	copy(data, id[:])
	return &Atom{Tag: tag, data: data}
}

// NewID4Atom creates a data atom with a 4-byte ID4 payload.
func NewID4Atom(tag ID4, v ID4) *Atom {
	data := make([]byte, 4)
	copy(data, v[:])
	return &Atom{Tag: tag, data: data}
}

// NewEmptyAtom creates a data atom with zero-length payload.
// Used for ping/pong and similar signaling atoms.
func NewEmptyAtom(tag ID4) *Atom {
	return &Atom{Tag: tag, data: []byte{}}
}

// ---------------------------------------------------------------------------
// Wire format reader
// ---------------------------------------------------------------------------

// ReadAtom reads a single atom (and all descendants for containers)
// from the reader. The reader must provide data in PCP wire format:
//
//	[4]byte  Tag
//	uint32   Length (MSB=1 → parent with N children; MSB=0 → N bytes of data)
//	...      Payload (children or data bytes)
func ReadAtom(r io.Reader) (*Atom, error) {
	var header [AtomHeaderSize]byte
	if _, err := io.ReadFull(r, header[:]); err != nil {
		return nil, fmt.Errorf("pcp: reading atom header: %w", err)
	}

	tag := ID4{header[0], header[1], header[2], header[3]}
	value := binary.LittleEndian.Uint32(header[4:8])

	if value&AtomParentBit != 0 {
		// Container atom
		numChildren, err := uint32ToInt(value & AtomParentMask)
		if err != nil {
			return nil, fmt.Errorf("pcp: invalid child count for %q: %w", tag, err)
		}
		var children []*Atom
		for i := 0; i < numChildren; i++ {
			child, err := ReadAtom(r)
			if err != nil {
				return nil, fmt.Errorf("pcp: reading child %d of %q: %w", i, tag, err)
			}
			children = append(children, child)
		}
		return &Atom{Tag: tag, isParent: true, children: children}, nil
	}

	// Data atom
	payloadLen, err := uint32ToInt(value)
	if err != nil {
		return nil, fmt.Errorf("pcp: invalid payload length for %q: %w", tag, err)
	}
	data := make([]byte, payloadLen)
	if value > 0 {
		if _, err := io.ReadFull(r, data); err != nil {
			return nil, fmt.Errorf("pcp: reading %d data bytes of %q: %w", value, tag, err)
		}
	}
	return &Atom{Tag: tag, data: data}, nil
}

// ---------------------------------------------------------------------------
// Wire format writer
// ---------------------------------------------------------------------------

// Write serializes the atom (and all descendants) to the writer in PCP
// wire format.
func (a *Atom) Write(w io.Writer) error {
	var header [AtomHeaderSize]byte
	copy(header[:4], a.Tag[:])

	if a.isParent {
		if uint64(len(a.children)) > uint64(AtomParentMask) {
			return fmt.Errorf("pcp: child count %d exceeds wire format limit for %q", len(a.children), a.Tag)
		}
		binary.LittleEndian.PutUint32(header[4:8], uint32(len(a.children))|AtomParentBit)
		if err := writeFull(w, header[:]); err != nil {
			return fmt.Errorf("pcp: writing header for %q: %w", a.Tag, err)
		}
		for _, child := range a.children {
			if err := child.Write(w); err != nil {
				return fmt.Errorf("pcp: writing child of %q: %w", a.Tag, err)
			}
		}
	} else {
		if uint64(len(a.data)) > uint64(AtomParentMask) {
			return fmt.Errorf("pcp: data length %d exceeds wire format limit for %q", len(a.data), a.Tag)
		}
		binary.LittleEndian.PutUint32(header[4:8], uint32(len(a.data)))
		if err := writeFull(w, header[:]); err != nil {
			return fmt.Errorf("pcp: writing header for %q: %w", a.Tag, err)
		}
		if len(a.data) > 0 {
			if err := writeFull(w, a.data); err != nil {
				return fmt.Errorf("pcp: writing %d data bytes for %q: %w", len(a.data), a.Tag, err)
			}
		}
	}
	return nil
}

func writeFull(w io.Writer, p []byte) error {
	for len(p) > 0 {
		n, err := w.Write(p)
		if err != nil {
			return err
		}
		if n <= 0 || n > len(p) {
			return io.ErrShortWrite
		}
		p = p[n:]
	}
	return nil
}

// ---------------------------------------------------------------------------
// Skip helper
// ---------------------------------------------------------------------------

// SkipAtom reads and discards a single atom (and all descendants) from
// the reader. This is useful for handling unknown tags gracefully.
func SkipAtom(r io.Reader) error {
	var header [AtomHeaderSize]byte
	if _, err := io.ReadFull(r, header[:]); err != nil {
		return fmt.Errorf("pcp: skipping atom header: %w", err)
	}
	value := binary.LittleEndian.Uint32(header[4:8])

	if value&AtomParentBit != 0 {
		numChildren, err := uint32ToInt(value & AtomParentMask)
		if err != nil {
			return fmt.Errorf("pcp: invalid child count in skipped atom: %w", err)
		}
		for i := 0; i < numChildren; i++ {
			if err := SkipAtom(r); err != nil {
				return fmt.Errorf("pcp: skipping child %d of %d: %w", i, numChildren, err)
			}
		}
		return nil
	}

	if value > 0 {
		_, err := io.CopyN(io.Discard, r, int64(value))
		if err != nil {
			return fmt.Errorf("pcp: skipping %d data bytes: %w", value, err)
		}
	}
	return nil
}

func uint32ToInt(v uint32) (int, error) {
	maxInt := int(^uint(0) >> 1)
	if v > uint32(maxInt) {
		return 0, fmt.Errorf("value %d overflows int", v)
	}
	return int(v), nil
}
