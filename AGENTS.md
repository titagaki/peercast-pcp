# AI Agent Instructions for peercast-pcp

## 1) Goal and Scope
- Build a robust, zero-dependency, low-level PeerCast Protocol (PCP) library in Go.
- Package name MUST be `pcp`.

## 2) Source of Truth
- Primary: `docs/PCP_SPEC.md`.
- If ambiguous: resolve by reading C++ headers (`docs/reference/pcp.h`, `docs/reference/atom.h`).
- In conflicts, C++ behavior is authoritative.

## 3) Protocol Rules (MUST)
- Endianness: all uint16/uint32 fields are Little Endian (`binary.LittleEndian`).
- Tag width: tags are exactly 4 bytes; short tags are null-padded (0x00), not space-padded.
- Strings (`STR`): null-terminated on wire; append `\0` on encode and trim trailing `\0` on decode.
- Packet format: 8-byte header (`[4]byte` tag + `uint32` length), followed by payload.

## 4) Go Implementation Rules (MUST)
- Use only Go standard library (no third-party dependencies).
- Base read/write logic on `io.Reader` / `io.Writer`; do not assume full-buffer input.
- Parse container tags recursively by treating payload as sub-packet stream.
- Unknown tags must be skipped by length and parsing must continue (forward compatibility).
- Wrap errors with context via `fmt.Errorf("...: %w", err)`.

## 5) Testing Rules
- Prefer table-driven tests for parser/serializer coverage.
- Constructor, roundtrip, and error-path tests may use standalone `TestXxx` style.
- Serialization tests must verify exact bytes.
- Parsing tests must verify decoded fields.
- Cover truncated input, wrong payload size, and EOF paths.