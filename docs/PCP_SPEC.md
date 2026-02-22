# PeerCast Protocol (PCP) Specification v1.0

## 1. General Architecture
PCP is a binary protocol used for P2P streaming. It operates over TCP and uses a **Tag-Length-Value (TLV)** packet structure. The protocol is strictly **Little Endian**.

### 1.1 Byte Order
All multi-byte numbers (uint16, uint32) are encoded in **Little Endian**.
- Example: `0x1234` -> `[0x34, 0x12]`

### 1.2 Basic Packet Structure (Header)
Every packet consists of an 8-byte header followed by a variable-length payload.

| Offset | Type     | Name   | Description |
|:-------|:---------|:-------|:------------|
| 0      | [4]byte  | Tag    | ASCII identifier (e.g., "helo", "chan"). Case-sensitive. |
| 4      | uint32   | Length | Length of the Payload in bytes. |
| 8      | []byte   | Payload| Data content (Atoms) or Sub-Packets (Containers). |

---

## 2. Data Types
Definitions of data types used within the Payload.

| Type Name | Go Type | Size | Description |
|:---|:---|:---|:---|
| **BYTE** | `uint8` | 1 | Single byte integer. |
| **SHORT**| `uint16`| 2 | 2-byte integer (Little Endian). |
| **INT** | `uint32`| 4 | 4-byte integer (Little Endian). |
| **ID** | `[16]byte`| 16 | Unique Identifier (GUID/UUID). |
| **STR** | `string`| Var | Null-terminated string (`\0`). Copilot must handle stripping/adding `\0`. |
| **RAW** | `[]byte`| Var | Binary data stream. |
| **PKT** | `[]Packet`| Var | **Container Type**. The payload consists of one or more PCP Packets. |

---

## 3. Protocol Flow & Handshake

### 3.1 Connection Start
1. **Client** sends Magic Bytes: `pcp\n` (`0x70 0x63 0x70 0x0a`).
   - This is treated as a packet with Tag="pcp\n" and Length=0.
2. **Server** validates Magic Bytes.
3. **Both** start exchanging root-level packets (usually starting with `helo` or `oleh`).

---

## 4. Packet Dictionary & Hierarchy
PCP packets are categorized into **Containers** (holding sub-packets) and **Atoms** (holding data).

### 4.1 Root Level Commands
These packets appear at the top level of the stream.

| Tag | Type | Payload Structure | Description |
|:---|:---|:---|:---|
| `helo` | **PKT** | See [4.2 HELO Structure] | Handshake: Node information exchange. |
| `oleh` | **PKT** | Same as `helo` | Handshake Reply (Reverse HELO). |
| `ping` | - | (Empty) | Keep-alive. Length=0. |
| `pong` | - | (Empty) | Reply to ping. Length=0. |
| `quit` | STR | "Reason String" | Disconnect notification. |
| `bcst` | **PKT** | See [4.3 BCST Structure] | Broadcast stream data (Group Header). |
| `push` | **PKT** | See [4.3 BCST Structure] | Relaying stream data. |

---

### 4.2 HELO / OLEH Structure (Handshake)
**Parent Tag:** `helo` or `oleh`
**Description:** Defines the node capabilities and session info.

| Sub-Tag | Type | Go Field Name | Description |
|:---|:---|:---|:---|
| `agnt` | STR | `Agent` | User Agent string. |
| `ostp` | INT | `OsType` | OS Type (Legacy, often unused). |
| `sid ` | ID | `SessionID` | 16-byte random session ID per client. |
| `port` | SHORT| `ListenPort` | The port number this node is listening on. |
| `ping` | SHORT| `PingPort` | Request to ping this specific port. |
| `rip ` | INT | `RemoteIP` | Global IPv4 address (as seen by the peer). |
| `ver ` | INT | `Version` | Client Version (e.g., 1218). |
| `bcid` | ID | `BroadcastID` | ID of the broadcast being relayed. |
| `dis ` | INT | `Disable` | If 1, connection is rejected. |

> **Note:** Tags with 3 letters like `sid` must be padded with a space: `"sid "`.

---

### 4.3 CHAN / INFO Structure (Metadata)
**Parent Tag:** `chan` (often wrapped in `info` or root)
**Description:** Describes channel metadata.

| Sub-Tag | Type | Go Field Name | Description |
|:---|:---|:---|:---|
| `id  ` | ID | `ChannelID` | Unique 16-byte Channel ID. |
| `name` | STR | `Name` | Channel Name. |
| `url ` | STR | `URL` | Related website URL. |
| `desc` | STR | `Description` | Channel description. |
| `genr` | STR | `Genre` | Content genre. |
| `type` | STR | `ContentType` | Stream type (e.g., "WMV", "FLV", "MP3"). |
| `bitr` | INT | `Bitrate` | Bitrate in kbps. |
| `sr  ` | INT | `SampleRate` | Audio sample rate (optional). |

---

### 4.4 BCST / PUSH Structure (Streaming)
**Parent Tag:** `bcst`
**Description:** Transports actual media data.

| Sub-Tag | Type | Go Field Name | Description |
|:---|:---|:---|:---|
| `id  ` | ID | `ChannelID` | Target Channel ID. |
| `pos ` | INT | `Position` | Stream position / sequence number. |
| `data` | RAW | `Data` | The actual media binary data. |

---

### 4.5 HOST Structure (Network Discovery)
**Parent Tag:** `host` (usually inside `lsth` list)

| Sub-Tag | Type | Go Field Name | Description |
|:---|:---|:---|:---|
| `ip  ` | INT | `IP` | IPv4 Address. |
| `port` | SHORT| `Port` | Port number. |
| `id  ` | ID | `NodeID` | Node ID. |
| `num ` | INT | `NumListeners`| Number of connected listeners. |
| `uptm` | INT | `UpTime` | Uptime seconds. |

---

## 5. Implementation Requirements for Copilot

1.  **Recursive Parsing**:
    - The parser MUST identify "Container Tags" (e.g., `helo`, `chan`, `bcst`).
    - For these tags, the `Payload` bytes must be treated as a new `io.Reader` and parsed recursively to extract Sub-Packets.

2.  **Unknown Tags**:
    - The library MUST handle unknown tags gracefully by skipping `Length` bytes.
    - Do not error on unknown tags; store them as generic `RawPacket` or discard them.

3.  **Strings**:
    - When writing `STR` type, append a null byte (`0x00`).
    - When reading `STR` type, read until the end of payload and trim the null byte.

4.  **Tag Padding**:
    - Tags are fixed 4 bytes. If a tag name is "id", it MUST be encoded as `{'i', 'd', ' ', ' '}`.

5.  **Go Struct Tags**:
    - Use struct tags to map PCP tags to fields.
    - Example:
      ```go
      type HeloPacket struct {
          Agent     string    `pcp:"agnt"`
          SessionID [16]byte  `pcp:"sid "`
          Port      uint16    `pcp:"port"`
      }
      ```