package pcp

import "testing"

func TestBroadcastState_InitPacketSettings(t *testing.T) {
	var chanID, bcid GnuID
	for i := range chanID {
		chanID[i] = byte(i + 1)
	}
	for i := range bcid {
		bcid[i] = byte(0xF0 + i)
	}

	bs := BroadcastState{
		ChanID:    chanID,
		BCID:      bcid,
		NumHops:   5,
		ForMe:     true,
		StreamPos: 12345,
		Group:     PCPBcstGroupRelays,
	}

	bs.InitPacketSettings()

	if bs.ForMe {
		t.Error("ForMe should be false after InitPacketSettings")
	}
	if bs.Group != 0 {
		t.Errorf("Group = %d, want 0", bs.Group)
	}
	if bs.NumHops != 0 {
		t.Errorf("NumHops = %d, want 0", bs.NumHops)
	}
	if !bs.BCID.IsEmpty() {
		t.Errorf("BCID should be cleared, got %v", bs.BCID)
	}
	if !bs.ChanID.IsEmpty() {
		t.Errorf("ChanID should be cleared, got %v", bs.ChanID)
	}
	// StreamPos is NOT reset by InitPacketSettings.
	if bs.StreamPos != 12345 {
		t.Errorf("StreamPos = %d, want 12345 (should be preserved)", bs.StreamPos)
	}
}

func TestBroadcastState_InitPacketSettings_AlreadyZero(t *testing.T) {
	var bs BroadcastState
	bs.InitPacketSettings()

	if bs.ForMe || bs.Group != 0 || bs.NumHops != 0 {
		t.Error("InitPacketSettings on zero state should remain zero")
	}
}
