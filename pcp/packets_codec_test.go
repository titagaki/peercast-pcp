package pcp

import (
	"bytes"
	"testing"
)

// ---------------------------------------------------------------------------
// HeloPacket
// ---------------------------------------------------------------------------

func TestHeloPacketBuildAndParse(t *testing.T) {
	var sid, bcid GnuID
	for i := range sid {
		sid[i] = byte(i + 1)
	}
	for i := range bcid {
		bcid[i] = byte(0xF0 + i)
	}

	h := HeloPacket{
		Agent:     "PeerCast/0.1218",
		OsType:    NewID4("lnux").Uint32(),
		SessionID: sid,
		Port:      7144,
		Ping:      7145,
		Pong:      7146,
		RemoteIP:  0xC0A80101,
		Version:   1218,
		BCID:      bcid,
		Disable:   0,
	}

	atom := h.BuildHeloAtom()
	if atom.Tag != PCPHelo {
		t.Fatalf("tag = %v, want PCPHelo", atom.Tag)
	}

	got, err := ParseHeloPacket(atom)
	if err != nil {
		t.Fatalf("ParseHeloPacket: %v", err)
	}
	if got != h {
		t.Errorf("ParseHeloPacket(BuildHeloAtom()):\ngot  %+v\nwant %+v", got, h)
	}
}

func TestHeloPacketBuildOleh(t *testing.T) {
	h := HeloPacket{Agent: "PeerCast/0.1218", Version: 1218}
	atom := h.BuildOlehAtom()
	if atom.Tag != PCPOleh {
		t.Fatalf("tag = %v, want PCPOleh", atom.Tag)
	}
	got, err := ParseHeloPacket(atom)
	if err != nil {
		t.Fatalf("ParseHeloPacket: %v", err)
	}
	if got != h {
		t.Errorf("roundtrip failed: got %+v, want %+v", got, h)
	}
}

func TestHeloPacketSkipsEmptyFields(t *testing.T) {
	h := HeloPacket{Agent: "Test", Version: 1218}
	atom := h.BuildHeloAtom()
	if n := atom.NumChildren(); n != 2 {
		t.Errorf("NumChildren = %d, want 2", n)
	}
}

func TestHeloPacketWireRoundtrip(t *testing.T) {
	var sid GnuID
	for i := range sid {
		sid[i] = byte(0x10 + i)
	}
	h := HeloPacket{
		Agent:     "TestAgent/1.0",
		SessionID: sid,
		Port:      7144,
		Version:   1218,
	}

	atom := h.BuildHeloAtom()
	var buf bytes.Buffer
	if err := atom.Write(&buf); err != nil {
		t.Fatal(err)
	}
	parsed, err := ReadAtom(&buf)
	if err != nil {
		t.Fatal(err)
	}
	got, err := ParseHeloPacket(parsed)
	if err != nil {
		t.Fatal(err)
	}
	if got != h {
		t.Errorf("wire roundtrip:\ngot  %+v\nwant %+v", got, h)
	}
}

// ---------------------------------------------------------------------------
// RootPacket
// ---------------------------------------------------------------------------

func TestRootPacketBuildAndParse(t *testing.T) {
	r := RootPacket{
		UpdateInterval: 120,
		CheckVersion:   1218,
		URL:            "http://example.com",
		Update:         []byte{0x01, 0x02},
		Next:           60,
	}

	atom := r.BuildAtom()
	if atom.Tag != PCPRoot {
		t.Fatalf("tag = %v, want PCPRoot", atom.Tag)
	}

	got, err := ParseRootPacket(atom)
	if err != nil {
		t.Fatalf("ParseRootPacket: %v", err)
	}
	if got.UpdateInterval != r.UpdateInterval || got.CheckVersion != r.CheckVersion ||
		got.URL != r.URL || got.Next != r.Next {
		t.Errorf("ParseRootPacket(BuildAtom()):\ngot  %+v\nwant %+v", got, r)
	}
	if !bytes.Equal(got.Update, r.Update) {
		t.Errorf("Update: got %v, want %v", got.Update, r.Update)
	}
}

func TestRootPacketSkipsEmptyFields(t *testing.T) {
	r := RootPacket{UpdateInterval: 120}
	atom := r.BuildAtom()
	if n := atom.NumChildren(); n != 1 {
		t.Errorf("NumChildren = %d, want 1", n)
	}
}

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

	got, err := ParseChanInfo(atom)
	if err != nil {
		t.Fatalf("ParseChanInfo: %v", err)
	}
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

	got, err := ParseChanInfo(atom)
	if err != nil {
		t.Fatalf("ParseChanInfo: %v", err)
	}
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

	got, err := ParseChanInfo(parsed)
	if err != nil {
		t.Fatalf("ParseChanInfo: %v", err)
	}
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

	got, err := ParseChanTrack(atom)
	if err != nil {
		t.Fatalf("ParseChanTrack: %v", err)
	}
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
// ChanPktData
// ---------------------------------------------------------------------------

func TestChanPktDataBuildAndParse(t *testing.T) {
	p := ChanPktData{
		Type:         NewID4("head"),
		Pos:          12345,
		Head:         []byte{0xDE, 0xAD},
		Data:         []byte{0xBE, 0xEF, 0x01},
		Meta:         []byte{0x42},
		Continuation: 1,
	}

	atom := p.BuildAtom()
	if atom.Tag != PCPChanPkt {
		t.Fatalf("tag = %v, want PCPChanPkt", atom.Tag)
	}

	got, err := ParseChanPktData(atom)
	if err != nil {
		t.Fatalf("ParseChanPktData: %v", err)
	}
	if got.Type != p.Type || got.Pos != p.Pos || got.Continuation != p.Continuation {
		t.Errorf("ParseChanPktData: scalar mismatch:\ngot  %+v\nwant %+v", got, p)
	}
	if !bytes.Equal(got.Head, p.Head) || !bytes.Equal(got.Data, p.Data) || !bytes.Equal(got.Meta, p.Meta) {
		t.Errorf("ParseChanPktData: bytes mismatch")
	}
}

func TestChanPktDataSkipsEmptyFields(t *testing.T) {
	p := ChanPktData{Pos: 100}
	atom := p.BuildAtom()
	if n := atom.NumChildren(); n != 1 {
		t.Errorf("NumChildren = %d, want 1", n)
	}
}

// ---------------------------------------------------------------------------
// ChanPacket
// ---------------------------------------------------------------------------

func TestChanPacketBuildAndParse(t *testing.T) {
	var cid GnuID
	for i := range cid {
		cid[i] = byte(0xA0 + i)
	}

	info := &ChanInfo{Name: "TestChan", Genre: "Variety", Bitrate: 1000}
	track := &ChanTrack{Title: "Song", Creator: "Artist"}

	cp := ChanPacket{
		ID:   cid,
		Info: info,
		Track: track,
	}

	atom := cp.BuildAtom()
	if atom.Tag != PCPChan {
		t.Fatalf("tag = %v, want PCPChan", atom.Tag)
	}

	got, err := ParseChanPacket(atom)
	if err != nil {
		t.Fatalf("ParseChanPacket: %v", err)
	}
	if got.ID != cp.ID {
		t.Errorf("ID: got %v, want %v", got.ID, cp.ID)
	}
	if got.Info == nil || got.Info.Name != "TestChan" {
		t.Errorf("Info: got %+v, want %+v", got.Info, info)
	}
	if got.Track == nil || got.Track.Title != "Song" {
		t.Errorf("Track: got %+v, want %+v", got.Track, track)
	}
}

func TestChanPacketWireRoundtrip(t *testing.T) {
	var cid GnuID
	for i := range cid {
		cid[i] = byte(i)
	}

	cp := ChanPacket{
		ID:   cid,
		Info: &ChanInfo{Name: "Wire", Type: "FLV", Bitrate: 500},
	}

	atom := cp.BuildAtom()
	var buf bytes.Buffer
	if err := atom.Write(&buf); err != nil {
		t.Fatal(err)
	}
	parsed, err := ReadAtom(&buf)
	if err != nil {
		t.Fatal(err)
	}
	got, err := ParseChanPacket(parsed)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != cid {
		t.Errorf("ID mismatch")
	}
	if got.Info == nil || got.Info.Name != "Wire" {
		t.Errorf("Info mismatch: %+v", got.Info)
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

	got, err := ParseHostPacket(atom)
	if err != nil {
		t.Fatalf("ParseHostPacket: %v", err)
	}
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
	got, err := ParseHostPacket(atom)
	if err != nil {
		t.Fatalf("ParseHostPacket: %v", err)
	}
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
	got, err := ParseHostPacket(atom)
	if err != nil {
		t.Fatalf("ParseHostPacket: %v", err)
	}
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

	got, err := ParseHostPacket(parsed)
	if err != nil {
		t.Fatalf("ParseHostPacket: %v", err)
	}
	if got != h {
		t.Errorf("wire roundtrip:\ngot  %+v\nwant %+v", got, h)
	}
}

// ---------------------------------------------------------------------------
// BcstPacket
// ---------------------------------------------------------------------------

func TestBcstPacketBuildAndParse(t *testing.T) {
	var from, dest, cid GnuID
	for i := range from {
		from[i] = byte(i + 1)
	}
	for i := range cid {
		cid[i] = byte(0xA0 + i)
	}

	b := BcstPacket{
		TTL:             7,
		Hops:            3,
		From:            from,
		Dest:            dest, // empty = broadcast
		Group:           PCPBcstGroupRelays,
		ChanID:          cid,
		Version:         1218,
		VersionVP:       27,
		VersionExPrefix: [2]byte{'S', 'T'},
		VersionExNumber: 5,
	}

	atom := b.BuildAtom()
	if atom.Tag != PCPBcst {
		t.Fatalf("tag = %v, want PCPBcst", atom.Tag)
	}

	got, err := ParseBcstPacket(atom)
	if err != nil {
		t.Fatalf("ParseBcstPacket: %v", err)
	}
	if got != b {
		t.Errorf("ParseBcstPacket(BuildAtom()):\ngot  %+v\nwant %+v", got, b)
	}
}

func TestBcstPacketSkipsEmptyFields(t *testing.T) {
	b := BcstPacket{TTL: 7, Group: PCPBcstGroupAll}
	atom := b.BuildAtom()
	if n := atom.NumChildren(); n != 2 {
		t.Errorf("NumChildren = %d, want 2", n)
	}
}

func TestBcstPacketWireRoundtrip(t *testing.T) {
	var from, cid GnuID
	for i := range from {
		from[i] = byte(0x30 + i)
	}
	for i := range cid {
		cid[i] = byte(0xB0 + i)
	}

	b := BcstPacket{
		TTL:     7,
		Hops:    1,
		From:    from,
		Group:   PCPBcstGroupRelays,
		ChanID:  cid,
		Version: 1218,
	}

	atom := b.BuildAtom()
	var buf bytes.Buffer
	if err := atom.Write(&buf); err != nil {
		t.Fatal(err)
	}
	parsed, err := ReadAtom(&buf)
	if err != nil {
		t.Fatal(err)
	}
	got, err := ParseBcstPacket(parsed)
	if err != nil {
		t.Fatal(err)
	}
	if got != b {
		t.Errorf("wire roundtrip:\ngot  %+v\nwant %+v", got, b)
	}
}

// ---------------------------------------------------------------------------
// PushPacket
// ---------------------------------------------------------------------------

func TestPushPacketBuildAndParse(t *testing.T) {
	var cid GnuID
	for i := range cid {
		cid[i] = byte(0xC0 + i)
	}

	p := PushPacket{
		IP:     0xC0A80101,
		Port:   7144,
		ChanID: cid,
	}

	atom := p.BuildAtom()
	if atom.Tag != PCPPush {
		t.Fatalf("tag = %v, want PCPPush", atom.Tag)
	}

	got, err := ParsePushPacket(atom)
	if err != nil {
		t.Fatalf("ParsePushPacket: %v", err)
	}
	if got != p {
		t.Errorf("ParsePushPacket(BuildAtom()):\ngot  %+v\nwant %+v", got, p)
	}
}

func TestPushPacketWireRoundtrip(t *testing.T) {
	var cid GnuID
	for i := range cid {
		cid[i] = byte(i)
	}
	p := PushPacket{IP: 0x0A000001, Port: 8144, ChanID: cid}

	atom := p.BuildAtom()
	var buf bytes.Buffer
	if err := atom.Write(&buf); err != nil {
		t.Fatal(err)
	}
	parsed, err := ReadAtom(&buf)
	if err != nil {
		t.Fatal(err)
	}
	got, err := ParsePushPacket(parsed)
	if err != nil {
		t.Fatal(err)
	}
	if got != p {
		t.Errorf("wire roundtrip:\ngot  %+v\nwant %+v", got, p)
	}
}

// ---------------------------------------------------------------------------
// GetPacket
// ---------------------------------------------------------------------------

func TestGetPacketBuildAndParse(t *testing.T) {
	var id GnuID
	for i := range id {
		id[i] = byte(0xD0 + i)
	}

	g := GetPacket{ID: id, Name: "channel-name"}

	atom := g.BuildAtom()
	if atom.Tag != PCPGet {
		t.Fatalf("tag = %v, want PCPGet", atom.Tag)
	}

	got, err := ParseGetPacket(atom)
	if err != nil {
		t.Fatalf("ParseGetPacket: %v", err)
	}
	if got != g {
		t.Errorf("ParseGetPacket(BuildAtom()):\ngot  %+v\nwant %+v", got, g)
	}
}

func TestGetPacketSkipsEmptyFields(t *testing.T) {
	g := GetPacket{Name: "test"}
	atom := g.BuildAtom()
	if n := atom.NumChildren(); n != 1 {
		t.Errorf("NumChildren = %d, want 1", n)
	}
}

// ---------------------------------------------------------------------------
// MesgPacket
// ---------------------------------------------------------------------------

func TestMesgPacketBuildAndParse(t *testing.T) {
	m := MesgPacket{ASCII: "hello", SJIS: "world"}

	atom := m.BuildAtom()
	if atom.Tag != PCPMesg {
		t.Fatalf("tag = %v, want PCPMesg", atom.Tag)
	}

	got, err := ParseMesgPacket(atom)
	if err != nil {
		t.Fatalf("ParseMesgPacket: %v", err)
	}
	if got != m {
		t.Errorf("ParseMesgPacket(BuildAtom()):\ngot  %+v\nwant %+v", got, m)
	}
}

func TestMesgPacketSkipsEmptyFields(t *testing.T) {
	m := MesgPacket{ASCII: "hello"}
	atom := m.BuildAtom()
	if n := atom.NumChildren(); n != 1 {
		t.Errorf("NumChildren = %d, want 1", n)
	}
}

// ---------------------------------------------------------------------------
// Parse error propagation
// ---------------------------------------------------------------------------

func TestParseHostPacket_MalformedField(t *testing.T) {
	// Port field with wrong size (4 bytes instead of 2).
	atom := NewParentAtom(PCPHost,
		NewIntAtom(PCPHostPort, 7144), // INT instead of SHORT
	)
	_, err := ParseHostPacket(atom)
	if err == nil {
		t.Fatal("expected error for malformed port field")
	}
}

func TestParseBcstPacket_MalformedField(t *testing.T) {
	// TTL field with wrong size (4 bytes instead of 1).
	atom := NewParentAtom(PCPBcst,
		NewIntAtom(PCPBcstTTL, 7), // INT instead of BYTE
	)
	_, err := ParseBcstPacket(atom)
	if err == nil {
		t.Fatal("expected error for malformed ttl field")
	}
}

func TestParseHeloPacket_MalformedField(t *testing.T) {
	// Version field with wrong size (2 bytes instead of 4).
	atom := NewParentAtom(PCPHelo,
		NewShortAtom(PCPHeloVersion, 1218), // SHORT instead of INT
	)
	_, err := ParseHeloPacket(atom)
	if err == nil {
		t.Fatal("expected error for malformed version field")
	}
}
