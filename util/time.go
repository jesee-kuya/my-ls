package util

import "time"

func FormatTime(t time.Time) string {
	return t.Format("Jan _2 15:04")
}
