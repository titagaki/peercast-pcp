package pcp

// Numeric constants extracted from pcp.h.
// These are true Go constants (int/uint types).

// ---------------------------------------------------------------------------
// Broadcast group flags — used in PCPBcstGroup ("grp") field.
// Multiple groups can be combined with bitwise OR.
// ---------------------------------------------------------------------------

const (
	PCPBcstGroupAll      = 0xff // All groups
	PCPBcstGroupRoot     = 1    // Root server
	PCPBcstGroupTrackers = 2    // Tracker nodes
	PCPBcstGroupRelays   = 4    // Relay nodes
)

// ---------------------------------------------------------------------------
// Error categories — base values for error classification.
// The actual error code is category + specific error.
// ---------------------------------------------------------------------------

const (
	PCPErrorQuit    = 1000
	PCPErrorBcst    = 2000
	PCPErrorRead    = 3000
	PCPErrorWrite   = 4000
	PCPErrorGeneral = 5000
)

// ---------------------------------------------------------------------------
// Specific error codes — added to an error category base to form
// the complete error value sent in a "quit" atom.
// Example: PCPErrorQuit + PCPErrorBanned = 1011
// ---------------------------------------------------------------------------

const (
	PCPErrorSkip             = 1  // Skip this node
	PCPErrorAlreadyConnected = 2  // Already connected to this node
	PCPErrorUnavailable      = 3  // Resource unavailable
	PCPErrorLoopback         = 4  // Loopback connection detected
	PCPErrorNotIdentified    = 5  // Node not identified
	PCPErrorBadResponse      = 6  // Bad response received
	PCPErrorBadAgent         = 7  // Bad user agent
	PCPErrorOffAir           = 8  // Channel is off-air
	PCPErrorShutdown         = 9  // Node is shutting down
	PCPErrorNoRoot           = 10 // No root server available
	PCPErrorBanned           = 11 // Node is banned
)

// ---------------------------------------------------------------------------
// Host flags (first byte) — used in PCPHostFlags1 ("flg1") field.
// Multiple flags can be combined with bitwise OR.
// ---------------------------------------------------------------------------

const (
	PCPHostFlags1Tracker = 0x01 // Node is a tracker
	PCPHostFlags1Relay   = 0x02 // Node is relaying
	PCPHostFlags1Direct  = 0x04 // Node accepts direct connections
	PCPHostFlags1Push    = 0x08 // Node requires push connections
	PCPHostFlags1Recv    = 0x10 // Node is receiving
	PCPHostFlags1CIN     = 0x20 // Node is control-in
	PCPHostFlags1Private = 0x40 // Node is on a private network
)

// ---------------------------------------------------------------------------
// Atom header constants
// ---------------------------------------------------------------------------

const (
	// AtomHeaderSize is the fixed size of an atom header (4-byte tag + 4-byte length/count).
	AtomHeaderSize = 8

	// AtomParentBit is the MSB of the length field, indicating a container (parent) atom.
	// When set, the lower 31 bits represent the number of child atoms.
	AtomParentBit = 0x80000000

	// AtomParentMask extracts the child count from a parent atom's length field.
	AtomParentMask = 0x7FFFFFFF

	// MaxAtomDataSize is the maximum payload size (in bytes) that ReadAtom
	// and SkipAtom will accept for a data atom. This guards against
	// excessive memory allocation from untrusted streams. (16 MiB)
	MaxAtomDataSize = 16 * 1024 * 1024
)
