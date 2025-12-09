package forum

import (
	"fmt"
	"strconv"
	"time"
)

func convertStringToTime(timeString string) (time.Time, error) {
	timestamp, err := strconv.ParseInt(timeString, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	t := time.Unix(timestamp, 0)
	return t, nil
}

func getCurrentTimestamp() string {
	return fmt.Sprintf("%d", time.Now().Unix())
}
