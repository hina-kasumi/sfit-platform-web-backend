package converter

import (
	"regexp"
	"strings"
	"time"
)

func ISO8601ToNumber(isoTime string) (time.Duration, error) {
	isoTime = strings.TrimPrefix(isoTime, "PT")

	// Tìm các phần thời gian (giờ, phút, giây)
	re := regexp.MustCompile(`(?:(\d+)H)?(?:(\d+)M)?(?:(\d+)S)?`)
	matches := re.FindStringSubmatch(isoTime)

	var h, m, s string
	if len(matches) == 4 {
		h = matches[1]
		m = matches[2]
		s = matches[3]
	}

	durationStr := ""
	if h != "" {
		durationStr += h + "h"
	}
	if m != "" {
		durationStr += m + "m"
	}
	if s != "" {
		durationStr += s + "s"
	}

	return time.ParseDuration(durationStr)
}
