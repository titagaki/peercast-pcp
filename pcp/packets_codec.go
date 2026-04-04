package pcp

// packets_codec.go provides Build/Parse functions that convert between
// the typed packet structs (packets.go) and the Atom tree representation.

// ---------------------------------------------------------------------------
// ChanInfo
// ---------------------------------------------------------------------------

// BuildAtom serializes a ChanInfo into a PCPChanInfo parent atom.
func (ci *ChanInfo) BuildAtom() *Atom {
	children := make([]*Atom, 0, 9)
	if ci.Name != "" {
		children = append(children, NewStringAtom(PCPChanInfoName, ci.Name))
	}
	if ci.URL != "" {
		children = append(children, NewStringAtom(PCPChanInfoURL, ci.URL))
	}
	if ci.Desc != "" {
		children = append(children, NewStringAtom(PCPChanInfoDesc, ci.Desc))
	}
	if ci.Comment != "" {
		children = append(children, NewStringAtom(PCPChanInfoComment, ci.Comment))
	}
	if ci.Genre != "" {
		children = append(children, NewStringAtom(PCPChanInfoGenre, ci.Genre))
	}
	if ci.Type != "" {
		children = append(children, NewStringAtom(PCPChanInfoType, ci.Type))
	}
	if ci.StreamType != "" {
		children = append(children, NewStringAtom(PCPChanInfoStreamType, ci.StreamType))
	}
	if ci.StreamExt != "" {
		children = append(children, NewStringAtom(PCPChanInfoStreamExt, ci.StreamExt))
	}
	if ci.Bitrate != 0 {
		children = append(children, NewIntAtom(PCPChanInfoBitrate, ci.Bitrate))
	}
	return NewParentAtom(PCPChanInfo, children...)
}

// ParseChanInfo extracts a ChanInfo from a PCPChanInfo parent atom.
func ParseChanInfo(a *Atom) ChanInfo {
	var ci ChanInfo
	for _, child := range a.Children() {
		switch child.Tag {
		case PCPChanInfoName:
			ci.Name = child.GetString()
		case PCPChanInfoURL:
			ci.URL = child.GetString()
		case PCPChanInfoDesc:
			ci.Desc = child.GetString()
		case PCPChanInfoComment:
			ci.Comment = child.GetString()
		case PCPChanInfoGenre:
			ci.Genre = child.GetString()
		case PCPChanInfoType:
			ci.Type = child.GetString()
		case PCPChanInfoStreamType:
			ci.StreamType = child.GetString()
		case PCPChanInfoStreamExt:
			ci.StreamExt = child.GetString()
		case PCPChanInfoBitrate:
			ci.Bitrate, _ = child.GetInt()
		}
	}
	return ci
}

// ---------------------------------------------------------------------------
// ChanTrack
// ---------------------------------------------------------------------------

// BuildAtom serializes a ChanTrack into a PCPChanTrack parent atom.
func (ct *ChanTrack) BuildAtom() *Atom {
	children := make([]*Atom, 0, 4)
	if ct.Title != "" {
		children = append(children, NewStringAtom(PCPChanTrackTitle, ct.Title))
	}
	if ct.Creator != "" {
		children = append(children, NewStringAtom(PCPChanTrackCreator, ct.Creator))
	}
	if ct.URL != "" {
		children = append(children, NewStringAtom(PCPChanTrackURL, ct.URL))
	}
	if ct.Album != "" {
		children = append(children, NewStringAtom(PCPChanTrackAlbum, ct.Album))
	}
	return NewParentAtom(PCPChanTrack, children...)
}

// ParseChanTrack extracts a ChanTrack from a PCPChanTrack parent atom.
func ParseChanTrack(a *Atom) ChanTrack {
	var ct ChanTrack
	for _, child := range a.Children() {
		switch child.Tag {
		case PCPChanTrackTitle:
			ct.Title = child.GetString()
		case PCPChanTrackCreator:
			ct.Creator = child.GetString()
		case PCPChanTrackURL:
			ct.URL = child.GetString()
		case PCPChanTrackAlbum:
			ct.Album = child.GetString()
		}
	}
	return ct
}

// ---------------------------------------------------------------------------
// HostPacket
// ---------------------------------------------------------------------------

// BuildAtom serializes a HostPacket into a PCPHost parent atom.
func (h *HostPacket) BuildAtom() *Atom {
	children := []*Atom{
		NewIDAtom(PCPHostID, h.ID),
		NewIntAtom(PCPHostIP, h.IP),
		NewShortAtom(PCPHostPort, h.Port),
		NewIntAtom(PCPHostNumListeners, h.NumListeners),
		NewIntAtom(PCPHostNumRelays, h.NumRelays),
		NewIntAtom(PCPHostUptime, h.Uptime),
		NewIntAtom(PCPHostOldPos, h.OldPos),
		NewIntAtom(PCPHostNewPos, h.NewPos),
		NewIDAtom(PCPHostChanID, h.ChanID),
		NewByteAtom(PCPHostFlags1, h.Flags1),
		NewIntAtom(PCPHostVersion, h.Version),
		NewIntAtom(PCPHostVersionVP, h.VersionVP),
		NewBytesAtom(PCPHostVersionExPrefix, h.VersionExPrefix[:]),
		NewShortAtom(PCPHostVersionExNumber, h.VersionExNumber),
	}

	if h.Tracker != 0 {
		children = append(children, NewIntAtom(PCPHostTracker, h.Tracker))
	}

	if h.UphostIP != 0 || h.UphostPort != 0 {
		children = append(children,
			NewIntAtom(PCPHostUphostIP, h.UphostIP),
			NewIntAtom(PCPHostUphostPort, h.UphostPort),
		)
		if h.UphostHops != 0 {
			children = append(children, NewIntAtom(PCPHostUphostHops, h.UphostHops))
		}
	}

	return NewParentAtom(PCPHost, children...)
}

// ParseHostPacket extracts a HostPacket from a PCPHost parent atom.
func ParseHostPacket(a *Atom) HostPacket {
	var h HostPacket
	for _, child := range a.Children() {
		switch child.Tag {
		case PCPHostID:
			h.ID, _ = child.GetID()
		case PCPHostIP:
			h.IP, _ = child.GetInt()
		case PCPHostPort:
			h.Port, _ = child.GetShort()
		case PCPHostNumListeners:
			h.NumListeners, _ = child.GetInt()
		case PCPHostNumRelays:
			h.NumRelays, _ = child.GetInt()
		case PCPHostUptime:
			h.Uptime, _ = child.GetInt()
		case PCPHostTracker:
			h.Tracker, _ = child.GetInt()
		case PCPHostChanID:
			h.ChanID, _ = child.GetID()
		case PCPHostVersion:
			h.Version, _ = child.GetInt()
		case PCPHostVersionVP:
			h.VersionVP, _ = child.GetInt()
		case PCPHostVersionExPrefix:
			if d := child.Data(); len(d) >= 2 {
				copy(h.VersionExPrefix[:], d[:2])
			}
		case PCPHostVersionExNumber:
			h.VersionExNumber, _ = child.GetShort()
		case PCPHostFlags1:
			h.Flags1, _ = child.GetByte()
		case PCPHostOldPos:
			h.OldPos, _ = child.GetInt()
		case PCPHostNewPos:
			h.NewPos, _ = child.GetInt()
		case PCPHostUphostIP:
			h.UphostIP, _ = child.GetInt()
		case PCPHostUphostPort:
			h.UphostPort, _ = child.GetInt()
		case PCPHostUphostHops:
			h.UphostHops, _ = child.GetInt()
		}
	}
	return h
}
