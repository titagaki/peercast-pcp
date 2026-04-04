package pcp

// PCP tag identifiers extracted from pcp.h.
//
// Go does not support const for array types ([4]byte), so these are
// defined as package-level variables. They should be treated as
// immutable constants — never reassign them.
//
// Several tag variables share the same wire value (e.g. "id", "port",
// "ip", "type", "cid", "vexp", "vexn"). They are disambiguated by their
// parent container context on the wire. The distinct Go variable names
// exist for code readability only.
//
// ID4 values are comparable with == and can be used directly in
// switch statements:
//
//	switch atom.Tag {
//	case PCPHelo:
//	    // handle helo
//	case PCPOleh:
//	    // handle oleh
//	}

// ---------------------------------------------------------------------------
// Connection
// ---------------------------------------------------------------------------

var (
	// PCPConnect is the magic bytes sent at connection start: "pcp\n"
	// (0x70 0x63 0x70 0x0a). Sent as an atom with Length=0.
	PCPConnect = NewID4("pcp\n")

	// PCPOK is the OK response atom.
	PCPOK = NewID4("ok")
)

// ---------------------------------------------------------------------------
// HELO / OLEH — Handshake
// ---------------------------------------------------------------------------

var (
	PCPHelo          = NewID4("helo") // Handshake request (container)
	PCPHeloAgent     = NewID4("agnt") // STR: User-Agent string
	PCPHeloOsType    = NewID4("ostp") // INT: OS type identifier
	PCPHeloSessionID = NewID4("sid")  // ID:  16-byte session ID
	PCPHeloPort      = NewID4("port") // SHORT: listen port
	PCPHeloPing      = NewID4("ping") // SHORT: request to ping this port
	PCPHeloPong      = NewID4("pong") // SHORT: ping reply port
	PCPHeloRemoteIP  = NewID4("rip")  // INT: remote IP as seen by peer
	PCPHeloVersion   = NewID4("ver")  // INT: client version number
	PCPHeloBCID      = NewID4("bcid") // ID:  broadcast ID
	PCPHeloDisable   = NewID4("dis")  // INT: if 1, connection rejected

	PCPOleh = NewID4("oleh") // Handshake reply (container, same sub-structure as helo)
)

// ---------------------------------------------------------------------------
// MODE
// ---------------------------------------------------------------------------

var (
	PCPMode       = NewID4("mode") // Mode container
	PCPModeGnut06 = NewID4("gn06") // Gnutella 0.6 mode
)

// ---------------------------------------------------------------------------
// ROOT — Root server information
// ---------------------------------------------------------------------------

var (
	PCPRoot         = NewID4("root") // Root container
	PCPRootUpdInt   = NewID4("uint") // INT: update interval (seconds)
	PCPRootCheckVer = NewID4("chkv") // INT: check version
	PCPRootURL      = NewID4("url")  // STR: root URL
	PCPRootUpdate   = NewID4("upd")  // Update data
	PCPRootNext     = NewID4("next") // Next root
)

// ---------------------------------------------------------------------------
// OS types
// ---------------------------------------------------------------------------

var (
	PCPOSLinux   = NewID4("lnux") // Linux
	PCPOSWindows = NewID4("w32")  // Windows
	PCPOSOSX     = NewID4("osx")  // macOS
	PCPOSWinamp  = NewID4("wamp") // Winamp
	PCPOSZaurus  = NewID4("zaur") // Sharp Zaurus
)

// ---------------------------------------------------------------------------
// GET — Request data
// ---------------------------------------------------------------------------

var (
	PCPGet     = NewID4("get")  // Get container
	PCPGetID   = NewID4("id")   // ID:  target ID
	PCPGetName = NewID4("name") // STR: target name
)

// ---------------------------------------------------------------------------
// HOST — Network discovery (usually inside a host list)
// ---------------------------------------------------------------------------

var (
	PCPHost                = NewID4("host") // Host container
	PCPHostID              = NewID4("id")   // ID:  node ID
	PCPHostIP              = NewID4("ip")   // INT: IPv4 address
	PCPHostPort            = NewID4("port") // SHORT: port number
	PCPHostNumListeners    = NewID4("numl") // INT: number of listeners
	PCPHostNumRelays       = NewID4("numr") // INT: number of relays
	PCPHostUptime          = NewID4("uptm") // INT: uptime in seconds
	PCPHostTracker         = NewID4("trkr") // INT: tracker flag
	PCPHostChanID          = NewID4("cid")  // ID:  channel ID
	PCPHostVersion         = NewID4("ver")  // INT: version
	PCPHostVersionVP       = NewID4("vevp") // INT: VP version
	PCPHostVersionExPrefix = NewID4("vexp") // RAW[2]: version extension prefix (2 ASCII bytes)
	PCPHostVersionExNumber = NewID4("vexn") // SHORT: version extension number
	PCPHostFlags1          = NewID4("flg1") // BYTE: host capability flags (see PCPHostFlags1* constants)
	PCPHostOldPos          = NewID4("oldp") // INT: oldest available stream position
	PCPHostNewPos          = NewID4("newp") // INT: newest stream position
	PCPHostUphostIP        = NewID4("upip") // INT or RAW[16]: upstream host IP
	PCPHostUphostPort      = NewID4("uppt") // INT: upstream host port
	PCPHostUphostHops      = NewID4("uphp") // INT: number of hops to upstream host
)

// ---------------------------------------------------------------------------
// QUIT
// ---------------------------------------------------------------------------

var (
	PCPQuit = NewID4("quit") // STR: disconnect with reason
)

// ---------------------------------------------------------------------------
// CHAN — Channel
// ---------------------------------------------------------------------------

var (
	PCPChan     = NewID4("chan") // Channel container
	PCPChanID   = NewID4("id")  // ID:  channel ID
	PCPChanBCID = NewID4("bcid") // ID:  broadcast ID
	PCPChanKey  = NewID4("key")  // RAW: channel key
)

// ---------------------------------------------------------------------------
// CHAN > PKT — Channel packets (stream data)
// ---------------------------------------------------------------------------

var (
	PCPChanPkt             = NewID4("pkt")  // Packet container
	PCPChanPktType         = NewID4("type") // ID4: packet type identifier
	PCPChanPktPos          = NewID4("pos")  // INT: stream position / sequence
	PCPChanPktData         = NewID4("data") // RAW: stream payload data
	PCPChanPktContinuation = NewID4("cont") // BYTE: continuation flag (non-zero if packet continues previous)
)

// ---------------------------------------------------------------------------
// CHAN > INFO — Channel metadata
// ---------------------------------------------------------------------------

var (
	PCPChanInfo           = NewID4("info") // Info container
	PCPChanInfoType       = NewID4("type") // STR: content type (e.g. "WMV", "FLV")
	PCPChanInfoStreamType = NewID4("styp") // STR: stream MIME type
	PCPChanInfoStreamExt  = NewID4("sext") // STR: stream file extension
	PCPChanInfoBitrate    = NewID4("bitr") // INT: bitrate in kbps
	PCPChanInfoGenre      = NewID4("gnre") // STR: genre (note: C++ uses "gnre", not "genr")
	PCPChanInfoName       = NewID4("name") // STR: channel name
	PCPChanInfoURL        = NewID4("url")  // STR: related website URL
	PCPChanInfoDesc       = NewID4("desc") // STR: channel description
	PCPChanInfoComment    = NewID4("cmnt") // STR: comment
)

// ---------------------------------------------------------------------------
// CHAN > TRACK — Current track information
// ---------------------------------------------------------------------------

var (
	PCPChanTrack        = NewID4("trck") // Track container
	PCPChanTrackTitle   = NewID4("titl") // STR: track title
	PCPChanTrackCreator = NewID4("crea") // STR: track creator / artist
	PCPChanTrackURL     = NewID4("url")  // STR: track URL
	PCPChanTrackAlbum   = NewID4("albm") // STR: album name
)

// ---------------------------------------------------------------------------
// MESG — Messages
// ---------------------------------------------------------------------------

var (
	PCPMesg      = NewID4("mesg") // Message container
	PCPMesgASCII = NewID4("asci") // STR: ASCII/SJIS message (deprecated)
	PCPMesgSJIS  = NewID4("sjis") // STR: Shift_JIS message (deprecated)
)

// ---------------------------------------------------------------------------
// BCST — Broadcast
// ---------------------------------------------------------------------------

var (
	PCPBcst                = NewID4("bcst") // Broadcast container
	PCPBcstTTL             = NewID4("ttl")  // BYTE: time-to-live
	PCPBcstHops            = NewID4("hops") // BYTE: hop count
	PCPBcstFrom            = NewID4("from") // ID:   source node ID
	PCPBcstDest            = NewID4("dest") // ID:   destination node ID
	PCPBcstGroup           = NewID4("grp")  // BYTE: target group (see PCPBcstGroup* constants)
	PCPBcstChanID          = NewID4("cid")  // ID:   channel ID
	PCPBcstVersion         = NewID4("vers") // INT:  version
	PCPBcstVersionVP       = NewID4("vrvp") // INT:    VP version
	PCPBcstVersionExPrefix = NewID4("vexp") // RAW[2]: version extension prefix (2 ASCII bytes)
	PCPBcstVersionExNumber = NewID4("vexn") // SHORT:  version extension number
)

// ---------------------------------------------------------------------------
// PUSH — Push relay request
// ---------------------------------------------------------------------------

var (
	PCPPush       = NewID4("push") // Push container
	PCPPushIP     = NewID4("ip")   // INT:   target IP
	PCPPushPort   = NewID4("port") // SHORT: target port
	PCPPushChanID = NewID4("cid")  // ID:    channel ID
)

// ---------------------------------------------------------------------------
// Miscellaneous
// ---------------------------------------------------------------------------

var (
	PCPSpkt      = NewID4("spkt") // Server packet
	PCPAtom      = NewID4("atom") // Generic atom
	PCPSessionID = NewID4("sid")  // ID: session ID (root-level)

	// Ping/Pong at root level (Length=0, empty payload).
	PCPPing = NewID4("ping")
	PCPPong = NewID4("pong")
)
