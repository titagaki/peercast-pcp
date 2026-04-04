package pcp

import (
	"bytes"
	"testing"
)

// ---------------------------------------------------------------------------
// ChanInfo
// ---------------------------------------------------------------------------

func TestChanInfoBuildAndParse(t *testing.T) {
	ci := ChanInfo{
		Name:       "Test Channel",
		URL:        "http://example.com",
		Desc:       "A test channel",
		Comment:    "hello",
		Genre:      "Variety",
		Type:       "FLV",
		StreamType: "video/x-flv",
		StreamExt:  ".flv",
		Bitrate:    1000,
	}

	atom := ci.BuildAtom()
	if atom.Tag != PCPChanInfo {
		t.Fatalf("tag = %v, want PCPChanInfo", atom.Tag)
	}

	got := ParseChanInfo(atom)
	if got != ci {
		t.Errorf("ParseChanInfo(BuildAtom()) = %+v, want %+v", got, ci)
	}
}

func TestChanInfoBuildSkipsEmptyFields(t *testing.T) {
	ci := ChanInfo{Name: "OnlyName", Bitrate: 500}
	atom := ci.BuildAtom()

	// Should only have Name + Bitrate children.
	if n := atom.NumChildren(); n != 2 {
		t.Errorf("NumChildren = %d, want 2", n)
	}

	got := ParseChanInfo(atom)
	if got.Name != "OnlyName" || got.Bitrate != 500 {
		t.Errorf("unexpected: %+v", got)
	}
}

func TestChanInfoRoundtripWire(t *testing.T) {
	ci := ChanInfo{
		Name:    "Wire Test",
		Genre:   "Tech",
		Type:    "MKV",
		Bitrate: 2500,
	}

	atom := ci.BuildAtom()
	var buf bytes.Buffer
	if err := atom.Write(&buf); err != nil {
		t.Fatal(err)
	}

	parsed, err := ReadAtom(&buf)
	if err != nil {
		t.Fatal(err)
	}

	got := ParseChanInfo(parsed)
	if got != ci {
		t.Errorf("wire roundtrip: got %+v, want %+v", got, ci)
	}
}

// ---------------------------------------------------------------------------
// ChanTrack
// ---------------------------------------------------------------------------

func TestChanTrackBuildAndParse(t *testing.T) {
	ct := ChanTrack{
		Title:   "Song Title",
		Creator: "Artist Name",
		URL:     "http://example.com/track",
		Album:   "Album Name",
	}

	atom := ct.BuildAtom()
	if atom.Tag != PCPChanTrack {
		t.Fatalf("tag = %v, want PCPChanTrack", atom.Tag)
	}

	got := ParseChanTrack(atom)
	if got != ct {
		t.Errorf("ParseChanTrack(BuildAtom()) = %+v, want %+v", got, ct)
	}
}

func TestChanTrackBuildSkipsEmptyFields(t *testing.T) {
	ct := ChanTrack{Title: "Only Title"}
	atom := ct.BuildAtom()
	if n := atom.NumChildren(); n != 1 {
		t.Errorf("NumChildren = %d, want 1", n)
	}
}

// ---------------------------------------------------------------------------
// HostPacket
// ---------------------------------------------------------------------------

func TestHostPacketBuildAndParse(t *testing.T) {
	var sid, cid GnuID
	for i := range sid {
		sid[i] = byte(i)
	}
	for i := range cid {
		cid[i] = byte(0xA0 + i)
	}

	h := HostPacket{
		ID:              sid,
		IP:              0xC0A80101, // 192.168.1.1
		Port:            7144,
		NumListeners:    5,
		NumRelays:       3,
		Uptime:          3600,
		ChanID:          cid,
		Flags1:          PCPHostFlags1Relay | PCPHostFlags1Direct | PCPHostFlags1Recv | PCPHostFlags1CIN,
		Version:         1218,
		VersionVP:       27,
		VersionExPrefix: [2]byte{'M', 'I'},
		VersionExNumber: 1,
		OldPos:          1000,
		NewPos:          5000,
	}

	atom := h.BuildAtom()
	if atom.Tag != PCPHost {
		t.Fatalf("tag = %v, want PCPHost", atom.Tag)
	}

	got := ParseHostPacket(atom)
	if got != h {
		t.Errorf("ParseHostPacket(BuildAtom()):\ngot  %+v\nwant %+v", got, h)
	}
}

func TestHostPacketWithTracker(t *testing.T) {
	h := HostPacket{
		Tracker: 1,
		Version: 1218,
	}

	atom := h.BuildAtom()
	got := ParseHostPacket(atom)
	if got.Tracker != 1 {
		t.Errorf("Tracker = %d, want 1", got.Tracker)
	}
}

func TestHostPacketWithUphost(t *testing.T) {
	h := HostPacket{
		UphostIP:   0x0A000001,
		UphostPort: 7144,
		UphostHops: 2,
		Version:    1218,
	}

	atom := h.BuildAtom()
	got := ParseHostPacket(atom)
	if got.UphostIP != h.UphostIP || got.UphostPort != h.UphostPort || got.UphostHops != h.UphostHops {
		t.Errorf("uphost: got IP=0x%08X port=%d hops=%d, want IP=0x%08X port=%d hops=%d",
			got.UphostIP, got.UphostPort, got.UphostHops,
			h.UphostIP, h.UphostPort, h.UphostHops)
	}
}

func TestHostPacketWireRoundtrip(t *testing.T) {
	var sid GnuID
	for i := range sid {
		sid[i] = byte(0x10 + i)
	}

	h := HostPacket{
		ID:              sid,
		IP:              0xAC100001,
		Port:            7144,
		NumListeners:    2,
		NumRelays:       1,
		Uptime:          120,
		Flags1:          PCPHostFlags1Relay | PCPHostFlags1Recv,
		Version:         1218,
		VersionVP:       27,
		VersionExPrefix: [2]byte{'Y', 'P'},
		VersionExNumber: 1,
	}

	atom := h.BuildAtom()
	var buf bytes.Buffer
	if err := atom.Write(&buf); err != nil {
		t.Fatal(err)
	}

	parsed, err := ReadAtom(&buf)
	if err != nil {
		t.Fatal(err)
	}

	got := ParseHostPacket(parsed)
	if got != h {
		t.Errorf("wire roundtrip:\ngot  %+v\nwant %+v", got, h)
	}
}
