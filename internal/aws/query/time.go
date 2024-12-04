package query

import "time"

func ToTime(in string) string {
	// in  - 2024-12-04 14:50:07.000
	// out - 14:50:07
	t, err := time.Parse("2006-01-02 15:04:05.999", in)
	if err != nil {
		return ""
	}
	return t.Format("15:04:05")
}
