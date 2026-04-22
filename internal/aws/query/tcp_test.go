package query

import (
	"reflect"
	"testing"
)

func TestToTcpFlagNames(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want []string
	}{
		{"empty string", "", nil},
		{"dash", "-", nil},
		{"not a number", "abc", nil},
		{"out of range", "256", nil},
		{"zero treated as ACK-only marker", "0", []string{"0"}},
		{"FIN only", "1", []string{"FIN"}},
		{"SYN only", "2", []string{"SYN"}},
		{"ACK only", "16", []string{"ACK"}},
		{"SYN+ACK", "18", []string{"SYN", "ACK"}},
		{"FIN+PSH+ACK", "25", []string{"FIN", "PSH", "ACK"}},
		{"all bits set", "255", []string{"FIN", "SYN", "RST", "PSH", "ACK", "URG", "RESERVED", "RESERVED"}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := ToTcpFlagNames(tc.in)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("ToTcpFlagNames(%q) = %v, want %v", tc.in, got, tc.want)
			}
		})
	}
}
