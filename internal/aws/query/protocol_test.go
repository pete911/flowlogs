package query

import (
	"strconv"
	"testing"
)

func TestProtocolFromKeywordToNumber(t *testing.T) {
	tests := []struct {
		in   string
		want int
	}{
		{"TCP", 6},
		{"tcp", 6},
		{"Tcp", 6},
		{"UDP", 17},
		{"ICMP", 1},
		{"HOPOPT", 0},
		{"not-a-protocol", -1},
	}

	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			if got := protocolFromKeywordToNumber(tc.in); got != tc.want {
				t.Errorf("protocolFromKeywordToNumber(%q) = %d, want %d", tc.in, got, tc.want)
			}
		})
	}
}

func TestProtocolFromNumberToKeyword(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"TCP", "6", "TCP"},
		{"UDP", "17", "UDP"},
		{"ICMP", "1", "ICMP"},
		{"reserved", "255", "Reserved"},
		{"dash passes through", "-", "-"},
		{"empty passes through", "", ""},
		{"out of range returns input", "999", "999"},
		{"non-numeric returns input", "abc", "abc"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := ProtocolFromNumberToKeyword(tc.in); got != tc.want {
				t.Errorf("ProtocolFromNumberToKeyword(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestProtocolRoundTrip(t *testing.T) {
	for _, kw := range []string{"TCP", "UDP", "ICMP"} {
		n := protocolFromKeywordToNumber(kw)
		if n < 0 {
			t.Errorf("%q: not found", kw)
			continue
		}
		back := ProtocolFromNumberToKeyword(strconv.Itoa(n))
		if back != kw {
			t.Errorf("round-trip %q -> %d -> %q", kw, n, back)
		}
	}
}
