package pcp

// packets_codec.go provides Build/Parse functions that convert between
// the typed packet structs (packets.go) and the Atom tree representation.

import "fmt"

// ---------------------------------------------------------------------------
// HeloPacket (used for both HELO and OLEH)
// ---------------------------------------------------------------------------

func (h *HeloPacket) buildAtom(tag ID4) *Atom {
	var children []*Atom
	if h.Agent != "" {
		children = append(children, NewStringAtom(PCPHeloAgent, h.Agent))
	}
	if h.OsType != 0 {
		children = append(children, NewIntAtom(PCPHeloOsType, h.OsType))
	}
	if !h.SessionID.IsEmpty() {
		children = append(children, NewIDAtom(PCPHeloSessionID, h.SessionID))
	}
	if h.Port != 0 {
		children = append(children, NewShortAtom(PCPHeloPort, h.Port))
	}
	if h.Ping != 0 {
		children = append(children, NewShortAtom(PCPHeloPing, h.Ping))
	}
	if h.Pong != 0 {
		children = append(children, NewShortAtom(PCPHeloPong, h.Pong))
	}
	if h.RemoteIP != 0 {
		children = append(children, NewIntAtom(PCPHeloRemoteIP, h.RemoteIP))
	}
	if h.Version != 0 {
		children = append(children, NewIntAtom(PCPHeloVersion, h.Version))
	}
	if !h.BCID.IsEmpty() {
		children = append(children, NewIDAtom(PCPHeloBCID, h.BCID))
	}
	if h.Disable != 0 {
		children = append(children, NewIntAtom(PCPHeloDisable, h.Disable))
	}
	return NewParentAtom(tag, children...)
}

// BuildHeloAtom serializes a HeloPacket into a PCPHelo parent atom.
func (h *HeloPacket) BuildHeloAtom() *Atom { return h.buildAtom(PCPHelo) }

// BuildOlehAtom serializes a HeloPacket into a PCPOleh parent atom.
func (h *HeloPacket) BuildOlehAtom() *Atom { return h.buildAtom(PCPOleh) }

// ParseHeloPacket extracts a HeloPacket from a PCPHelo or PCPOleh parent atom.
func ParseHeloPacket(a *Atom) (HeloPacket, error) {
	var h HeloPacket
	for _, child := range a.Children() {
		switch child.Tag {
		case PCPHeloAgent:
			h.Agent = child.GetString()
		case PCPHeloOsType:
			v, err := child.GetInt()
			if err != nil {
				return h, fmt.Errorf("pcp: parsing helo ostype: %w", err)
			}
			h.OsType = v
		case PCPHeloSessionID:
			v, err := child.GetID()
			if err != nil {
				return h, fmt.Errorf("pcp: parsing helo session id: %w", err)
			}
			h.SessionID = v
		case PCPHeloPort:
			v, err := child.GetShort()
			if err != nil {
				return h, fmt.Errorf("pcp: parsing helo port: %w", err)
			}
			h.Port = v
		case PCPHeloPing:
			v, err := child.GetShort()
			if err != nil {
				return h, fmt.Errorf("pcp: parsing helo ping: %w", err)
			}
			h.Ping = v
		case PCPHeloPong:
			v, err := child.GetShort()
			if err != nil {
				return h, fmt.Errorf("pcp: parsing helo pong: %w", err)
			}
			h.Pong = v
		case PCPHeloRemoteIP:
			v, err := child.GetInt()
			if err != nil {
				return h, fmt.Errorf("pcp: parsing helo remote ip: %w", err)
			}
			h.RemoteIP = v
		case PCPHeloVersion:
			v, err := child.GetInt()
			if err != nil {
				return h, fmt.Errorf("pcp: parsing helo version: %w", err)
			}
			h.Version = v
		case PCPHeloBCID:
			v, err := child.GetID()
			if err != nil {
				return h, fmt.Errorf("pcp: parsing helo bcid: %w", err)
			}
			h.BCID = v
		case PCPHeloDisable:
			v, err := child.GetInt()
			if err != nil {
				return h, fmt.Errorf("pcp: parsing helo disable: %w", err)
			}
			h.Disable = v
		}
	}
	return h, nil
}

// ---------------------------------------------------------------------------
// RootPacket
// ---------------------------------------------------------------------------

// BuildAtom serializes a RootPacket into a PCPRoot parent atom.
func (r *RootPacket) BuildAtom() *Atom {
	var children []*Atom
	if r.UpdateInterval != 0 {
		children = append(children, NewIntAtom(PCPRootUpdInt, r.UpdateInterval))
	}
	if r.CheckVersion != 0 {
		children = append(children, NewIntAtom(PCPRootCheckVer, r.CheckVersion))
	}
	if r.URL != "" {
		children = append(children, NewStringAtom(PCPRootURL, r.URL))
	}
	if len(r.Update) > 0 {
		children = append(children, NewBytesAtom(PCPRootUpdate, r.Update))
	}
	if r.Next != 0 {
		children = append(children, NewIntAtom(PCPRootNext, r.Next))
	}
	return NewParentAtom(PCPRoot, children...)
}

// ParseRootPacket extracts a RootPacket from a PCPRoot parent atom.
func ParseRootPacket(a *Atom) (RootPacket, error) {
	var r RootPacket
	for _, child := range a.Children() {
		switch child.Tag {
		case PCPRootUpdInt:
			v, err := child.GetInt()
			if err != nil {
				return r, fmt.Errorf("pcp: parsing root update interval: %w", err)
			}
			r.UpdateInterval = v
		case PCPRootCheckVer:
			v, err := child.GetInt()
			if err != nil {
				return r, fmt.Errorf("pcp: parsing root check version: %w", err)
			}
			r.CheckVersion = v
		case PCPRootURL:
			r.URL = child.GetString()
		case PCPRootUpdate:
			r.Update = child.Data()
		case PCPRootNext:
			v, err := child.GetInt()
			if err != nil {
				return r, fmt.Errorf("pcp: parsing root next: %w", err)
			}
			r.Next = v
		}
	}
	return r, nil
}

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
func ParseChanInfo(a *Atom) (ChanInfo, error) {
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
			v, err := child.GetInt()
			if err != nil {
				return ci, fmt.Errorf("pcp: parsing chan info bitrate: %w", err)
			}
			ci.Bitrate = v
		}
	}
	return ci, nil
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
func ParseChanTrack(a *Atom) (ChanTrack, error) {
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
	return ct, nil
}

// ---------------------------------------------------------------------------
// ChanPktData
// ---------------------------------------------------------------------------

// BuildAtom serializes a ChanPktData into a PCPChanPkt parent atom.
func (p *ChanPktData) BuildAtom() *Atom {
	var children []*Atom
	var emptyID4 ID4
	if p.Type != emptyID4 {
		children = append(children, NewID4Atom(PCPChanPktType, p.Type))
	}
	if p.Pos != 0 {
		children = append(children, NewIntAtom(PCPChanPktPos, p.Pos))
	}
	if len(p.Data) > 0 {
		children = append(children, NewBytesAtom(PCPChanPktData, p.Data))
	}
	if p.Continuation != 0 {
		children = append(children, NewByteAtom(PCPChanPktContinuation, p.Continuation))
	}
	return NewParentAtom(PCPChanPkt, children...)
}

// ParseChanPktData extracts a ChanPktData from a PCPChanPkt parent atom.
func ParseChanPktData(a *Atom) (ChanPktData, error) {
	var p ChanPktData
	for _, child := range a.Children() {
		switch child.Tag {
		case PCPChanPktType:
			v, err := child.GetID4()
			if err != nil {
				return p, fmt.Errorf("pcp: parsing chan pkt type: %w", err)
			}
			p.Type = v
		case PCPChanPktPos:
			v, err := child.GetInt()
			if err != nil {
				return p, fmt.Errorf("pcp: parsing chan pkt pos: %w", err)
			}
			p.Pos = v
		case PCPChanPktData:
			p.Data = child.Data()
		case PCPChanPktContinuation:
			v, err := child.GetByte()
			if err != nil {
				return p, fmt.Errorf("pcp: parsing chan pkt continuation: %w", err)
			}
			p.Continuation = v
		}
	}
	return p, nil
}

// ---------------------------------------------------------------------------
// ChanPacket
// ---------------------------------------------------------------------------

// BuildAtom serializes a ChanPacket into a PCPChan parent atom.
func (c *ChanPacket) BuildAtom() *Atom {
	var children []*Atom
	if !c.ID.IsEmpty() {
		children = append(children, NewIDAtom(PCPChanID, c.ID))
	}
	if !c.BroadcastID.IsEmpty() {
		children = append(children, NewIDAtom(PCPChanBCID, c.BroadcastID))
	}
	if len(c.Key) > 0 {
		children = append(children, NewBytesAtom(PCPChanKey, c.Key))
	}
	if c.Info != nil {
		children = append(children, c.Info.BuildAtom())
	}
	if c.Track != nil {
		children = append(children, c.Track.BuildAtom())
	}
	if c.Pkt != nil {
		children = append(children, c.Pkt.BuildAtom())
	}
	return NewParentAtom(PCPChan, children...)
}

// ParseChanPacket extracts a ChanPacket from a PCPChan parent atom.
func ParseChanPacket(a *Atom) (ChanPacket, error) {
	var c ChanPacket
	for _, child := range a.Children() {
		switch child.Tag {
		case PCPChanID:
			v, err := child.GetID()
			if err != nil {
				return c, fmt.Errorf("pcp: parsing chan id: %w", err)
			}
			c.ID = v
		case PCPChanBCID:
			v, err := child.GetID()
			if err != nil {
				return c, fmt.Errorf("pcp: parsing chan bcid: %w", err)
			}
			c.BroadcastID = v
		case PCPChanKey:
			c.Key = child.Data()
		case PCPChanInfo:
			v, err := ParseChanInfo(child)
			if err != nil {
				return c, fmt.Errorf("pcp: parsing chan info: %w", err)
			}
			c.Info = &v
		case PCPChanTrack:
			v, err := ParseChanTrack(child)
			if err != nil {
				return c, fmt.Errorf("pcp: parsing chan track: %w", err)
			}
			c.Track = &v
		case PCPChanPkt:
			v, err := ParseChanPktData(child)
			if err != nil {
				return c, fmt.Errorf("pcp: parsing chan pkt: %w", err)
			}
			c.Pkt = &v
		}
	}
	return c, nil
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
func ParseHostPacket(a *Atom) (HostPacket, error) {
	var h HostPacket
	for _, child := range a.Children() {
		switch child.Tag {
		case PCPHostID:
			v, err := child.GetID()
			if err != nil {
				return h, fmt.Errorf("pcp: parsing host id: %w", err)
			}
			h.ID = v
		case PCPHostIP:
			v, err := child.GetInt()
			if err != nil {
				return h, fmt.Errorf("pcp: parsing host ip: %w", err)
			}
			h.IP = v
		case PCPHostPort:
			v, err := child.GetShort()
			if err != nil {
				return h, fmt.Errorf("pcp: parsing host port: %w", err)
			}
			h.Port = v
		case PCPHostNumListeners:
			v, err := child.GetInt()
			if err != nil {
				return h, fmt.Errorf("pcp: parsing host num listeners: %w", err)
			}
			h.NumListeners = v
		case PCPHostNumRelays:
			v, err := child.GetInt()
			if err != nil {
				return h, fmt.Errorf("pcp: parsing host num relays: %w", err)
			}
			h.NumRelays = v
		case PCPHostUptime:
			v, err := child.GetInt()
			if err != nil {
				return h, fmt.Errorf("pcp: parsing host uptime: %w", err)
			}
			h.Uptime = v
		case PCPHostTracker:
			v, err := child.GetInt()
			if err != nil {
				return h, fmt.Errorf("pcp: parsing host tracker: %w", err)
			}
			h.Tracker = v
		case PCPHostChanID:
			v, err := child.GetID()
			if err != nil {
				return h, fmt.Errorf("pcp: parsing host chan id: %w", err)
			}
			h.ChanID = v
		case PCPHostVersion:
			v, err := child.GetInt()
			if err != nil {
				return h, fmt.Errorf("pcp: parsing host version: %w", err)
			}
			h.Version = v
		case PCPHostVersionVP:
			v, err := child.GetInt()
			if err != nil {
				return h, fmt.Errorf("pcp: parsing host version vp: %w", err)
			}
			h.VersionVP = v
		case PCPHostVersionExPrefix:
			if d := child.Data(); len(d) >= 2 {
				copy(h.VersionExPrefix[:], d[:2])
			}
		case PCPHostVersionExNumber:
			v, err := child.GetShort()
			if err != nil {
				return h, fmt.Errorf("pcp: parsing host version ex number: %w", err)
			}
			h.VersionExNumber = v
		case PCPHostFlags1:
			v, err := child.GetByte()
			if err != nil {
				return h, fmt.Errorf("pcp: parsing host flags1: %w", err)
			}
			h.Flags1 = v
		case PCPHostOldPos:
			v, err := child.GetInt()
			if err != nil {
				return h, fmt.Errorf("pcp: parsing host old pos: %w", err)
			}
			h.OldPos = v
		case PCPHostNewPos:
			v, err := child.GetInt()
			if err != nil {
				return h, fmt.Errorf("pcp: parsing host new pos: %w", err)
			}
			h.NewPos = v
		case PCPHostUphostIP:
			v, err := child.GetInt()
			if err != nil {
				return h, fmt.Errorf("pcp: parsing host uphost ip: %w", err)
			}
			h.UphostIP = v
		case PCPHostUphostPort:
			v, err := child.GetInt()
			if err != nil {
				return h, fmt.Errorf("pcp: parsing host uphost port: %w", err)
			}
			h.UphostPort = v
		case PCPHostUphostHops:
			v, err := child.GetInt()
			if err != nil {
				return h, fmt.Errorf("pcp: parsing host uphost hops: %w", err)
			}
			h.UphostHops = v
		}
	}
	return h, nil
}

// ---------------------------------------------------------------------------
// BcstPacket
// ---------------------------------------------------------------------------

// BuildAtom serializes a BcstPacket into a PCPBcst parent atom.
func (b *BcstPacket) BuildAtom() *Atom {
	var children []*Atom
	if b.TTL != 0 {
		children = append(children, NewByteAtom(PCPBcstTTL, b.TTL))
	}
	if b.Hops != 0 {
		children = append(children, NewByteAtom(PCPBcstHops, b.Hops))
	}
	if !b.From.IsEmpty() {
		children = append(children, NewIDAtom(PCPBcstFrom, b.From))
	}
	if !b.Dest.IsEmpty() {
		children = append(children, NewIDAtom(PCPBcstDest, b.Dest))
	}
	if b.Group != 0 {
		children = append(children, NewByteAtom(PCPBcstGroup, b.Group))
	}
	if !b.ChanID.IsEmpty() {
		children = append(children, NewIDAtom(PCPBcstChanID, b.ChanID))
	}
	if b.Version != 0 {
		children = append(children, NewIntAtom(PCPBcstVersion, b.Version))
	}
	if b.VersionVP != 0 {
		children = append(children, NewIntAtom(PCPBcstVersionVP, b.VersionVP))
	}
	if b.VersionExPrefix != [2]byte{} {
		children = append(children, NewBytesAtom(PCPBcstVersionExPrefix, b.VersionExPrefix[:]))
	}
	if b.VersionExNumber != 0 {
		children = append(children, NewShortAtom(PCPBcstVersionExNumber, b.VersionExNumber))
	}
	return NewParentAtom(PCPBcst, children...)
}

// ParseBcstPacket extracts a BcstPacket from a PCPBcst parent atom.
func ParseBcstPacket(a *Atom) (BcstPacket, error) {
	var b BcstPacket
	for _, child := range a.Children() {
		switch child.Tag {
		case PCPBcstTTL:
			v, err := child.GetByte()
			if err != nil {
				return b, fmt.Errorf("pcp: parsing bcst ttl: %w", err)
			}
			b.TTL = v
		case PCPBcstHops:
			v, err := child.GetByte()
			if err != nil {
				return b, fmt.Errorf("pcp: parsing bcst hops: %w", err)
			}
			b.Hops = v
		case PCPBcstFrom:
			v, err := child.GetID()
			if err != nil {
				return b, fmt.Errorf("pcp: parsing bcst from: %w", err)
			}
			b.From = v
		case PCPBcstDest:
			v, err := child.GetID()
			if err != nil {
				return b, fmt.Errorf("pcp: parsing bcst dest: %w", err)
			}
			b.Dest = v
		case PCPBcstGroup:
			v, err := child.GetByte()
			if err != nil {
				return b, fmt.Errorf("pcp: parsing bcst group: %w", err)
			}
			b.Group = v
		case PCPBcstChanID:
			v, err := child.GetID()
			if err != nil {
				return b, fmt.Errorf("pcp: parsing bcst chan id: %w", err)
			}
			b.ChanID = v
		case PCPBcstVersion:
			v, err := child.GetInt()
			if err != nil {
				return b, fmt.Errorf("pcp: parsing bcst version: %w", err)
			}
			b.Version = v
		case PCPBcstVersionVP:
			v, err := child.GetInt()
			if err != nil {
				return b, fmt.Errorf("pcp: parsing bcst version vp: %w", err)
			}
			b.VersionVP = v
		case PCPBcstVersionExPrefix:
			if d := child.Data(); len(d) >= 2 {
				copy(b.VersionExPrefix[:], d[:2])
			}
		case PCPBcstVersionExNumber:
			v, err := child.GetShort()
			if err != nil {
				return b, fmt.Errorf("pcp: parsing bcst version ex number: %w", err)
			}
			b.VersionExNumber = v
		}
	}
	return b, nil
}

// ---------------------------------------------------------------------------
// PushPacket
// ---------------------------------------------------------------------------

// BuildAtom serializes a PushPacket into a PCPPush parent atom.
func (p *PushPacket) BuildAtom() *Atom {
	var children []*Atom
	if p.IP != 0 {
		children = append(children, NewIntAtom(PCPPushIP, p.IP))
	}
	if p.Port != 0 {
		children = append(children, NewShortAtom(PCPPushPort, p.Port))
	}
	if !p.ChanID.IsEmpty() {
		children = append(children, NewIDAtom(PCPPushChanID, p.ChanID))
	}
	return NewParentAtom(PCPPush, children...)
}

// ParsePushPacket extracts a PushPacket from a PCPPush parent atom.
func ParsePushPacket(a *Atom) (PushPacket, error) {
	var p PushPacket
	for _, child := range a.Children() {
		switch child.Tag {
		case PCPPushIP:
			v, err := child.GetInt()
			if err != nil {
				return p, fmt.Errorf("pcp: parsing push ip: %w", err)
			}
			p.IP = v
		case PCPPushPort:
			v, err := child.GetShort()
			if err != nil {
				return p, fmt.Errorf("pcp: parsing push port: %w", err)
			}
			p.Port = v
		case PCPPushChanID:
			v, err := child.GetID()
			if err != nil {
				return p, fmt.Errorf("pcp: parsing push chan id: %w", err)
			}
			p.ChanID = v
		}
	}
	return p, nil
}

// ---------------------------------------------------------------------------
// GetPacket
// ---------------------------------------------------------------------------

// BuildAtom serializes a GetPacket into a PCPGet parent atom.
func (g *GetPacket) BuildAtom() *Atom {
	var children []*Atom
	if !g.ID.IsEmpty() {
		children = append(children, NewIDAtom(PCPGetID, g.ID))
	}
	if g.Name != "" {
		children = append(children, NewStringAtom(PCPGetName, g.Name))
	}
	return NewParentAtom(PCPGet, children...)
}

// ParseGetPacket extracts a GetPacket from a PCPGet parent atom.
func ParseGetPacket(a *Atom) (GetPacket, error) {
	var g GetPacket
	for _, child := range a.Children() {
		switch child.Tag {
		case PCPGetID:
			v, err := child.GetID()
			if err != nil {
				return g, fmt.Errorf("pcp: parsing get id: %w", err)
			}
			g.ID = v
		case PCPGetName:
			g.Name = child.GetString()
		}
	}
	return g, nil
}

// ---------------------------------------------------------------------------
// MesgPacket
// ---------------------------------------------------------------------------

// BuildAtom serializes a MesgPacket into a PCPMesg parent atom.
func (m *MesgPacket) BuildAtom() *Atom {
	var children []*Atom
	if m.ASCII != "" {
		children = append(children, NewStringAtom(PCPMesgASCII, m.ASCII))
	}
	if m.SJIS != "" {
		children = append(children, NewStringAtom(PCPMesgSJIS, m.SJIS))
	}
	return NewParentAtom(PCPMesg, children...)
}

// ParseMesgPacket extracts a MesgPacket from a PCPMesg parent atom.
func ParseMesgPacket(a *Atom) (MesgPacket, error) {
	var m MesgPacket
	for _, child := range a.Children() {
		switch child.Tag {
		case PCPMesgASCII:
			m.ASCII = child.GetString()
		case PCPMesgSJIS:
			m.SJIS = child.GetString()
		}
	}
	return m, nil
}
