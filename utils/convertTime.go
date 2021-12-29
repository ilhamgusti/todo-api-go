package utils

import (
	"strings"
	"time"
)

func ConvertTimeISO(value string) time.Time {
	// convert iso-8601 into rfc-3339 format
	//"2015-12-23 00:00:00"
	rfc3339t := strings.Replace(value, " ", "T", 1) + "Z"

	// parse rfc-3339 datetime
	t, err := time.Parse(time.RFC3339, rfc3339t)
	if err != nil {
		panic(err)
	}

	return t
}
