package pcp

import (
	"testing"
)

// ---------------------------------------------------------------------------
// ID4 — construction and conversion
// ---------------------------------------------------------------------------

func TestNewID4_FourChars(t *testing.T) {
	id := NewID4("helo")
	want := ID4{'h', 'e', 'l', 'o'}
	if id != want {
		t.Errorf("NewID4(%q): got %v, want %v", "helo", id, want)
	}
}

func TestNewID4_ShortString(t *testing.T) {
	// Strings shorter than 4 chars must be null-padded on the right.
	id := NewID4("ok")
	want := ID4{'o', 'k', 0, 0}
	if id != want {
		t.Errorf("NewID4(%q): got %v, want %v", "ok", id, want)
	}
}

func TestNewID4_LongString(t *testing.T) {
	// Only first 4 bytes should be used.
	id := NewID4("toolong")
	want := ID4{'t', 'o', 'o', 'l'}
	if id != want {
		t.Errorf("NewID4(%q): got %v, want %v", "toolong", id, want)
	}
}

func TestNewID4_Empty(t *testing.T) {
	id := NewID4("")
	want := ID4{0, 0, 0, 0}
	if id != want {
		t.Errorf("NewID4(%q): got %v, want %v", "", id, want)
	}
}

func TestID4_Uint32_RoundTrip(t *testing.T) {
	id := NewID4("helo")
	v := id.Uint32()
	got := ID4FromUint32(v)
	if got != id {
		t.Errorf("Uint32/ID4FromUint32 roundtrip: got %v, want %v", got, id)
	}
}

func TestID4FromUint32_Zero(t *testing.T) {
	id := ID4FromUint32(0)
	want := ID4{0, 0, 0, 0}
	if id != want {
		t.Errorf("ID4FromUint32(0): got %v, want %v", id, want)
	}
}

func TestID4_Equality(t *testing.T) {
	a := NewID4("helo")
	b := NewID4("helo")
	c := NewID4("oleh")
	if a != b {
		t.Error("identical ID4 values should be equal")
	}
	if a == c {
		t.Error("different ID4 values should not be equal")
	}
}

// ---------------------------------------------------------------------------
// ID4.String
// ---------------------------------------------------------------------------

func TestID4_String_Printable(t *testing.T) {
	id := NewID4("helo")
	if s := id.String(); s != "helo" {
		t.Errorf("ID4.String(): got %q, want %q", s, "helo")
	}
}

func TestID4_String_ShortTag(t *testing.T) {
	// Null bytes should be omitted from the string representation.
	id := NewID4("ok")
	if s := id.String(); s != "ok" {
		t.Errorf("ID4.String(): got %q, want %q", s, "ok")
	}
}

func TestID4_String_WithNewline(t *testing.T) {
	// "pcp\n" is the connect magic; newline should appear as \n.
	id := NewID4("pcp\n")
	want := "pcp\\n"
	if s := id.String(); s != want {
		t.Errorf("ID4.String(): got %q, want %q", s, want)
	}
}

func TestID4_String_NonPrintable(t *testing.T) {
	id := ID4{0x01, 0x02, 0x03, 0x04}
	s := id.String()
	want := `\x01\x02\x03\x04`
	if s != want {
		t.Errorf("ID4.String() non-printable: got %q, want %q", s, want)
	}
}

// ---------------------------------------------------------------------------
// GnuID
// ---------------------------------------------------------------------------

func TestGnuID_IsEmpty_Zero(t *testing.T) {
	var id GnuID
	if !id.IsEmpty() {
		t.Error("zero GnuID should be empty")
	}
}

func TestGnuID_IsEmpty_NonZero(t *testing.T) {
	var id GnuID
	id[0] = 1
	if id.IsEmpty() {
		t.Error("non-zero GnuID should not be empty")
	}
}

func TestGnuID_Clear(t *testing.T) {
	var id GnuID
	for i := range id {
		id[i] = byte(i + 1)
	}
	id.Clear()
	if !id.IsEmpty() {
		t.Error("GnuID after Clear() should be empty")
	}
}

func TestGnuID_String(t *testing.T) {
	var id GnuID
	for i := range id {
		id[i] = byte(i)
	}
	s := id.String()
	want := "000102030405060708090a0b0c0d0e0f"
	if s != want {
		t.Errorf("GnuID.String(): got %q, want %q", s, want)
	}
}

func TestGnuID_String_Zero(t *testing.T) {
	var id GnuID
	want := "00000000000000000000000000000000"
	if s := id.String(); s != want {
		t.Errorf("GnuID.String() zero: got %q, want %q", s, want)
	}
}

// ---------------------------------------------------------------------------
// Tag variables — verify key tags are well-formed
// ---------------------------------------------------------------------------

func TestTagVariables(t *testing.T) {
	tests := []struct {
		name    string
		id      ID4
		wantStr string
	}{
		{"PCPHelo", PCPHelo, "helo"},
		{"PCPOleh", PCPOleh, "oleh"},
		{"PCPQuit", PCPQuit, "quit"},
		{"PCPPing", PCPPing, "ping"},
		{"PCPPong", PCPPong, "pong"},
		{"PCPConnect", PCPConnect, "pcp\\n"},
		{"PCPOK", PCPOK, "ok"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.id.String(); got != tt.wantStr {
				t.Errorf("%s.String(): got %q, want %q", tt.name, got, tt.wantStr)
			}
		})
	}
}
