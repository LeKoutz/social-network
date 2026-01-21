package utils

import (
	"errors"
	"fmt"
	"log"
	"runtime"
	"strconv"
	"time"
)

func LogDebug(v any) {
	log.Printf("Debug: %#v", v)
}

func ConvertStringToTime(timeString string) (time.Time, error) {
	timestamp, err := strconv.ParseInt(timeString, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return ConvertInt64ToTime(timestamp), nil
}

func ConvertTimeToString(t time.Time) string {
	return t.String()
}

func ConvertInt64ToTime(i int64) time.Time {
	return time.Unix(i, 0)
}

func GetCurrentTimestamp() string {
	return fmt.Sprintf("%d", time.Now().Unix())
}

func GetFunctionName() error {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return errors.New("")
	}
	f := runtime.FuncForPC(pc)
	if f == nil {
		return errors.New("")
	}
	return errors.New(f.Name())
}
