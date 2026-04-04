package pcp

import (
	"bytes"
	"errors"
	"io"
	"math"
	"strings"
	"testing"
)

type shortWriter struct {
	w        io.Writer
	maxChunk int
}

func (sw *shortWriter) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	n := sw.maxChunk
	if n > len(p) {
		n = len(p)
	}
	return sw.w.Write(p[:n])
}

// ---------------------------------------------------------------------------
// Constructors — data atoms
// ---------------------------------------------------------------------------

func TestNewIntAtom(t *testing.T) {
	a := NewIntAtom(PCPHeloVersion, 1234)
	if a.Tag != PCPHeloVersion {
		t.Errorf("Tag: got %v, want %v", a.Tag, PCPHeloVersion)
	}
	if a.IsParent() {
		t.Error("expected data atom, got parent")
	}
	v, err := a.GetInt()
	if err != nil {
		t.Fatalf("GetInt: %v", err)
	}
	if v != 1234 {
		t.Errorf("GetInt: got %d, want 1234", v)
	}
}

func TestNewShortAtom(t *testing.T) {
	a := NewShortAtom(PCPHeloPort, 7144)
	v, err := a.GetShort()
	if err != nil {
		t.Fatalf("GetShort: %v", err)
	}
	if v != 7144 {
		t.Errorf("GetShort: got %d, want 7144", v)
	}
}

func TestNewByteAtom(t *testing.T) {
	a := NewByteAtom(PCPBcstTTL, 7)
	v, err := a.GetByte()
	if err != nil {
		t.Fatalf("GetByte: %v", err)
	}
	if v != 7 {
		t.Errorf("GetByte: got %d, want 7", v)
	}
}

func TestNewStringAtom(t *testing.T) {
	a := NewStringAtom(PCPHeloAgent, "TestAgent/1.0")
	s := a.GetString()
	if s != "TestAgent/1.0" {
		t.Errorf("GetString: got %q, want %q", s, "TestAgent/1.0")
	}
	// raw data must have trailing null byte
	if last := a.Data()[len(a.Data())-1]; last != 0 {
		t.Errorf("trailing null byte missing, last byte is %d", last)
	}
}

func TestNewStringAtom_EmptyString(t *testing.T) {
	a := NewStringAtom(PCPHeloAgent, "")
	if a.GetString() != "" {
		t.Errorf("expected empty string, got %q", a.GetString())
	}
}

func TestGetString_TrimsAllTrailingNulls(t *testing.T) {
	a := NewBytesAtom(PCPHeloAgent, []byte{'A', 'B', 0, 0, 0})
	if got := a.GetString(); got != "AB" {
		t.Errorf("GetString(): got %q, want %q", got, "AB")
	}
}

func TestNewBytesAtom(t *testing.T) {
	src := []byte{0xde, 0xad, 0xbe, 0xef}
	a := NewBytesAtom(PCPChanKey, src)
	if !bytes.Equal(a.Data(), src) {
		t.Errorf("Data mismatch: got %v, want %v", a.Data(), src)
	}
	// mutation of src must not affect atom
	src[0] = 0x00
	if a.Data()[0] != 0xde {
		t.Error("NewBytesAtom should copy data, not reference")
	}
}

func TestNewIDAtom(t *testing.T) {
	var id GnuID
	for i := range id {
		id[i] = byte(i + 1)
	}
	a := NewIDAtom(PCPHeloSessionID, id)
	got, err := a.GetID()
	if err != nil {
		t.Fatalf("GetID: %v", err)
	}
	if got != id {
		t.Errorf("GetID: got %v, want %v", got, id)
	}
}

func TestNewID4Atom(t *testing.T) {
	tag := NewID4("test")
	val := NewID4("ver\x00")
	a := NewID4Atom(tag, val)
	got, err := a.GetID4()
	if err != nil {
		t.Fatalf("GetID4: %v", err)
	}
	if got != val {
		t.Errorf("GetID4: got %v, want %v", got, val)
	}
}

func TestNewEmptyAtom(t *testing.T) {
	a := NewEmptyAtom(PCPPing)
	if a.IsParent() {
		t.Error("expected data atom")
	}
	if len(a.Data()) != 0 {
		t.Errorf("expected empty data, got len %d", len(a.Data()))
	}
}

// ---------------------------------------------------------------------------
// Constructors — parent atom
// ---------------------------------------------------------------------------

func TestNewParentAtom(t *testing.T) {
	child1 := NewIntAtom(PCPHeloVersion, 1200)
	child2 := NewStringAtom(PCPHeloAgent, "PeerCast")
	parent := NewParentAtom(PCPHelo, child1, child2)

	if !parent.IsParent() {
		t.Error("expected parent atom")
	}
	if parent.NumChildren() != 2 {
		t.Errorf("NumChildren: got %d, want 2", parent.NumChildren())
	}
	if parent.Data() != nil {
		t.Error("parent atom Data() should be nil")
	}
}

// ---------------------------------------------------------------------------
// Child lookup
// ---------------------------------------------------------------------------

func TestFindChild(t *testing.T) {
	ver := NewIntAtom(PCPHeloVersion, 1200)
	agent := NewStringAtom(PCPHeloAgent, "PeerCast")
	parent := NewParentAtom(PCPHelo, ver, agent)

	found := parent.FindChild(PCPHeloVersion)
	if found == nil {
		t.Fatal("FindChild: expected to find PCPHeloVersion")
	}
	v, _ := found.GetInt()
	if v != 1200 {
		t.Errorf("FindChild: got %d, want 1200", v)
	}
}

func TestFindChild_NotFound(t *testing.T) {
	parent := NewParentAtom(PCPHelo, NewIntAtom(PCPHeloVersion, 1))
	if parent.FindChild(PCPHeloPort) != nil {
		t.Error("FindChild: should return nil for missing tag")
	}
}

func TestFindChildren(t *testing.T) {
	ip1 := NewIntAtom(PCPHostIP, 0x01020304)
	ip2 := NewIntAtom(PCPHostIP, 0x05060708)
	port := NewShortAtom(PCPHostPort, 7144)
	parent := NewParentAtom(PCPHost, ip1, port, ip2)

	results := parent.FindChildren(PCPHostIP)
	if len(results) != 2 {
		t.Errorf("FindChildren: got %d, want 2", len(results))
	}
}

func TestFindChildren_None(t *testing.T) {
	parent := NewParentAtom(PCPHelo, NewIntAtom(PCPHeloVersion, 1))
	results := parent.FindChildren(PCPHeloPort)
	if len(results) != 0 {
		t.Errorf("FindChildren: got %d, want 0", len(results))
	}
}

// ---------------------------------------------------------------------------
// Type-safe getter error cases
// ---------------------------------------------------------------------------

func TestGetInt_WrongSize(t *testing.T) {
	a := NewByteAtom(PCPBcstTTL, 1)
	_, err := a.GetInt()
	if err == nil {
		t.Error("GetInt: expected error for 1-byte atom")
	}
}

func TestGetShort_WrongSize(t *testing.T) {
	a := NewIntAtom(PCPHeloVersion, 1)
	_, err := a.GetShort()
	if err == nil {
		t.Error("GetShort: expected error for 4-byte atom")
	}
}

func TestGetByte_WrongSize(t *testing.T) {
	a := NewIntAtom(PCPHeloVersion, 1)
	_, err := a.GetByte()
	if err == nil {
		t.Error("GetByte: expected error for 4-byte atom")
	}
}

func TestGetID_WrongSize(t *testing.T) {
	a := NewIntAtom(PCPHeloVersion, 1)
	_, err := a.GetID()
	if err == nil {
		t.Error("GetID: expected error for 4-byte atom")
	}
}

func TestGetID4_WrongSize(t *testing.T) {
	a := NewByteAtom(PCPBcstTTL, 1)
	_, err := a.GetID4()
	if err == nil {
		t.Error("GetID4: expected error for 1-byte atom")
	}
}

// ---------------------------------------------------------------------------
// Wire format roundtrip
// ---------------------------------------------------------------------------

func atomRoundtrip(t *testing.T, original *Atom) *Atom {
	t.Helper()
	var buf bytes.Buffer
	if err := original.Write(&buf); err != nil {
		t.Fatalf("Write: %v", err)
	}
	got, err := ReadAtom(&buf)
	if err != nil {
		t.Fatalf("ReadAtom: %v", err)
	}
	return got
}

func TestRoundtrip_IntAtom(t *testing.T) {
	orig := NewIntAtom(PCPHeloVersion, 0xDEADBEEF)
	got := atomRoundtrip(t, orig)
	if got.Tag != orig.Tag {
		t.Errorf("Tag: got %v, want %v", got.Tag, orig.Tag)
	}
	v, _ := got.GetInt()
	if v != 0xDEADBEEF {
		t.Errorf("value: got %#x, want 0xDEADBEEF", v)
	}
}

func TestRoundtrip_ShortAtom(t *testing.T) {
	orig := NewShortAtom(PCPHeloPort, 7144)
	got := atomRoundtrip(t, orig)
	v, _ := got.GetShort()
	if v != 7144 {
		t.Errorf("value: got %d, want 7144", v)
	}
}

func TestRoundtrip_ByteAtom(t *testing.T) {
	orig := NewByteAtom(PCPBcstTTL, 5)
	got := atomRoundtrip(t, orig)
	v, _ := got.GetByte()
	if v != 5 {
		t.Errorf("value: got %d, want 5", v)
	}
}

func TestRoundtrip_StringAtom(t *testing.T) {
	orig := NewStringAtom(PCPHeloAgent, "PeerCast/0.1218")
	got := atomRoundtrip(t, orig)
	if s := got.GetString(); s != "PeerCast/0.1218" {
		t.Errorf("value: got %q, want %q", s, "PeerCast/0.1218")
	}
}

func TestRoundtrip_EmptyAtom(t *testing.T) {
	orig := NewEmptyAtom(PCPPing)
	got := atomRoundtrip(t, orig)
	if got.Tag != PCPPing {
		t.Errorf("Tag: got %v, want %v", got.Tag, PCPPing)
	}
	if len(got.Data()) != 0 {
		t.Errorf("Data: expected empty, got len %d", len(got.Data()))
	}
}

func TestRoundtrip_IDAtom(t *testing.T) {
	var id GnuID
	for i := range id {
		id[i] = byte(i + 0x10)
	}
	orig := NewIDAtom(PCPHeloSessionID, id)
	got := atomRoundtrip(t, orig)
	g, _ := got.GetID()
	if g != id {
		t.Errorf("GnuID mismatch: got %v, want %v", g, id)
	}
}

func TestRoundtrip_ParentAtom(t *testing.T) {
	child1 := NewIntAtom(PCPHeloVersion, 1200)
	child2 := NewStringAtom(PCPHeloAgent, "PeerCast")
	child3 := NewShortAtom(PCPHeloPort, 7144)
	orig := NewParentAtom(PCPHelo, child1, child2, child3)

	got := atomRoundtrip(t, orig)

	if !got.IsParent() {
		t.Error("expected parent atom after roundtrip")
	}
	if got.NumChildren() != 3 {
		t.Errorf("NumChildren: got %d, want 3", got.NumChildren())
	}
	// verify individual children
	ver := got.FindChild(PCPHeloVersion)
	if ver == nil {
		t.Fatal("FindChild PCPHeloVersion: not found after roundtrip")
	}
	v, _ := ver.GetInt()
	if v != 1200 {
		t.Errorf("version: got %d, want 1200", v)
	}
}

func TestRoundtrip_NestedParent(t *testing.T) {
	inner := NewParentAtom(PCPChanInfo, NewStringAtom(PCPChanInfoName, "TestChannel"))
	outer := NewParentAtom(PCPChan, NewIDAtom(PCPChanID, GnuID{}), inner)

	got := atomRoundtrip(t, outer)
	if got.NumChildren() != 2 {
		t.Fatalf("outer NumChildren: got %d, want 2", got.NumChildren())
	}
	info := got.FindChild(PCPChanInfo)
	if info == nil {
		t.Fatal("nested info atom not found")
	}
	name := info.FindChild(PCPChanInfoName)
	if name == nil {
		t.Fatal("info name atom not found")
	}
	if s := name.GetString(); s != "TestChannel" {
		t.Errorf("channel name: got %q, want %q", s, "TestChannel")
	}
}

// ---------------------------------------------------------------------------
// ReadAtom — truncated input errors
// ---------------------------------------------------------------------------

func TestReadAtom_EOF(t *testing.T) {
	_, err := ReadAtom(bytes.NewReader(nil))
	if !errors.Is(err, io.EOF) && !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Errorf("expected EOF/UnexpectedEOF, got %v", err)
	}
}

func TestReadAtom_TruncatedHeader(t *testing.T) {
	_, err := ReadAtom(bytes.NewReader([]byte{0x70, 0x63, 0x70})) // only 3 bytes
	if err == nil {
		t.Error("expected error for truncated header")
	}
	if !strings.Contains(err.Error(), "reading atom header") {
		t.Errorf("expected contextual error, got %v", err)
	}
}

func TestReadAtom_TruncatedData(t *testing.T) {
	// Write a valid int atom then slice off half the payload
	a := NewIntAtom(PCPHeloVersion, 42)
	var buf bytes.Buffer
	_ = a.Write(&buf)
	truncated := buf.Bytes()[:buf.Len()-2] // remove last 2 bytes of 4-byte payload
	_, err := ReadAtom(bytes.NewReader(truncated))
	if err == nil {
		t.Error("expected error for truncated data")
	}
}

// ---------------------------------------------------------------------------
// SkipAtom
// ---------------------------------------------------------------------------

func TestSkipAtom_DataAtom(t *testing.T) {
	a := NewIntAtom(PCPHeloVersion, 99)
	var buf bytes.Buffer
	_ = a.Write(&buf)
	if err := SkipAtom(&buf); err != nil {
		t.Fatalf("SkipAtom: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected buffer drained, %d bytes remain", buf.Len())
	}
}

func TestSkipAtom_ParentAtom(t *testing.T) {
	parent := NewParentAtom(PCPHelo,
		NewIntAtom(PCPHeloVersion, 1200),
		NewStringAtom(PCPHeloAgent, "PeerCast"),
	)
	var buf bytes.Buffer
	_ = parent.Write(&buf)
	if err := SkipAtom(&buf); err != nil {
		t.Fatalf("SkipAtom parent: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected buffer drained, %d bytes remain", buf.Len())
	}
}

func TestSkipAtom_MultipeAtoms(t *testing.T) {
	// Write two atoms; skip only the first; the second should remain readable.
	first := NewIntAtom(PCPHeloVersion, 1)
	second := NewStringAtom(PCPHeloAgent, "PeerCast")
	var buf bytes.Buffer
	_ = first.Write(&buf)
	_ = second.Write(&buf)

	if err := SkipAtom(&buf); err != nil {
		t.Fatalf("SkipAtom: %v", err)
	}
	got, err := ReadAtom(&buf)
	if err != nil {
		t.Fatalf("ReadAtom after skip: %v", err)
	}
	if got.GetString() != "PeerCast" {
		t.Errorf("after skip: got %q, want %q", got.GetString(), "PeerCast")
	}
}

func TestWrite_ShortWriter(t *testing.T) {
	orig := NewParentAtom(PCPHelo,
		NewIntAtom(PCPHeloVersion, 1200),
		NewStringAtom(PCPHeloAgent, "PeerCast"),
	)

	var buf bytes.Buffer
	sw := &shortWriter{w: &buf, maxChunk: 1}
	if err := orig.Write(sw); err != nil {
		t.Fatalf("Write with short writer: %v", err)
	}

	got, err := ReadAtom(&buf)
	if err != nil {
		t.Fatalf("ReadAtom: %v", err)
	}
	if !got.IsParent() || got.NumChildren() != 2 {
		t.Fatalf("decoded atom mismatch: isParent=%v children=%d", got.IsParent(), got.NumChildren())
	}
}

func TestSkipAtom_TruncatedData_Context(t *testing.T) {
	a := NewIntAtom(PCPHeloVersion, 42)
	var buf bytes.Buffer
	_ = a.Write(&buf)
	truncated := buf.Bytes()[:buf.Len()-1]

	err := SkipAtom(bytes.NewReader(truncated))
	if err == nil {
		t.Fatal("expected error for truncated skip data")
	}
	if !strings.Contains(err.Error(), "skipping 4 data bytes") {
		t.Errorf("expected contextual error, got %v", err)
	}
}

func TestReadAtom_ExceedsMaxDepth(t *testing.T) {
	// Build a deeply nested atom that exceeds maxReadDepth.
	inner := NewEmptyAtom(PCPPing)
	for i := 0; i < maxReadDepth+2; i++ {
		inner = NewParentAtom(PCPHelo, inner)
	}
	var buf bytes.Buffer
	if err := inner.Write(&buf); err != nil {
		t.Fatal(err)
	}
	_, err := ReadAtom(&buf)
	if err == nil {
		t.Fatal("expected error for excessive nesting depth")
	}
	if !strings.Contains(err.Error(), "nesting depth exceeds") {
		t.Errorf("expected depth error, got %v", err)
	}
}

func TestSkipAtom_ExceedsMaxDepth(t *testing.T) {
	inner := NewEmptyAtom(PCPPing)
	for i := 0; i < maxReadDepth+2; i++ {
		inner = NewParentAtom(PCPHelo, inner)
	}
	var buf bytes.Buffer
	if err := inner.Write(&buf); err != nil {
		t.Fatal(err)
	}
	err := SkipAtom(&buf)
	if err == nil {
		t.Fatal("expected error for excessive nesting depth")
	}
	if !strings.Contains(err.Error(), "nesting depth exceeds") {
		t.Errorf("expected depth error, got %v", err)
	}
}

func TestReadAtom_ExceedsMaxDataSize(t *testing.T) {
	// Craft a header claiming a payload larger than MaxAtomDataSize.
	var header [8]byte
	copy(header[:4], []byte("test"))
	// Set length to MaxAtomDataSize + 1
	size := uint32(MaxAtomDataSize + 1)
	header[4] = byte(size)
	header[5] = byte(size >> 8)
	header[6] = byte(size >> 16)
	header[7] = byte(size >> 24)

	_, err := ReadAtom(bytes.NewReader(header[:]))
	if err == nil {
		t.Fatal("expected error for oversized data atom")
	}
	if !strings.Contains(err.Error(), "exceeds maximum") {
		t.Errorf("expected max size error, got %v", err)
	}
}

func TestUint32ToInt(t *testing.T) {
	v, err := uint32ToInt(123)
	if err != nil {
		t.Fatalf("uint32ToInt(123): %v", err)
	}
	if v != 123 {
		t.Fatalf("uint32ToInt(123): got %d, want 123", v)
	}

	maxInt := int(^uint(0) >> 1)
	if uint64(maxInt) < math.MaxUint32 {
		_, err = uint32ToInt(math.MaxUint32)
		if err == nil {
			t.Fatal("uint32ToInt(MaxUint32): expected overflow error on this architecture")
		}
	}
}
