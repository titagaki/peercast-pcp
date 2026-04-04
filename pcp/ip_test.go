package pcp

import (
	"net"
	"testing"
)

func TestIPv4ToUint32(t *testing.T) {
	tests := []struct {
		ip   net.IP
		want uint32
	}{
		{net.IPv4(192, 168, 1, 1), 0xC0A80101},
		{net.IPv4(10, 0, 0, 1), 0x0A000001},
		{net.IPv4(255, 255, 255, 255), 0xFFFFFFFF},
		{net.IPv4(0, 0, 0, 0), 0},
		{nil, 0},
	}
	for _, tt := range tests {
		got := IPv4ToUint32(tt.ip)
		if got != tt.want {
			t.Errorf("IPv4ToUint32(%v) = 0x%08X, want 0x%08X", tt.ip, got, tt.want)
		}
	}
}

func TestIPv4FromUint32(t *testing.T) {
	tests := []struct {
		v    uint32
		want net.IP
	}{
		{0xC0A80101, net.IP{192, 168, 1, 1}},
		{0x0A000001, net.IP{10, 0, 0, 1}},
		{0, net.IP{0, 0, 0, 0}},
	}
	for _, tt := range tests {
		got := IPv4FromUint32(tt.v)
		if !got.Equal(tt.want) {
			t.Errorf("IPv4FromUint32(0x%08X) = %v, want %v", tt.v, got, tt.want)
		}
	}
}

func TestIPv4Roundtrip(t *testing.T) {
	ip := net.IPv4(172, 16, 254, 3).To4()
	got := IPv4FromUint32(IPv4ToUint32(ip))
	if !got.Equal(ip) {
		t.Errorf("roundtrip: got %v, want %v", got, ip)
	}
}

func TestDecodeIPv4(t *testing.T) {
	// Wire bytes for 192.168.1.1: little-endian of 0xC0A80101 = [01, 01, A8, C0]
	data := []byte{0x01, 0x01, 0xA8, 0xC0}
	got := DecodeIPv4(data)
	want := net.IP{192, 168, 1, 1}
	if !got.Equal(want) {
		t.Errorf("DecodeIPv4(%x) = %v, want %v", data, got, want)
	}
}

func TestDecodeIPv4_InvalidLength(t *testing.T) {
	if got := DecodeIPv4([]byte{1, 2, 3}); got != nil {
		t.Errorf("DecodeIPv4(3 bytes) = %v, want nil", got)
	}
	if got := DecodeIPv4(nil); got != nil {
		t.Errorf("DecodeIPv4(nil) = %v, want nil", got)
	}
}
