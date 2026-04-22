package query

import "testing"

func TestToTime(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"valid with ms", "2024-12-04 14:50:07.000", "14:50:07"},
		{"valid no ms", "2024-12-04 14:50:07", "14:50:07"},
		{"empty string", "", ""},
		{"malformed", "not-a-date", ""},
		{"wrong format (iso)", "2024-12-04T14:50:07Z", ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := ToTime(tc.in); got != tc.want {
				t.Errorf("ToTime(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestToPathName(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"1", "vpc"},
		{"7", "vpc endpoint"},
		{"8", "internet gateway"},
		{"99", ""},
		{"", ""},
	}

	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			if got := ToPathName(tc.in); got != tc.want {
				t.Errorf("ToPathName(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}
