package pcp

// Packet struct definitions for each PCP message type.
//
// These structs provide a typed view over the atom tree. The `pcp` struct
// tag indicates the wire-format tag name used when encoding/decoding.
//
// All structures are derived from the C++ source (pcp.h / atom.h) as the
// authoritative reference. Fields present in PCP_SPEC.md but absent from
// the C++ source are omitted. Fields present in C++ but absent from the
// spec are included (marked as such).

// ---------------------------------------------------------------------------
// HELO / OLEH — Handshake
// ---------------------------------------------------------------------------

// HeloPacket represents a "helo" or "oleh" handshake.
// Both directions use the same sub-atom structure.
type HeloPacket struct {
	Agent     string `pcp:"agnt"` // User-Agent string
	OsType    uint32 `pcp:"ostp"` // OS type identifier
	SessionID GnuID  `pcp:"sid"`  // 16-byte random session ID
	Port      uint16 `pcp:"port"` // Listen port
	Ping      uint16 `pcp:"ping"` // Request to ping this port
	Pong      uint16 `pcp:"pong"` // Ping reply port (not in spec)
	RemoteIP  uint32 `pcp:"rip"`  // Global IPv4 address as seen by peer
	Version   uint32 `pcp:"ver"`  // Client version number
	BCID      GnuID  `pcp:"bcid"` // Broadcast ID
	Disable   uint32 `pcp:"dis"`  // If 1, connection is rejected
}

// ---------------------------------------------------------------------------
// ROOT — Root server information
// ---------------------------------------------------------------------------

// RootPacket represents information from or about the root server.
type RootPacket struct {
	UpdateInterval uint32 `pcp:"uint"` // Update interval in seconds
	CheckVersion   uint32 `pcp:"chkv"` // Version to check against
	URL            string `pcp:"url"`  // Root server URL
	Update         []byte `pcp:"upd"`  // Update data
	Next           []byte `pcp:"next"` // Next root information
}

// ---------------------------------------------------------------------------
// CHAN — Channel
// ---------------------------------------------------------------------------

// ChanPacket represents a "chan" container with channel state.
type ChanPacket struct {
	ID          GnuID        `pcp:"id"`   // 16-byte Channel ID
	BroadcastID GnuID        `pcp:"bcid"` // Broadcast ID (not in spec)
	Key         []byte       `pcp:"key"`  // Channel key (not in spec)
	Info        *ChanInfo    `pcp:"info"` // Channel metadata (container)
	Track       *ChanTrack   `pcp:"trck"` // Current track info (container, not in spec)
	Pkt         *ChanPktData `pcp:"pkt"`  // Stream data packet (container)
}

// ---------------------------------------------------------------------------
// CHAN > INFO — Channel metadata
// ---------------------------------------------------------------------------

// ChanInfo represents channel metadata inside a "chan" > "info" container.
type ChanInfo struct {
	Type       string `pcp:"type"` // Content type (e.g. "WMV", "FLV", "MP3")
	StreamType string `pcp:"styp"` // Stream MIME type (not in spec)
	StreamExt  string `pcp:"sext"` // Stream file extension (not in spec)
	Bitrate    uint32 `pcp:"bitr"` // Bitrate in kbps
	Genre      string `pcp:"gnre"` // Genre (note: wire tag is "gnre", not "genr" as in spec)
	Name       string `pcp:"name"` // Channel name
	URL        string `pcp:"url"`  // Related website URL
	Desc       string `pcp:"desc"` // Channel description
	Comment    string `pcp:"cmnt"` // Comment (not in spec)
}

// ---------------------------------------------------------------------------
// CHAN > TRACK — Current track
// ---------------------------------------------------------------------------

// ChanTrack represents track metadata inside a "chan" > "trck" container.
// This entire structure is absent from PCP_SPEC.md but present in C++.
type ChanTrack struct {
	Title   string `pcp:"titl"` // Track title
	Creator string `pcp:"crea"` // Track creator / artist
	URL     string `pcp:"url"`  // Track URL
	Album   string `pcp:"albm"` // Album name
}

// ---------------------------------------------------------------------------
// CHAN > PKT — Stream data packet
// ---------------------------------------------------------------------------

// ChanPktData represents stream data inside a "chan" > "pkt" container.
type ChanPktData struct {
	Type         ID4    `pcp:"type"` // Packet type (head/data/meta)
	Pos          uint32 `pcp:"pos"`  // Stream position / sequence number
	Head         []byte `pcp:"head"` // Stream header data
	Data         []byte `pcp:"data"` // Stream payload data
	Meta         []byte `pcp:"meta"` // Metadata
	Continuation []byte `pcp:"cont"` // Continuation data
}

// ---------------------------------------------------------------------------
// HOST — Network discovery
// ---------------------------------------------------------------------------

// HostPacket represents a node in the P2P network, usually found
// inside a host list container.
//
// PCP_SPEC.md only documents ip, port, id, num (listeners), and uptm.
// The C++ implementation has many additional fields listed below.
type HostPacket struct {
	ID              GnuID  `pcp:"id"`   // 16-byte Node ID
	IP              uint32 `pcp:"ip"`   // IPv4 address
	Port            uint16 `pcp:"port"` // Port number
	NumListeners    uint32 `pcp:"numl"` // Number of connected listeners
	NumRelays       uint32 `pcp:"numr"` // Number of relays (not in spec)
	Uptime          uint32 `pcp:"uptm"` // Uptime in seconds
	Tracker         uint32 `pcp:"trkr"` // Tracker flag (not in spec)
	ChanID          GnuID  `pcp:"cid"`  // Channel ID (not in spec)
	Version         uint32 `pcp:"ver"`  // Version number (not in spec)
	VersionVP       uint32 `pcp:"vevp"` // VP version (not in spec)
	VersionExPrefix ID4    `pcp:"vexp"` // Version extension prefix (not in spec)
	VersionExNumber uint32 `pcp:"vexn"` // Version extension number (not in spec)
	Flags1          uint32 `pcp:"flg1"` // Host capability flags (not in spec)
	OldPos          uint32 `pcp:"oldp"` // Old stream position (not in spec)
	NewPos          uint32 `pcp:"newp"` // New stream position (not in spec)
	UphostIP        uint32 `pcp:"upip"` // Upstream host IP (not in spec)
	UphostPort      uint16 `pcp:"uppt"` // Upstream host port (not in spec)
	UphostHops      uint32 `pcp:"uphp"` // Upstream host hops (not in spec)
}

// ---------------------------------------------------------------------------
// BCST — Broadcast
// ---------------------------------------------------------------------------

// BcstPacket represents a broadcast container that transports channel
// updates across the P2P network.
//
// PCP_SPEC.md describes bcst as carrying id/pos/data directly, but in
// the actual C++ implementation, bcst wraps routing metadata and embeds
// other containers (chan, host, etc.) as children.
type BcstPacket struct {
	TTL             byte   `pcp:"ttl"`  // Time-to-live (decremented per hop)
	Hops            byte   `pcp:"hops"` // Number of hops traversed
	From            GnuID  `pcp:"from"` // Source node ID
	Dest            GnuID  `pcp:"dest"` // Destination node ID (empty = broadcast)
	Group           byte   `pcp:"grp"`  // Target group (PCPBcstGroup* constants)
	ChanID          GnuID  `pcp:"cid"`  // Channel ID
	Version         uint32 `pcp:"vers"` // Broadcaster version
	VersionVP       uint32 `pcp:"vrvp"` // VP version
	VersionExPrefix ID4    `pcp:"vexp"` // Version extension prefix
	VersionExNumber uint32 `pcp:"vexn"` // Version extension number
}

// ---------------------------------------------------------------------------
// PUSH — Push relay request
// ---------------------------------------------------------------------------

// PushPacket requests a node to initiate a connection back to the
// requester for relaying.
type PushPacket struct {
	IP     uint32 `pcp:"ip"`   // Target IPv4 address
	Port   uint16 `pcp:"port"` // Target port
	ChanID GnuID  `pcp:"cid"`  // Channel ID to relay
}

// ---------------------------------------------------------------------------
// GET — Data request
// ---------------------------------------------------------------------------

// GetPacket represents a "get" request for channel data.
type GetPacket struct {
	ID   GnuID  `pcp:"id"`   // Target ID
	Name string `pcp:"name"` // Target name
}

// ---------------------------------------------------------------------------
// MESG — Message
// ---------------------------------------------------------------------------

// MesgPacket represents a text message.
// ASCII and SJIS formats are deprecated in favor of UTF-8.
type MesgPacket struct {
	ASCII string `pcp:"asci"` // ASCII / Shift_JIS message (deprecated)
	SJIS  string `pcp:"sjis"` // Shift_JIS message (deprecated)
}

// ---------------------------------------------------------------------------
// BroadcastState — Runtime state during broadcast processing
// ---------------------------------------------------------------------------

// BroadcastState tracks the transient state of a broadcast exchange,
// mirroring the C++ BroadcastState class used during atom processing.
type BroadcastState struct {
	ChanID    GnuID  // Channel being processed
	BCID      GnuID  // Broadcast session ID
	NumHops   int    // Accumulated hop count
	ForMe     bool   // True if this broadcast is targeted at this node
	StreamPos uint32 // Current stream byte position
	Group     int    // Target group bitmask
}

// InitPacketSettings resets transient per-packet state fields,
// matching the C++ BroadcastState::initPacketSettings() method.
func (bs *BroadcastState) InitPacketSettings() {
	bs.ForMe = false
	bs.Group = 0
	bs.NumHops = 0
	bs.BCID.Clear()
	bs.ChanID.Clear()
}
