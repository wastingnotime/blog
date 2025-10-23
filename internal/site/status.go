package site

import "time"

func SagaStatus(last *time.Time) string {
	if last == nil || last.IsZero() {
		return "Hiatus"
	}
	days := int(time.Since(*last).Hours() / 24)
	switch {
	case days <= 30:
		return "Now Airing"
	case days <= 90:
		return "Active"
	case days <= 180:
		return "Paused"
	default:
		return "Hiatus"
	}
}
