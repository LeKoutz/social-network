package forum

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

func LogDebug(v any) {
	log.Printf("Debug: %#v", v)
}

func convertStringToTime(timeString string) (time.Time, error) {
	timestamp, err := strconv.ParseInt(timeString, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return convertInt64ToTime(timestamp), nil
}

func convertTimeToString(t time.Time) string {
	return t.String()
}

func convertInt64ToTime(i int64) time.Time {
	return time.Unix(i, 0)
}

func getCurrentTimestamp() string {
	return fmt.Sprintf("%d", time.Now().Unix())
}
