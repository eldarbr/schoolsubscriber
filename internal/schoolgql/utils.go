package schoolgql

import "time"

func FormatTimeToStr(t time.Time) string {
	t = t.UTC()

	return t.Format(time.RFC3339)
}

func FormatStrToTime(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s) //nolint:wrapcheck // too verbal.
}
