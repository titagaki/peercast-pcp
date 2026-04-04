package pcp

import (
	"encoding/binary"
	"net"
)

// IPv4ToUint32 converts a net.IP (IPv4) to the uint32 representation used
// in PCP atom fields (e.g. PCPHostIP, PCPHeloRemoteIP).
//
// The returned value has the first octet in the most significant byte
// (network byte order as an integer). For example, 192.168.1.1 → 0xC0A80101.
//
// This uint32 is intended to be written with NewIntAtom, which stores it
// as little-endian on the wire — matching the PCP wire convention.
//
// Returns 0 if ip is nil or not a valid IPv4 address.
func IPv4ToUint32(ip net.IP) uint32 {
	ip4 := ip.To4()
	if ip4 == nil {
		return 0
	}
	return binary.BigEndian.Uint32(ip4)
}

// IPv4FromUint32 converts a uint32 (as returned by Atom.GetInt on an IP field)
// back to a net.IP.
//
// The uint32 is interpreted with the first octet in the most significant byte.
// For example, 0xC0A80101 → 192.168.1.1.
func IPv4FromUint32(v uint32) net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, v)
	return ip
}

// DecodeIPv4 converts raw PCP wire bytes (from Atom.Data) directly to a net.IP
// without going through GetInt. The 4 bytes are in little-endian uint32 order,
// so they are reversed to produce a standard net.IP.
//
// Returns nil if data is not exactly 4 bytes.
func DecodeIPv4(data []byte) net.IP {
	if len(data) != 4 {
		return nil
	}
	return net.IP{data[3], data[2], data[1], data[0]}
}
